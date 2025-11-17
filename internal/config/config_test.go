package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	assert.NotNil(t, cfg)
	assert.Equal(t, "github-copilot", cfg.Providers.Primary)
	assert.Equal(t, "openai", cfg.Providers.Fallback)
	assert.True(t, cfg.Providers.GitHubCopilot.Enabled)
	assert.Equal(t, "gpt-4o", cfg.Providers.OpenAI.Model)
	assert.Equal(t, 0.7, cfg.Providers.OpenAI.Temperature)
	assert.True(t, cfg.Providers.AutoFallback)
}

func TestLoadConfig_NewFile(t *testing.T) {
	tmpDir := t.TempDir()
	cfgPath := filepath.Join(tmpDir, "config.yaml")

	cfg, err := LoadConfig(cfgPath)
	require.NoError(t, err)
	assert.NotNil(t, cfg)

	// Verify default config was saved
	_, err = os.Stat(cfgPath)
	assert.False(t, os.IsNotExist(err))

	// Check some default values
	assert.Equal(t, "github-copilot", cfg.Providers.Primary)
}

func TestLoadConfig_ExistingFile(t *testing.T) {
	tmpDir := t.TempDir()
	cfgPath := filepath.Join(tmpDir, "config.yaml")

	// Create a custom config file
	customConfigContent := `
providers:
  primary: anthropic
  fallback: github-copilot
  openai:
    model: gpt-3.5-turbo
    max_tokens: 2000
`
	err := os.WriteFile(cfgPath, []byte(customConfigContent), 0644)
	require.NoError(t, err)

	cfg, err := LoadConfig(cfgPath)
	require.NoError(t, err)
	assert.NotNil(t, cfg)

	// Verify custom values
	assert.Equal(t, "anthropic", cfg.Providers.Primary)
	assert.Equal(t, "github-copilot", cfg.Providers.Fallback)
	assert.Equal(t, "gpt-3.5-turbo", cfg.Providers.OpenAI.Model)
	assert.Equal(t, 2000, cfg.Providers.OpenAI.MaxTokens)
	// Verify default values for fields not specified in custom config
	assert.Equal(t, 60*time.Second, cfg.Providers.GitHubCopilot.Timeout)
}

func TestLoadConfig_InvalidFile(t *testing.T) {
	tmpDir := t.TempDir()
	cfgPath := filepath.Join(tmpDir, "config.yaml")

	// Create an invalid config file
	invalidConfigContent := `
providers:
  primary: anthropic
  fallback: github-copilot
  openai:
    model: gpt-3.5-turbo
    max_tokens: 2000
  invalid_field: [
`
	err := os.WriteFile(cfgPath, []byte(invalidConfigContent), 0644)
	require.NoError(t, err)

	_, err = LoadConfig(cfgPath)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to unmarshal config")
}

func TestSaveConfig(t *testing.T) {
	tmpDir := t.TempDir()
	cfgPath := filepath.Join(tmpDir, "config.yaml")

	cfg := DefaultConfig()
	cfg.Providers.Primary = "test-primary"
	cfg.Providers.OpenAI.Model = "test-model"

	err := SaveConfig(cfgPath, cfg)
	require.NoError(t, err)

	// Load and verify
	loadedCfg, err := LoadConfig(cfgPath)
	require.NoError(t, err)
	assert.Equal(t, "test-primary", loadedCfg.Providers.Primary)
	assert.Equal(t, "test-model", loadedCfg.Providers.OpenAI.Model)
}

func TestSaveConfig_DirectoryCreation(t *testing.T) {
	tmpDir := t.TempDir()
	nestedDir := filepath.Join(tmpDir, "nested", "sub")
	cfgPath := filepath.Join(nestedDir, "config.yaml")

	cfg := DefaultConfig()
	err := SaveConfig(cfgPath, cfg)
	require.NoError(t, err)

	// Verify directory was created
	_, err = os.Stat(nestedDir)
	assert.False(t, os.IsNotExist(err))
}
