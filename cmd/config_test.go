package cmd_test

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/joelklabo/copilot-research/cmd" // Import the cmd package
	"github.com/joelklabo/copilot-research/internal/config"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func TestConfigShowCmd(t *testing.T) {
	// Create a temporary config directory and file
	tempDir := t.TempDir()
	tempConfigFile := filepath.Join(tempDir, "config.yaml")

	// Set cmd.CfgFile to the temporary config file
	oldCfgFile := cmd.CfgFile
	cmd.CfgFile = tempConfigFile
	defer func() { cmd.CfgFile = oldCfgFile }()

	// Initialize config with default values and save it
	cmd.AppConfig = config.DefaultConfig()
	err := config.SaveConfig(tempConfigFile, cmd.AppConfig)
	assert.NoError(t, err)

	// Redirect stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Execute the command
	cmd.ConfigShowCmd.Run(cmd.ConfigShowCmd, []string{})

	// Restore stdout
	w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	_, err = buf.ReadFrom(r)
	assert.NoError(t, err)

	// Expected output (YAML representation of default config)
	expectedConfig, err := yaml.Marshal(config.DefaultConfig())
	assert.NoError(t, err)

	// Compare actual output with expected output
	assert.Equal(t, strings.TrimSpace(string(expectedConfig)), strings.TrimSpace(buf.String()))
}
