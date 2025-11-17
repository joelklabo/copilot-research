package provider

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/sashabaranov/go-openai"
)

// OpenAIProvider implements the AIProvider interface for OpenAI
type OpenAIProvider struct {
	client  *openai.Client
	model   string
	timeout time.Duration
	apiKey  string
}

// NewOpenAIProvider creates a new OpenAI provider
func NewOpenAIProvider(model string, timeout time.Duration) *OpenAIProvider {
	apiKey := os.Getenv("OPENAI_API_KEY")
	
	var client *openai.Client
	if apiKey != "" {
		client = openai.NewClient(apiKey)
	}
	
	return &OpenAIProvider{
		client:  client,
		model:   model,
		timeout: timeout,
		apiKey:  apiKey,
	}
}

// Name returns the provider name
func (o *OpenAIProvider) Name() string {
	return "openai"
}

// Query executes a query using OpenAI API
func (o *OpenAIProvider) Query(ctx context.Context, prompt string, opts QueryOptions) (*Response, error) {
	// Check authentication first
	if !o.IsAuthenticated() {
		return nil, fmt.Errorf("not authenticated: please set OPENAI_API_KEY environment variable")
	}
	
	// Create context with timeout
	queryCtx, cancel := context.WithTimeout(ctx, o.timeout)
	defer cancel()
	
	// Prepare request
	model := o.model
	if opts.Model != "" {
		model = opts.Model
	}
	
	maxTokens := 4000
	if opts.MaxTokens > 0 {
		maxTokens = opts.MaxTokens
	}
	
	temperature := float32(0.7)
	if opts.Temperature > 0 {
		temperature = float32(opts.Temperature)
	}
	
	topP := float32(1.0)
	if opts.TopP > 0 {
		topP = float32(opts.TopP)
	}
	
	// Execute request
	start := time.Now()
	req := openai.ChatCompletionRequest{
		Model: model,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleUser,
				Content: prompt,
			},
		},
		MaxTokens:   maxTokens,
		Temperature: temperature,
		TopP:        topP,
	}
	
	resp, err := o.client.CreateChatCompletion(queryCtx, req)
	duration := time.Since(start)
	
	if err != nil {
		// Check for timeout
		if queryCtx.Err() == context.DeadlineExceeded {
			return nil, fmt.Errorf("query timeout after %v", o.timeout)
		}
		
		// Check for rate limiting
		if isRateLimitError(err) {
			return nil, fmt.Errorf("OpenAI rate limit exceeded: %w", err)
		}
		
		return nil, fmt.Errorf("OpenAI API error: %w", err)
	}
	
	// Parse response
	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("no response from OpenAI")
	}
	
	content := resp.Choices[0].Message.Content
	
	return &Response{
		Content:  content,
		Provider: "openai",
		Model:    resp.Model,
		Duration: duration,
		TokensUsed: TokenUsage{
			Prompt:     resp.Usage.PromptTokens,
			Completion: resp.Usage.CompletionTokens,
			Total:      resp.Usage.TotalTokens,
		},
		Metadata: map[string]interface{}{
			"finish_reason": resp.Choices[0].FinishReason,
		},
	}, nil
}

// IsAuthenticated checks if the provider is authenticated
func (o *OpenAIProvider) IsAuthenticated() bool {
	return o.apiKey != ""
}

// RequiresAuth returns authentication information
func (o *OpenAIProvider) RequiresAuth() AuthInfo {
	if o.IsAuthenticated() {
		return AuthInfo{
			Type:         "apikey",
			IsConfigured: true,
		}
	}
	
	return AuthInfo{
		Type:         "apikey",
		IsConfigured: false,
		HelpURL:      "https://platform.openai.com/api-keys",
		Instructions: `OpenAI API key required.

Get your API key:
  1. Visit https://platform.openai.com/api-keys
  2. Create a new API key
  3. Set it in your environment:
     export OPENAI_API_KEY=sk-...
     
Or add to config:
  copilot-research config set providers.openai.api_key sk-...

Pricing: https://openai.com/pricing`,
	}
}

// Capabilities returns the provider's capabilities
func (o *OpenAIProvider) Capabilities() ProviderCapabilities {
	return ProviderCapabilities{
		Streaming:      true,
		FunctionCall:   true,
		MaxTokens:      128000, // GPT-4-turbo
		SupportsImages: true,   // GPT-4-vision
	}
}

// isRateLimitError checks if an error is a rate limit error
func isRateLimitError(err error) bool {
	// OpenAI SDK wraps rate limit errors
	// Check if error message contains rate limit keywords
	if err == nil {
		return false
	}
	errMsg := err.Error()
	// Use strings package for substring check
	return len(errMsg) > 0 && (
		findSubstring(errMsg, "rate limit") || 
		findSubstring(errMsg, "429"))
}

// findSubstring checks if substr is in s
func findSubstring(s, substr string) bool {
	if len(substr) > len(s) {
		return false
	}
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
