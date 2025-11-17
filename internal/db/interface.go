package db

// DB defines the interface for database operations
type DB interface {
	// Sessions
	SaveSession(session *ResearchSession) error
	GetSession(id int64) (*ResearchSession, error)
	ListSessions(limit, offset int) ([]*ResearchSession, error)
	SearchSessions(query string) ([]*ResearchSession, error)

	// Patterns
	SavePattern(pattern *LearnedPattern) error
	GetPattern(name string) (*LearnedPattern, error)
	IncrementPattern(name string) error

	// Stats
	GetTotalSessions() (int, error)
	GetModeStats() (map[string]int, error)

	// Cleanup
	Close() error
}

// MockDB is a mock implementation of the DB interface for testing
type MockDB struct {
	SaveSessionFunc    func(session *ResearchSession) error
	GetSessionFunc     func(id int64) (*ResearchSession, error)
	ListSessionsFunc   func(limit, offset int) ([]*ResearchSession, error)
	SearchSessionsFunc func(query string) ([]*ResearchSession, error)
	SavePatternFunc    func(pattern *LearnedPattern) error
	GetPatternFunc     func(name string) (*LearnedPattern, error)
	IncrementPatternFunc func(name string) error
	GetTotalSessionsFunc func() (int, error)
	GetModeStatsFunc   func() (map[string]int, error)
	CloseFunc          func() error
}

// SaveSession calls SaveSessionFunc
func (m *MockDB) SaveSession(session *ResearchSession) error {
	if m.SaveSessionFunc != nil {
		return m.SaveSessionFunc(session)
	}
	return nil
}

// GetSession calls GetSessionFunc
func (m *MockDB) GetSession(id int64) (*ResearchSession, error) {
	if m.GetSessionFunc != nil {
		return m.GetSessionFunc(id)
	}
	return nil, nil
}

// ListSessions calls ListSessionsFunc
func (m *MockDB) ListSessions(limit, offset int) ([]*ResearchSession, error) {
	if m.ListSessionsFunc != nil {
		return m.ListSessionsFunc(limit, offset)
	}
	return nil, nil
}

// SearchSessions calls SearchSessionsFunc
func (m *MockDB) SearchSessions(query string) ([]*ResearchSession, error) {
	if m.SearchSessionsFunc != nil {
		return m.SearchSessionsFunc(query)
	}
	return nil, nil
}

// SavePattern calls SavePatternFunc
func (m *MockDB) SavePattern(pattern *LearnedPattern) error {
	if m.SavePatternFunc != nil {
		return m.SavePatternFunc(pattern)
	}
	return nil
}

// GetPattern calls GetPatternFunc
func (m *MockDB) GetPattern(name string) (*LearnedPattern, error) {
	if m.GetPatternFunc != nil {
		return m.GetPatternFunc(name)
	}
	return nil, nil
}

// IncrementPattern calls IncrementPatternFunc
func (m *MockDB) IncrementPattern(name string) error {
	if m.IncrementPatternFunc != nil {
		return m.IncrementPatternFunc(name)
	}
	return nil
}

// GetTotalSessions calls GetTotalSessionsFunc
func (m *MockDB) GetTotalSessions() (int, error) {
	if m.GetTotalSessionsFunc != nil {
		return m.GetTotalSessionsFunc()
	}
	return 0, nil
}

// GetModeStats calls GetModeStatsFunc
func (m *MockDB) GetModeStats() (map[string]int, error) {
	if m.GetModeStatsFunc != nil {
		return m.GetModeStatsFunc()
	}
	return nil, nil
}

// Close calls CloseFunc
func (m *MockDB) Close() error {
	if m.CloseFunc != nil {
		return m.CloseFunc()
	}
	return nil
}
