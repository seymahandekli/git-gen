package models

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"
)

const (
	apiEndpoint = "https://api.openai.com/v1/chat/completions"
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

type OpenAi struct{}

func NewOpenAi() *OpenAi {
	return &OpenAi{}
}

func (o *OpenAi) ExecPrompt(systemPrompt string, userPrompt string, modelConfig ModelConfig) (*ModelResponse, error) {
	// Create the request body
	request := openAiPromptRequest{
		Model: modelConfig.PromptModel,
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
		MaxTokens: modelConfig.PromptMaxTokens,
	}

	body, err := json.MarshalIndent(request, "", "  ") // Use json.MarshalIndent for pretty printing
	if err != nil {
		return nil, err
	}

	// Create the HTTP request
	req, err := http.NewRequest("POST", apiEndpoint, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+modelConfig.ApiKey)

	// Send the request
	client := &http.Client{
		Timeout: time.Duration(modelConfig.PromptRequestTimeoutSeconds) * time.Second,
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
