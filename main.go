package main

import "log/slog"

func main() {
	ReadEnvVars()

	err := CloneRepostories()
	if err != nil {
		slog.Error("there was problem when cloning, stop execution", "err", err)
		panic(err)
	}
	slog.Info("successfully cloned the repositories")

	err = Promote()
	if err != nil {
		slog.Error("there was problem when promoting, stop execution", "err", err)
		panic(err)
	}
	slog.Info("successfully promoted the helm chart")

	err = PushEnvs()
	if err != nil {
		slog.Error("there was problem when pushing, stop execution", "err", err)
	}
	slog.Info("successfully promoted the helm chart")
}
