package provider

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/liushuangls/go-anthropic"
)

// AnthropicProvider implements the AIProvider interface for Anthropic Claude
type AnthropicProvider struct {
	client  *anthropic.Client
	model   string
	timeout time.Duration
	apiKey  string
}

// NewAnthropicProvider creates a new Anthropic provider
func NewAnthropicProvider(model string, timeout time.Duration, apiKeyEnv string) *AnthropicProvider {
	apiKey := os.Getenv(apiKeyEnv)

	var client *anthropic.Client
	if apiKey != "" {
		client = anthropic.NewClient(apiKey)
	}

	return &AnthropicProvider{
		client:  client,
		model:   model,
		timeout: timeout,
		apiKey:  apiKey,
	}
}

// Name returns the provider name
func (a *AnthropicProvider) Name() string {
	return "anthropic"
}

// Query executes a query using Anthropic API
func (a *AnthropicProvider) Query(ctx context.Context, prompt string, opts QueryOptions) (*Response, error) {
	// Check authentication first
	if !a.IsAuthenticated() {
		return nil, fmt.Errorf("not authenticated: please set %s environment variable", a.apiKey)
	}

	// Create context with timeout
	queryCtx, cancel := context.WithTimeout(ctx, a.timeout)
	defer cancel()

	// Prepare request
	model := a.model
	if opts.Model != "" {
		model = opts.Model
	}

	maxTokens := 4000 // Default max tokens for Anthropic
	if opts.MaxTokens > 0 {
		maxTokens = opts.MaxTokens
	}

	// Anthropic API uses messages format
	// Corrected: Wrap prompt string in anthropic.MessageContent
	messages := []anthropic.Message{
		{
			Role: anthropic.RoleUser,
			Content: []anthropic.MessageContent{
				{
					Type: "text",
					Text: &prompt, // Changed from prompt to &prompt
				},
			},
		},
	}

	start := time.Now()
	resp, err := a.client.CreateMessages(queryCtx, anthropic.MessagesRequest{
		Model:     model,
		Messages:  messages,
		MaxTokens: maxTokens,
		// Temperature and TopP are not directly exposed in QueryOptions for Anthropic yet
	})
	duration := time.Since(start)

	if err != nil {
		// Check for timeout
		if queryCtx.Err() == context.DeadlineExceeded {
			return nil, fmt.Errorf("query timeout after %v", a.timeout)
		}
		// TODO: Add more specific error handling for rate limits, etc.
		return nil, fmt.Errorf("Anthropic API error: %w", err)
	}

	// Parse response
	if len(resp.Content) == 0 {
		return nil, fmt.Errorf("no response from Anthropic")
	}

	content := ""
	for _, block := range resp.Content {
		if block.Type == "text" {
			content += block.Text
		}
	}

	return &Response{
		Content:  content,
		Provider: "anthropic",
		Model:    resp.Model,
		Duration: duration,
		TokensUsed: TokenUsage{
			Prompt:     resp.Usage.InputTokens,
			Completion: resp.Usage.OutputTokens,
			Total:      resp.Usage.InputTokens + resp.Usage.OutputTokens,
		},
		Metadata: map[string]interface{}{
			"stop_reason": resp.StopReason,
		},
	}, nil
}

// IsAuthenticated checks if the provider is authenticated
func (a *AnthropicProvider) IsAuthenticated() bool {
	return a.apiKey != ""
}

// RequiresAuth returns authentication information
func (a *AnthropicProvider) RequiresAuth() AuthInfo {
	if a.IsAuthenticated() {
		return AuthInfo{
			Type:         "apikey",
			IsConfigured: true,
		}
	}

	return AuthInfo{
		Type:         "apikey",
		IsConfigured: false,
		HelpURL:      "https://console.anthropic.com/",
		Instructions: `Anthropic API key required.

Get your API key:
  1. Visit https://console.anthropic.com/
  2. Create an account
  3. Generate API key
  4. Set in environment:
     export ANTHROPIC_API_KEY=sk-ant-...

Or add to config:
  copilot-research config set providers.anthropic.api_key sk-ant-...`,
	}
}

// Capabilities returns the provider's capabilities
func (a *AnthropicProvider) Capabilities() ProviderCapabilities {
	return ProviderCapabilities{
		Streaming:      true,
		FunctionCall:   false, // Anthropic does not directly support function calling in this SDK version
		MaxTokens:      200000, // Claude 3
		SupportsImages: true,
	}
}
