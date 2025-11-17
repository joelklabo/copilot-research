package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHistoryCommand(t *testing.T) {
	assert.NotNil(t, researchHistoryCmd)
	assert.Contains(t, researchHistoryCmd.Use, "history")
	assert.NotEmpty(t, researchHistoryCmd.Short)
}

func TestHistoryCommand_Flags(t *testing.T) {
	flags := []string{"search", "mode", "id", "clear", "limit"}
	
	for _, flagName := range flags {
		t.Run(flagName, func(t *testing.T) {
			flag := researchHistoryCmd.Flags().Lookup(flagName)
			assert.NotNil(t, flag, "Flag %s should exist", flagName)
		})
	}
}

func TestFormatSession(t *testing.T) {
	// Test session formatting
	result := formatSessionSummary(1, "Test query", "quick", "2025-11-17")
	assert.Contains(t, result, "Test query")
	assert.Contains(t, result, "quick")
}

func TestFormatDuration(t *testing.T) {
	tests := []struct {
		name     string
		seconds  int64
		expected string
	}{
		{"seconds", 45, "45s"},
		{"minutes", 90, "1m 30s"},
		{"hours", 3665, "1h 1m"},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatDuration(tt.seconds)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestTruncateString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		maxLen   int
		expected string
	}{
		{"short", "Hello", 10, "Hello"},
		{"exact", "Hello", 5, "Hello"},
		{"long", "Hello World", 8, "Hello..."},
		{"empty", "", 5, ""},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := truncateString(tt.input, tt.maxLen)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestValidateClearConfirmation(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"yes", "yes", true},
		{"y", "y", true},
		{"Y", "Y", true},
		{"YES", "YES", true},
		{"no", "no", false},
		{"n", "n", false},
		{"other", "maybe", false},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := validateClearConfirmation(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}
