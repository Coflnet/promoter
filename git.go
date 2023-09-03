package main

import (
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
)

var repo map[string]*git.Repository

func CloneRepostories() error {
	repo = make(map[string]*git.Repository)
	config.RepositoryFolders = make(map[string]string)

	repos := strings.Split(config.GitRepository, ",")

	for i, repo := range repos {
		config.RepositoryFolders[fmt.Sprintf("dir-%d", i)] = repo
	}

	for key, url := range config.RepositoryFolders {
		err := CloneRepository(key, url)
		if err != nil {
			return err
		}
	}

	return nil
}

func CloneRepository(folder string, repositoryUrl string) error {
	var err error
	username := config.GitUsername
	token := config.GitToken
	auth := &http.BasicAuth{Username: username, Password: token}

	// check if the repository folder exists and delete it
	if _, err := os.Stat(folder); !os.IsNotExist(err) {
		slog.Warn("delete repository folder because it already exists")
		err = os.RemoveAll(folder)
		if err != nil {
			return err
		}
	}

	repository, err := git.PlainClone(folder, false, &git.CloneOptions{
		URL:      repositoryUrl,
		Auth:     auth,
		Progress: os.Stdout,
	})

	if err != nil {
		return err
	}

	repo[folder] = repository

	return nil
}

func PipelineStart() time.Time {
	startString := os.Getenv("PIPELINE_START")

	// format as RFC3339
	start, err := time.Parse(time.RFC3339, startString)

	if err != nil {
		return time.Time{}
	}

	slog.Debug("pipeline start calculated", "start", start.String())
	return start
}

func PushEnvs() error {
	for folder, _ := range config.RepositoryFolders {
		err := PushEnv(folder)

		if err != nil {
			return err
		}
	}
	return nil
}

func PushEnv(folder string) error {
	repository := repo[folder]

	if repository == nil {
		return fmt.Errorf("the repository variable is nil")
	}

	worktree, err := repository.Worktree()
	if err != nil {
		return fmt.Errorf("error when getting the worktree, %w", err)
	}

	stateChanged := stateChanged(worktree)
	if !stateChanged {
		slog.Debug("the state of the repository has not changed")
		return nil
	}

	worktree.Pull(&git.PullOptions{})

	worktree.Add(".")

	_, err = worktree.Commit(fmt.Sprintf("[CI] promote %s", config.Filename), &git.CommitOptions{
		All: true,
		Committer: &object.Signature{
			Name:  "coflnet-bot",
			Email: "ci@coflnet.com",
			When:  time.Now(),
		},
		Author: &object.Signature{
			Name:  "coflnet-bot",
			Email: "ci@coflnet.com",
			When:  time.Now(),
		},
	})

	if err != nil {
		return err
	}

	username := config.GitUsername
	token := config.GitToken
	auth := &http.BasicAuth{Username: username, Password: token}

	err = repository.Push(&git.PushOptions{
		RemoteName: "origin",
		Progress:   os.Stdout,
		Auth:       auth,
	})

	if err != nil {
		return err
	}

	return nil
}

func stateChanged(worktree *git.Worktree) bool {
	status, err := worktree.Status()

	if err != nil {
		return false
	}

	return !status.IsClean()
}
