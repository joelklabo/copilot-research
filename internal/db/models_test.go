package db

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestResearchSessionStruct(t *testing.T) {
	// Test that ResearchSession struct exists and has correct fields
	now := time.Now()
	qualityScore := 5
	
	session := ResearchSession{
		ID:           1,
		Query:        "Test query",
		Mode:         "deep",
		PromptUsed:   "default",
		Result:       "Test result",
		QualityScore: &qualityScore,
		CreatedAt:    now,
	}
	
	assert.Equal(t, int64(1), session.ID)
	assert.Equal(t, "Test query", session.Query)
	assert.Equal(t, "deep", session.Mode)
	assert.Equal(t, "default", session.PromptUsed)
	assert.Equal(t, "Test result", session.Result)
	assert.Equal(t, 5, *session.QualityScore)
	assert.Equal(t, now, session.CreatedAt)
}

func TestResearchSessionJSON(t *testing.T) {
	// Test that ResearchSession can be marshaled/unmarshaled to JSON
	now := time.Now()
	qualityScore := 4
	
	original := ResearchSession{
		ID:           123,
		Query:        "Swift concurrency",
		Mode:         "quick",
		PromptUsed:   "default",
		Result:       "Result content",
		QualityScore: &qualityScore,
		CreatedAt:    now,
	}
	
	// Marshal
	jsonData, err := json.Marshal(original)
	require.NoError(t, err)
	
	// Unmarshal
	var decoded ResearchSession
	err = json.Unmarshal(jsonData, &decoded)
	require.NoError(t, err)
	
	// Verify
	assert.Equal(t, original.ID, decoded.ID)
	assert.Equal(t, original.Query, decoded.Query)
	assert.Equal(t, original.Mode, decoded.Mode)
	assert.Equal(t, original.PromptUsed, decoded.PromptUsed)
	assert.Equal(t, original.Result, decoded.Result)
	assert.Equal(t, *original.QualityScore, *decoded.QualityScore)
}

func TestLearnedPatternStruct(t *testing.T) {
	// Test that LearnedPattern struct exists and has correct fields
	now := time.Now()
	lastUsed := now.Add(-time.Hour)
	
	pattern := LearnedPattern{
		ID:           1,
		PatternName:  "test-pattern",
		Description:  "Test description",
		SuccessCount: 10,
		LastUsed:     lastUsed,
		CreatedAt:    now,
	}
	
	assert.Equal(t, int64(1), pattern.ID)
	assert.Equal(t, "test-pattern", pattern.PatternName)
	assert.Equal(t, "Test description", pattern.Description)
	assert.Equal(t, 10, pattern.SuccessCount)
	assert.Equal(t, lastUsed, pattern.LastUsed)
	assert.Equal(t, now, pattern.CreatedAt)
}

func TestLearnedPatternJSON(t *testing.T) {
	// Test that LearnedPattern can be marshaled/unmarshaled to JSON
	now := time.Now()
	
	original := LearnedPattern{
		ID:           456,
		PatternName:  "research-pattern",
		Description:  "A successful research pattern",
		SuccessCount: 25,
		LastUsed:     now,
		CreatedAt:    now.Add(-24 * time.Hour),
	}
	
	// Marshal
	jsonData, err := json.Marshal(original)
	require.NoError(t, err)
	
	// Unmarshal
	var decoded LearnedPattern
	err = json.Unmarshal(jsonData, &decoded)
	require.NoError(t, err)
	
	// Verify
	assert.Equal(t, original.ID, decoded.ID)
	assert.Equal(t, original.PatternName, decoded.PatternName)
	assert.Equal(t, original.Description, decoded.Description)
	assert.Equal(t, original.SuccessCount, decoded.SuccessCount)
}

func TestSearchHistoryStruct(t *testing.T) {
	// Test that SearchHistory struct exists and has correct fields
	now := time.Now()
	
	history := SearchHistory{
		ID:        1,
		SessionID: 100,
		Query:     "test search",
		CreatedAt: now,
	}
	
	assert.Equal(t, int64(1), history.ID)
	assert.Equal(t, int64(100), history.SessionID)
	assert.Equal(t, "test search", history.Query)
	assert.Equal(t, now, history.CreatedAt)
}

func TestSearchHistoryJSON(t *testing.T) {
	// Test that SearchHistory can be marshaled/unmarshaled to JSON
	now := time.Now()
	
	original := SearchHistory{
		ID:        789,
		SessionID: 123,
		Query:     "iOS development",
		CreatedAt: now,
	}
	
	// Marshal
	jsonData, err := json.Marshal(original)
	require.NoError(t, err)
	
	// Unmarshal
	var decoded SearchHistory
	err = json.Unmarshal(jsonData, &decoded)
	require.NoError(t, err)
	
	// Verify
	assert.Equal(t, original.ID, decoded.ID)
	assert.Equal(t, original.SessionID, decoded.SessionID)
	assert.Equal(t, original.Query, decoded.Query)
}

func TestQualityScoreOptional(t *testing.T) {
	// Test that QualityScore can be nil (optional field)
	session := ResearchSession{
		ID:           1,
		Query:        "Test",
		Mode:         "quick",
		PromptUsed:   "default",
		Result:       "Result",
		QualityScore: nil, // Optional field
		CreatedAt:    time.Now(),
	}
	
	assert.Nil(t, session.QualityScore)
	
	// Can also have a value
	score := 3
	session.QualityScore = &score
	assert.NotNil(t, session.QualityScore)
	assert.Equal(t, 3, *session.QualityScore)
}
