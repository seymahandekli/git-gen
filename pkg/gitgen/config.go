package gitgen

type Config struct {
	PlatformApiKey              string
	SourceRef                   string
	DestinationRef              string
	Platform                    string
	Model                       string
	PromptMaxTokens             int64
	PromptRequestTimeoutSeconds int64
}

type ConfigOption func(*Config)

func WithPlatformApiKey(apiKey string) ConfigOption {
	return func(c *Config) {
		c.PlatformApiKey = apiKey
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

func WithPlatform(platform string) ConfigOption {
	return func(c *Config) {
		c.Platform = platform
	}
}

func WithModel(model string) ConfigOption {
	return func(c *Config) {
		c.Model = model
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
		PlatformApiKey:              "",
		SourceRef:                   "HEAD",
		DestinationRef:              "",
		Platform:                    "openai",
		Model:                       "",
		PromptMaxTokens:             3500,
		PromptRequestTimeoutSeconds: 3600,
	}

	for _, opt := range opts {
		opt(config)
	}

	return config
}
