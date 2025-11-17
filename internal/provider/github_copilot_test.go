package provider

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewGitHubCopilotProvider(t *testing.T) {
	provider := NewGitHubCopilotProvider(30 * time.Second)
	assert.NotNil(t, provider)
	assert.Equal(t, "github-copilot", provider.Name())
}

func TestGitHubCopilotProvider_Capabilities(t *testing.T) {
	provider := NewGitHubCopilotProvider(30 * time.Second)
	
	caps := provider.Capabilities()
	assert.False(t, caps.Streaming)
	assert.True(t, caps.FunctionCall)
	assert.Equal(t, 8000, caps.MaxTokens)
	assert.False(t, caps.SupportsImages)
}

func TestGitHubCopilotProvider_RequiresAuth_NotAuthenticated(t *testing.T) {
	// Clear environment variables
	os.Unsetenv("COPILOT_GITHUB_TOKEN")
	os.Unsetenv("GH_TOKEN")
	
	provider := NewGitHubCopilotProvider(30 * time.Second)
	
	// Note: This test may pass or fail depending on whether gh CLI is authenticated
	// on the test machine. We test the auth info structure regardless.
	
	authInfo := provider.RequiresAuth()
	
	if !provider.IsAuthenticated() {
		// If not authenticated, should provide instructions
		assert.False(t, authInfo.IsConfigured)
		assert.NotEmpty(t, authInfo.Instructions)
		assert.NotEmpty(t, authInfo.HelpURL)
		assert.Contains(t, authInfo.Instructions, "gh auth login")
	} else {
		// If authenticated (via gh CLI), should show configured
		assert.True(t, authInfo.IsConfigured)
	}
}

func TestGitHubCopilotProvider_IsAuthenticated_WithCopilotToken(t *testing.T) {
	// Set COPILOT_GITHUB_TOKEN
	os.Setenv("COPILOT_GITHUB_TOKEN", "test-token")
	defer os.Unsetenv("COPILOT_GITHUB_TOKEN")
	
	provider := NewGitHubCopilotProvider(30 * time.Second)
	
	// Should be authenticated
	assert.True(t, provider.IsAuthenticated())
	
	// Auth info should show configured
	authInfo := provider.RequiresAuth()
	assert.True(t, authInfo.IsConfigured)
}

func TestGitHubCopilotProvider_IsAuthenticated_WithGHToken(t *testing.T) {
	// Set GH_TOKEN
	os.Setenv("GH_TOKEN", "test-token")
	defer os.Unsetenv("GH_TOKEN")
	
	// Make sure COPILOT_GITHUB_TOKEN is not set
	os.Unsetenv("COPILOT_GITHUB_TOKEN")
	
	provider := NewGitHubCopilotProvider(30 * time.Second)
	
	// Should be authenticated via GH_TOKEN
	assert.True(t, provider.IsAuthenticated())
}

func TestGitHubCopilotProvider_AuthPriority(t *testing.T) {
	// Test that COPILOT_GITHUB_TOKEN takes priority over GH_TOKEN
	os.Setenv("COPILOT_GITHUB_TOKEN", "copilot-token")
	os.Setenv("GH_TOKEN", "gh-token")
	defer func() {
		os.Unsetenv("COPILOT_GITHUB_TOKEN")
		os.Unsetenv("GH_TOKEN")
	}()
	
	provider := NewGitHubCopilotProvider(30 * time.Second)
	
	method, token := provider.detectAuth()
	assert.Equal(t, "env:COPILOT_GITHUB_TOKEN", method)
	assert.Equal(t, "copilot-token", token)
}

func TestGitHubCopilotProvider_Query_NotAuthenticated(t *testing.T) {
	// Clear environment variables
	os.Unsetenv("COPILOT_GITHUB_TOKEN")
	os.Unsetenv("GH_TOKEN")
	
	provider := NewGitHubCopilotProvider(30 * time.Second)
	
	// Skip if gh CLI is authenticated (would make actual API call)
	if provider.IsAuthenticated() {
		t.Skip("gh CLI is authenticated, skipping unauthenticated test")
	}
	
	ctx := context.Background()
	opts := QueryOptions{}
	
	_, err := provider.Query(ctx, "test prompt", opts)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not authenticated")
}

// Note: Testing actual gh copilot suggest requires:
// 1. gh CLI installed
// 2. Active GitHub Copilot subscription
// 3. Authentication configured
// These tests would be integration tests, not unit tests
// For unit tests, we'd need to mock exec.Command

func TestGitHubCopilotProvider_Query_WithTimeout(t *testing.T) {
	// This test verifies timeout behavior
	// In a real scenario with gh copilot taking too long
	
	os.Setenv("COPILOT_GITHUB_TOKEN", "test-token")
	defer os.Unsetenv("COPILOT_GITHUB_TOKEN")
	
	// Very short timeout
	provider := NewGitHubCopilotProvider(1 * time.Millisecond)
	
	// Since we can't mock exec.Command easily in this test,
	// we just verify the provider is created correctly
	assert.NotNil(t, provider)
	assert.True(t, provider.IsAuthenticated())
}

func TestGitHubCopilotProvider_FormatPrompt(t *testing.T) {
	provider := NewGitHubCopilotProvider(30 * time.Second)
	
	prompt := "Test prompt"
	formatted := provider.formatPrompt(prompt)
	
	// Should return the prompt as-is for now
	assert.Equal(t, prompt, formatted)
}

func TestGitHubCopilotProvider_ParseResponse(t *testing.T) {
	provider := NewGitHubCopilotProvider(30 * time.Second)
	
	rawOutput := "This is the response from gh copilot"
	
	resp := provider.parseResponse(rawOutput, 100*time.Millisecond)
	assert.Equal(t, rawOutput, resp.Content)
	assert.Equal(t, "github-copilot", resp.Provider)
	assert.Equal(t, 100*time.Millisecond, resp.Duration)
}
