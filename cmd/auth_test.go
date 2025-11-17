package cmd

import (
	"testing"

	"github.com/spf13/cobra" // Added this import
	"github.com/stretchr/testify/assert"
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
