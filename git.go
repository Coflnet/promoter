package main

import (
	"fmt"
	"os"
	"time"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/rs/zerolog/log"
)

var repository *git.Repository

func CloneRepository() error {

	var err error
	username := config.GitUsername
	token := config.GitToken
	url := config.GitRepository
	auth := &http.BasicAuth{Username: username, Password: token}
	config.RepositoryFolder = "k8s"

	repository, err = git.PlainClone(config.RepositoryFolder, false, &git.CloneOptions{
		URL:      url,
		Auth:     auth,
		Progress: os.Stdout,
	})

	if err != nil {
		log.Fatal().Err(err).Msgf("can not clone the kube repository")
		return err
	}

	return nil
}

func PushEnv() error {
	if repository == nil {
		log.Fatal().Msgf("the repository variable is nil")
	}

	worktree, err := repository.Worktree()
	if err != nil {
		log.Fatal().Err(err).Msgf("error when getting the worktree")
		return err
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
		log.Fatal().Err(err).Msgf("something went wrong when committing")
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
		log.Fatal().Err(err).Msgf("something went wrong when pushing")
		return err
	}

	return nil
}
