package cmd

import (
	"context"
	"io"
	"os"
	"testing"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/joelklabo/copilot-research/internal/config"
	"github.com/joelklabo/copilot-research/internal/provider"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAuthCommand(t *testing.T) {
	assert.NotNil(t, authCmd)
	assert.Equal(t, "auth", authCmd.Use)
	assert.Contains(t, authCmd.Short, "authentication")
	assert.NotEmpty(t, authCmd.Long)
}

func TestAuthSubcommands(t *testing.T) {
	tests := []struct {
		name    string
		command *cobra.Command
		use     string
		short   string
	}{
		{
			name:    "status",
			command: authStatusCommand,
			use:     "status",
			short:   "Show authentication status for all providers",
		},
		{
			name:    "login",
			command: authLoginCommand,
			use:     "login [provider]",
			short:   "Interactively authenticate with a provider",
		},
		{
			name:    "test",
			command: authTestCommand,
			use:     "test [provider]",
			short:   "Test connectivity and authentication for a provider",
		},
		{
			name:    "logout",
			command: authLogoutCommand,
			use:     "logout [provider]",
			short:   "Clear authentication credentials for a provider",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.NotNil(t, tt.command)
			assert.Equal(t, tt.use, tt.command.Use)
			assert.Equal(t, tt.short, tt.command.Short)
		})
	}
}

// MockProvider is a mock implementation of the AIProvider interface for testing
type MockProvider struct {
	name          string
	authenticated bool
	authInfo      provider.AuthInfo
	queryFunc     func(ctx context.Context, prompt string, opts provider.QueryOptions) (*provider.Response, error)
	capabilities  provider.ProviderCapabilities
}

func (m *MockProvider) Name() string { return m.name }
func (m *MockProvider) Query(ctx context.Context, prompt string, opts provider.QueryOptions) (*provider.Response, error) {
	if m.queryFunc != nil {
		return m.queryFunc(ctx, prompt, opts)
	}
	return &provider.Response{Content: "Mock response"}, nil
}
func (m *MockProvider) IsAuthenticated() bool { return m.authenticated }
func (m *MockProvider) RequiresAuth() provider.AuthInfo { return m.authInfo }
func (m *MockProvider) Capabilities() provider.ProviderCapabilities {
	return provider.ProviderCapabilities{MaxTokens: 1000} // Hardcoded for now
}

func TestRunAuthStatus(t *testing.T) {
	// Dummy usage to satisfy linter for time import
	var _ time.Duration

	// Save original global variables and defer their restoration
	oldAppConfig := AppConfig
	oldAppProviderManager := AppProviderManager
	defer func() {
		AppConfig = oldAppConfig
		AppProviderManager = oldAppProviderManager
	}()

	// Create mock config
	mockConfig := config.DefaultConfig()
	mockConfig.Providers.Primary = "mock-authenticated"
	mockConfig.Providers.Fallback = "mock-unauthenticated"
	AppConfig = mockConfig

	// Create mock providers
	mockAuthProvider := &MockProvider{
		name:          "mock-authenticated",
		authenticated: true,
		authInfo: provider.AuthInfo{
			Type:         "cli",
			IsConfigured: true,
		},
		capabilities: provider.ProviderCapabilities{MaxTokens: 1000},
	}
	mockUnauthProvider := &MockProvider{
		name:          "mock-unauthenticated",
		authenticated: false,
		authInfo: provider.AuthInfo{
			Type:         "apikey",
			IsConfigured: false,
			Instructions: "Please set MOCK_API_KEY",
			HelpURL:      "http://mock.com/help",
		},
		capabilities: provider.ProviderCapabilities{MaxTokens: 500},
	}
	// Removed mockErrorProvider as it was unused

	// Create mock provider factory and manager
	mockFactory := provider.NewProviderFactory()
	mockFactory.Register(mockAuthProvider.Name(), mockAuthProvider)
	mockFactory.Register(mockUnauthProvider.Name(), mockUnauthProvider)

	// Updated call to NewProviderManager
	AppProviderManager = provider.NewProviderManager(
		mockFactory,
		mockConfig.Providers.Primary,
		mockConfig.Providers.Fallback,
		mockConfig.Providers.AutoFallback,
		mockConfig.Providers.NotifyFallback,
	)

	// Capture stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := runAuthStatus(authStatusCommand, []string{})
	require.NoError(t, err)

	w.Close()
	out, _ := io.ReadAll(r)
	os.Stdout = oldStdout // Restore stdout

	output := string(out)

	// Assertions for output content
	assert.Contains(t, output, lipgloss.NewStyle().Bold(true).Render("Authentication Status"))
	assert.Contains(t, output, lipgloss.NewStyle().Bold(true).Render("Provider"))
	assert.Contains(t, output, lipgloss.NewStyle().Bold(true).Render("Status"))
	assert.Contains(t, output, lipgloss.NewStyle().Bold(true).Render("Method"))

	assert.Contains(t, output, "mock-authenticated")
	assert.Contains(t, output, lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("42")).Render("✅ Authenticated")) // SuccessStyle
	assert.Contains(t, output, "CLI Tool")

	assert.Contains(t, output, "mock-unauthenticated")
	assert.Contains(t, output, lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("196")).Render("❌ Not Configured")) // ErrorStyle
	assert.Contains(t, output, "apikey") // Changed from "API Key (Env/Config)" to "apikey"

	assert.Contains(t, output, "Primary: mock-authenticated")
	assert.Contains(t, output, "Fallback: mock-unauthenticated")

	assert.Contains(t, output, lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("205")).Render("Authentication Required")) // TitleStyle
	assert.Contains(t, output, "To authenticate mock-unauthenticated:\nPlease set MOCK_API_KEY")
}