package provider

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

// GitHubCopilotProvider implements the AIProvider interface for GitHub Copilot
type GitHubCopilotProvider struct {
	timeout    time.Duration
	authMethod string
	token      string
}

// NewGitHubCopilotProvider creates a new GitHub Copilot provider
func NewGitHubCopilotProvider(timeout time.Duration) *GitHubCopilotProvider {
	return &GitHubCopilotProvider{
		timeout: timeout,
	}
}

// Name returns the provider name
func (g *GitHubCopilotProvider) Name() string {
	return "github-copilot"
}

// Query executes a query using gh copilot suggest
func (g *GitHubCopilotProvider) Query(ctx context.Context, prompt string, opts QueryOptions) (*Response, error) {
	// Check authentication first
	if !g.IsAuthenticated() {
		return nil, fmt.Errorf("not authenticated: please run 'gh auth login' or set COPILOT_GITHUB_TOKEN")
	}
	
	// Format the prompt
	formattedPrompt := g.formatPrompt(prompt)
	
	// Create context with timeout
	queryCtx, cancel := context.WithTimeout(ctx, g.timeout)
	defer cancel()
	
	// Execute gh copilot suggest
	start := time.Now()
	cmd := exec.CommandContext(queryCtx, "gh", "copilot", "suggest", formattedPrompt)
	
	// Set environment if we have a token
	if g.token != "" {
		cmd.Env = append(os.Environ(), fmt.Sprintf("GH_TOKEN=%s", g.token))
	}
	
	output, err := cmd.CombinedOutput()
	duration := time.Since(start)
	
	if err != nil {
		// Check for timeout
		if queryCtx.Err() == context.DeadlineExceeded {
			return nil, fmt.Errorf("query timeout after %v", g.timeout)
		}
		
		// Parse error message for helpful feedback
		errorMsg := string(output)
		if strings.Contains(errorMsg, "not authenticated") || strings.Contains(errorMsg, "authentication") {
			return nil, fmt.Errorf("GitHub Copilot authentication failed: %w", err)
		}
		if strings.Contains(errorMsg, "subscription") {
			return nil, fmt.Errorf("GitHub Copilot subscription required: %w", err)
		}
		
		return nil, fmt.Errorf("gh copilot suggest failed: %w\nOutput: %s", err, errorMsg)
	}
	
	// Parse and return response
	return g.parseResponse(string(output), duration), nil
}

// IsAuthenticated checks if the provider is authenticated
func (g *GitHubCopilotProvider) IsAuthenticated() bool {
	method, token := g.detectAuth()
	g.authMethod = method
	g.token = token
	return method != "none"
}

// RequiresAuth returns authentication information
func (g *GitHubCopilotProvider) RequiresAuth() AuthInfo {
	if g.IsAuthenticated() {
		return AuthInfo{
			Type:         g.authMethod,
			IsConfigured: true,
		}
	}
	
	return AuthInfo{
		Type:         "oauth-device-flow",
		IsConfigured: false,
		HelpURL:      "https://github.com/features/copilot",
		Instructions: `GitHub Copilot authentication required.

Please authenticate using one of these methods:

1. GitHub CLI (recommended):
   gh auth login
   
2. Personal Access Token:
   export COPILOT_GITHUB_TOKEN=ghp_your_token_here
   
3. Set GH_TOKEN:
   export GH_TOKEN=ghp_your_token_here

Note: You need an active GitHub Copilot subscription.
Get one at https://github.com/features/copilot

Once authenticated, run your command again.`,
	}
}

// Capabilities returns the provider's capabilities
func (g *GitHubCopilotProvider) Capabilities() ProviderCapabilities {
	return ProviderCapabilities{
		Streaming:      false,
		FunctionCall:   true, // Via MCP
		MaxTokens:      8000,
		SupportsImages: false,
	}
}

// detectAuth checks authentication in priority order
func (g *GitHubCopilotProvider) detectAuth() (string, string) {
	// 1. Check COPILOT_GITHUB_TOKEN
	if token := os.Getenv("COPILOT_GITHUB_TOKEN"); token != "" {
		return "env:COPILOT_GITHUB_TOKEN", token
	}
	
	// 2. Check GH_TOKEN
	if token := os.Getenv("GH_TOKEN"); token != "" {
		return "env:GH_TOKEN", token
	}
	
	// 3. Check gh CLI authentication
	cmd := exec.Command("gh", "auth", "status")
	if err := cmd.Run(); err == nil {
		return "gh-cli", ""
	}
	
	return "none", ""
}

// formatPrompt formats the prompt for gh copilot
func (g *GitHubCopilotProvider) formatPrompt(prompt string) string {
	// For now, return as-is
	// Could add additional formatting or preprocessing here
	return prompt
}

// parseResponse parses the gh copilot output into a Response
func (g *GitHubCopilotProvider) parseResponse(output string, duration time.Duration) *Response {
	// gh copilot suggest returns markdown output
	// We'll clean it up and structure it
	
	content := strings.TrimSpace(output)
	
	return &Response{
		Content:  content,
		Provider: "github-copilot",
		Model:    "gpt-4", // GitHub Copilot uses GPT-4
		Duration: duration,
		TokensUsed: TokenUsage{
			// gh copilot doesn't provide token usage info
			// We could estimate based on content length
			Total: len(content) / 4, // Rough estimate: 4 chars per token
		},
		Metadata: map[string]interface{}{
			"auth_method": g.authMethod,
		},
	}
}
