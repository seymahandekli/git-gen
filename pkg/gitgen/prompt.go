package gitgen

import (
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var (
	//go:embed prompts/commit-message.txt
	PromptForCommit string

	//go:embed prompts/code-review.txt
	PromptForCodeReview string

	//go:embed prompts/test-scenario.txt
	PromptForTestScenario string

	//go:embed prompts/test.txt
	PromptForTest string
)

type PromptAttachment struct {
	Filename string
	Type     string
	Content  string
}

type Prompt struct {
	ActionType ActionType

	Diff        string
	Attachments []PromptAttachment
}

func (p *Prompt) GetSystemPrompt() string {
	if p.ActionType == ActionCommitMessage {
		return PromptForCommit
	}

	if p.ActionType == ActionCodeReview {
		return PromptForCodeReview
	}

	if p.ActionType == ActionTestScenario {
		return PromptForTestScenario
	}

	return PromptForTest
}

func (p *Prompt) GetUserPrompt() string {
	var builder strings.Builder

	// Add file contents to prompt
	for _, attachment := range p.Attachments {
		builder.WriteString(fmt.Sprintf("\n<file>\n%s\n```%s\n%s\n```\n</file>\n",
			attachment.Filename,
			attachment.Type,
			attachment.Content))
	}

	builder.WriteString(fmt.Sprintf("\n<diff>\n```\n%s\n```\n</diff>\n", p.Diff))

	return builder.String()
}

func (p *Prompt) Attach(filepath string) error {
	content, err := readFileContent(filepath)
	if err != nil {
		return err
	}

	attachment := PromptAttachment{
		Filename: filepath,
		Type:     getFileExtension(filepath),
		Content:  content,
	}

	p.Attachments = append(p.Attachments, attachment)

	return nil
}

// readFileContent reads and returns the content of a file
func readFileContent(path string) (string, error) {
	content, err := os.ReadFile(path)

	if err != nil {
		return "", err
	}

	return string(content), nil
}

// getFileExtension returns the file extension without the dot
func getFileExtension(path string) string {
	ext := filepath.Ext(path)

	if ext != "" {
		return ext[1:] // Remove the leading dot
	}

	return ""
}
