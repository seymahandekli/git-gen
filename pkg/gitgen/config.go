package gitgen

import (
	"errors"
	"fmt"
	"os"
)

type Config struct {
	OpenApiKey                  string
	PromptModel                 string
	PromptMaxTokens             int64
	PromptRequestTimeoutSeconds int64
}

var (
	ErrOpenApiKeyIsNeeded = errors.New("open api key is needed")
)

func InitConfig() (Config, error) {
	openApiKey, openApiKeyOk := os.LookupEnv("OPENAI_API_KEY")
	if !openApiKeyOk {
		return Config{}, fmt.Errorf("environment variable OPENAI_API_KEY must be set: %w", ErrOpenApiKeyIsNeeded)
	}

	config := Config{
		OpenApiKey:                  openApiKey,
		PromptModel:                 "gpt-4o",
		PromptMaxTokens:             3500,
		PromptRequestTimeoutSeconds: 3600,
	}

	return config, nil
}
