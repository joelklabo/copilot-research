package cmd_test

import (
	"path/filepath"
	"testing"

	"github.com/joelklabo/copilot-research/cmd" // Import the cmd package
	"github.com/stretchr/testify/assert"
)

func TestRootCommand(t *testing.T) {
	assert.NotNil(t, cmd.RootCmd)
	assert.Equal(t, "copilot-research", cmd.RootCmd.Use)
	assert.Contains(t, cmd.RootCmd.Short, "research")
	assert.NotEmpty(t, cmd.RootCmd.Long)
}

func TestGlobalFlags(t *testing.T) {
	tests := []struct {
		name     string
		flagName string
	}{
		{"config", "config"},
		{"output", "output"},
		{"quiet", "quiet"},
		{"json", "json"},
		{"mode", "mode"},
		{"prompt", "prompt"},
		{"no-store", "no-store"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			flag := cmd.RootCmd.PersistentFlags().Lookup(tt.flagName)
			assert.NotNil(t, flag, "Flag %s should exist", tt.flagName)
		})
	}
}

func TestVersionFlag(t *testing.T) {
	assert.NotEmpty(t, cmd.RootCmd.Version)
}

func TestGetKnowledgeDir(t *testing.T) {
	dir := cmd.GetKnowledgeDir()
	assert.NotEmpty(t, dir)
	assert.Contains(t, dir, ".copilot-research")
	assert.Contains(t, dir, "knowledge")
}

func TestGetConfigFile(t *testing.T) {
	// Save original CfgFile and restore after test
	originalCfgFile := cmd.CfgFile
	defer func() { cmd.CfgFile = originalCfgFile }()

	// Reset and initialize
	cmd.CfgFile = ""
	cmd.InitConfig() // Call the exported InitConfig
	
	assert.NotEmpty(t, cmd.CfgFile)
	assert.Contains(t, cmd.CfgFile, ".copilot-research")
	assert.Contains(t, cmd.CfgFile, "config.yaml")
}

func TestConfigFileCustomPath(t *testing.T) {
	// Save original CfgFile and restore after test
	originalCfgFile := cmd.CfgFile
	defer func() { cmd.CfgFile = originalCfgFile }()

	// Create a temporary directory for the custom config file
	tmpDir := t.TempDir()
	customPath := filepath.Join(tmpDir, "custom_config.yaml")

	// Set custom path
	cmd.CfgFile = customPath
	cmd.InitConfig() // Call the exported InitConfig
	
	// Should now be the custom path
	assert.Equal(t, customPath, cmd.CfgFile)
}