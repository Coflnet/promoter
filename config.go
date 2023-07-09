package main

import (
	"os"

	_ "github.com/joho/godotenv/autoload"
	"github.com/rs/zerolog/log"
)

var config Config

type Config struct {
	GitRepository     string
	GitUsername       string
	GitToken          string
	Filename          string
	NewTag            string
	ImageName         string
	RepositoryFolder  string
	RepositoryFolders map[string]string
}

func Init() {

	ReadEnvVars()
}

func ReadEnvVars() {
	config.GitRepository = os.Getenv("GIT_REPOSITORY")
	if config.GitRepository == "" {
		log.Panic().Msgf("GIT_REPOSITORY env var is not set")
	}
	log.Info().Msgf("using git repository: %s", config.GitRepository)

	config.GitUsername = os.Getenv("GIT_USERNAME")
	if config.GitUsername == "" {
		log.Panic().Msgf("GIT_USERNAME env var is not set")
	}
	log.Info().Msgf("using username: %s", config.GitUsername)

	config.GitToken = os.Getenv("GIT_TOKEN")
	if config.GitToken == "" {
		log.Panic().Msgf("GIT_TOKEN is not set")
	}

	config.Filename = os.Getenv("FILENAME")
	if config.Filename == "" {
		log.Panic().Msgf("FILENAME env var is not set")
	}
	log.Info().Msgf("using filename: %s", config.Filename)

	config.NewTag = os.Getenv("NEW_TAG")
	if config.NewTag == "" {
		log.Panic().Msgf("NEW_TAG env var is not set")
	}
	log.Info().Msgf("using tag: %s", config.NewTag)

	config.ImageName = os.Getenv("IMAGE_NAME")
	if config.ImageName == "" {
		log.Panic().Msgf("IMAGE_NAME env var is not set")
	}
	log.Info().Msgf("using image name: %s", config.ImageName)
}
