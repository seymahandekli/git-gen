package models

import (
	"context"
	"fmt"
	"log"

	"github.com/ollama/ollama/api"
)

type OllamaAi struct{}

func NewOllamaAi() *OllamaAi {
	return &OllamaAi{}
}

func (o *OllamaAi) ExecPrompt(systemPrompt string, userPrompt string, modelConfig ModelConfig) (*ModelResponse, error) {
	client, err := api.ClientFromEnvironment()
	if err != nil {
		log.Fatal(err)
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
	ctx := context.Background()
	req := &api.ChatRequest{
		Model:    "llama3",
		Messages: request,
	}
	respFunc := func(resp api.ChatResponse) error {
		fmt.Print(resp.Message.Content)
		return nil
	}

	err = client.Chat(ctx, req, respFunc)
	if err != nil {
		log.Fatal(err)
	}

	return &ModelResponse{Content: "response"}, nil

}
