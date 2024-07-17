package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"

	"github.com/urfave/cli/v3"
)

const (
	apiEndpoint = "https://api.openai.com/v1/chat/completions"
)

type PromptType int

const (
	PromptCommitMessage PromptType = iota
	PromptCodeReview
)

func runDiff() (error, string, string) {
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
		return err, "", ""
	}

	// Convert the output to a string
	return nil, stdout.String(), stderr.String()
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
	err, stdout, stderr := runDiff()
	if err != nil {
		fmt.Println("Error:", err)
		fmt.Println("Stderr:", stderr)

		return
	}

	result := fmt.Sprintf("~~~diff\n%s~~~\n\n%s", stdout, prompt)

	// Print the output
	fmt.Println(result)

	apiKey, apiOk := os.LookupEnv("OPENAI_API_KEY")
	if !apiOk {
		log.Fatal("Environment variable OPENAI_API_KEY must be set")
		return
	}

	// Create the request body
	body, err := json.Marshal(map[string]interface{}{
		"model":      "gpt-4o",
		"messages":   []interface{}{map[string]interface{}{"role": "system", "content": result}},
		"max_tokens": 3500,
	})
	if err != nil {
		log.Fatalf("Error while marshalling request body: %v", err)
		return
	}

	// Create the HTTP request
	req, err := http.NewRequest("POST", apiEndpoint, bytes.NewBuffer(body))
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	// Send the request
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error while sending request: %v", err)
		return
	}
	defer res.Body.Close()

	var data map[string]interface{}

	err = json.NewDecoder(res.Body).Decode(&data)
	if err != nil {
		fmt.Println("Error while decoding JSON response:", err)
		return
	}

	// Extract the content from the JSON response
	content := data
	fmt.Println(content)
}

