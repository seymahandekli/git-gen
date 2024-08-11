package gitgen

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"

	_ "embed"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"

	"github.com/seymahandekli/git-gen/pkg/models"
)

//go:generate stringer -type=PromptType
type PromptType int

var (
	//go:embed prompts/commit-message.txt
	PromptForCommit string

	//go:embed prompts/code-review.txt
	PromptForCodeReview string

	//go:embed prompts/test-case.txt
	PromptForTestCase string
)

const (
	PromptCommitMessage PromptType = iota
	PromptCodeReview
	PromptTestCase
)

var (
	ErrUnknownPlatform = errors.New("unknown platform")
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

func GetPrompt(promptType PromptType) string {
	if promptType == PromptCommitMessage {
		return PromptForCommit
	}

	if promptType == PromptCodeReview {
		return PromptForCodeReview
	}

	return PromptForTestCase
}

func Do(promptType PromptType, config Config) (string, error) {
	systemPrompt := GetPrompt(promptType)

	// Run the git diff command
	userPrompt, err := runDiffOnCli(config)
	if err != nil {
		return "", err
	}

	log.Printf("System Prompt:\n%s\n\n", systemPrompt)
	// log.Printf("User Prompt:\n%s\n\n", userPrompt)
	log.Printf("User Prompt Length:\n%d\n\n", len(userPrompt))

	modelConfig := models.ModelConfig{
		PlatformApiKey:              config.PlatformApiKey,
		Platform:                    config.Platform,
		Model:                       config.Model,
		PromptMaxTokens:             config.PromptMaxTokens,
		PromptRequestTimeoutSeconds: config.PromptRequestTimeoutSeconds,
	}

	var runtime models.Model

	switch modelConfig.Platform {
	case "openai":
		runtime = models.NewOpenAi(modelConfig)
	case "ollama":
		runtime, err = models.NewOllamaAi(modelConfig)

		if err != nil {
			return "", err
		}
	default:
		return "", fmt.Errorf("unknown platform %s - %w", modelConfig.Platform, ErrUnknownPlatform)
	}

	response, err := runtime.ExecPrompt(context.Background(), systemPrompt, userPrompt)
	if err != nil {
		return "", err
	}

	return response.Content, nil
}
