package config

import (
	"fmt" // Added this import
	"os"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v3"
)

// Config holds the entire application configuration
type Config struct {
	Providers ProviderConfig `yaml:"providers"`
}

// ProviderConfig holds configuration for AI providers
type ProviderConfig struct {
	Primary  string `yaml:"primary"`
	Fallback string `yaml:"fallback"`

	GitHubCopilot GitHubCopilotConfig `yaml:"github-copilot"`
	OpenAI        OpenAIConfig        `yaml:"openai"`
	Anthropic     AnthropicConfig     `yaml:"anthropic"`

	AutoFallback   bool `yaml:"auto_fallback"`
	NotifyFallback bool `yaml:"notify_fallback"`
}

// GitHubCopilotConfig holds configuration for the GitHub Copilot provider
type GitHubCopilotConfig struct {
	Enabled  bool          `yaml:"enabled"`
	AuthType string        `yaml:"auth_type"` // cli, pat, oauth
	Timeout  time.Duration `yaml:"timeout"`
}

// OpenAIConfig holds configuration for the OpenAI provider
type OpenAIConfig struct {
	Enabled     bool          `yaml:"enabled"`
	AuthType    string        `yaml:"auth_type"` // apikey
	APIKeyEnv   string        `yaml:"api_key_env"`
	Model       string        `yaml:"model"`
	Temperature float64       `yaml:"temperature"`
	MaxTokens   int           `yaml:"max_tokens"`
	Timeout     time.Duration `yaml:"timeout"`
}

// AnthropicConfig holds configuration for the Anthropic provider
type AnthropicConfig struct {
	Enabled   bool          `yaml:"enabled"`
	AuthType  string        `yaml:"auth_type"` // apikey
	APIKeyEnv string        `yaml:"api_key_env"`
	Model     string        `yaml:"model"`
	Timeout   time.Duration `yaml:"timeout"`
}

// DefaultConfig returns a new Config with sensible defaults
func DefaultConfig() *Config {
	return &Config{
		Providers: ProviderConfig{
			Primary:  "github-copilot",
			Fallback: "openai",

			GitHubCopilot: GitHubCopilotConfig{
				Enabled:  true,
				AuthType: "cli",
				Timeout:  60 * time.Second,
			},
			OpenAI: OpenAIConfig{
				Enabled:     true,
				AuthType:    "apikey",
				APIKeyEnv:   "OPENAI_API_KEY",
				Model:       "gpt-4o", // Updated to a more recent model
				Temperature: 0.7,
				MaxTokens:   4000,
				Timeout:     30 * time.Second,
			},
			Anthropic: AnthropicConfig{
				Enabled:   false,
				AuthType:  "apikey",
				APIKeyEnv: "ANTHROPIC_API_KEY",
				Model:     "claude-3-5-sonnet",
				Timeout:   30 * time.Second,
			},
			AutoFallback:   true,
			NotifyFallback: true,
		},
	}
}

// LoadConfig loads configuration from the specified path
func LoadConfig(path string) (*Config, error) {
	cfg := DefaultConfig()

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			// If config file doesn't exist, save default and return it
			if err := SaveConfig(path, cfg); err != nil {
				return nil, fmt.Errorf("failed to save default config: %w", err)
			}
			return cfg, nil
		}
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return cfg, nil
}

// SaveConfig saves the configuration to the specified path
func SaveConfig(path string, cfg *Config) error {
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	// Ensure directory exists
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}