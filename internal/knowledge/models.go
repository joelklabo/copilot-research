package knowledge

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

// Knowledge represents a piece of learned information
type Knowledge struct {
	ID         string    `json:"id" yaml:"id"`                   // SHA-256 hash of topic+content
	Topic      string    `json:"topic" yaml:"topic"`             // e.g., "swift-concurrency"
	Content    string    `json:"content" yaml:"content"`         // Markdown content
	Source     string    `json:"source" yaml:"source"`           // URL or "learned" or "manual"
	Confidence float64   `json:"confidence" yaml:"confidence"`   // 0.0 to 1.0
	Tags       []string  `json:"tags" yaml:"tags"`               // Topic tags
	CreatedAt  time.Time `json:"created_at" yaml:"created_at"`   // Created timestamp
	UpdatedAt  time.Time `json:"updated_at" yaml:"updated_at"`   // Last updated
	Version    int       `json:"version" yaml:"version"`         // Incremented on update
}

// GenerateID creates a unique ID from topic and content
func (k *Knowledge) GenerateID() string {
	data := k.Topic + k.Content
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}

// Rule represents a user-defined preference or exclusion
type Rule struct {
	ID          string    `json:"id" yaml:"id"`                                       // UUID
	Type        string    `json:"type" yaml:"type"`                                   // "exclude", "prefer", "always", "never"
	Pattern     string    `json:"pattern" yaml:"pattern"`                             // What to match (regex)
	Replacement string    `json:"replacement,omitempty" yaml:"replacement,omitempty"` // Optional replacement
	Reason      string    `json:"reason" yaml:"reason"`                               // Why this rule exists
	CreatedAt   time.Time `json:"created_at" yaml:"created_at"`                       // When created
}

// KnowledgeMetadata tracks overall knowledge base state
type KnowledgeMetadata struct {
	Version     string    `json:"version" yaml:"version"`           // Knowledge base version
	LastSync    time.Time `json:"last_sync" yaml:"last_sync"`       // Last sync time
	TotalTopics int       `json:"total_topics" yaml:"total_topics"` // Number of topics
	TotalRules  int       `json:"total_rules" yaml:"total_rules"`   // Number of rules
}

// Manifest represents the central registry of knowledge
type Manifest struct {
	Version  string            `yaml:"version"`
	Updated  time.Time         `yaml:"updated"`
	Topics   []ManifestTopic   `yaml:"topics"`
	Metadata KnowledgeMetadata `yaml:"metadata"`
}

// ManifestTopic represents a topic entry in the manifest
type ManifestTopic struct {
	Name       string    `yaml:"name"`
	File       string    `yaml:"file"`
	Version    int       `yaml:"version"`
	UpdatedAt  time.Time `yaml:"updated_at"`
	Confidence float64   `yaml:"confidence"`
	Tags       []string  `yaml:"tags"`
}

// Frontmatter represents the YAML frontmatter in knowledge files
type Frontmatter struct {
	Topic      string    `yaml:"topic"`
	Version    int       `yaml:"version"`
	Confidence float64   `yaml:"confidence"`
	Tags       []string  `yaml:"tags"`
	Source     string    `yaml:"source"`
	CreatedAt  time.Time `yaml:"created"`
	UpdatedAt  time.Time `yaml:"updated"`
}

// Save writes knowledge to a markdown file with YAML frontmatter
func (k *Knowledge) Save(filename string) error {
	fm := Frontmatter{
		Topic:      k.Topic,
		Version:    k.Version,
		Confidence: k.Confidence,
		Tags:       k.Tags,
		Source:     k.Source,
		CreatedAt:  k.CreatedAt,
		UpdatedAt:  k.UpdatedAt,
	}

	fmBytes, err := yaml.Marshal(fm)
	if err != nil {
		return fmt.Errorf("failed to marshal frontmatter: %w", err)
	}

	content := fmt.Sprintf("---\n%s---\n\n%s\n", string(fmBytes), k.Content)
	
	if err := os.WriteFile(filename, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

// ParseKnowledge reads and parses a knowledge markdown file
func ParseKnowledge(filename string) (*Knowledge, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	content := string(data)
	
	// Split frontmatter and content
	parts := splitFrontmatter(content)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid file format")
	}

	var fm Frontmatter
	if err := yaml.Unmarshal([]byte(parts[0]), &fm); err != nil {
		return nil, fmt.Errorf("failed to parse frontmatter: %w", err)
	}

	k := &Knowledge{
		Topic:      fm.Topic,
		Content:    parts[1],
		Source:     fm.Source,
		Confidence: fm.Confidence,
		Tags:       fm.Tags,
		CreatedAt:  fm.CreatedAt,
		UpdatedAt:  fm.UpdatedAt,
		Version:    fm.Version,
	}
	k.ID = k.GenerateID()

	return k, nil
}

func splitFrontmatter(content string) []string {
	// Find frontmatter delimiters
	lines := []string{}
	currentLine := ""
	for _, c := range content {
		if c == '\n' {
			lines = append(lines, currentLine)
			currentLine = ""
		} else {
			currentLine += string(c)
		}
	}
	if currentLine != "" {
		lines = append(lines, currentLine)
	}

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
