package gitgen

type Config struct {
	OpenApiKey                  string
	SourceRef                   string
	DestinationRef              string
	PromptModel                 string
	PromptMaxTokens             int64
	PromptRequestTimeoutSeconds int64
}

func DefaultConfig() Config {
	config := Config{
		OpenApiKey:                  "",
		SourceRef:                   "HEAD",
		DestinationRef:              "",
		PromptModel:                 "gpt-4o",
		PromptMaxTokens:             3500,
		PromptRequestTimeoutSeconds: 3600,
		
	}

	return config
}
