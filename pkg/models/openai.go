package models

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

const (
	apiEndpoint        = "https://api.openai.com/v1/chat/completions"
	openaiDefaultModel = "gpt-4o"
)

var (
	ErrPlatformApiKeyIsRequired = errors.New("OpenAI platform requires PLATFORM_API_KEY is specified")
)

type openAiPromptRequestMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type openAiPromptRequest struct {
	Model     string                       `json:"model"`
	Messages  []openAiPromptRequestMessage `json:"messages"`
	MaxTokens int64                        `json:"max_tokens"`
}

type openAiPromptResponse struct {
	Id      string `json:"id"`
	Object  string `json:"object"`
	Created int    `json:"created"`
	Model   string `json:"model"`
	Usage   struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
	Choices []struct {
		Index   int `json:"index"`
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
		Logprobs     *bool  `json:"logprobs"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
}

type OpenAi struct {
	modelConfig ModelConfig
}

func NewOpenAi(modelConfig ModelConfig) *OpenAi {
	return &OpenAi{
		modelConfig: modelConfig,
	}
}

func (o *OpenAi) ExecPrompt(ctx context.Context, systemPrompt string, userPrompt string) (*ModelResponse, error) {
	if o.modelConfig.PlatformApiKey == "" {
		return nil, ErrPlatformApiKeyIsRequired
	}

	var targetModel string = openaiDefaultModel

	// if model is specified by user
	if o.modelConfig.Model != "" {
		targetModel = o.modelConfig.Model
	}

	// Create the request body
	request := openAiPromptRequest{
		Model: targetModel,
		Messages: []openAiPromptRequestMessage{
			{
				Role:    "system",
				Content: systemPrompt,
			},
			{
				Role:    "user",
				Content: userPrompt,
			},
		},
		MaxTokens: o.modelConfig.PromptMaxTokens,
	}

	body, err := json.MarshalIndent(request, "", "  ") // Use json.MarshalIndent for pretty printing
	if err != nil {
		return nil, err
	}

	// Create the HTTP request
	req, err := http.NewRequestWithContext(ctx, "POST", apiEndpoint, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+o.modelConfig.PlatformApiKey)

	// Send the request
	client := &http.Client{
		Timeout: time.Duration(o.modelConfig.PromptRequestTimeoutSeconds) * time.Second,
	}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var data openAiPromptResponse

	err = json.NewDecoder(res.Body).Decode(&data)
	if err != nil {
		return nil, err
	}

	response := ModelResponse{
		Content: data.Choices[0].Message.Content,
	}

	return &response, nil
}
