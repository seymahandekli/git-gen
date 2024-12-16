package platforms

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

const (
	anthropicApiEndpoint = "https://api.anthropic.com/v1/messages"
	// Using Claude 3.5 Sonnet as default model based on docs
	anthropicDefaultModel = "claude-3-5-sonnet-20241022"
	anthropicApiVersion   = "2023-06-01"
)

var (
	ErrAnthropicApiKeyIsRequired = errors.New("Anthropic platform requires PLATFORM_API_KEY is specified")
)

// Request structures
type anthropicMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type anthropicRequest struct {
	Model     string             `json:"model"`
	Messages  []anthropicMessage `json:"messages"`
	MaxTokens int64              `json:"max_tokens"`
	System    string             `json:"system,omitempty"`
}

// Response structures
type anthropicContentBlock struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type anthropicResponse struct {
	Id      string                  `json:"id"`
	Type    string                  `json:"type"`
	Role    string                  `json:"role"`
	Content []anthropicContentBlock `json:"content"`
	Model   string                  `json:"model"`
	Usage   struct {
		InputTokens  int `json:"input_tokens"`
		OutputTokens int `json:"output_tokens"`
	} `json:"usage"`
	StopReason   string  `json:"stop_reason"`
	StopSequence *string `json:"stop_sequence"`
}

type anthropicErrorDetail struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

type anthropicErrorResponse struct {
	Type  string               `json:"type"`
	Error anthropicErrorDetail `json:"error"`
}

type Anthropic struct {
	platformConfig PlatformConfig
}

func NewAnthropic(platformConfig PlatformConfig) *Anthropic {
	return &Anthropic{
		platformConfig: platformConfig,
	}
}

func (a *Anthropic) ExecPrompt(ctx context.Context, promptSource PromptGenerator) (*ModelResponse, error) {
	if a.platformConfig.ApiKey == "" {
		return nil, ErrAnthropicApiKeyIsRequired
	}

	targetModel := anthropicDefaultModel
	if a.platformConfig.Model != "" {
		targetModel = a.platformConfig.Model
	}

	// Create request payload
	payload := anthropicRequest{
		Model: targetModel,
		Messages: []anthropicMessage{
			{
				Role:    "user",
				Content: promptSource.GetUserPrompt(),
			},
		},
		MaxTokens: a.platformConfig.PromptMaxTokens,
	}

	// Add system prompt if provided
	if systemPrompt := promptSource.GetSystemPrompt(); systemPrompt != "" {
		payload.System = systemPrompt
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	// Create HTTP request
	req, err := http.NewRequestWithContext(ctx, "POST", anthropicApiEndpoint, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	// Set required headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", a.platformConfig.ApiKey)
	req.Header.Set("anthropic-version", anthropicApiVersion)

	// Send request
	client := &http.Client{
		Timeout: time.Duration(a.platformConfig.PromptRequestTimeoutSeconds) * time.Second,
	}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		var errResp anthropicErrorResponse
		if err := json.NewDecoder(res.Body).Decode(&errResp); err != nil {
			return nil, fmt.Errorf("anthropic API error (status %d): failed to decode error response: %w",
				res.StatusCode, err)
		}
		return nil, fmt.Errorf("anthropic API error: %s - %s",
			errResp.Error.Type, errResp.Error.Message)
	}

	var data anthropicResponse
	err = json.NewDecoder(res.Body).Decode(&data)
	if err != nil {
		return nil, fmt.Errorf("failed to decode successful response: %w", err)
	}

	// Check for error type in successful response
	if data.Type == "error" {
		return nil, fmt.Errorf("anthropic API returned error type in response")
	}

	// Extract text content from the first content block
	var content string
	if len(data.Content) > 0 {
		content = data.Content[0].Text
	}

	response := ModelResponse{
		Content: content,
	}

	return &response, nil
}
