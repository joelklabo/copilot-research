package knowledge

import (
	"bytes"
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

// ParseKnowledgeFile parses a markdown file with YAML frontmatter
func ParseKnowledgeFile(data []byte) (*Knowledge, error) {
	// Check for frontmatter delimiters
	if !bytes.HasPrefix(data, []byte("---\n")) {
		return nil, fmt.Errorf("no frontmatter found")
	}

	// Split into frontmatter and content
	parts := bytes.SplitN(data[4:], []byte("\n---\n"), 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid frontmatter format")
	}

	// Parse frontmatter
	var fm Frontmatter
	if err := yaml.Unmarshal(parts[0], &fm); err != nil {
		return nil, fmt.Errorf("failed to parse frontmatter: %w", err)
	}

	// Extract content
	content := strings.TrimSpace(string(parts[1]))

	// Create Knowledge struct
	k := &Knowledge{
		Topic:      fm.Topic,
		Content:    content,
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

// SerializeKnowledge converts a Knowledge struct to markdown with frontmatter
func SerializeKnowledge(k *Knowledge) ([]byte, error) {
	fm := Frontmatter{
		Topic:      k.Topic,
		Version:    k.Version,
		Confidence: k.Confidence,
		Tags:       k.Tags,
		Source:     k.Source,
		CreatedAt:  k.CreatedAt,
		UpdatedAt:  k.UpdatedAt,
	}

	fmData, err := yaml.Marshal(&fm)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal frontmatter: %w", err)
	}

	var buf bytes.Buffer
	buf.WriteString("---\n")
	buf.Write(fmData)
	buf.WriteString("---\n\n")
	buf.WriteString(k.Content)
	buf.WriteString("\n")

	return buf.Bytes(), nil
}
