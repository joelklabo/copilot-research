package db

import (
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestDB(t *testing.T) (*SQLiteDB, string) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")
	
	db, err := NewSQLiteDB(dbPath)
	require.NoError(t, err, "should create database successfully")
	
	return db, dbPath
}

func TestNewSQLiteDB(t *testing.T) {
	// Test database creation
	db, dbPath := setupTestDB(t)
	defer db.Close()
	
	// Verify database file was created
	_, err := os.Stat(dbPath)
	assert.NoError(t, err, "database file should exist")
}

func TestSaveAndGetSession(t *testing.T) {
	db, _ := setupTestDB(t)
	defer db.Close()
	
	// Create a session
	session := &ResearchSession{
		Query:      "Test query",
		Mode:       "quick",
		PromptUsed: "default",
		Result:     "Test result",
		CreatedAt:  time.Now(),
	}
	
	// Save it
	err := db.SaveSession(session)
	require.NoError(t, err, "should save session successfully")
	assert.Greater(t, session.ID, int64(0), "should assign an ID")
	
	// Retrieve it
	retrieved, err := db.GetSession(session.ID)
	require.NoError(t, err, "should retrieve session successfully")
	assert.Equal(t, session.Query, retrieved.Query)
	assert.Equal(t, session.Mode, retrieved.Mode)
	assert.Equal(t, session.PromptUsed, retrieved.PromptUsed)
	assert.Equal(t, session.Result, retrieved.Result)
}

func TestSaveSessionWithQualityScore(t *testing.T) {
	db, _ := setupTestDB(t)
	defer db.Close()
	
	score := 5
	session := &ResearchSession{
		Query:        "Test with score",
		Mode:         "deep",
		PromptUsed:   "default",
		Result:       "Result",
		QualityScore: &score,
		CreatedAt:    time.Now(),
	}
	
	err := db.SaveSession(session)
	require.NoError(t, err)
	
	retrieved, err := db.GetSession(session.ID)
	require.NoError(t, err)
	require.NotNil(t, retrieved.QualityScore)
	assert.Equal(t, 5, *retrieved.QualityScore)
}

func TestListSessions(t *testing.T) {
	db, _ := setupTestDB(t)
	defer db.Close()
	
	// Create multiple sessions
	for i := 1; i <= 5; i++ {
		session := &ResearchSession{
			Query:      "Query " + string(rune('0'+i)),
			Mode:       "quick",
			PromptUsed: "default",
			Result:     "Result",
			CreatedAt:  time.Now(),
		}
		err := db.SaveSession(session)
		require.NoError(t, err)
	}
	
	// List first 3
	sessions, err := db.ListSessions(3, 0)
	require.NoError(t, err)
	assert.Len(t, sessions, 3, "should return 3 sessions")
	
	// List next 2
	sessions, err = db.ListSessions(2, 3)
	require.NoError(t, err)
	assert.Len(t, sessions, 2, "should return 2 sessions")
}

func TestSearchSessions(t *testing.T) {
	db, _ := setupTestDB(t)
	defer db.Close()
	
	// Create sessions with different queries
	sessions := []*ResearchSession{
		{Query: "Swift concurrency", Mode: "deep", PromptUsed: "default", Result: "R1", CreatedAt: time.Now()},
		{Query: "iOS 26 features", Mode: "quick", PromptUsed: "default", Result: "R2", CreatedAt: time.Now()},
		{Query: "Swift actor model", Mode: "deep", PromptUsed: "default", Result: "R3", CreatedAt: time.Now()},
	}
	
	for _, s := range sessions {
		err := db.SaveSession(s)
		require.NoError(t, err)
	}
	
	// Search for "Swift"
	results, err := db.SearchSessions("Swift")
	require.NoError(t, err)
	assert.Len(t, results, 2, "should find 2 sessions containing 'Swift'")
}

func TestSaveAndGetPattern(t *testing.T) {
	db, _ := setupTestDB(t)
	defer db.Close()
	
	pattern := &LearnedPattern{
		PatternName:  "test-pattern",
		Description:  "Test description",
		SuccessCount: 10,
		LastUsed:     time.Now(),
		CreatedAt:    time.Now(),
	}
	
	err := db.SavePattern(pattern)
	require.NoError(t, err)
	assert.Greater(t, pattern.ID, int64(0), "should assign an ID")
	
	retrieved, err := db.GetPattern("test-pattern")
	require.NoError(t, err)
	assert.Equal(t, pattern.PatternName, retrieved.PatternName)
	assert.Equal(t, pattern.Description, retrieved.Description)
	assert.Equal(t, pattern.SuccessCount, retrieved.SuccessCount)
}

func TestIncrementPattern(t *testing.T) {
	db, _ := setupTestDB(t)
	defer db.Close()
	
	pattern := &LearnedPattern{
		PatternName:  "increment-test",
		Description:  "Test",
		SuccessCount: 5,
		LastUsed:     time.Now(),
		CreatedAt:    time.Now(),
	}
	
	err := db.SavePattern(pattern)
	require.NoError(t, err)
	
	// Increment
	err = db.IncrementPattern("increment-test")
	require.NoError(t, err)
	
	// Verify
	retrieved, err := db.GetPattern("increment-test")
	require.NoError(t, err)
	assert.Equal(t, 6, retrieved.SuccessCount, "success count should be incremented")
}

func TestGetTotalSessions(t *testing.T) {
	db, _ := setupTestDB(t)
	defer db.Close()
	
	// Initially should be 0
	total, err := db.GetTotalSessions()
	require.NoError(t, err)
	assert.Equal(t, 0, total)
	
	// Add some sessions
	for i := 0; i < 3; i++ {
		session := &ResearchSession{
			Query:      "Query",
			Mode:       "quick",
			PromptUsed: "default",
			Result:     "Result",
			CreatedAt:  time.Now(),
		}
		err := db.SaveSession(session)
		require.NoError(t, err)
	}
	
	// Should now be 3
	total, err = db.GetTotalSessions()
	require.NoError(t, err)
	assert.Equal(t, 3, total)
}

func TestGetModeStats(t *testing.T) {
	db, _ := setupTestDB(t)
	defer db.Close()
	
	// Create sessions with different modes
	modes := []string{"quick", "quick", "deep", "quick", "deep", "compare"}
	for _, mode := range modes {
		session := &ResearchSession{
			Query:      "Query",
			Mode:       mode,
			PromptUsed: "default",
			Result:     "Result",
			CreatedAt:  time.Now(),
		}
		err := db.SaveSession(session)
		require.NoError(t, err)
	}
	
	// Get stats
	stats, err := db.GetModeStats()
	require.NoError(t, err)
	
	assert.Equal(t, 3, stats["quick"])
	assert.Equal(t, 2, stats["deep"])
	assert.Equal(t, 1, stats["compare"])
}

func TestConcurrentAccess(t *testing.T) {
	db, _ := setupTestDB(t)
	defer db.Close()
	
	var wg sync.WaitGroup
	numGoroutines := 10
	
	// Concurrent writes
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			session := &ResearchSession{
				Query:      "Concurrent query",
				Mode:       "quick",
				PromptUsed: "default",
				Result:     "Result",
				CreatedAt:  time.Now(),
			}
			err := db.SaveSession(session)
			assert.NoError(t, err)
		}(i)
	}
	
	wg.Wait()
	
	// Verify all sessions were saved
	total, err := db.GetTotalSessions()
	require.NoError(t, err)
	assert.Equal(t, numGoroutines, total, "all concurrent writes should succeed")
}

func TestClose(t *testing.T) {
	db, _ := setupTestDB(t)
	
	// Close should succeed
	err := db.Close()
	assert.NoError(t, err)
	
	// Operations after close should fail
	session := &ResearchSession{
		Query:      "Query",
		Mode:       "quick",
		PromptUsed: "default",
		Result:     "Result",
		CreatedAt:  time.Now(),
	}
	err = db.SaveSession(session)
	assert.Error(t, err, "operations should fail after close")
}

func TestGetSessionNotFound(t *testing.T) {
	db, _ := setupTestDB(t)
	defer db.Close()
	
	// Try to get non-existent session
	_, err := db.GetSession(9999)
	assert.Error(t, err, "should return error for non-existent session")
}

func TestGetPatternNotFound(t *testing.T) {
	db, _ := setupTestDB(t)
	defer db.Close()
	
	// Try to get non-existent pattern
	_, err := db.GetPattern("non-existent")
	assert.Error(t, err, "should return error for non-existent pattern")
}
