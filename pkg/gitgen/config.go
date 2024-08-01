package gitgen

type Config struct {
	OpenAiKey                   string
	SourceRef                   string
	DestinationRef              string
	PromptModel                 string
	OllamaAiModel			   string
	PromptMaxTokens             int64
	PromptRequestTimeoutSeconds int64
}

type ConfigOption func(*Config)

func WithOpenAiKey(apiKey string) ConfigOption {
	return func(c *Config) {
		c.OpenAiKey = apiKey
	}
}

func WithSourceRef(ref string) ConfigOption {
	return func(c *Config) {
		c.SourceRef = ref
	}
}

func WithDestinationRef(ref string) ConfigOption {
	return func(c *Config) {
		c.DestinationRef = ref
	}
}

func WithPromptModel(model string) ConfigOption {
	return func(c *Config) {
		c.PromptModel = model
	}
}

func WithOllamaAiModel(model string) ConfigOption {
	return func(c *Config) {
		c.OllamaAiModel = model
	}
}

func WithPromptMaxTokens(tokens int64) ConfigOption {
	return func(c *Config) {
		c.PromptMaxTokens = tokens
	}
}

func WithPromptRequestTimeoutSeconds(timeout int64) ConfigOption {
	return func(c *Config) {
		c.PromptRequestTimeoutSeconds = timeout
	}
}

func NewConfig(opts ...ConfigOption) *Config {
	config := &Config{
		OpenAiKey:                   "",
		SourceRef:                   "HEAD",
		DestinationRef:              "",
		PromptModel:                 "gpt-4o",
		OllamaAiModel:               "llama3",
		PromptMaxTokens:             3500,
		PromptRequestTimeoutSeconds: 3600,
	}

	for _, opt := range opts {
		opt(config)
	}

	return config
}
