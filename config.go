package main

import (
	"fmt"
	"log/slog"
	"os"

	_ "github.com/joho/godotenv/autoload"
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
	config.GitRepository = mustReadEnv("GIT_REPOSITORY")
	slog.Info(fmt.Sprintf("using git repository: %s", config.GitRepository))

	config.GitUsername = mustReadEnv("GIT_USERNAME")
	slog.Info(fmt.Sprintf("using username: %s", config.GitUsername))

	config.GitToken = mustReadEnv("GIT_TOKEN")
    slog.Info(fmt.Sprintf("using token: %s", config.GitToken))

	config.Filename = mustReadEnv("FILENAME")
    slog.Info(fmt.Sprintf("using filename: %s", config.Filename))

	config.NewTag = mustReadEnv("NEW_TAG")
    slog.Info(fmt.Sprintf("using new tag: %s", config.NewTag))

	config.ImageName = mustReadEnv("IMAGE_NAME")
    slog.Info(fmt.Sprintf("using image name: %s", config.ImageName))
}

func mustReadEnv(key string) string {
    value := os.Getenv(key)
    if value == "" {
        panic(fmt.Sprintf("env var %s is not set", key))
    }
    return value
}

