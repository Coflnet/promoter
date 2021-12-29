package main

import (
	"os"

	_ "github.com/joho/godotenv/autoload"
	"github.com/rs/zerolog/log"
)

var config Config

type Config struct {
	GitRepository    string
	GitUsername      string
	GitToken         string
	Filename         string
	NewTag           string
	ImageName        string
	RepositoryFolder string
}

func Init() {

	ReadEnvVars()
}

func ReadEnvVars() {
	config.GitRepository = os.Getenv("GIT_REPOSITORY")
	if config.GitRepository == "" {
		log.Fatal().Msgf("GIT_REPOSITORY env var is not set")
	}

	config.GitUsername = os.Getenv("GIT_USERNAME")
	if config.GitUsername == "" {
		log.Fatal().Msgf("GIT_USERNAME env var is not set")
	}

	config.GitToken = os.Getenv("GIT_TOKEN")
	if config.GitToken == "" {
		log.Fatal().Msgf("GIT_TOKEN is not set")
	}

	config.Filename = os.Getenv("FILENAME")
	if config.Filename == "" {
		log.Fatal().Msgf("FILENAME env var is not set")
	}

	config.NewTag = os.Getenv("NEW_TAG")
	if config.NewTag == "" {
		log.Fatal().Msgf("NEW_TAG env var is not set")
	}

	config.ImageName = os.Getenv("IMAGE_NAME")
	if config.ImageName == "" {
		log.Fatal().Msgf("IMAGE_NAME env var is not set")
	}
}
