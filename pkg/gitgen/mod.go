package gitgen

import (
	"bytes"
	"fmt"
	"os/exec"
)

type PromptType int

const (
	PromptCommitMessage PromptType = iota
	PromptCodeReview
)

func runDiff() (string, string, error) {
	// Define the Git command
	cmd := exec.Command("git", "diff", "HEAD")

	// Create buffers to capture the output and error
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// Run the command
	err := cmd.Run()
	if err != nil {
		return "", "", err
	}

	// Convert the output to a string
	return stdout.String(), stderr.String(), nil
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

func Do(promptType PromptType, maxTokens int64) (string, error) {
	config, err := InitConfig()
	config.PromptMaxTokens = maxTokens

	if err != nil {
		return "", err
	}

	prompt := getPrompt(promptType)

	// Run the git diff command
	stdout, _, err := runDiff()
	if err != nil {
		return "", err
	}

	result := fmt.Sprintf("~~~diff\n%s~~~\n\n%s", stdout, prompt)

	content, err := ExecPrompt(result, config)
	if err != nil {
		return "", err
	}

	return content.Choices[0].Message.Content, nil
}
