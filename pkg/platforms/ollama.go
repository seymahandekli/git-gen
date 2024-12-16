package platforms

import (
	"context"
	"fmt"

	"github.com/ollama/ollama/api"
)

const (
	ollamaDefaultModel = "llama3"
)

type Ollama struct {
	platformConfig PlatformConfig

	client *api.Client
}

func NewOllama(platformConfig PlatformConfig) (*Ollama, error) {
	client, err := api.ClientFromEnvironment()
	if err != nil {
		return nil, fmt.Errorf("failed to create Ollama API client: %w", err)
	}

	return &Ollama{
		platformConfig: platformConfig,
		client:         client,
	}, nil
}

func (o *Ollama) ExecPrompt(ctx context.Context, promptSource PromptGenerator) (*ModelResponse, error) {
	var targetModel string = ollamaDefaultModel

	// if model is specified by user
	if o.platformConfig.Model != "" {
		targetModel = o.platformConfig.Model
	}

	payload := &api.ChatRequest{
		Model: targetModel,
		Messages: []api.Message{
			{
				Role:    "system",
				Content: promptSource.GetSystemPrompt(),
			},
			{
				Role:    "user",
				Content: promptSource.GetUserPrompt(),
			},
		},
	}

	respFunc := func(resp api.ChatResponse) error {
		fmt.Print(resp.Message.Content)
		return nil
	}

	if err := o.client.Chat(ctx, payload, respFunc); err != nil {
		return nil, fmt.Errorf("failed to execute chat request: %w", err)
	}

	return &ModelResponse{Content: "response"}, nil
}
