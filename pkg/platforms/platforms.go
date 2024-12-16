package platforms

import (
	"context"
	"errors"
	"fmt"
)

type Platform string

const (
	PlatformOpenAI    Platform = "openai"
	PlatformOllama    Platform = "ollama"
	PlatformAnthropic Platform = "anthropic"
)

var (
	ErrUnknownPlatform = errors.New("unknown platform")
)

type PlatformConfig struct {
	ApiKey                      string
	Model                       string
	PromptMaxTokens             int64
	PromptRequestTimeoutSeconds int64
}

type ModelRequest struct {
	SystemPrompt string
	UserPrompt   string
}

type ModelResponse struct {
	Content string
}

type PromptGenerator interface {
	GetSystemPrompt() string
	GetUserPrompt() string
}

type PromptExecutor interface {
	ExecPrompt(ctx context.Context, promptSource PromptGenerator) (*ModelResponse, error)
}

func NewPromptExecutor(platform Platform, platformConfig PlatformConfig) (PromptExecutor, error) {
	switch platform {
	case PlatformOpenAI:
		return NewOpenAi(platformConfig), nil
	case PlatformOllama:
		return NewOllama(platformConfig)
	case PlatformAnthropic:
		return NewAnthropic(platformConfig), nil
	default:
		return nil, fmt.Errorf("unknown platform %s - %w", platform, ErrUnknownPlatform)
	}
}
