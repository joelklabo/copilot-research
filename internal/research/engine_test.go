package research

import (
	"context"
	"testing"
	"time"

	"github.com/joelklabo/copilot-research/internal/db"
	"github.com/joelklabo/copilot-research/internal/prompts"
	"github.com/joelklabo/copilot-research/internal/provider"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// MockProvider for testing
type MockProvider struct {
	name           string
	authenticated  bool
	queryResponse  *provider.Response
	queryError     error
	queryCalled    bool
}

func (m *MockProvider) Name() string {
	return m.name
}

func (m *MockProvider) Query(ctx context.Context, prompt string, opts provider.QueryOptions) (*provider.Response, error) {
	m.queryCalled = true
	if m.queryError != nil {
		return nil, m.queryError
	}
	return m.queryResponse, nil
}

func (m *MockProvider) IsAuthenticated() bool {
	return m.authenticated
}

func (m *MockProvider) RequiresAuth() provider.AuthInfo {
	return provider.AuthInfo{IsConfigured: m.authenticated}
}

func (m *MockProvider) Capabilities() provider.ProviderCapabilities {
	return provider.ProviderCapabilities{}
}

func TestNewEngine(t *testing.T) {
	// Create temp database
	database, err := db.NewSQLiteDB(":memory:")
	require.NoError(t, err)
	defer database.Close()

	// Create prompt loader
	loader := prompts.NewPromptLoader("../../prompts")

	// Create provider manager
	factory := provider.NewProviderFactory()
	mockProvider := &MockProvider{
		name:          "test",
		authenticated: true,
		queryResponse: &provider.Response{
			Content:  "Test response",
			Provider: "test",
			Model:    "test-model",
			Duration: 100 * time.Millisecond,
		},
	}
	err = factory.Register("test", mockProvider) // Added error check
	require.NoError(t, err)
	providerMgr := provider.NewProviderManager(factory, "test", "", false, false) // Updated

	// Create engine
	engine := NewEngine(database, loader, providerMgr)
	assert.NotNil(t, engine)
}

func TestEngine_Research_FullFlow(t *testing.T) {
	// Create temp database
	database, err := db.NewSQLiteDB(":memory:")
	require.NoError(t, err)
	defer database.Close()

	// Create prompt loader
	loader := prompts.NewPromptLoader("../../prompts")

	// Create provider manager with mock
	factory := provider.NewProviderFactory()
	mockProvider := &MockProvider{
		name:          "test",
		authenticated: true,
		queryResponse: &provider.Response{
			Content:  "Test response about Swift actors",
			Provider: "test",
			Model:    "test-model",
			Duration: 100 * time.Millisecond,
		},
	}
	err = factory.Register("test", mockProvider) // Added error check
	require.NoError(t, err)
	providerMgr := provider.NewProviderManager(factory, "test", "", false, false) // Updated

	// Create engine
	engine := NewEngine(database, loader, providerMgr)

	// Research options
	opts := ResearchOptions{
		Query:      "How do Swift actors work?",
		Mode:       "quick",
		PromptName: "default",
		NoStore:    false,
	}

	// Progress channel
	progress := make(chan string, 10)
	go func() {
		// Drain progress channel
		for range progress {
		}
	}()

	// Execute research
	ctx := context.Background()
	result, err := engine.Research(ctx, opts, progress)

	// Verify
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, opts.Query, result.Query)
	assert.Equal(t, opts.Mode, result.Mode)
	assert.Equal(t, "Test response about Swift actors", result.Content)
	assert.GreaterOrEqual(t, result.Duration.Nanoseconds(), int64(0))
	assert.Greater(t, result.SessionID, int64(0))

	// Verify provider was called
	assert.True(t, mockProvider.queryCalled)

	// Verify session was stored in database
	session, err := database.GetSession(result.SessionID)
	require.NoError(t, err)
	assert.Equal(t, opts.Query, session.Query)
	assert.Equal(t, opts.Mode, session.Mode)
	assert.Equal(t, "Test response about Swift actors", session.Result)

	close(progress)
}

func TestEngine_Research_NoStore(t *testing.T) {
	// Create temp database
	database, err := db.NewSQLiteDB(":memory:")
	require.NoError(t, err)
	defer database.Close()

	// Create prompt loader
	loader := prompts.NewPromptLoader("../../prompts")

	// Create provider manager with mock
	factory := provider.NewProviderFactory()
	mockProvider := &MockProvider{
		name:          "test",
		authenticated: true,
		queryResponse: &provider.Response{
			Content:  "Test response",
			Provider: "test",
			Model:    "test-model",
			Duration: 100 * time.Millisecond,
		},
	}
	err = factory.Register("test", mockProvider) // Added error check
	require.NoError(t, err)
	providerMgr := provider.NewProviderManager(factory, "test", "", false, false) // Updated

	// Create engine
	engine := NewEngine(database, loader, providerMgr)

	// Research options with NoStore
	opts := ResearchOptions{
		Query:      "Test query",
		Mode:       "quick",
		PromptName: "default",
		NoStore:    true,
	}

	// Progress channel
	progress := make(chan string, 10)
	go func() {
		for range progress {
		}
	}()

	// Execute research
	ctx := context.Background()
	result, err := engine.Research(ctx, opts, progress)

	// Verify
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, int64(0), result.SessionID) // No session ID when NoStore is true

	// Verify no session was stored
	sessions, err := database.ListSessions(10, 0)
	require.NoError(t, err)
	assert.Equal(t, 0, len(sessions))

	close(progress)
}

func TestEngine_Research_ProgressEvents(t *testing.T) {
	// Create temp database
	database, err := db.NewSQLiteDB(":memory:")
	require.NoError(t, err)
	defer database.Close()

	// Create prompt loader
	loader := prompts.NewPromptLoader("../../prompts")

	// Create provider manager with mock
	factory := provider.NewProviderFactory()
	mockProvider := &MockProvider{
		name:          "test",
		authenticated: true,
		queryResponse: &provider.Response{
			Content:  "Test response",
			Provider: "test",
			Model:    "test-model",
			Duration: 100 * time.Millisecond,
		},
	}
	err = factory.Register("test", mockProvider) // Added error check
	require.NoError(t, err)
	providerMgr := provider.NewProviderManager(factory, "test", "", false, false) // Updated

	// Create engine
	engine := NewEngine(database, loader, providerMgr)

	// Research options
	opts := ResearchOptions{
		Query:      "Test query",
		Mode:       "quick",
		PromptName: "default",
		NoStore:    false,
	}

	// Collect progress events
	progress := make(chan string, 10)
	var events []string
	done := make(chan struct{})
	go func() {
		for msg := range progress {
			events = append(events, msg)
		}
		close(done)
	}()

	// Execute research
	ctx := context.Background()
	_, err = engine.Research(ctx, opts, progress)
	require.NoError(t, err)

	close(progress)
	<-done

	// Verify progress events were sent
	assert.Greater(t, len(events), 0)
	// Should contain expected progress messages
	hasLoadingPrompt := false
	hasQuerying := false
	for _, event := range events {
		if event == "Loading prompt..." {
			hasLoadingPrompt = true
		}
		if event == "Querying AI provider..." {
			hasQuerying = true
		}
	}
	assert.True(t, hasLoadingPrompt, "Expected 'Loading prompt...' event")
	assert.True(t, hasQuerying, "Expected 'Querying AI provider...' event")
}

func TestEngine_Research_ContextCancellation(t *testing.T) {
	// Create temp database
	database, err := db.NewSQLiteDB(":memory:")
	require.NoError(t, err)
	defer database.Close()

	// Create prompt loader
	loader := prompts.NewPromptLoader("../../prompts")

	// Create provider manager with mock that simulates slow query
	factory := provider.NewProviderFactory()
	mockProvider := &MockProvider{
		name:          "test",
		authenticated: true,
		queryError:    context.Canceled,
	}
	err = factory.Register("test", mockProvider) // Added error check
	require.NoError(t, err)
	providerMgr := provider.NewProviderManager(factory, "test", "", false, false) // Updated

	// Create engine
	engine := NewEngine(database, loader, providerMgr)

	// Research options
	opts := ResearchOptions{
		Query:      "Test query",
		Mode:       "quick",
		PromptName: "default",
		NoStore:    false,
	}

	// Progress channel
	progress := make(chan string, 10)
	go func() {
		for range progress {
		}
	}()

	// Create cancelable context
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	// Execute research
	result, err := engine.Research(ctx, opts, progress)

	// Verify error occurred
	assert.Error(t, err)
	assert.Nil(t, result)

	close(progress)
}

func TestEngine_Research_ProviderError(t *testing.T) {
	// Create temp database
	database, err := db.NewSQLiteDB(":memory:")
	require.NoError(t, err)
	defer database.Close()

	// Create prompt loader
	loader := prompts.NewPromptLoader("../../prompts")

	// Create provider manager with mock that returns error
	factory := provider.NewProviderFactory()
	mockProvider := &MockProvider{
		name:          "test",
		authenticated: true,
		queryError:    assert.AnError,
	}
	err = factory.Register("test", mockProvider) // Added error check
	require.NoError(t, err)
	providerMgr := provider.NewProviderManager(factory, "test", "", false, false) // Updated

	// Create engine
	engine := NewEngine(database, loader, providerMgr)

	// Research options
	opts := ResearchOptions{
		Query:      "Test query",
		Mode:       "quick",
		PromptName: "default",
		NoStore:    false,
	}

	// Progress channel
	progress := make(chan string, 10)
	go func() {
		for range progress {
		}
	}()

	// Execute research
	ctx := context.Background()
	result, err := engine.Research(ctx, opts, progress)

	// Verify error occurred
	assert.Error(t, err)
	assert.Nil(t, result)

	close(progress)
}
