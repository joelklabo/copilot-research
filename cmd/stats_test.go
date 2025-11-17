package cmd

import (
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/joelklabo/copilot-research/internal/config"
	"github.com/joelklabo/copilot-research/internal/db"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRunStats(t *testing.T) {
	// Save original global variables and defer their restoration
	oldAppConfig := AppConfig
	defer func() {
		AppConfig = oldAppConfig
	}()

	// Create mock config
	mockConfig := config.DefaultConfig()
	AppConfig = mockConfig

	// Create a mock database
	mockDB := &db.MockDB{
		GetTotalSessionsFunc: func() (int, error) {
			return 3, nil
		},
		GetModeStatsFunc: func() (map[string]int, error) {
			return map[string]int{
				"quick": 2,
				"deep":  1,
			}, nil
		},
		CloseFunc: func() error {
			return nil
		},
	}

	// Create a dummy db file for os.Stat to work
	tmpDir := t.TempDir()
	dummyDbPath := filepath.Join(tmpDir, ".copilot-research", "research.db")
	err := os.MkdirAll(filepath.Dir(dummyDbPath), 0755)
	require.NoError(t, err)
	err = os.WriteFile(dummyDbPath, []byte("dummy db content"), 0644)
	require.NoError(t, err)

	// Capture stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Execute the command with the mock DB
	err = _runStats(mockDB, dummyDbPath)
	require.NoError(t, err)

	w.Close()
	out, _ := io.ReadAll(r)
	os.Stdout = oldStdout // Restore stdout

	output := string(out)

	// Assertions for output content
	assert.Contains(t, output, "Research Statistics")
	assert.Contains(t, output, "Total Sessions: 3")
	assert.Contains(t, output, "Database Size:") // Will check formatBytes output
	assert.Contains(t, output, "Mode Usage:")
	assert.Contains(t, output, "quick   2 (67%)")
	assert.Contains(t, output, "deep    1 (33%)")
}

func TestFormatBytes(t *testing.T) {
	tests := []struct {
		name string
		input int64
		expected string
	}{
		{"bytes", 100, "100 B"},
		{"kilobytes", 1024, "1.0 KB"},
		{"megabytes", 1024 * 1024, "1.0 MB"},
		{"gigabytes", 1024 * 1024 * 1024, "1.0 GB"},
		{"terabytes", 1024 * 1024 * 1024 * 1024, "1.0 TB"},
		{"petabytes", 1024 * 1024 * 1024 * 1024 * 1024, "1.0 PB"},
		{"large bytes", 1500, "1.5 KB"},
		{"large megabytes", 2.5 * 1024 * 1024, "2.5 MB"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, formatBytes(tt.input))
		})
	}
}