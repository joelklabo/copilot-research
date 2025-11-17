package knowledge

import (
	"testing"
	"time"
)

func TestKnowledgeGenerateID(t *testing.T) {
	k1 := &Knowledge{
		Topic:   "swift-concurrency",
		Content: "Test content",
	}

	k2 := &Knowledge{
		Topic:   "swift-concurrency",
		Content: "Test content",
	}

	k3 := &Knowledge{
		Topic:   "swift-concurrency",
		Content: "Different content",
	}

	id1 := k1.GenerateID()
	id2 := k2.GenerateID()
	id3 := k3.GenerateID()

	// Same topic and content should generate same ID
	if id1 != id2 {
		t.Errorf("Expected same IDs for identical knowledge, got %s and %s", id1, id2)
	}

	// Different content should generate different ID
	if id1 == id3 {
		t.Errorf("Expected different IDs for different content, got %s", id1)
	}

	// ID should be 64 characters (SHA-256 hex)
	if len(id1) != 64 {
		t.Errorf("Expected ID length 64, got %d", len(id1))
	}
}

func TestManifestTopicStruct(t *testing.T) {
	now := time.Now()
	mt := ManifestTopic{
		Name:       "test-topic",
		File:       "topics/test-topic.md",
		Version:    1,
		UpdatedAt:  now,
		Confidence: 0.95,
		Tags:       []string{"test", "go"},
	}

	if mt.Name != "test-topic" {
		t.Errorf("Expected name 'test-topic', got %s", mt.Name)
	}

	if mt.Confidence != 0.95 {
		t.Errorf("Expected confidence 0.95, got %f", mt.Confidence)
	}

	if len(mt.Tags) != 2 {
		t.Errorf("Expected 2 tags, got %d", len(mt.Tags))
	}
}

func TestRuleStruct(t *testing.T) {
	now := time.Now()
	rule := Rule{
		ID:          "test-rule-1",
		Type:        "exclude",
		Pattern:     "MVC|Model View Controller",
		Replacement: "",
		Reason:      "Using MV architecture",
		CreatedAt:   now,
	}

	if rule.Type != "exclude" {
		t.Errorf("Expected type 'exclude', got %s", rule.Type)
	}

	if rule.Pattern != "MVC|Model View Controller" {
		t.Errorf("Expected pattern 'MVC|Model View Controller', got %s", rule.Pattern)
	}
}
