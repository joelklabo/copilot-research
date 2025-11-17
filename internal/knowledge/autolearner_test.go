package knowledge

import (
	"testing"
	"time"

	"github.com/joelklabo/copilot-research/internal/research"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// MockKnowledgeManager for testing AutoLearner
type MockKnowledgeManager struct {
	addCalled bool
	addedKnowledge []*Knowledge
}

func (m *MockKnowledgeManager) Add(k *Knowledge) error {
	m.addCalled = true
	m.addedKnowledge = append(m.addedKnowledge, k)
	return nil
}

func (m *MockKnowledgeManager) Update(id string, k *Knowledge) error { return nil }
func (m *MockKnowledgeManager) Get(id string) (*Knowledge, error) { return nil, nil }
func (m *MockKnowledgeManager) Delete(id string) error { return nil }
func (m *MockKnowledgeManager) List() ([]*Knowledge, error) { return nil, nil }
func (m *MockKnowledgeManager) Search(query string) ([]*Knowledge, error) { return nil, nil }
func (m *MockKnowledgeManager) Deduplicate(topicPrefix string) error { return nil }
func (m *MockKnowledgeManager) Consolidate() error { return nil }
func (m *MockKnowledgeManager) GetRelevantKnowledge(query string, maxSize int) (string, error) { return "", nil }
func (m *MockKnowledgeManager) History(topic string) ([]GitCommit, error) { return nil, nil }
func (m *MockKnowledgeManager) Diff(from, to string) (string, error) { return "", nil }
func (m *MockKnowledgeManager) Commit(message string) error { return nil }


func TestNewAutoLearner(t *testing.T) {
	km := &MockKnowledgeManager{}
	al := NewAutoLearner(km)
	assert.NotNil(t, al)
	assert.Equal(t, km, al.km)
}

func TestAutoLearner_AnalyzeResult_Basic(t *testing.T) {
	km := &MockKnowledgeManager{}
	al := NewAutoLearner(km)

	testResult := &research.ResearchResult{
		Query:   "How to use Go modules",
		Mode:    "quick",
		Content: "Go modules are the dependency management system for Go.",
		Duration: 10 * time.Second,
		SessionID: 1,
	}

	knowledgeEntries, err := al.AnalyzeResult(testResult)
	require.NoError(t, err)
	assert.NotNil(t, knowledgeEntries)
	assert.Len(t, knowledgeEntries, 1)

	entry := knowledgeEntries[0]
	assert.Equal(t, testResult.Query, entry.Topic)
	assert.Equal(t, testResult.Content, entry.Content)
	assert.Equal(t, "auto-learned", entry.Source)
	assert.Equal(t, 0.7, entry.Confidence)
	assert.Contains(t, entry.Tags, "auto-learned")
	assert.Contains(t, entry.Tags, testResult.Mode)
}

func TestAutoLearner_AnalyzeResult_EmptyContent(t *testing.T) {
	km := &MockKnowledgeManager{}
	al := NewAutoLearner(km)

	testResult := &research.ResearchResult{
		Query:   "Empty query",
		Mode:    "quick",
		Content: "", // Empty content
	}

	knowledgeEntries, err := al.AnalyzeResult(testResult)
	assert.Error(t, err)
	assert.Nil(t, knowledgeEntries)
	assert.Contains(t, err.Error(), "research result is empty or nil")
}

func TestAutoLearner_AnalyzeResult_NilResult(t *testing.T) {
	km := &MockKnowledgeManager{}
	al := NewAutoLearner(km)

	knowledgeEntries, err := al.AnalyzeResult(nil)
	assert.Error(t, err)
	assert.Nil(t, knowledgeEntries)
	assert.Contains(t, err.Error(), "research result is empty or nil")
}
