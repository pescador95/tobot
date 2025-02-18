package services

import (
	"fmt"
	"os"
	"strings"

	"github.com/bwmarrin/discordgo"
)

type DiscordService struct {
	Session       *discordgo.Session
	GitHubService *GitHubService
	ChannelID     string
}

type UserInfo struct {
	UserID    string `json:userId`
	ChannelId string `json:channelId`
}

func NewDiscordService(session *discordgo.Session, githubService *GitHubService) *DiscordService {
	return &DiscordService{
		Session:       session,
		GitHubService: githubService,
		ChannelID:     os.Getenv("DISCORD_CHANNEL_ID"),
	}
}

func (s *DiscordService) RegisterHandlers() {
	s.Session.AddHandler(s.messageCreate)
}

func (s *DiscordService) messageCreate(session *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == session.State.User.ID {
		return
	}

	if m.ChannelID != s.ChannelID {
		return
	}

	handleMessage(session, m, s)
}

func handleMessage(session *discordgo.Session, m *discordgo.MessageCreate, s *DiscordService) {
	args := strings.Split(m.Content, " ")
	if len(args) < 1 {
		sendChannel(session, m, " Comando inválido. Uso: !<comando> [parâmetros]")
		return
	}

	prefix := args[0]
	switch prefix {
	case "!msg":
		handleMsg(session, m, s, args)
	case "!whoami":
		handleWhoAmI(session, m, s, args)
	case "!buildstatus":
		handleBuildStatus(session, m, s, args)
	case "!branches":
		handleBranches(session, m, s, args)
	case "!issues":
		handleIssues(session, m, s, args)
	case "!pullrequests":
		handlePullRequests(session, m, s, args)
	case "!commits":
		handleCommits(session, m, s, args)
	case "!comandos":
		handleCommands(session, m, s)
	default:
		sendChannel(session, m, "Comando desconhecido. Digite !comandos para verificar os comandos disponíveis...")
	}
}

func handleMsg(session *discordgo.Session, m *discordgo.MessageCreate, s *DiscordService, args []string) {
	if len(m.Mentions) == 0 {
		session.ChannelMessageSend(m.ChannelID, "Por favor, mencione um usuário para enviar a mensagem.")
		sendChannel(session, m, "Uso: !msg <suamensagem> <@userMention>")
		return
	}

	if len(args) < 3 {
		sendChannel(session, m, "Uso: !msg <suamensagem> <@userMention>")
		return
	}

	userID := m.Mentions[0].ID

	channel, err := session.UserChannelCreate(userID)
	if err != nil {
		fmt.Println("Erro ao criar canal de DM:", err)
		return
	}

	message := strings.Join(args[1:len(args)-1], " ")

	fullMessage := message

	var userInfo = &UserInfo{
		UserID:    userID,
		ChannelId: channel.ID,
	}

	session.ChannelMessageSend(userInfo.ChannelId, fullMessage)
}

func handleWhoAmI(session *discordgo.Session, m *discordgo.MessageCreate, s *DiscordService, args []string) {
	userID := m.Author.ID
	channel, err := session.UserChannelCreate(userID)

	if err != nil {
		return
	}

	message := "userID: " + userID + "\n" + "channelId: " + channel.ID

	var userInfo = &UserInfo{
		UserID:    userID,
		ChannelId: channel.ID,
	}

	session.ChannelMessageSend(userInfo.ChannelId, message)
}

func handleCommands(session *discordgo.Session, m *discordgo.MessageCreate, s *DiscordService) {
	sendChannel(session, m, "Comandos disponíveis: !buildstatus, !branches,!issues, !pullrequests, !commits, !whoami, !msg")
}

func handleBuildStatus(session *discordgo.Session, m *discordgo.MessageCreate, s *DiscordService, args []string) {
	if len(args) != 3 {
		sendChannel(session, m, "Uso: !buildstatus <username> <repo>")
		return
	}

	username := args[1]
	repo := args[2]
	status, err := s.GitHubService.GetLatestBuildStatus(username, repo)
	if err != nil {
		sendReplyChannel(session, m, "Erro ao obter status da build.")
		return
	}

	sendReplyChannel(session, m, "Status da última build: "+status)
}

func handleBranches(session *discordgo.Session, m *discordgo.MessageCreate, s *DiscordService, args []string) {
	if len(args) != 3 {
		sendChannel(session, m, "Uso: !branches <username> <repo>")
		return
	}

	username := args[1]
	repo := args[2]
	branches, err := s.GitHubService.GetBranches(username, repo)
	if err != nil {
		sendReplyChannel(session, m, "Erro ao obter branches.")
		return
	}

	var branchesList []string
	for _, branch := range branches {
		branchesList = append(branchesList, *branch.Name)
	}
	sendReplyChannel(session, m, "Branches: "+strings.Join(branchesList, "\n "))
}

func handleIssues(session *discordgo.Session, m *discordgo.MessageCreate, s *DiscordService, args []string) {
	if len(args) != 3 {
		sendChannel(session, m, "Uso: !issues <username> <repo>")
		return
	}

	username := args[1]
	repo := args[2]
	issues, err := s.GitHubService.GetIssues(username, repo)
	if err != nil {
		sendReplyChannel(session, m, "Erro ao obter issues.")
		return
	}

	var issuesList []string
	for _, issue := range issues {
		issuesList = append(issuesList, *issue.Title)
	}
	sendReplyChannel(session, m, "Issues: "+strings.Join(issuesList, ", "))
}

func handlePullRequests(session *discordgo.Session, m *discordgo.MessageCreate, s *DiscordService, args []string) {
	if len(args) != 3 {
		sendChannel(session, m, "Uso: !pullrequests <username> <repo>")
		return
	}

	username := args[1]
	repo := args[2]
	pullRequests, err := s.GitHubService.GetPullRequests(username, repo)
	if err != nil {
		sendReplyChannel(session, m, "Erro ao obter pull requests.")
		return
	}

	var prsList []string
	for _, pr := range pullRequests {
		prsList = append(prsList, *pr.Title)
	}
	sendReplyChannel(session, m, "Pull Requests: "+strings.Join(prsList, ", "))
}

func handleCommits(session *discordgo.Session, m *discordgo.MessageCreate, s *DiscordService, args []string) {
	if len(args) != 4 {
		sendChannel(session, m, "Uso: !commits <username> <repo> <branch>")
		return
	}

	username := args[1]
	repo := args[2]
	branch := args[3]
	commits, err := s.GitHubService.GetCommits(username, repo, branch)
	if err != nil {
		sendReplyChannel(session, m, "Erro ao obter commits.")
		return
	}

	var commitsList []string
	for _, commit := range commits {
		author := "Author: " + *commit.Commit.Author.Name
		date := "Date: " + commit.Commit.Author.Date.Format("2006-01-02 15:04:05")
		message := "Message: " + *commit.Commit.Message
		commitsList = append(commitsList, " \n "+author+"\n"+date+" \n "+message+" \n ")
	}

	sendReplyChannel(session, m, "Commits: "+strings.Join(commitsList, "\n "))
}

func getUserMention(m *discordgo.MessageCreate) string {
	return "<@" + m.Author.ID + ">\n"
}

func getUserMentionById(id string) string {
	return "<@" + id + ">\n"
}

func sendDM(session *discordgo.Session, m *discordgo.MessageCreate, message string) {
	if message == "" {
		return
	}
	userID := m.Author.ID
	channel, err := session.UserChannelCreate(userID)
	if err != nil {
		fmt.Println("Erro ao criar canal de DM:", err)
		return
	}
	session.ChannelMessageSend(channel.ID, message)
}

func sendChannel(session *discordgo.Session, m *discordgo.MessageCreate, message string) {
	if message == "" {
		return
	}
	userMention := getUserMention(m)
	session.ChannelMessageSend(m.ChannelID, userMention+message)
}

func sendReplyChannel(session *discordgo.Session, m *discordgo.MessageCreate, message string) {
	if message == "" {
		return
	}

	userMention := getUserMention(m)

	msg := &discordgo.MessageSend{
		Content: userMention + " " + message,
		Reference: &discordgo.MessageReference{
			MessageID: m.ID,
			ChannelID: m.ChannelID,
			GuildID:   m.GuildID,
		},
	}

	session.ChannelMessageSendComplex(m.ChannelID, msg)
}
