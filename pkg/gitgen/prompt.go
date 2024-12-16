package gitgen

import _ "embed"

var (
	//go:embed prompts/commit-message.txt
	PromptForCommit string

	//go:embed prompts/code-review.txt
	PromptForCodeReview string

	//go:embed prompts/test-case.txt
	PromptForTestCase string
)

type Prompt struct {
	systemPrompt string
	userPrompt   string
}

func GetPromptFor(actionType ActionType) string {
	if actionType == ActionCommitMessage {
		return PromptForCommit
	}

	if actionType == ActionCodeReview {
		return PromptForCodeReview
	}

	return PromptForTestCase
}

func NewPrompt(actionType ActionType) *Prompt {
	return &Prompt{
		systemPrompt: GetPromptFor(actionType),
		userPrompt:   "",
	}
}

func (p *Prompt) SetUserPrompt(userPrompt string) {
	p.userPrompt = userPrompt
}

func (p *Prompt) GetSystemPrompt() string {
	return p.systemPrompt
}

func (p *Prompt) GetUserPrompt() string {
	return p.userPrompt
}
