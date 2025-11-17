package knowledge

import (
	"testing"
	"time"
)

func TestParseKnowledgeFile(t *testing.T) {
	markdown := `---
topic: swift-concurrency
version: 1
confidence: 0.9
tags: [swift, concurrency, actors]
source: https://docs.swift.org/
created: 2025-11-17T12:00:00Z
updated: 2025-11-17T14:00:00Z
---

# Swift Concurrency

Swift concurrency provides structured concurrency with async/await.

## Key Features

- Actors for safe mutable state
- Async/await syntax
- Task groups for parallel work
`

	k, err := ParseKnowledgeFile([]byte(markdown))
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	if k.Topic != "swift-concurrency" {
		t.Errorf("Expected topic 'swift-concurrency', got %s", k.Topic)
	}

	if k.Version != 1 {
		t.Errorf("Expected version 1, got %d", k.Version)
	}

	if k.Confidence != 0.9 {
		t.Errorf("Expected confidence 0.9, got %f", k.Confidence)
	}

	if len(k.Tags) != 3 {
		t.Errorf("Expected 3 tags, got %d", len(k.Tags))
	}

	if k.Source != "https://docs.swift.org/" {
		t.Errorf("Expected source 'https://docs.swift.org/', got %s", k.Source)
	}

	if !Contains(k.Content, "Swift Concurrency") {
		t.Errorf("Content missing expected text")
	}

	if k.ID == "" {
		t.Errorf("Expected ID to be generated")
	}
}

func TestSerializeKnowledge(t *testing.T) {
	now := time.Now()
	k := &Knowledge{
		Topic:      "test-topic",
		Content:    "# Test\n\nThis is test content.",
		Source:     "manual",
		Confidence: 0.85,
		Tags:       []string{"test", "example"},
		CreatedAt:  now,
		UpdatedAt:  now,
		Version:    1,
	}

	data, err := SerializeKnowledge(k)
	if err != nil {
		t.Fatalf("Failed to serialize: %v", err)
	}

	// Parse it back
	k2, err := ParseKnowledgeFile(data)
	if err != nil {
		t.Fatalf("Failed to parse serialized data: %v", err)
	}

	if k2.Topic != k.Topic {
		t.Errorf("Expected topic %s, got %s", k.Topic, k2.Topic)
	}

	if k2.Confidence != k.Confidence {
		t.Errorf("Expected confidence %f, got %f", k.Confidence, k2.Confidence)
	}

	if len(k2.Tags) != len(k.Tags) {
		t.Errorf("Expected %d tags, got %d", len(k.Tags), len(k2.Tags))
	}
}

func TestParseInvalidFrontmatter(t *testing.T) {
	tests := []struct {
		name    string
		content string
	}{
		{
			name:    "no frontmatter",
			content: "# Just content\n\nNo frontmatter here.",
		},
		{
			name: "incomplete frontmatter",
			content: `---
topic: test
---`,
		},
		{
			name: "invalid yaml",
			content: `---
topic: test
invalid yaml: [unclosed
---

Content`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ParseKnowledgeFile([]byte(tt.content))
			if err == nil {
				t.Errorf("Expected error for %s, got nil", tt.name)
			}
		})
	}
}

// Helper function
func Contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && ContainsHelper(s, substr))
}

func ContainsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
