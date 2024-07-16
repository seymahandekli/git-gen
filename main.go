package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/go-resty/resty/v2"
)

const (
	apiEndpoint = "https://api.openai.com/v1/chat/completions"
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

func main() {
	// Run the git diff command
	err, stdout, stderr := runDiff()
	if err != nil {
		fmt.Println("Error:", err)
		fmt.Println("Stderr:", stderr)

		return
	}

	prompt := "please generate a git commit message with a simple explanation from the changes stated above which is an output of a git diff command. all response of this message should be wrapped in a markdown format because it will be shared in a text-only terminal interface."
	result := fmt.Sprintf("~~~diff\n%s~~~\n\n%s", stdout, prompt)

	// Print the output
	fmt.Println(result)

	apiKey, apiOk := os.LookupEnv("OPENAI_API_KEY")
	if !apiOk {
		log.Fatal("Environment variable OPENAI_API_KEY must be set")
		return
	}

	client := resty.New()

	response, err := client.R().
		SetAuthToken(apiKey).
		SetHeader("Content-Type", "application/json").
		SetBody(map[string]interface{}{
			"model":      "gpt-4o",
			"messages":   []interface{}{map[string]interface{}{"role": "system", "content": result}},
			"max_tokens": 50,
		}).
		Post(apiEndpoint)

	if err != nil {
		log.Fatalf("Error while sending send the request: %v", err)
	}

	body := response.Body()

	var data map[string]interface{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		fmt.Println("Error while decoding JSON response:", err)
		return
	}

	// Extract the content from the JSON response
	content := data["choices"].([]interface{})[0].(map[string]interface{})["message"].(map[string]interface{})["content"].(string)
	fmt.Println(content)
}
