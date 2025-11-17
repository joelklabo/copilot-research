package db

import "time"

// ResearchSession represents a single research query and its result
type ResearchSession struct {
	ID           int64     `json:"id"`
	Query        string    `json:"query"`
	Mode         string    `json:"mode"`
	PromptUsed   string    `json:"prompt_used"`
	Result       string    `json:"result"`
	QualityScore *int      `json:"quality_score,omitempty"` // Optional user rating
	CreatedAt    time.Time `json:"created_at"`
}

// LearnedPattern tracks successful research patterns and strategies
type LearnedPattern struct {
	ID           int64     `json:"id"`
	PatternName  string    `json:"pattern_name"`
	Description  string    `json:"description"`
	SuccessCount int       `json:"success_count"`
	LastUsed     time.Time `json:"last_used"`
	CreatedAt    time.Time `json:"created_at"`
}

// SearchHistory maintains a log of all search queries
type SearchHistory struct {
	ID        int64     `json:"id"`
	SessionID int64     `json:"session_id"`
	Query     string    `json:"query"`
	CreatedAt time.Time `json:"created_at"`
}
