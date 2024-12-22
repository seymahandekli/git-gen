package platforms

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const (
	geminiApiEndpoint    = "https://generativelanguage.googleapis.com/v1beta/openai/chat/completions"
	geminiDefaultModel   = "gemini-1.5-pro"
	maxScannerBufferSize = 10 * 1024 * 1024 // 10 MB buffer
)

var (
	ErrGeminiApiKeyIsRequired = errors.New("Gemini platform requires PLATFORM_API_KEY to be specified")
)

type geminiPromptRequestMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type geminiPromptRequest struct {
	Model    string                       `json:"model"`
	Messages []geminiPromptRequestMessage `json:"messages"`
	Stream   bool                         `json:"stream"`
}

type geminiPromptResponseChunk struct {
	Choices []struct {
		Delta struct {
			Content string `json:"content"`
		} `json:"delta"`
	} `json:"choices"`
}

type Gemini struct {
	platformConfig PlatformConfig
}

func NewGemini(platformConfig PlatformConfig) *Gemini {
	return &Gemini{
		platformConfig: platformConfig,
	}
}

func (g *Gemini) ExecPrompt(ctx context.Context, promptSource PromptGenerator) (*ModelResponse, error) {
	if strings.TrimSpace(g.platformConfig.ApiKey) == "" {
		return nil, ErrGeminiApiKeyIsRequired
	}

	model := geminiDefaultModel
	if g.platformConfig.Model != "" {
		model = g.platformConfig.Model
	}

	payload := geminiPromptRequest{
		Model: model,
		Messages: []geminiPromptRequestMessage{
			{Role: "system", Content: promptSource.GetSystemPrompt()},
			{Role: "user", Content: promptSource.GetUserPrompt()},
		},
		Stream: true,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", geminiApiEndpoint, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+g.platformConfig.ApiKey)

	client := &http.Client{
		Timeout: 60 * time.Second,
	}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusTooManyRequests {
		retryAfter := res.Header.Get("Retry-After")
		return nil, fmt.Errorf("Rate limit exceeded. Retry after %s seconds", retryAfter)
	}

	if res.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(res.Body)
		return nil, fmt.Errorf("Gemini API request failed with status code %d. Response: %s", res.StatusCode, string(body))
	}

	scanner := bufio.NewScanner(res.Body)
	buf := make([]byte, maxScannerBufferSize)
	scanner.Buffer(buf, maxScannerBufferSize)

	var responseBuilder strings.Builder
	for scanner.Scan() {
		line := scanner.Text()

		// Skip empty lines or lines that start with non-data prefixes
		if len(line) == 0 || !strings.HasPrefix(line, "data: ") {
			continue
		}

		// Remove the "data: " prefix
		line = strings.TrimPrefix(line, "data: ")

		// Check for the "[DONE]" marker
		if line == "[DONE]" {
			break // End of stream
		}

		// Parse the JSON chunk
		var chunk geminiPromptResponseChunk
		if err := json.Unmarshal([]byte(line), &chunk); err != nil {
			return nil, fmt.Errorf("Failed to parse chunk: %v. Line: %s", err, line)
		}

		// Append the content to the response builder
		if len(chunk.Choices) > 0 && chunk.Choices[0].Delta.Content != "" {
			responseBuilder.WriteString(chunk.Choices[0].Delta.Content)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	// Return the final model response
	return &ModelResponse{
		Content: responseBuilder.String(),
	}, nil

}
