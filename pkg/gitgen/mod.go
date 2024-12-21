package gitgen

import (
	"context"
	"log"
	"os"
	"os/exec"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/seymahandekli/git-gen/pkg/platforms"
)

//go:generate stringer -type=ActionType
type ActionType int

const (
	ActionCommitMessage ActionType = iota
	ActionCodeReview
	ActionTestScenario
	ActionTest
)

func runDiffOnCli(config Config) (string, error) {
	// Define the Git command
	cmdArgs := []string{
		"diff",
		"--patch",
		"--minimal",
		"--diff-algorithm=minimal",
		"--ignore-all-space",
		"--ignore-blank-lines",
		"--no-ext-diff",
		"--no-color",
		"--unified=10",
		config.SourceRef,
	}
	if config.DestinationRef != "" {
		cmdArgs = append(cmdArgs, config.DestinationRef)
	}

	cmd := exec.Command("git", cmdArgs...)

	// cmd.Env = os.Environ()

	// var newEnv []string
	// for _, e := range cmd.Env {
	// 	if e[:18] != "GIT_EXTERNAL_DIFF=" {
	// 		newEnv = append(newEnv, e)
	// 	}
	// }
	// cmd.Env = newEnv

	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	// Convert the output to a string
	return string(output), nil
}

func runDiffWithGoGit(config Config) (string, error) {
	workingDir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	repo, err := git.PlainOpenWithOptions(workingDir, &git.PlainOpenOptions{DetectDotGit: true})
	if err != nil {
		return "", err
	}

	srcRefName := plumbing.ReferenceName(config.SourceRef)
	if err := srcRefName.Validate(); err != nil {
		return "", err
	}
	srcRef, err := repo.Reference(srcRefName, true)
	if err != nil {
		return "", err
	}
	srcCommit, err := repo.CommitObject(srcRef.Hash())
	if err != nil {
		return "", err
	}
	srcTree, err := srcCommit.Tree()
	if err != nil {
		return "", err
	}

	var destRef *plumbing.Reference

	if config.DestinationRef != "" {
		destRefName := plumbing.ReferenceName(config.DestinationRef)
		if err := destRefName.Validate(); err != nil {
			return "", err
		}
		destRef, err = repo.Reference(destRefName, true)
		if err != nil {
			return "", err
		}
	} else {
		destRef, err = repo.Storer.Reference(plumbing.HEAD)
		if err != nil {
			return "", err
		}
	}

	destCommit, err := repo.CommitObject(destRef.Hash())
	if err != nil {
		return "", err
	}
	destTree, err := destCommit.Tree()
	if err != nil {
		return "", err
	}

	patch, err := destTree.Diff(srcTree)
	if err != nil {
		return "", err
	}

	return patch.String(), nil
}

func Do(actionType ActionType, config Config) (string, error) {
	// Run the git diff command
	diff, err := runDiffOnCli(config)
	if err != nil {
		return "", err
	}

	platformConfig := platforms.PlatformConfig{
		ApiKey:                      config.PlatformApiKey,
		Model:                       config.Model,
		PromptMaxTokens:             config.PromptMaxTokens,
		PromptRequestTimeoutSeconds: config.PromptRequestTimeoutSeconds,
	}

	// Convert string to Platform type
	platform := platforms.Platform(config.Platform)
	runtime, err := platforms.NewPromptExecutor(platform, platformConfig)
	if err != nil {
		return "", err
	}
	prompt := &Prompt{
		ActionType:  actionType,
		Diff:        diff,
		Attachments: []PromptAttachment{},
	}

	log.Printf("System Prompt:\n%s\n\n", prompt.GetSystemPrompt())
	// log.Printf("User Prompt:\n%s\n\n", prompt.GetUserPrompt())
	log.Printf("User Prompt Length:\n%d\n\n", len(prompt.GetUserPrompt()))

	response, err := runtime.ExecPrompt(context.Background(), prompt)
	if err != nil {
		return "", err
	}

	return response.Content, nil
}
