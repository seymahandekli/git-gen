package models

import "context"

type ModelConfig struct {
	PlatformApiKey              string
	Platform                    string
	Model                       string
	PromptMaxTokens             int64
	PromptRequestTimeoutSeconds int64
}

type ModelResponse struct {
	Content string
}

type Model interface {
	ExecPrompt(ctx context.Context, systemPrompt string, userPrompt string) (*ModelResponse, error)
}
