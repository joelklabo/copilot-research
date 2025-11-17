package cmd

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRootCommand(t *testing.T) {
	assert.NotNil(t, rootCmd)
	assert.Equal(t, "copilot-research", rootCmd.Use)
	assert.Contains(t, rootCmd.Short, "research")
	assert.NotEmpty(t, rootCmd.Long)
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
			flag := rootCmd.PersistentFlags().Lookup(tt.flagName)
			assert.NotNil(t, flag, "Flag %s should exist", tt.flagName)
		})
	}
}

func TestVersionFlag(t *testing.T) {
	assert.NotEmpty(t, rootCmd.Version)
}

func TestGetKnowledgeDir(t *testing.T) {
	dir := GetKnowledgeDir()
	assert.NotEmpty(t, dir)
	assert.Contains(t, dir, ".copilot-research")
	assert.Contains(t, dir, "knowledge")
}

func TestGetConfigFile(t *testing.T) {
	// Save original cfgFile and restore after test
	originalCfgFile := cfgFile
	defer func() { cfgFile = originalCfgFile }()

	// Reset and initialize
	cfgFile = ""
	initConfig()
	
	assert.NotEmpty(t, cfgFile)
	assert.Contains(t, cfgFile, ".copilot-research")
	assert.Contains(t, cfgFile, "config.yaml")
}

func TestConfigFileCustomPath(t *testing.T) {
	// Save original cfgFile and restore after test
	originalCfgFile := cfgFile
	defer func() { cfgFile = originalCfgFile }()

	// Create a temporary directory for the custom config file
	tmpDir := t.TempDir()
	customPath := filepath.Join(tmpDir, "custom_config.yaml")

	// Set custom path
	cfgFile = customPath
	initConfig()
	
	// Should now be the custom path
	assert.Equal(t, customPath, cfgFile)
}