package gitgen

import (
	"fmt"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

type PromptType int

const (
	PromptCommitMessage PromptType = iota
	PromptCodeReview
)

func runDiff(config Config) (string, error) {
	repo, err := git.PlainOpen(".")
	if err != nil {
		return "", err
	}

	srcRef, err := repo.ResolveRevision(plumbing.Revision(config.SourceRef))
	if err != nil {
		return "", err
	}

	var destRef *plumbing.Hash
	if config.DestinationRef != "" {
		destRef, err = repo.ResolveRevision(plumbing.Revision(config.DestinationRef))
	} else {
		destRef, err = repo.ResolveRevision(plumbing.Revision("HEAD"))
	}
	if err != nil {
		return "", err
	}

	srcCommit, err := repo.CommitObject(*srcRef)
	if err != nil {
		return "", err
	}

	destCommit, err := repo.CommitObject(*destRef)
	if err != nil {
		return "", err
	}

	patch, err := destCommit.Patch(srcCommit)
	if err != nil {
		return "", err
	}

	return patch.String(), nil
}

func getPrompt(promptType PromptType) string {
	var prompt string
	switch promptType {
	case PromptCommitMessage:
		prompt = "please generate a git commit message with a simple explanation from the changes stated above which is an output of a git diff command. all response of this message should be wrapped in a markdown format because it will be shared in a text-only terminal interface."

	case PromptCodeReview:
		prompt = "please perform a efficient and concise code review which points out crucial improvements could be changed on the target code. the target code is stated above which is an output of a git diff command. all response of this message should be wrapped in a markdown format because it will be shared in a text-only terminal interface."

	}

	return prompt
}

func Do(promptType PromptType, config Config) (string, error) {
	systemPrompt := getPrompt(promptType)

	// Run the git diff command
	userPrompt, err := runDiff(config)
	if err != nil {
		return "", err
	}

	fmt.Println("System Prompt:")
	fmt.Println(systemPrompt)
	fmt.Println("User Prompt:")
	fmt.Println(userPrompt)

	content, err := ExecPrompt(systemPrompt, userPrompt, config)
	if err != nil {
		return "", err
	}

	fmt.Println("OpenAI Response:")
	return content.Choices[0].Message.Content, nil
}
