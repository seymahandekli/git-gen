package gitgen

type Config struct {
	OpenApiKey                  string
	PromptModel                 string
	PromptMaxTokens             int64
	PromptRequestTimeoutSeconds int64
	SourceRef                   string
	DestinationRef              string
}

func DefaultConfig() Config {
	config := Config{
		OpenApiKey:                  "",
		PromptModel:                 "gpt-4o",
		PromptMaxTokens:             3500,
		PromptRequestTimeoutSeconds: 3600,
		SourceRef:                   "HEAD~",
		DestinationRef:              "",
	}

	return config
}
