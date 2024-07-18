package gitgen

type Config struct {
	OpenApiKey                  string
	PromptModel                 string
	PromptMaxTokens             int64
	PromptRequestTimeoutSeconds int64
}

func DefaultConfig() Config {
	config := Config{
		OpenApiKey:                  "",
		PromptModel:                 "gpt-4o",
		PromptMaxTokens:             3500,
		PromptRequestTimeoutSeconds: 3600,
	}

	return config
}
