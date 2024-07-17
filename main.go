package main

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/urfave/cli/v3"
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

// help - urfave/cli
func main() {
	config, err := InitConfig()

	if err != nil {
		log.Fatal("Configuration problem", err)
		return
	}

	var prompt string

	cmd := &cli.Command{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "prompt",
				Value:       "commit",
				Destination: &prompt,
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			if prompt == "review" {
				prompt = getPrompt(PromptCodeReview)
			} else {
				prompt = getPrompt(PromptCommitMessage)
			}
			return nil
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
		return
	}

	// Run the git diff command
	stdout, stderr, err := runDiff()
	if err != nil {
		fmt.Println("Error:", err)
		fmt.Println("Stderr:", stderr)

		return
	}

	result := fmt.Sprintf("~~~diff\n%s~~~\n\n%s", stdout, prompt)

	// Print the output
	fmt.Println(result)

	content, err := ExecPrompt(result, config)
	if err != nil {
		log.Fatal("Prompt execution problem", err)
		return
	}

	fmt.Println(content.Choices[0].Message.Content)
}
