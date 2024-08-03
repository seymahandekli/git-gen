package models

import (
	"context"
	"fmt"

	"github.com/ollama/ollama/api"
)

const (
	ollamaDefaultModel = "llama3"
)

type OllamaAi struct {
	modelConfig ModelConfig

	client *api.Client
}

func NewOllamaAi(modelConfig ModelConfig) (*OllamaAi, error) {
	client, err := api.ClientFromEnvironment()
	if err != nil {
		return nil, fmt.Errorf("failed to create Ollama API client: %w", err)
	}

	return &OllamaAi{
		modelConfig: modelConfig,
		client:      client,
	}, nil
}

func (o *OllamaAi) ExecPrompt(ctx context.Context, systemPrompt string, userPrompt string) (*ModelResponse, error) {
	var targetModel string = ollamaDefaultModel

	// if model is specified by user
	if o.modelConfig.Model != "" {
		targetModel = o.modelConfig.Model
	}

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

	req := &api.ChatRequest{
		Model:    targetModel,
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
