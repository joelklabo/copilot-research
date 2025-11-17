package provider

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Test that we can create a provider factory
func TestNewProviderFactory(t *testing.T) {
	factory := NewProviderFactory()
	assert.NotNil(t, factory)
	
	// Should be empty initially
	list := factory.List()
	assert.Empty(t, list)
}

// Test registering a provider
func TestProviderFactory_Register(t *testing.T) {
	factory := NewProviderFactory()
	
	// Create a mock provider
	mockProvider := &MockProvider{
		name: "test-provider",
	}
	
	// Register it
	err := factory.Register("test-provider", mockProvider)
	require.NoError(t, err)
	
	// Should appear in list
	list := factory.List()
	assert.Contains(t, list, "test-provider")
}

// Test registering duplicate provider fails
func TestProviderFactory_RegisterDuplicate(t *testing.T) {
	factory := NewProviderFactory()
	
	mockProvider := &MockProvider{name: "test"}
	
	err := factory.Register("test", mockProvider)
	require.NoError(t, err)
	
	// Try to register again
	err = factory.Register("test", mockProvider)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "already registered")
}

// Test getting a registered provider
func TestProviderFactory_Get(t *testing.T) {
	factory := NewProviderFactory()
	
	mockProvider := &MockProvider{name: "test"}
	factory.Register("test", mockProvider)
	
	// Get the provider
	provider, err := factory.Get("test")
	require.NoError(t, err)
	assert.Equal(t, "test", provider.Name())
}

// Test getting non-existent provider fails
func TestProviderFactory_GetNonExistent(t *testing.T) {
	factory := NewProviderFactory()
	
	_, err := factory.Get("nonexistent")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

// Test provider interface methods
func TestProviderInterface(t *testing.T) {
	provider := &MockProvider{
		name:            "test",
		authenticated:   true,
		capabilities: ProviderCapabilities{
			Streaming:      true,
			FunctionCall:   false,
			MaxTokens:      4096,
			SupportsImages: false,
		},
	}
	
	// Test Name
	assert.Equal(t, "test", provider.Name())
	
	// Test IsAuthenticated
	assert.True(t, provider.IsAuthenticated())
	
	// Test Capabilities
	caps := provider.Capabilities()
	assert.True(t, caps.Streaming)
	assert.False(t, caps.FunctionCall)
	assert.Equal(t, 4096, caps.MaxTokens)
	
	// Test Query
	ctx := context.Background()
	opts := QueryOptions{
		MaxTokens:   100,
		Temperature: 0.7,
	}
	
	resp, err := provider.Query(ctx, "test prompt", opts)
	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotEmpty(t, resp.Content)
}

// Test RequiresAuth for unauthenticated provider
func TestProviderInterface_RequiresAuth(t *testing.T) {
	provider := &MockProvider{
		name:          "test",
		authenticated: false,
	}
	
	assert.False(t, provider.IsAuthenticated())
	
	authInfo := provider.RequiresAuth()
	assert.False(t, authInfo.IsConfigured)
	assert.NotEmpty(t, authInfo.Instructions)
}

// Test Query with context cancellation
func TestProviderInterface_QueryWithCancellation(t *testing.T) {
	provider := &MockProvider{
		name:          "test",
		authenticated: true,
		queryDelay:    100 * time.Millisecond,
	}
	
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()
	
	opts := QueryOptions{}
	_, err := provider.Query(ctx, "test", opts)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "context")
}

// Test Response structure
func TestResponse(t *testing.T) {
	resp := &Response{
		Content:  "Test response",
		Provider: "test-provider",
		TokensUsed: TokenUsage{
			Prompt:     10,
			Completion: 20,
			Total:      30,
		},
		Metadata: map[string]interface{}{
			"model": "test-model",
		},
	}
	
	assert.Equal(t, "Test response", resp.Content)
	assert.Equal(t, "test-provider", resp.Provider)
	assert.Equal(t, 30, resp.TokensUsed.Total)
	assert.Equal(t, "test-model", resp.Metadata["model"])
}

// Test ProviderManager with fallback
func TestProviderManager_QueryWithFallback(t *testing.T) {
	factory := NewProviderFactory()
	
	// Primary provider that will fail
	primaryProvider := &MockProvider{
		name:          "primary",
		authenticated: false,
	}
	
	// Fallback provider that will succeed
	fallbackProvider := &MockProvider{
		name:          "fallback",
		authenticated: true,
	}
	
	factory.Register("primary", primaryProvider)
	factory.Register("fallback", fallbackProvider)
	
	// Create manager
	manager := NewProviderManager(factory, "primary", "fallback")
	
	// Query should use fallback since primary is not authenticated
	ctx := context.Background()
	opts := QueryOptions{}
	
	resp, err := manager.Query(ctx, "test prompt", opts)
	require.NoError(t, err)
	assert.Equal(t, "fallback", resp.Provider)
}

// Test ProviderManager when both fail
func TestProviderManager_QueryBothFail(t *testing.T) {
	factory := NewProviderFactory()
	
	primaryProvider := &MockProvider{
		name:          "primary",
		authenticated: false,
	}
	
	fallbackProvider := &MockProvider{
		name:          "fallback",
		authenticated: false,
	}
	
	factory.Register("primary", primaryProvider)
	factory.Register("fallback", fallbackProvider)
	
	manager := NewProviderManager(factory, "primary", "fallback")
	
	ctx := context.Background()
	opts := QueryOptions{}
	
	_, err := manager.Query(ctx, "test", opts)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "all providers failed")
}

// Test ProviderManager with primary success
func TestProviderManager_PrimarySuccess(t *testing.T) {
	factory := NewProviderFactory()
	
	primaryProvider := &MockProvider{
		name:          "primary",
		authenticated: true,
	}
	
	fallbackProvider := &MockProvider{
		name:          "fallback",
		authenticated: true,
	}
	
	factory.Register("primary", primaryProvider)
	factory.Register("fallback", fallbackProvider)
	
	manager := NewProviderManager(factory, "primary", "fallback")
	
	ctx := context.Background()
	opts := QueryOptions{}
	
	resp, err := manager.Query(ctx, "test", opts)
	require.NoError(t, err)
	assert.Equal(t, "primary", resp.Provider)
}

// Test CheckAuthentication
func TestProviderManager_CheckAuthentication(t *testing.T) {
	factory := NewProviderFactory()
	
	authenticatedProvider := &MockProvider{
		name:          "auth",
		authenticated: true,
	}
	
	unauthenticatedProvider := &MockProvider{
		name:          "unauth",
		authenticated: false,
	}
	
	factory.Register("auth", authenticatedProvider)
	factory.Register("unauth", unauthenticatedProvider)
	
	manager := NewProviderManager(factory, "auth", "unauth")
	
	authenticated, unauthenticated := manager.CheckAuthentication()
	
	assert.Contains(t, authenticated, "auth")
	assert.Contains(t, unauthenticated, "unauth")
}

// Mock provider for testing
type MockProvider struct {
	name          string
	authenticated bool
	capabilities  ProviderCapabilities
	queryDelay    time.Duration
}

func (m *MockProvider) Name() string {
	return m.name
}

func (m *MockProvider) Query(ctx context.Context, prompt string, opts QueryOptions) (*Response, error) {
	// Simulate delay if set
	if m.queryDelay > 0 {
		select {
		case <-time.After(m.queryDelay):
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}
	
	return &Response{
		Content:  "Mock response for: " + prompt,
		Provider: m.name,
		TokensUsed: TokenUsage{
			Prompt:     len(prompt),
			Completion: 50,
			Total:      len(prompt) + 50,
		},
	}, nil
}

func (m *MockProvider) IsAuthenticated() bool {
	return m.authenticated
}

func (m *MockProvider) RequiresAuth() AuthInfo {
	if m.authenticated {
		return AuthInfo{
			IsConfigured: true,
		}
	}
	
	return AuthInfo{
		Type:         "test",
		IsConfigured: false,
		HelpURL:      "https://test.com/auth",
		Instructions: "Test authentication instructions",
	}
}

func (m *MockProvider) Capabilities() ProviderCapabilities {
	return m.capabilities
}
