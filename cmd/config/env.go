package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnv() error {
	err := godotenv.Load()
	if err != nil {
		return fmt.Errorf("erro ao carregar o arquivo .env: %v", err)
	}
	return nil
}

func GetDiscordToken() string {
	return os.Getenv("DISCORD_TOKEN")
}

func GetGitHubToken() string {
	return os.Getenv("GITHUB_TOKEN")
}

func GetChannelId() string {
	return os.Getenv("DISCORD_CHANNEL_ID")
}

func GetEnv(key string) string {
	return os.Getenv(key)
}
