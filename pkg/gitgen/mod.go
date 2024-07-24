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

func runDiff(config Config) (string, string, error) {
	// Define the Git command

	cmd := exec.Command("git", "diff", config.SourceRef, config.DestinationRef)

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

func Do(promptType PromptType, config Config) (string, error) {
	systemPrompt := getPrompt(promptType)

	// Run the git diff command
	userPrompt, _, err := runDiff(config)
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
