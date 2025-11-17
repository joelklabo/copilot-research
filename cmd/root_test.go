package cmd

import (
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
	// Reset and initialize
	cfgFile = ""
	initConfig()
	
	assert.NotEmpty(t, cfgFile)
	assert.Contains(t, cfgFile, ".copilot-research")
	assert.Contains(t, cfgFile, "config.yaml")
}

func TestConfigFileCustomPath(t *testing.T) {
	// Set custom path
	cfgFile = "/custom/path/config.yaml"
	initConfig()
	
	// Should not change
	assert.Equal(t, "/custom/path/config.yaml", cfgFile)
}
