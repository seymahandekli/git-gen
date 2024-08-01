package models

import (
	"context"
	"fmt"

	"github.com/ollama/ollama/api"
)

type OllamaAi struct {
	client *api.Client
}

func NewOllamaAi() *OllamaAi {
	client, err := api.ClientFromEnvironment()
	if err != nil {
		println(fmt.Errorf("failed to create OllamaAi client: %w", err))
	}
	return &OllamaAi{client: client}

}

func (o *OllamaAi) ExecPrompt(systemPrompt string, userPrompt string, modelConfig ModelConfig) (*ModelResponse, error) {
	request := []api.Message{
		{
			Role:    "system",
			Content: systemPrompt,
		},
		{
			Role:    "user",
			Content: userPrompt,
		},
	}
	ctx := context.Background()
	req := &api.ChatRequest{
		Model:    modelConfig.OllamaAiModel,
		Messages: request,
	}

	respFunc := func(resp api.ChatResponse) error {
		fmt.Print(resp.Message.Content)
		return nil
	}

	if err := o.client.Chat(ctx, req, respFunc); err != nil {
		return nil, fmt.Errorf("failed to execute chat request: %w", err)
	}

	return &ModelResponse{Content: "response"}, nil

}
