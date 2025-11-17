package provider

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewOpenAIProvider(t *testing.T) {
	provider := NewOpenAIProvider("gpt-4", 60*time.Second)
	assert.NotNil(t, provider)
	assert.Equal(t, "openai", provider.Name())
}

func TestOpenAIProvider_Name(t *testing.T) {
	provider := NewOpenAIProvider("gpt-4", 60*time.Second)
	assert.Equal(t, "openai", provider.Name())
}

func TestOpenAIProvider_Capabilities(t *testing.T) {
	provider := NewOpenAIProvider("gpt-4", 60*time.Second)
	caps := provider.Capabilities()
	
	assert.True(t, caps.Streaming)
	assert.True(t, caps.FunctionCall)
	assert.Equal(t, 128000, caps.MaxTokens)
	assert.True(t, caps.SupportsImages)
}

func TestOpenAIProvider_IsAuthenticated_WithAPIKey(t *testing.T) {
	// Set API key in environment
	os.Setenv("OPENAI_API_KEY", "sk-test-key")
	defer os.Unsetenv("OPENAI_API_KEY")
	
	provider := NewOpenAIProvider("gpt-4", 60*time.Second)
	assert.True(t, provider.IsAuthenticated())
}

func TestOpenAIProvider_IsAuthenticated_NoAPIKey(t *testing.T) {
	// Make sure no API key is set
	os.Unsetenv("OPENAI_API_KEY")
	
	provider := NewOpenAIProvider("gpt-4", 60*time.Second)
	assert.False(t, provider.IsAuthenticated())
}

func TestOpenAIProvider_RequiresAuth_Authenticated(t *testing.T) {
	os.Setenv("OPENAI_API_KEY", "sk-test-key")
	defer os.Unsetenv("OPENAI_API_KEY")
	
	provider := NewOpenAIProvider("gpt-4", 60*time.Second)
	authInfo := provider.RequiresAuth()
	
	assert.True(t, authInfo.IsConfigured)
}

func TestOpenAIProvider_RequiresAuth_NotAuthenticated(t *testing.T) {
	os.Unsetenv("OPENAI_API_KEY")
	
	provider := NewOpenAIProvider("gpt-4", 60*time.Second)
	authInfo := provider.RequiresAuth()
	
	assert.False(t, authInfo.IsConfigured)
	assert.Equal(t, "apikey", authInfo.Type)
	assert.Contains(t, authInfo.HelpURL, "openai.com")
	assert.Contains(t, authInfo.Instructions, "OPENAI_API_KEY")
	assert.Contains(t, authInfo.Instructions, "https://platform.openai.com/api-keys")
}

func TestOpenAIProvider_Query_NotAuthenticated(t *testing.T) {
	os.Unsetenv("OPENAI_API_KEY")
	
	provider := NewOpenAIProvider("gpt-4", 60*time.Second)
	
	ctx := context.Background()
	_, err := provider.Query(ctx, "test prompt", QueryOptions{})
	
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not authenticated")
}

func TestOpenAIProvider_Query_WithTimeout(t *testing.T) {
	os.Setenv("OPENAI_API_KEY", "sk-test-key")
	defer os.Unsetenv("OPENAI_API_KEY")
	
	provider := NewOpenAIProvider("gpt-4", 1*time.Millisecond)
	
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
	defer cancel()
	
	_, err := provider.Query(ctx, "test prompt", QueryOptions{})
	
	// Should timeout or fail (we don't have real API key)
	assert.Error(t, err)
}

func TestOpenAIProvider_ModelSelection(t *testing.T) {
	tests := []struct {
		name  string
		model string
	}{
		{"gpt-4", "gpt-4"},
		{"gpt-4-turbo", "gpt-4-turbo-preview"},
		{"gpt-3.5-turbo", "gpt-3.5-turbo"},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			provider := NewOpenAIProvider(tt.model, 60*time.Second)
			assert.NotNil(t, provider)
		})
	}
}

func TestOpenAIProvider_QueryOptions(t *testing.T) {
	os.Setenv("OPENAI_API_KEY", "sk-test-key")
	defer os.Unsetenv("OPENAI_API_KEY")
	
	provider := NewOpenAIProvider("gpt-4", 60*time.Second)
	
	opts := QueryOptions{
		MaxTokens:   1000,
		Temperature: 0.7,
		TopP:        0.9,
		Model:       "gpt-4",
	}
	
	ctx := context.Background()
	_, err := provider.Query(ctx, "test", opts)
	
	// Will fail with invalid API key, but we're testing it accepts the options
	assert.Error(t, err)
}

func TestOpenAIProvider_ErrorMessages(t *testing.T) {
	// Test that error messages are helpful
	os.Unsetenv("OPENAI_API_KEY")
	
	provider := NewOpenAIProvider("gpt-4", 60*time.Second)
	authInfo := provider.RequiresAuth()
	
	// Should have clear instructions
	assert.Contains(t, authInfo.Instructions, "export OPENAI_API_KEY")
	assert.Contains(t, authInfo.Instructions, "https://platform.openai.com/api-keys")
	assert.Contains(t, authInfo.Instructions, "Pricing")
}
