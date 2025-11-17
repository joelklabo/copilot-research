package knowledge

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

const (
	knowledgeDirName = ".copilot-research"
	knowledgeSubDir  = "knowledge"
)

// GetKnowledgeDir returns the path to the knowledge directory
func GetKnowledgeDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}
	return filepath.Join(home, knowledgeDirName, knowledgeSubDir), nil
}

// InitKnowledgeDir creates the knowledge directory structure and initializes Git
func InitKnowledgeDir() error {
	dir, err := GetKnowledgeDir()
	if err != nil {
		return err
	}

	// Create directory structure
	dirs := []string{
		dir,
		filepath.Join(dir, "topics"),
		filepath.Join(dir, "patterns"),
		filepath.Join(dir, "rules"),
	}

	for _, d := range dirs {
		if err := os.MkdirAll(d, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", d, err)
		}
	}

	// Initialize Git repository if not already initialized
	gitDir := filepath.Join(dir, ".git")
	if _, err := os.Stat(gitDir); os.IsNotExist(err) {
		cmd := exec.Command("git", "init")
		cmd.Dir = dir
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to initialize git repository: %w", err)
		}

		// Create .gitignore
		gitignore := filepath.Join(dir, ".gitignore")
		content := `# OS files
.DS_Store
Thumbs.db

# Temp files
*.tmp
*.swp
*~

# IDE
.vscode/
.idea/
`
		if err := os.WriteFile(gitignore, []byte(content), 0644); err != nil {
			return fmt.Errorf("failed to create .gitignore: %w", err)
		}

		// Initial commit
		cmd = exec.Command("git", "add", ".")
		cmd.Dir = dir
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to git add: %w", err)
		}

		cmd = exec.Command("git", "commit", "-m", "Initial commit: Initialize knowledge base")
		cmd.Dir = dir
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to initial commit: %w", err)
		}
	}

	// Create initial MANIFEST.yaml if it doesn't exist
	manifestPath := filepath.Join(dir, "MANIFEST.yaml")
	if _, err := os.Stat(manifestPath); os.IsNotExist(err) {
		manifest := &Manifest{
			Version: "1.0.0",
			Topics:  []ManifestTopic{},
			Metadata: KnowledgeMetadata{
				Version:     "1.0.0",
				TotalTopics: 0,
				TotalRules:  0,
			},
		}
		if err := SaveManifest(dir, manifest); err != nil {
			return fmt.Errorf("failed to create initial manifest: %w", err)
		}
	}

	// Create initial rules file if it doesn't exist
	rulesPath := filepath.Join(dir, "rules", "preferences.yaml")
	if _, err := os.Stat(rulesPath); os.IsNotExist(err) {
		rulesContent := `# User Preferences and Rules
# 
# Rule types:
#   - exclude: Remove matching content
#   - prefer: Replace with preferred alternative
#   - always: Always include when condition matches
#   - never: Never include matching content
#
rules: []
`
		if err := os.WriteFile(rulesPath, []byte(rulesContent), 0644); err != nil {
			return fmt.Errorf("failed to create rules file: %w", err)
		}
	}

	return nil
}

// EnsureKnowledgeDir ensures the knowledge directory exists and is initialized
func EnsureKnowledgeDir() (string, error) {
	dir, err := GetKnowledgeDir()
	if err != nil {
		return "", err
	}

	// Check if directory exists
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := InitKnowledgeDir(); err != nil {
			return "", err
		}
	}

	return dir, nil
}
