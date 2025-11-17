package provider

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewAnthropicProvider(t *testing.T) {
	// Test with API key set
	os.Setenv("ANTHROPIC_API_KEY", "test-api-key")
	p := NewAnthropicProvider("claude-3-opus-20240229", 30*time.Second, "ANTHROPIC_API_KEY")
	assert.NotNil(t, p)
	assert.NotNil(t, p.client)
	assert.Equal(t, "claude-3-opus-20240229", p.model)
	assert.Equal(t, 30*time.Second, p.timeout)
	assert.Equal(t, "test-api-key", p.apiKey)
	os.Unsetenv("ANTHROPIC_API_KEY")

	// Test without API key set
	p = NewAnthropicProvider("claude-3-opus-20240229", 30*time.Second, "ANTHROPIC_API_KEY")
	assert.NotNil(t, p)
	assert.Nil(t, p.client)
	assert.Empty(t, p.apiKey)
}

func TestAnthropicProvider_Name(t *testing.T) {
	p := NewAnthropicProvider("claude-3-opus-20240229", 30*time.Second, "ANTHROPIC_API_KEY")
	assert.Equal(t, "anthropic", p.Name())
}

func TestAnthropicProvider_IsAuthenticated(t *testing.T) {
	// Authenticated
	os.Setenv("ANTHROPIC_API_KEY", "test-api-key")
	p := NewAnthropicProvider("claude-3-opus-20240229", 30*time.Second, "ANTHROPIC_API_KEY")
	assert.True(t, p.IsAuthenticated())
	os.Unsetenv("ANTHROPIC_API_KEY")

	// Not authenticated
	p = NewAnthropicProvider("claude-3-opus-20240229", 30*time.Second, "ANTHROPIC_API_KEY")
	assert.False(t, p.IsAuthenticated())
}

func TestAnthropicProvider_RequiresAuth(t *testing.T) {
	// Authenticated
	os.Setenv("ANTHROPIC_API_KEY", "test-api-key")
	p := NewAnthropicProvider("claude-3-opus-20240229", 30*time.Second, "ANTHROPIC_API_KEY")
	authInfo := p.RequiresAuth()
	assert.True(t, authInfo.IsConfigured)
	assert.Equal(t, "apikey", authInfo.Type)
	os.Unsetenv("ANTHROPIC_API_KEY")

	// Not authenticated
	p = NewAnthropicProvider("claude-3-opus-20240229", 30*time.Second, "ANTHROPIC_API_KEY")
	authInfo = p.RequiresAuth()
	assert.False(t, authInfo.IsConfigured)
	assert.Equal(t, "apikey", authInfo.Type)
	assert.Contains(t, authInfo.Instructions, "Anthropic API key required")
	assert.Contains(t, authInfo.HelpURL, "anthropic.com")
}

func TestAnthropicProvider_Capabilities(t *testing.T) {
	p := NewAnthropicProvider("claude-3-opus-20240229", 30*time.Second, "ANTHROPIC_API_KEY")
	caps := p.Capabilities()
	assert.True(t, caps.Streaming)
	assert.False(t, caps.FunctionCall) // As per implementation note
	assert.True(t, caps.SupportsImages)
	assert.Greater(t, caps.MaxTokens, 0)
}
