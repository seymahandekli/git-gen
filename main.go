package main

import (
	"bytes"
	"fmt"
	"os/exec"
)

func main() {
	// Define the Git command
	cmd := exec.Command("git", "diff")

	// Create buffers to capture the output and error
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	// Run the command
	err := cmd.Run()
	if err != nil {
		fmt.Println("Error:", err)
		fmt.Println("Stderr:", stderr.String())
		return
	}

	// Convert the output to a string
	output := out.String()
	prompt := "\nplease generate a git commit message with a simple explanation from the changes stated above which is an output of a git diff command. all response of this message should be wrapped in a markdown format because it will be shared in a text-only terminal interface."

	// Print the output
	fmt.Println(output,prompt)
}
