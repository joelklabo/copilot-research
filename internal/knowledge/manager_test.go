package knowledge

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewKnowledgeManager(t *testing.T) {
	tmpDir := t.TempDir()
	
	km, err := NewKnowledgeManager(tmpDir)
	require.NoError(t, err)
	require.NotNil(t, km)
	
	// Should have initialized git repo
	gitDir := filepath.Join(tmpDir, ".git")
	_, err = os.Stat(gitDir)
	assert.NoError(t, err, "Git repo should be initialized")
}

func TestKnowledgeManager_Add(t *testing.T) {
	tmpDir := t.TempDir()
	km, err := NewKnowledgeManager(tmpDir)
	require.NoError(t, err)
	
	knowledge := &Knowledge{
		Topic:      "swift-concurrency",
		Content:    "Swift 6 introduces strict concurrency checking",
		Source:     "test",
		Confidence: 0.9,
		Tags:       []string{"swift", "concurrency"},
	}
	
	err = km.Add(knowledge)
	require.NoError(t, err)
	
	// Should have created the file
	filePath := filepath.Join(tmpDir, "swift-concurrency.md")
	_, err = os.Stat(filePath)
	assert.NoError(t, err, "Knowledge file should exist")
	
	// Should have committed to git
	commits, err := km.History("swift-concurrency")
	require.NoError(t, err)
	assert.Greater(t, len(commits), 0, "Should have at least one commit")
}

func TestKnowledgeManager_Get(t *testing.T) {
	tmpDir := t.TempDir()
	km, err := NewKnowledgeManager(tmpDir)
	require.NoError(t, err)
	
	original := &Knowledge{
		Topic:      "testing",
		Content:    "Use Swift Testing framework",
		Source:     "test",
		Confidence: 0.95,
		Tags:       []string{"testing"},
	}
	
	err = km.Add(original)
	require.NoError(t, err)
	
	retrieved, err := km.Get("testing")
	require.NoError(t, err)
	assert.Equal(t, original.Topic, retrieved.Topic)
	assert.Equal(t, original.Content, retrieved.Content)
	assert.Equal(t, original.Confidence, retrieved.Confidence)
}

func TestKnowledgeManager_Update(t *testing.T) {
	tmpDir := t.TempDir()
	km, err := NewKnowledgeManager(tmpDir)
	require.NoError(t, err)
	
	original := &Knowledge{
		Topic:      "mvvm",
		Content:    "MVVM is a pattern",
		Source:     "test",
		Confidence: 0.8,
	}
	
	err = km.Add(original)
	require.NoError(t, err)
	
	updated := &Knowledge{
		Topic:      "mvvm",
		Content:    "MVVM is outdated, use MV instead",
		Source:     "test-update",
		Confidence: 0.9,
	}
	
	err = km.Update("mvvm", updated)
	require.NoError(t, err)
	
	retrieved, err := km.Get("mvvm")
	require.NoError(t, err)
	assert.Equal(t, updated.Content, retrieved.Content)
	assert.Greater(t, retrieved.Version, original.Version)
}

func TestKnowledgeManager_Delete(t *testing.T) {
	tmpDir := t.TempDir()
	km, err := NewKnowledgeManager(tmpDir)
	require.NoError(t, err)
	
	knowledge := &Knowledge{
		Topic:   "temp",
		Content: "Temporary content",
		Source:  "test",
	}
	
	err = km.Add(knowledge)
	require.NoError(t, err)
	
	err = km.Delete("temp")
	require.NoError(t, err)
	
	_, err = km.Get("temp")
	assert.Error(t, err, "Should not find deleted knowledge")
}

func TestKnowledgeManager_List(t *testing.T) {
	tmpDir := t.TempDir()
	km, err := NewKnowledgeManager(tmpDir)
	require.NoError(t, err)
	
	topics := []string{"topic1", "topic2", "topic3"}
	for _, topic := range topics {
		err = km.Add(&Knowledge{
			Topic:   topic,
			Content: "Content for " + topic,
			Source:  "test",
		})
		require.NoError(t, err)
	}
	
	list, err := km.List()
	require.NoError(t, err)
	assert.Len(t, list, 3)
}

func TestKnowledgeManager_Search(t *testing.T) {
	tmpDir := t.TempDir()
	km, err := NewKnowledgeManager(tmpDir)
	require.NoError(t, err)
	
	err = km.Add(&Knowledge{
		Topic:   "swift-async",
		Content: "async/await in Swift",
		Source:  "test",
		Tags:    []string{"swift", "concurrency"},
	})
	require.NoError(t, err)
	
	err = km.Add(&Knowledge{
		Topic:   "swiftui-views",
		Content: "SwiftUI view hierarchy",
		Source:  "test",
		Tags:    []string{"swiftui"},
	})
	require.NoError(t, err)
	
	results, err := km.Search("swift")
	require.NoError(t, err)
	assert.Greater(t, len(results), 0, "Should find swift-related knowledge")
}

func TestKnowledgeManager_Deduplicate(t *testing.T) {
	tmpDir := t.TempDir()
	km, err := NewKnowledgeManager(tmpDir)
	require.NoError(t, err)
	
	// Add very similar content that should be deduped
	err = km.Add(&Knowledge{
		Topic:      "testing-duplicate",
		Content:    "Swift Testing is the new testing framework for Swift developers and it provides modern testing capabilities",
		Source:     "source1",
		Confidence: 0.8,
	})
	require.NoError(t, err)
	
	err = km.Add(&Knowledge{
		Topic:      "testing-duplicate-2",
		Content:    "Swift Testing is the new testing framework for Swift developers and it provides modern testing capabilities",
		Source:     "source2",
		Confidence: 0.9,
	})
	require.NoError(t, err)
	
	initialCount, _ := km.List()
	
	err = km.Deduplicate("testing")
	require.NoError(t, err)
	
	afterCount, _ := km.List()
	assert.Less(t, len(afterCount), len(initialCount), "Should have fewer entries after dedup")
}

func TestKnowledgeManager_Consolidate(t *testing.T) {
	tmpDir := t.TempDir()
	km, err := NewKnowledgeManager(tmpDir)
	require.NoError(t, err)
	
	// Add multiple entries with different topics
	for i := 0; i < 5; i++ {
		err = km.Add(&Knowledge{
			Topic:   fmt.Sprintf("swift-feature%d", i),
			Content: "Some content about Swift",
			Source:  "test",
		})
		require.NoError(t, err)
	}
	
	err = km.Consolidate()
	require.NoError(t, err)
	
	// Should have committed consolidation if anything was consolidated
	// For now just verify no error
}

func TestKnowledgeManager_GetRelevantKnowledge(t *testing.T) {
	tmpDir := t.TempDir()
	km, err := NewKnowledgeManager(tmpDir)
	require.NoError(t, err)
	
	err = km.Add(&Knowledge{
		Topic:   "swiftui-state",
		Content: "@State is for local view state",
		Source:  "test",
		Tags:    []string{"swiftui", "state"},
	})
	require.NoError(t, err)
	
	// Search for "state" which should match the topic and tags
	relevant, err := km.GetRelevantKnowledge("state", 1000)
	require.NoError(t, err)
	assert.Contains(t, relevant, "@State")
}

func TestKnowledgeManager_ThreadSafety(t *testing.T) {
	tmpDir := t.TempDir()
	km, err := NewKnowledgeManager(tmpDir)
	require.NoError(t, err)
	
	// Run concurrent operations
	done := make(chan bool)
	
	for i := 0; i < 10; i++ {
		go func(id int) {
			k := &Knowledge{
				Topic:   filepath.Join("concurrent", fmt.Sprintf("%d", id)), // Use fmt.Sprintf to convert int to string
				Content: "Concurrent content",
				Source:  "test",
			}
			err := km.Add(k) // Added error check
			require.NoError(t, err)
			done <- true
		}(i)
	}
	
	for i := 0; i < 10; i++ {
		<-done
	}
	
	list, err := km.List()
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(list), 10)
}