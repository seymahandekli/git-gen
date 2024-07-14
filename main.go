package main

import (
	"bytes"
	"fmt"
	"os/exec"
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
}
