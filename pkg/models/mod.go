package models

type ModelConfig struct {
	ApiKey                      string
	PromptModel                 string
	OllamaAiModel               string
	PromptMaxTokens             int64
	PromptRequestTimeoutSeconds int64
}

type ModelResponse struct {
	Content string
}

type Model interface {
	ExecPrompt(systemPrompt string, userPrompt string, modelConfig ModelConfig) (*ModelResponse, error)
}
