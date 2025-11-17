-- Research Sessions Table
-- Stores all research queries and their results
CREATE TABLE IF NOT EXISTS research_sessions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    query TEXT NOT NULL,
    mode TEXT NOT NULL,
    prompt_used TEXT NOT NULL,
    result TEXT NOT NULL,
    quality_score INTEGER,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Index for fast lookups by creation date (most recent first)
CREATE INDEX IF NOT EXISTS idx_sessions_created ON research_sessions(created_at DESC);

-- Index for searching by query text
CREATE INDEX IF NOT EXISTS idx_sessions_query ON research_sessions(query);

-- Index for filtering by mode
CREATE INDEX IF NOT EXISTS idx_sessions_mode ON research_sessions(mode);

-- Learned Patterns Table
-- Tracks successful research patterns and strategies
CREATE TABLE IF NOT EXISTS learned_patterns (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    pattern_name TEXT UNIQUE NOT NULL,
    description TEXT,
    success_count INTEGER DEFAULT 0,
    last_used DATETIME,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Index for fast pattern lookups by name
CREATE INDEX IF NOT EXISTS idx_patterns_name ON learned_patterns(pattern_name);

-- Index for finding most successful patterns
CREATE INDEX IF NOT EXISTS idx_patterns_success ON learned_patterns(success_count DESC);

-- Search History Table
-- Maintains a log of all search queries
CREATE TABLE IF NOT EXISTS search_history (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    session_id INTEGER,
    query TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (session_id) REFERENCES research_sessions(id) ON DELETE CASCADE
);

-- Index for finding history by session
CREATE INDEX IF NOT EXISTS idx_history_session ON search_history(session_id);

-- Index for temporal queries
CREATE INDEX IF NOT EXISTS idx_history_created ON search_history(created_at DESC);
