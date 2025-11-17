package cmd

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/joelklabo/copilot-research/internal/knowledge"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestKnowledgeListCommand(t *testing.T) {
	// Setup temp directory
	tempDir := t.TempDir()
	km, err := knowledge.NewKnowledgeManager(tempDir)
	require.NoError(t, err)

	// Add some test knowledge
	k1 := &knowledge.Knowledge{
		Topic:      "swift-concurrency",
		Content:    "Swift concurrency content",
		Source:     "manual",
		Confidence: 0.9,
		Tags:       []string{"swift", "concurrency"},
	}
	err = km.Add(k1)
	require.NoError(t, err)

	k2 := &knowledge.Knowledge{
		Topic:      "swiftui-patterns",
		Content:    "SwiftUI patterns content",
		Source:     "manual",
		Confidence: 0.8,
		Tags:       []string{"swiftui"},
	}
	err = km.Add(k2)
	require.NoError(t, err)

	// Test list functionality
	list, err := km.List()
	require.NoError(t, err)
	assert.Len(t, list, 2)
}

func TestKnowledgeShowCommand(t *testing.T) {
	// Setup temp directory
	tempDir := t.TempDir()
	km, err := knowledge.NewKnowledgeManager(tempDir)
	require.NoError(t, err)

	// Add test knowledge
	k := &knowledge.Knowledge{
		Topic:      "test-topic",
		Content:    "Test content here",
		Source:     "manual",
		Confidence: 0.95,
		Tags:       []string{"test"},
	}
	err = km.Add(k)
	require.NoError(t, err)

	// Test show functionality
	retrieved, err := km.Get("test-topic")
	require.NoError(t, err)
	assert.Equal(t, "test-topic", retrieved.Topic)
	assert.Equal(t, "Test content here", retrieved.Content)
	assert.Equal(t, 0.95, retrieved.Confidence)
}

func TestKnowledgeSearchCommand(t *testing.T) {
	// Setup temp directory
	tempDir := t.TempDir()
	km, err := knowledge.NewKnowledgeManager(tempDir)
	require.NoError(t, err)

	// Add test knowledge
	k1 := &knowledge.Knowledge{
		Topic:      "swift-actors",
		Content:    "Actor isolation in Swift",
		Source:     "manual",
		Confidence: 0.9,
		Tags:       []string{"swift", "actors"},
	}
	err = km.Add(k1)
	require.NoError(t, err)

	k2 := &knowledge.Knowledge{
		Topic:      "kotlin-coroutines",
		Content:    "Kotlin coroutines content",
		Source:     "manual",
		Confidence: 0.8,
		Tags:       []string{"kotlin"},
	}
	err = km.Add(k2)
	require.NoError(t, err)

	// Test search by topic
	results, err := km.Search("swift")
	require.NoError(t, err)
	assert.Len(t, results, 1)
	assert.Equal(t, "swift-actors", results[0].Topic)

	// Test search by content
	results, err = km.Search("isolation")
	require.NoError(t, err)
	assert.Len(t, results, 1)
	assert.Equal(t, "swift-actors", results[0].Topic)
}

func TestKnowledgeHistoryCommand(t *testing.T) {
	// Setup temp directory
	tempDir := t.TempDir()
	km, err := knowledge.NewKnowledgeManager(tempDir)
	require.NoError(t, err)

	// Add and update knowledge
	k := &knowledge.Knowledge{
		Topic:      "test-topic",
		Content:    "Initial content",
		Source:     "manual",
		Confidence: 0.9,
		Tags:       []string{"test"},
	}
	err = km.Add(k)
	require.NoError(t, err)

	// Wait a bit to ensure different timestamps
	time.Sleep(100 * time.Millisecond)

	// Update it
	k.Content = "Updated content"
	err = km.Update("test-topic", k)
	require.NoError(t, err)

	// Test history
	history, err := km.History("test-topic")
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(history), 2, "Should have at least 2 commits")
}

func TestKnowledgeConsolidateCommand(t *testing.T) {
	// Setup temp directory
	tempDir := t.TempDir()
	km, err := knowledge.NewKnowledgeManager(tempDir)
	require.NoError(t, err)

	// Add some knowledge
	k1 := &knowledge.Knowledge{
		Topic:      "swift/actors",
		Content:    "Actors content",
		Source:     "manual",
		Confidence: 0.9,
		Tags:       []string{"swift"},
	}
	err = km.Add(k1)
	require.NoError(t, err)

	// Test consolidate runs without error
	err = km.Consolidate()
	assert.NoError(t, err)
}

func TestKnowledgeAddCommand(t *testing.T) {
	// Setup temp directory
	tempDir := t.TempDir()
	km, err := knowledge.NewKnowledgeManager(tempDir)
	require.NoError(t, err)

	// Create a temporary file with content
	content := `---
topic: new-topic
version: 1
confidence: 0.8
tags: [test, new]
source: manual
created: 2025-11-17T12:00:00Z
updated: 2025-11-17T12:00:00Z
---

# New Topic

This is new content.
`
	tempFile := filepath.Join(t.TempDir(), "new-topic.md")
	err = os.WriteFile(tempFile, []byte(content), 0644)
	require.NoError(t, err)

	// Parse and add
	k, err := knowledge.ParseKnowledge(tempFile)
	require.NoError(t, err)

	err = km.Add(k)
	require.NoError(t, err)

	// Verify it was added
	retrieved, err := km.Get("new-topic")
	require.NoError(t, err)
	assert.Equal(t, "new-topic", retrieved.Topic)
}

func TestKnowledgeRulesCommand(t *testing.T) {
	// Setup temp directory
	tempDir := t.TempDir()
	km, err := knowledge.NewKnowledgeManager(tempDir)
	require.NoError(t, err)

	// Create rule engine
	re, err := knowledge.NewRuleEngine(km)
	require.NoError(t, err)

	// Test adding rule
	rule := knowledge.Rule{
		Type:    "exclude",
		Pattern: "MVC",
		Reason:  "Using MV architecture",
	}
	err = re.AddRule(rule)
	require.NoError(t, err)

	// Test listing rules
	rules := re.ListRules()
	assert.Len(t, rules, 1)
	assert.Equal(t, "exclude", rules[0].Type)

	// Test removing rule
	err = re.RemoveRule(rules[0].ID)
	require.NoError(t, err)

	// Verify removal
	rules = re.ListRules()
	assert.Len(t, rules, 0)
}
