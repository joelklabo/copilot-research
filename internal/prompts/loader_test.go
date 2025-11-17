package prompts

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDefaultPromptExists(t *testing.T) {
	// This test verifies that the default prompt template exists
	// and has the required structure
	
	promptPath := filepath.Join("..", "..", "prompts", "default.md")
	
	// Test that file exists
	_, err := os.Stat(promptPath)
	require.NoError(t, err, "default.md should exist in prompts directory")
	
	// Read the file
	content, err := os.ReadFile(promptPath)
	require.NoError(t, err, "should be able to read default.md")
	
	contentStr := string(content)
	
	// Test that it contains frontmatter
	assert.True(t, strings.HasPrefix(contentStr, "---"), "should start with frontmatter delimiter")
	assert.Contains(t, contentStr, "name:", "should have name field in frontmatter")
	assert.Contains(t, contentStr, "description:", "should have description field")
	assert.Contains(t, contentStr, "version:", "should have version field")
	
	// Test that it has template variables
	assert.Contains(t, contentStr, "{{query}}", "should have {{query}} template variable")
	assert.Contains(t, contentStr, "{{mode}}", "should have {{mode}} template variable")
	
	// Test that it has key sections
	assert.Contains(t, contentStr, "### Overview", "should have Overview section")
	assert.Contains(t, contentStr, "### Key Concepts", "should have Key Concepts section")
	assert.Contains(t, contentStr, "### Examples", "should have Examples section")
	assert.Contains(t, contentStr, "### Best Practices", "should have Best Practices section")
	assert.Contains(t, contentStr, "### Resources", "should have Resources section")
}

func TestDefaultPromptFormat(t *testing.T) {
	// Test that the prompt produces properly formatted markdown
	promptPath := filepath.Join("..", "..", "prompts", "default.md")
	content, err := os.ReadFile(promptPath)
	require.NoError(t, err)
	
	contentStr := string(content)
	
	// Should have markdown headers
	assert.Contains(t, contentStr, "##", "should use markdown headers")
	assert.Contains(t, contentStr, "###", "should use sub-headers")
	
	// Should encourage structured output
	assert.Contains(t, contentStr, "Markdown", "should mention Markdown format")
	assert.Contains(t, contentStr, "structure", "should emphasize structure")
}

func TestDefaultPromptVariables(t *testing.T) {
	// Verify all expected template variables are present
	promptPath := filepath.Join("..", "..", "prompts", "default.md")
	content, err := os.ReadFile(promptPath)
	require.NoError(t, err)
	
	contentStr := string(content)
	
	// Required variables
	requiredVars := []string{
		"{{query}}",
		"{{mode}}",
	}
	
	for _, v := range requiredVars {
		assert.Contains(t, contentStr, v, "should contain variable %s", v)
	}
}
