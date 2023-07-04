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

	// check if the repository folder exists and delete it
	if _, err := os.Stat(config.RepositoryFolder); !os.IsNotExist(err) {
		log.Warn().Msgf("delete repository folder because it already exists")
		err = os.RemoveAll(config.RepositoryFolder)
		if err != nil {
			log.Panic().Err(err).Msgf("could not delete repository folder")
		}
	}

	repository, err = git.PlainClone(config.RepositoryFolder, false, &git.CloneOptions{
		URL:      url,
		Auth:     auth,
		Progress: os.Stdout,
	})

	if err != nil {
		log.Panic().Err(err).Msgf("can not clone the kube repository")
		return err
	}

	commits, err := repository.CommitObjects()
	if err != nil {
		log.Panic().Err(err).Msgf("can not get the commits")
		return err
	}

	pipelineStart := PipelineStart()
	err = commits.ForEach(func(commit *object.Commit) error {
		// get timestamp of commit
		timestamp := commit.Committer.When

		if timestamp.After(pipelineStart) {
			return fmt.Errorf("commit %s is newer than pipeline start, fail", commit.Hash.String())
		}

		return nil
	})

	if err != nil {
		log.Panic().Err(err).Msgf("error while checking commits")
		return err
	}

	log.Info().Msgf("have not found a commit that is newer than the pipeline start, continue")

	return nil
}

func PipelineStart() time.Time {
	startString := os.Getenv("PIPELINE_START")

	// format as RFC3339
	start, err := time.Parse(time.RFC3339, startString)

	if err != nil {
		return time.Time{}
	}

	log.Info().Msgf("pipeline start: %s", start.String())
	return start
}

func PushEnv() error {
	if repository == nil {
		log.Panic().Msgf("the repository variable is nil")
	}

	worktree, err := repository.Worktree()
	if err != nil {
		log.Panic().Err(err).Msgf("error when getting the worktree")
		return err
	}

	worktree.Pull(&git.PullOptions{})

	worktree.Add(".")

	// _, err = worktree.Commit(fmt.Sprintf("[CI] promote %s", config.Filename), &git.CommitOptions{
	// 	All: true,
	// 	Committer: &object.Signature{
	// 		Name:  "coflnet-bot",
	// 		Email: "ci@coflnet.com",
	// 		When:  time.Now(),
	// 	},
	// 	Author: &object.Signature{
	// 		Name:  "coflnet-bot",
	// 		Email: "ci@coflnet.com",
	// 		When:  time.Now(),
	// 	},
	// })

	// if err != nil {
	// 	log.Panic().Err(err).Msgf("something went wrong when committing")
	// 	return err
	// }

	username := config.GitUsername
	token := config.GitToken
	auth := &http.BasicAuth{Username: username, Password: token}

	err = repository.Push(&git.PushOptions{
		RemoteName: "origin",
		Progress:   os.Stdout,
		Auth:       auth,
	})

	if err != nil {
		log.Panic().Err(err).Msgf("something went wrong when pushing")
		return err
	}

	return nil
}
