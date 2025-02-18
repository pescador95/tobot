package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/pescador95/go/cmd/api/services"
	"github.com/pescador95/go/cmd/config"
)

func main() {

	err := config.LoadEnv()
	if err != nil {
		fmt.Println(err)
		return
	}
	discordToken := config.GetDiscordToken()
	if discordToken == "" {
		fmt.Println("DISCORD_TOKEN não encontrado no arquivo .env")
		return
	}

	dg, err := discordgo.New("Bot " + discordToken)
	if err != nil {
		fmt.Println("Erro ao criar a sessão do Discord,", err)
		return
	}

	gitHubToken := config.GetGitHubToken()
	if gitHubToken == "" {
		fmt.Println("GITHUB_TOKEN não encontrado no arquivo .env")
		return
	}

	githubService := services.NewGitHubService(gitHubToken)

	discordService := services.NewDiscordService(dg, githubService)

	discordService.RegisterHandlers()

	err = dg.Open()
	if err != nil {
		fmt.Println("Erro ao abrir a conexão,", err)
		return
	}

	fmt.Println("Bot está rodando. Pressione CTRL-C para sair.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	dg.Close()
}
