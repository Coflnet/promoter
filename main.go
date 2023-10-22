package main

import (
	"log/slog"
	"os"
)

func main() {
	SetupDebugLogLevel()
	ReadEnvVars()

	err := CloneRepostories()
	if err != nil {
		slog.Error("there was problem when cloning, stop execution", "err", err)
		os.Exit(1)
	}
	slog.Info("successfully cloned the repositories")

	err = Promote()
	if err != nil {
		slog.Error("there was problem when promoting, stop execution", "err", err)
		os.Exit(1)
	}
	slog.Info("successfully promoted the helm chart")

	err = PushEnvs()
	if err != nil {
		slog.Error("there was problem when pushing, stop execution", "err", err)
		os.Exit(1)
	}

	slog.Info("successfully promoted the helm chart")
}

func SetupDebugLogLevel() {
	level := slog.LevelInfo
	if os.Getenv("DEBUG") == "true" {
		level = slog.LevelDebug
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
	}))

	slog.SetDefault(logger)
}
