package prompts

import (
	"fmt"
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

func TestPromptLoader_Load(t *testing.T) {
	// Create a temp directory for test prompts
	tempDir := t.TempDir()
	
	// Create a test prompt file
	testPrompt := `---
name: test
description: Test prompt
version: 1.0.0
---

Test prompt content with {{query}} variable.`
	
	err := os.WriteFile(filepath.Join(tempDir, "test.md"), []byte(testPrompt), 0644)
	require.NoError(t, err)
	
	// Create loader
	loader := NewPromptLoader(tempDir)
	
	// Test loading
	prompt, err := loader.Load("test")
	require.NoError(t, err)
	assert.Equal(t, "test", prompt.Name)
	assert.Equal(t, "Test prompt", prompt.Description)
	assert.Equal(t, "1.0.0", prompt.Version)
	assert.Contains(t, prompt.Template, "{{query}}")
}

func TestPromptLoader_LoadDefault(t *testing.T) {
	// Load default from the actual prompts directory
	promptsDir := filepath.Join("..", "..", "prompts")
	loader := NewPromptLoader(promptsDir)
	
	prompt, err := loader.Load("default")
	require.NoError(t, err)
	assert.Equal(t, "default", prompt.Name)
	assert.Contains(t, prompt.Template, "{{query}}")
	assert.Contains(t, prompt.Template, "{{mode}}")
}

func TestPromptLoader_Cache(t *testing.T) {
	tempDir := t.TempDir()
	
	testPrompt := `---
name: cached
description: Cached prompt
version: 1.0.0
---

Cached content`
	
	err := os.WriteFile(filepath.Join(tempDir, "cached.md"), []byte(testPrompt), 0644)
	require.NoError(t, err)
	
	loader := NewPromptLoader(tempDir)
	
	// Load first time
	prompt1, err := loader.Load("cached")
	require.NoError(t, err)
	
	// Load second time (should be from cache)
	prompt2, err := loader.Load("cached")
	require.NoError(t, err)
	
	// Should be same instance
	assert.Equal(t, prompt1.Name, prompt2.Name)
}

func TestPromptLoader_Render(t *testing.T) {
	tempDir := t.TempDir()
	
	testPrompt := `---
name: render-test
description: Render test
version: 1.0.0
---

Query: {{query}}
Mode: {{mode}}
User: {{user}}`
	
	err := os.WriteFile(filepath.Join(tempDir, "render-test.md"), []byte(testPrompt), 0644)
	require.NoError(t, err)
	
	loader := NewPromptLoader(tempDir)
	prompt, err := loader.Load("render-test")
	require.NoError(t, err)
	
	// Test rendering with variables
	vars := map[string]string{
		"query": "How do actors work?",
		"mode":  "deep",
		"user":  "Alice",
	}
	
	rendered := loader.Render(prompt, vars)
	
	assert.Contains(t, rendered, "Query: How do actors work?")
	assert.Contains(t, rendered, "Mode: deep")
	assert.Contains(t, rendered, "User: Alice")
	assert.NotContains(t, rendered, "{{")
}

func TestPromptLoader_List(t *testing.T) {
	tempDir := t.TempDir()
	
	// Create multiple test prompts
	prompts := []string{"quick", "deep", "compare"}
	for _, name := range prompts {
		content := fmt.Sprintf(`---
name: %s
description: %s prompt
version: 1.0.0
---

Content for %s`, name, name, name)
		err := os.WriteFile(filepath.Join(tempDir, name+".md"), []byte(content), 0644)
		require.NoError(t, err)
	}
	
	loader := NewPromptLoader(tempDir)
	names, err := loader.List()
	require.NoError(t, err)
	
	// Should include default and all created prompts
	assert.Contains(t, names, "default")
	assert.Contains(t, names, "quick")
	assert.Contains(t, names, "deep")
	assert.Contains(t, names, "compare")
}

func TestPromptLoader_Reload(t *testing.T) {
	tempDir := t.TempDir()
	
	testPrompt := `---
name: reload-test
description: Reload test
version: 1.0.0
---

Original content`
	
	filename := filepath.Join(tempDir, "reload-test.md")
	err := os.WriteFile(filename, []byte(testPrompt), 0644)
	require.NoError(t, err)
	
	loader := NewPromptLoader(tempDir)
	
	// Load first time
	prompt1, err := loader.Load("reload-test")
	require.NoError(t, err)
	assert.Contains(t, prompt1.Template, "Original content")
	
	// Update file
	updatedPrompt := `---
name: reload-test
description: Reload test
version: 2.0.0
---

Updated content`
	
	err = os.WriteFile(filename, []byte(updatedPrompt), 0644)
	require.NoError(t, err)
	
	// Reload cache
	loader.Reload()
	
	// Load again (should get updated version)
	prompt2, err := loader.Load("reload-test")
	require.NoError(t, err)
	assert.Contains(t, prompt2.Template, "Updated content")
	assert.Equal(t, "2.0.0", prompt2.Version)
}

func TestPromptLoader_MissingPrompt(t *testing.T) {
	loader := NewPromptLoader(t.TempDir())
	
	_, err := loader.Load("nonexistent")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "nonexistent")
}

func TestPromptLoader_InvalidFrontmatter(t *testing.T) {
	tempDir := t.TempDir()
	
	// Prompt with invalid YAML
	invalidPrompt := `---
name: invalid
description: 
  - this
  - is
  - wrong: format
---

Content`
	
	err := os.WriteFile(filepath.Join(tempDir, "invalid.md"), []byte(invalidPrompt), 0644)
	require.NoError(t, err)
	
	loader := NewPromptLoader(tempDir)
	_, err = loader.Load("invalid")
	assert.Error(t, err)
}

func TestPromptLoader_MissingFrontmatter(t *testing.T) {
	tempDir := t.TempDir()
	
	// Prompt without frontmatter
	noFrontmatter := `Just content without frontmatter`
	
	err := os.WriteFile(filepath.Join(tempDir, "no-fm.md"), []byte(noFrontmatter), 0644)
	require.NoError(t, err)
	
	loader := NewPromptLoader(tempDir)
	_, err = loader.Load("no-fm")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "frontmatter")
}
