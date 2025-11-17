package db

import (
	"database/sql"
	_ "embed"
	"fmt"
	"sync"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// Compile-time check that SQLiteDB implements the DB interface
var _ DB = (*SQLiteDB)(nil)

//go:embed schema.sql
var schemaSQL string

// SQLiteDB implements database operations for SQLite
type SQLiteDB struct {
	db *sql.DB
	mu sync.RWMutex
}

// NewSQLiteDB creates a new SQLite database connection
func NewSQLiteDB(path string) (DB, error) {
	db, err := sql.Open("sqlite3", path+"?_journal_mode=WAL&_timeout=5000")
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Set connection pool settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Initialize schema
	if _, err := db.Exec(schemaSQL); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to initialize schema: %w", err)
	}

	return &SQLiteDB{db: db}, nil
}

// SaveSession saves a research session to the database
func (s *SQLiteDB) SaveSession(session *ResearchSession) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	query := `
		INSERT INTO research_sessions (query, mode, prompt_used, result, quality_score, created_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`

	result, err := s.db.Exec(
		query,
		session.Query,
		session.Mode,
		session.PromptUsed,
		session.Result,
		session.QualityScore,
		session.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to save session: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get insert ID: %w", err)
	}

	session.ID = id
	return nil
}

// GetSession retrieves a session by ID
func (s *SQLiteDB) GetSession(id int64) (*ResearchSession, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	query := `
		SELECT id, query, mode, prompt_used, result, quality_score, created_at
		FROM research_sessions
		WHERE id = ?
	`

	session := &ResearchSession{}
	err := s.db.QueryRow(query, id).Scan(
		&session.ID,
		&session.Query,
		&session.Mode,
		&session.PromptUsed,
		&session.Result,
		&session.QualityScore,
		&session.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("session not found: %d", id)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get session: %w", err)
	}

	return session, nil
}

// ListSessions retrieves sessions with pagination
func (s *SQLiteDB) ListSessions(limit, offset int) ([]*ResearchSession, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	query := `
		SELECT id, query, mode, prompt_used, result, quality_score, created_at
		FROM research_sessions
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`

	rows, err := s.db.Query(query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list sessions: %w", err)
	}
	defer rows.Close()

	var sessions []*ResearchSession
	for rows.Next() {
		session := &ResearchSession{}
		err := rows.Scan(
			&session.ID,
			&session.Query,
			&session.Mode,
			&session.PromptUsed,
			&session.Result,
			&session.QualityScore,
			&session.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan session: %w", err)
		}
		sessions = append(sessions, session)
	}

	return sessions, nil
}

// SearchSessions finds sessions matching a query string
func (s *SQLiteDB) SearchSessions(query string) ([]*ResearchSession, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	sql := `
		SELECT id, query, mode, prompt_used, result, quality_score, created_at
		FROM research_sessions
		WHERE query LIKE ?
		ORDER BY created_at DESC
	`

	rows, err := s.db.Query(sql, "%"+query+"%")
	if err != nil {
		return nil, fmt.Errorf("failed to search sessions: %w", err)
	}
	defer rows.Close()

	var sessions []*ResearchSession
	for rows.Next() {
		session := &ResearchSession{}
		err := rows.Scan(
			&session.ID,
			&session.Query,
			&session.Mode,
			&session.PromptUsed,
			&session.Result,
			&session.QualityScore,
			&session.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan session: %w", err)
		}
		sessions = append(sessions, session)
	}

	return sessions, nil
}

// SavePattern saves a learned pattern to the database
func (s *SQLiteDB) SavePattern(pattern *LearnedPattern) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	query := `
		INSERT INTO learned_patterns (pattern_name, description, success_count, last_used, created_at)
		VALUES (?, ?, ?, ?, ?)
		ON CONFLICT(pattern_name) DO UPDATE SET
			description = excluded.description,
			success_count = excluded.success_count,
			last_used = excluded.last_used
	`

	result, err := s.db.Exec(
		query,
		pattern.PatternName,
		pattern.Description,
		pattern.SuccessCount,
		pattern.LastUsed,
		pattern.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to save pattern: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get insert ID: %w", err)
	}

	pattern.ID = id
	return nil
}

// GetPattern retrieves a pattern by name
func (s *SQLiteDB) GetPattern(name string) (*LearnedPattern, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	query := `
		SELECT id, pattern_name, description, success_count, last_used, created_at
		FROM learned_patterns
		WHERE pattern_name = ?
	`

	pattern := &LearnedPattern{}
	err := s.db.QueryRow(query, name).Scan(
		&pattern.ID,
		&pattern.PatternName,
		&pattern.Description,
		&pattern.SuccessCount,
		&pattern.LastUsed,
		&pattern.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("pattern not found: %s", name)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get pattern: %w", err)
	}

	return pattern, nil
}

// IncrementPattern increments the success count for a pattern
func (s *SQLiteDB) IncrementPattern(name string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	query := `
		UPDATE learned_patterns
		SET success_count = success_count + 1,
		    last_used = ?
		WHERE pattern_name = ?
	`

	result, err := s.db.Exec(query, time.Now(), name)
	if err != nil {
		return fmt.Errorf("failed to increment pattern: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("pattern not found: %s", name)
	}

	return nil
}

// GetTotalSessions returns the total number of research sessions
func (s *SQLiteDB) GetTotalSessions() (int, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var count int
	err := s.db.QueryRow("SELECT COUNT(*) FROM research_sessions").Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to get total sessions: %w", err)
	}

	return count, nil
}

// GetModeStats returns statistics about mode usage
func (s *SQLiteDB) GetModeStats() (map[string]int, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	query := `
		SELECT mode, COUNT(*) as count
		FROM research_sessions
		GROUP BY mode
	`

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get mode stats: %w", err)
	}
	defer rows.Close()

	stats := make(map[string]int)
	for rows.Next() {
		var mode string
		var count int
		if err := rows.Scan(&mode, &count); err != nil {
			return nil, fmt.Errorf("failed to scan mode stats: %w", err)
		}
		stats[mode] = count
	}

	return stats, nil
}

// GetTopQueries returns the most common queries
func (s *SQLiteDB) GetTopQueries(limit int) ([]QueryCount, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	query := `
		SELECT query, COUNT(*) as count
		FROM research_sessions
		GROUP BY query
		ORDER BY count DESC
		LIMIT ?
	`

	rows, err := s.db.Query(query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get top queries: %w", err)
	}
	defer rows.Close()

	var topQueries []QueryCount
	for rows.Next() {
		var qc QueryCount
		if err := rows.Scan(&qc.Query, &qc.Count); err != nil {
			return nil, fmt.Errorf("failed to scan top query: %w", err)
		}
		topQueries = append(topQueries, qc)
	}

	return topQueries, nil
}

// Close closes the database connection
func (s *SQLiteDB) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.db != nil {
		return s.db.Close()
	}
	return nil
}
