package gitgen

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"
)

const (
	apiEndpoint = "https://api.openai.com/v1/chat/completions"
)

type PromptRequestMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type PromptRequest struct {
	Model     string                 `json:"model"`
	Messages  []PromptRequestMessage `json:"messages"`
	MaxTokens int64                  `json:"max_tokens"`
}

type PromptResponse struct {
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

func ExecPrompt(systemPrompt string, userPrompt string, config Config) (*PromptResponse, error) {
	// Create the request body
	request := PromptRequest{
		Model: config.PromptModel,
		Messages: []PromptRequestMessage{
			{
				Role:    "system",
				Content: systemPrompt,
			},
			{
				Role:    "user",
				Content: userPrompt,
			},
		},
		MaxTokens: config.PromptMaxTokens,
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
	req.Header.Set("Authorization", "Bearer "+config.OpenAiKey)

	// Send the request
	client := &http.Client{
		Timeout: time.Duration(config.PromptRequestTimeoutSeconds) * time.Second,
	}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var data PromptResponse

	err = json.NewDecoder(res.Body).Decode(&data)
	if err != nil {
		return nil, err
	}

	return &data, nil
}
