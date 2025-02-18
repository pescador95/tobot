package services

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/google/go-github/v50/github"
)

type GitHubService struct {
	Token string
}

func NewGitHubService(token string) *GitHubService {
	return &GitHubService{Token: token}
}

func (s *GitHubService) GetLatestBuildStatus(username, repo string) (string, error) {
	url := "https://api.github.com/repos/" + username + "/" + repo + "/actions/runs"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "token "+s.Token)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	var result map[string]interface{}
	json.Unmarshal(body, &result)

	if runs, ok := result["workflow_runs"].([]interface{}); ok && len(runs) > 0 {
		if conclusion, ok := runs[0].(map[string]interface{})["conclusion"].(string); ok {
			return conclusion, nil
		}
	}

	return "Nenhum status disponível", nil
}

func (s *GitHubService) GetBranches(username, repo string) ([]*github.Branch, error) {
	url := "https://api.github.com/repos/" + username + "/" + repo + "/branches"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "token "+s.Token)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	var branches []*github.Branch
	if err := json.Unmarshal(body, &branches); err != nil {
		return nil, err
	}

	return branches, nil
}

func (s *GitHubService) GetIssues(username, repo string) ([]*github.Issue, error) {
	url := "https://api.github.com/repos/" + username + "/" + repo + "/issues"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "token "+s.Token)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	var issues []*github.Issue
	if err := json.Unmarshal(body, &issues); err != nil {
		return nil, err
	}

	return issues, nil
}

func (s *GitHubService) GetPullRequests(username, repo string) ([]*github.PullRequest, error) {
	url := "https://api.github.com/repos/" + username + "/" + repo + "/pulls"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "token "+s.Token)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	var pullRequests []*github.PullRequest
	if err := json.Unmarshal(body, &pullRequests); err != nil {
		return nil, err
	}

	return pullRequests, nil
}

func (s *GitHubService) GetCommits(username, repo, branch string) ([]*github.RepositoryCommit, error) {
	url := "https://api.github.com/repos/" + username + "/" + repo + "/commits?sha=" + branch
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "token "+s.Token)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	var commits []*github.RepositoryCommit
	if err := json.Unmarshal(body, &commits); err != nil {
		return nil, err
	}

	return commits, nil
}
