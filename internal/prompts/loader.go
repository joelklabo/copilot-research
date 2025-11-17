package prompts

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"gopkg.in/yaml.v3"
)

// Embedded default prompt content (loaded at compile time)
var embeddedDefault string

// Prompt represents a loaded prompt template
type Prompt struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	Version     string `yaml:"version"`
	Mode        string `yaml:"mode,omitempty"`
	Template    string `yaml:"-"` // The template content (not in frontmatter)
}

// PromptLoader loads and manages prompt templates
type PromptLoader struct {
	promptsDir string
	cache      map[string]*Prompt
	mu         sync.RWMutex
}

// NewPromptLoader creates a new prompt loader
func NewPromptLoader(promptsDir string) *PromptLoader {
	return &PromptLoader{
		promptsDir: promptsDir,
		cache:      make(map[string]*Prompt),
	}
}

// Load loads a prompt by name
func (l *PromptLoader) Load(name string) (*Prompt, error) {
	// Check cache first
	l.mu.RLock()
	if cached, exists := l.cache[name]; exists {
		l.mu.RUnlock()
		return cached, nil
	}
	l.mu.RUnlock()

	// Try to load from file
	filename := filepath.Join(l.promptsDir, name+".md")
	prompt, err := l.loadFromFile(filename)
	if err != nil {
		// Fall back to embedded default if loading "default"
		if name == "default" {
			return l.loadEmbeddedDefault()
		}
		return nil, fmt.Errorf("failed to load prompt '%s': %w", name, err)
	}

	// Cache it
	l.mu.Lock()
	l.cache[name] = prompt
	l.mu.Unlock()

	return prompt, nil
}

// loadFromFile loads a prompt from a file
func (l *PromptLoader) loadFromFile(filename string) (*Prompt, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	return parsePrompt(string(data))
}

// loadEmbeddedDefault loads the embedded default prompt
func (l *PromptLoader) loadEmbeddedDefault() (*Prompt, error) {
	// If embedded default is not set, load from relative path
	if embeddedDefault == "" {
		// Try loading from the repository root's prompts directory
		data, err := os.ReadFile(filepath.Join(l.promptsDir, "default.md"))
		if err != nil {
			// As last resort, return a minimal default
			return &Prompt{
				Name:        "default",
				Description: "Built-in default prompt",
				Version:     "1.0.0",
				Template:    "Research Query: {{query}}\n\nPlease provide a comprehensive answer.",
			}, nil
		}
		return parsePrompt(string(data))
	}

	return parsePrompt(embeddedDefault)
}

// parsePrompt parses a prompt file with YAML frontmatter
func parsePrompt(content string) (*Prompt, error) {
	// Split frontmatter and template content
	parts := splitFrontmatter(content)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid prompt format: missing frontmatter")
	}

	// Parse frontmatter
	var prompt Prompt
	if err := yaml.Unmarshal([]byte(parts[0]), &prompt); err != nil {
		return nil, fmt.Errorf("failed to parse frontmatter: %w", err)
	}

	// Set template content
	prompt.Template = strings.TrimSpace(parts[1])

	// Validate required fields
	if prompt.Name == "" {
		return nil, fmt.Errorf("prompt name is required")
	}

	return &prompt, nil
}

// splitFrontmatter splits content into frontmatter and body
func splitFrontmatter(content string) []string {
	lines := strings.Split(content, "\n")

	start := -1
	end := -1
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "---" {
			if start == -1 {
				start = i
			} else {
				end = i
				break
			}
		}
	}

	if start == -1 || end == -1 {
		return []string{content}
	}

	frontmatter := strings.Join(lines[start+1:end], "\n")
	body := strings.TrimSpace(strings.Join(lines[end+1:], "\n"))

	return []string{frontmatter, body}
}

// Render renders a prompt template with variables
func (l *PromptLoader) Render(prompt *Prompt, vars map[string]string) string {
	result := prompt.Template

	// Replace all variables
	for key, value := range vars {
		placeholder := fmt.Sprintf("{{%s}}", key)
		result = strings.ReplaceAll(result, placeholder, value)
	}

	return result
}

// List returns all available prompt names
func (l *PromptLoader) List() ([]string, error) {
	l.mu.RLock()
	defer l.mu.RUnlock()

	names := make([]string, 0)

	// Always include default
	names = append(names, "default")

	// List files in prompts directory
	entries, err := os.ReadDir(l.promptsDir)
	if err != nil {
		// If directory doesn't exist, just return default
		if os.IsNotExist(err) {
			return names, nil
		}
		return nil, fmt.Errorf("failed to read prompts directory: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if strings.HasSuffix(name, ".md") {
			// Remove .md extension
			promptName := strings.TrimSuffix(name, ".md")
			// Don't duplicate default
			if promptName != "default" {
				names = append(names, promptName)
			}
		}
	}

	return names, nil
}

// Reload clears the cache and forces reload on next access
func (l *PromptLoader) Reload() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.cache = make(map[string]*Prompt)
}
