# Implementation Plan - Copilot Research CLI

**Project**: Copilot Research - Beautiful CLI for AI-powered research  
**Language**: Go with Bubble Tea + Lipgloss  
**Status**: 🚧 In Progress  
**Started**: 2025-11-17

---

## Critical Issues

*None currently*

---

## Tasks

### Phase 1: Foundation & Setup

#### Task 1.1: Initialize Go Module and Dependencies
**Status**: ✅ Complete  
**Estimated**: 10 minutes  
**Priority**: P0 (Blocker)  
**Commit**: 46c982b

**Description**:
Initialize Go module and install core dependencies including Bubble Tea, Lipgloss, Cobra, and SQLite driver.

**Acceptance Criteria**:
- `go.mod` created with correct module path
- All dependencies listed and downloaded
- `go.sum` checksums verified
- Can build empty project

**Dependencies**: None

**Implementation Notes**:
```bash
go mod init github.com/joelklabo/copilot-research
go get github.com/charmbracelet/bubbletea@latest
go get github.com/charmbracelet/lipgloss@latest
go get github.com/charmbracelet/bubbles@latest
go get github.com/spf13/cobra@latest
go get github.com/mattn/go-sqlite3@latest
go get github.com/stretchr/testify@latest
```

**Tests**:
- Verify `go build` succeeds
- Verify `go mod tidy` has no changes

**Commit Template**:
```
[Setup] Initialize Go module and dependencies

Added all required dependencies for CLI tool:
- Bubble Tea (TUI framework)
- Lipgloss (styling)
- Bubbles (components)
- Cobra (CLI parsing)
- SQLite3 (database)
- Testify (testing)

Tests: Build verification
```

---

#### Task 1.2: Create Project Structure
**Status**: ✅ Complete  
**Estimated**: 15 minutes  
**Priority**: P0 (Blocker)  
**Commit**: 09b6619

**Description**:
Create all directories and placeholder files for the project structure. Set up `.gitignore` and basic configuration files.

**Acceptance Criteria**:
- All directories created (`cmd`, `internal/*`, `prompts`, etc.)
- `.gitignore` includes `tmp/`, binaries, OS files
- `Makefile` with basic targets (build, test, clean, install)
- `.editorconfig` for consistent formatting
- GitHub Actions workflow directory ready

**Dependencies**: Task 1.1

**Implementation Notes**:
```
internal/
├── research/
├── ui/
├── db/
├── prompts/
└── config/
```

**`.gitignore`**:
```
# Binaries
copilot-research
*.exe
*.dll
*.so
*.dylib

# Test files
*.test
*.out
coverage.txt

# Temp
tmp/
*.tmp

# Database
*.db
*.db-shm
*.db-wal

# OS
.DS_Store
Thumbs.db

# IDE
.idea/
.vscode/
*.swp
```

**Tests**:
- All directories exist
- Can create files in each directory
- `.gitignore` patterns work

**Commit Template**:
```
[Setup] Create project structure and configuration

Created directory layout:
- internal/ for private packages
- prompts/ for template files
- docs/ for documentation

Added .gitignore, Makefile, and .editorconfig.

Tests: Directory structure validation
```

---

#### Task 1.3: Default Prompt Template
**Status**: ✅ Complete  
**Estimated**: 20 minutes  
**Priority**: P0 (Blocker)  
**Commit**: 7f9c26b

**Description**:
Create the default prompt template that will be used to query GitHub Copilot. This should be optimized for research-style queries and produce high-quality markdown output.

**Acceptance Criteria**:
- `prompts/default.md` exists
- Prompt produces clean markdown output
- Includes instructions for citations and structure
- Tested with real `gh copilot` queries
- Has template variables: `{{query}}`, `{{mode}}`

**Dependencies**: Task 1.2

**Implementation Notes**:

The prompt should:
- Request structured markdown output
- Ask for clear sections (Overview, Details, Examples, References)
- Emphasize accuracy and citations
- Support different research modes
- Be concise but comprehensive

Template structure:
```markdown
---
name: default
description: Default research prompt for comprehensive queries
version: 1.0.0
---

You are an expert research assistant specializing in {{mode}} research.

[Detailed instructions...]

Research Query: {{query}}
```

**Tests**:
- Parse template successfully
- Replace variables correctly
- Test with `gh copilot suggest` manually

**Commit Template**:
```
[Prompts] Add default research prompt template

Created default.md with structured research instructions.
Optimized for markdown output with clear sections.

Supports template variables:
- {{query}} - User's research question
- {{mode}} - Research mode (quick, deep, compare, etc.)

Tests: Manual testing with gh copilot
```

---

### Phase 2: Core Database Layer

#### Task 2.1: Database Models and Schema
**Status**: ✅ Complete  
**Estimated**: 30 minutes  
**Priority**: P0  
**Commit**: 963cdd7

**Description**:
Define Go structs for data models and SQLite schema for storing research sessions, learned patterns, and statistics.

**Acceptance Criteria**:
- `internal/db/models.go` with all structs
- `internal/db/schema.sql` with CREATE statements
- Proper indexing for common queries
- Timestamps and metadata fields

**Dependencies**: Task 1.1

**Models Needed**:
```go
type ResearchSession struct {
    ID          int64
    Query       string
    Mode        string
    PromptUsed  string
    Result      string
    QualityScore *int  // Optional user rating
    CreatedAt   time.Time
}

type LearnedPattern struct {
    ID           int64
    PatternName  string
    Description  string
    SuccessCount int
    LastUsed     time.Time
    CreatedAt    time.Time
}

type SearchHistory struct {
    ID        int64
    SessionID int64
    Query     string
    CreatedAt time.Time
}
```

**Schema**:
```sql
CREATE TABLE IF NOT EXISTS research_sessions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    query TEXT NOT NULL,
    mode TEXT NOT NULL,
    prompt_used TEXT NOT NULL,
    result TEXT NOT NULL,
    quality_score INTEGER,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_sessions_created ON research_sessions(created_at DESC);
CREATE INDEX idx_sessions_query ON research_sessions(query);

CREATE TABLE IF NOT EXISTS learned_patterns (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    pattern_name TEXT UNIQUE NOT NULL,
    description TEXT,
    success_count INTEGER DEFAULT 0,
    last_used DATETIME,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS search_history (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    session_id INTEGER,
    query TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (session_id) REFERENCES research_sessions(id)
);
```

**Tests**:
- Struct field types match SQLite columns
- Can marshal/unmarshal JSON
- Timestamp handling works correctly

**Commit Template**:
```
[DB] Define data models and schema

Created structs for:
- ResearchSession (stores queries and results)
- LearnedPattern (tracks successful strategies)
- SearchHistory (query log)

Schema includes proper indexes for performance.

Tests: Model validation, JSON marshaling
```

---

#### Task 2.2: SQLite Database Implementation
**Status**: ✅ Complete  
**Estimated**: 45 minutes  
**Priority**: P0  
**Commit**: 051b589

**Description**:
Implement SQLite database layer with connection management, schema initialization, and CRUD operations for all models.

**Acceptance Criteria**:
- `internal/db/sqlite.go` with DB struct and methods
- Connection pooling and WAL mode enabled
- Auto-initialize schema on first run
- Proper error handling and logging
- Thread-safe operations

**Dependencies**: Task 2.1

**Interface**:
```go
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
```

**Implementation Notes**:
- Use `database/sql` with `mattn/go-sqlite3`
- Enable WAL mode: `PRAGMA journal_mode=WAL`
- Use prepared statements for performance
- Handle migration/schema updates gracefully

**Tests**:
- Create database successfully
- Save and retrieve sessions
- Search functionality works
- Concurrent access doesn't corrupt
- Close properly releases resources

**Commit Template**:
```
[DB] Implement SQLite database layer

Created DB implementation with:
- Connection management (WAL mode enabled)
- Auto schema initialization
- CRUD operations for all models
- Search and stats queries

Thread-safe with proper error handling.

Tests: CRUD operations, concurrent access, search
```

---

### Phase 3: Knowledge Management System

#### Task 3.1: Knowledge Base Structure and Models
**Status**: ✅ Complete  
**Estimated**: 30 minutes  
**Priority**: P0  
**Commit**: 119af21

**Description**:
Create the knowledge base directory structure and Go models for managing research knowledge with Git-based versioning.

**Acceptance Criteria**:
- Create `~/.copilot-research/knowledge/` directory structure
- Define Knowledge, Rule, and KnowledgeMetadata structs
- Create knowledge schema (markdown files with frontmatter)
- Initialize Git repository in knowledge directory
- Create `.gitignore` for knowledge repo

**Dependencies**: Task 2.2

**Directory Structure**:
```
~/.copilot-research/knowledge/
├── .git/                  # Git repo for versioning
├── topics/                # Topic-based knowledge
│   ├── swift-concurrency.md
│   ├── swiftui-patterns.md
│   └── ...
├── patterns/              # Learned patterns
│   ├── common-errors.md
│   └── best-practices.md
├── rules/                 # User preferences and rules
│   ├── preferences.yaml
│   └── exclusions.yaml
└── MANIFEST.yaml         # Central registry
```

**Models**:
```go
type Knowledge struct {
    ID          string    // SHA-256 of content
    Topic       string    // e.g., "swift-concurrency"
    Content     string    // Markdown content
    Source      string    // URL or "learned" or "manual"
    Confidence  float64   // 0.0 to 1.0
    Tags        []string
    CreatedAt   time.Time
    UpdatedAt   time.Time
    Version     int       // Incremented on update
}

type Rule struct {
    ID          string
    Type        string    // "exclude", "prefer", "always", "never"
    Pattern     string    // What to match
    Replacement string    // Optional
    Reason      string    // Why this rule exists
    CreatedAt   time.Time
}

type KnowledgeMetadata struct {
    Version     string
    LastSync    time.Time
    TotalTopics int
    TotalRules  int
}
```

**Markdown Frontmatter Format**:
```yaml
---
topic: swift-concurrency
version: 3
confidence: 0.95
tags: [swift, concurrency, actors]
source: https://example.com/swift-concurrency
created: 2025-11-17T12:00:00Z
updated: 2025-11-17T14:00:00Z
---

# Swift Concurrency

[Content...]
```

**Tests**:
- Directory creation works
- Git initialization succeeds
- Can parse frontmatter correctly
- Can serialize/deserialize models

**Commit Template**:
```
[Knowledge] Create knowledge base structure and models

Created knowledge management system foundation:
- Directory structure in ~/.copilot-research/knowledge/
- Git-based versioning
- Knowledge, Rule, and Metadata models
- Markdown frontmatter format
- MANIFEST.yaml schema

Tests: Directory creation, parsing, serialization
```

---

#### Task 3.2: Knowledge Manager Implementation
**Status**: ✅ Complete  
**Estimated**: 60 minutes  
**Priority**: P0  
**Commit**: 6165308

**Description**:
Implement KnowledgeManager that handles CRUD operations, Git operations, deduplication, and consolidation of knowledge files.

**Acceptance Criteria**:
- CRUD operations for knowledge entries
- Git commits with descriptive messages
- Automatic deduplication on write
- Consolidation pass after updates
- Thread-safe operations
- Load knowledge into prompt context

**Dependencies**: Task 3.1

**Implementation**:
```go
type KnowledgeManager struct {
    baseDir    string
    gitRepo    *git.Repository
    cache      map[string]*Knowledge
    mu         sync.RWMutex
}

func (km *KnowledgeManager) Add(k *Knowledge) error
func (km *KnowledgeManager) Update(id string, k *Knowledge) error
func (km *KnowledgeManager) Get(id string) (*Knowledge, error)
func (km *KnowledgeManager) Search(query string) ([]*Knowledge, error)
func (km *KnowledgeManager) List() ([]*Knowledge, error)
func (km *KnowledgeManager) Delete(id string) error

// Git operations
func (km *KnowledgeManager) Commit(message string) error
func (km *KnowledgeManager) History(topic string) ([]GitCommit, error)
func (km *KnowledgeManager) Diff(from, to string) (string, error)

// Consolidation
func (km *KnowledgeManager) Consolidate() error
func (km *KnowledgeManager) Deduplicate(topic string) error

// Loading for prompts
func (km *KnowledgeManager) GetRelevantKnowledge(query string, maxSize int) (string, error)
```

**Consolidation Strategy**:
1. Group by topic
2. Merge duplicate/similar content
3. Keep highest confidence version
4. Preserve unique information
5. Update version numbers
6. Git commit changes

**Deduplication Algorithm**:
- Calculate similarity score (cosine similarity of embeddings)
- If similarity > 0.9, merge entries
- Keep newer version if timestamps differ
- Combine tags and sources

**Git Commit Messages**:
- "Add: {topic} - {summary}"
- "Update: {topic} - {summary}"
- "Consolidate: Merged {n} entries in {topic}"
- "Remove: {topic} - {reason}"

**Tests**:
- CRUD operations work
- Git commits created
- Deduplication removes duplicates
- Consolidation reduces file size
- Thread-safe concurrent access
- Relevant knowledge retrieval

**Commit Template**:
```
[Knowledge] Implement knowledge manager with Git versioning

Created KnowledgeManager with:
- CRUD operations for knowledge entries
- Automatic Git commits with descriptive messages
- Deduplication on write
- Consolidation pass for cleanup
- Thread-safe operations
- Relevant knowledge retrieval for prompts

Tests: CRUD, Git ops, dedup, consolidation, concurrency
```

---

#### Task 3.3: Rule System Implementation
**Status**: ⬜ Not Started  
**Estimated**: 45 minutes  
**Priority**: P1

**Description**:
Implement rule system for user preferences, exclusions, and content filtering.

**Acceptance Criteria**:
- Add/remove/list rules
- Apply rules to knowledge content
- Persist rules in YAML
- Rule validation
- Pattern matching (regex support)

**Dependencies**: Task 3.2

**Rule Types**:
```yaml
rules:
  - type: exclude
    pattern: "Model View Controller|MVC"
    reason: "Using MV architecture instead"
    
  - type: prefer
    pattern: "Swift Testing"
    over: "XCTest"
    reason: "Modern testing framework"
    
  - type: always_mention
    pattern: "actor isolation"
    when: "swift.*concurrency"
    
  - type: never_mention
    pattern: "Objective-C"
    reason: "Swift-only codebase"
```

**Implementation**:
```go
type RuleEngine struct {
    rules []Rule
    km    *KnowledgeManager
}

func (re *RuleEngine) AddRule(rule Rule) error
func (re *RuleEngine) RemoveRule(id string) error
func (re *RuleEngine) ListRules() []Rule
func (re *RuleEngine) Apply(content string) (string, error)
func (re *RuleEngine) Validate(rule Rule) error
```

**Apply Algorithm**:
1. Load all rules
2. For each rule:
   - If type == "exclude", remove matching content
   - If type == "prefer", replace patterns
   - If type == "always_mention", ensure inclusion
   - If type == "never_mention", filter out
3. Return filtered content

**Tests**:
- Rule CRUD operations
- Pattern matching works
- Content filtering correct
- Rule validation catches errors
- Rules persist across restarts

**Commit Template**:
```
[Knowledge] Implement rule system for preferences

Created RuleEngine with:
- Add/remove/list rules
- Multiple rule types (exclude, prefer, always, never)
- Pattern matching with regex
- Content filtering
- YAML persistence

Tests: Rule operations, pattern matching, filtering
```

---

#### Task 3.4: Knowledge CLI Commands
**Status**: ⬜ Not Started  
**Estimated**: 45 minutes  
**Priority**: P1

**Description**:
Add CLI commands for managing knowledge base, viewing history, and editing rules.

**Acceptance Criteria**:
- `knowledge list` - Show all topics
- `knowledge show <topic>` - Display knowledge
- `knowledge add <topic>` - Add new knowledge
- `knowledge edit <topic>` - Edit in $EDITOR
- `knowledge search <query>` - Search content
- `knowledge history <topic>` - Show Git history
- `knowledge consolidate` - Run consolidation
- `knowledge rules` - Manage rules

**Dependencies**: Task 3.3

**Commands**:
```bash
copilot-research knowledge list
copilot-research knowledge show swift-concurrency
copilot-research knowledge add new-topic
copilot-research knowledge edit swift-concurrency
copilot-research knowledge search "actor isolation"
copilot-research knowledge history swift-concurrency
copilot-research knowledge consolidate
copilot-research knowledge stats

# Rule management
copilot-research knowledge rules list
copilot-research knowledge rules add --exclude "MVC"
copilot-research knowledge rules remove <id>
```

**Output Examples**:
```
$ copilot-research knowledge list

Knowledge Base (12 topics)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
Topic                      Version  Updated
swift-concurrency          3        2 hours ago
swiftui-patterns          5        1 day ago
networking-best-practices  2        3 days ago

$ copilot-research knowledge show swift-concurrency

Swift Concurrency (v3)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
Confidence: 95%
Tags: swift, concurrency, actors
Source: https://docs.swift.org/...
Updated: 2 hours ago

[Content displayed with syntax highlighting]
```

**Tests**:
- All commands execute correctly
- Output formatted nicely
- Editor launches for edit
- Search returns relevant results
- History shows commits

**Commit Template**:
```
[CLI] Add knowledge management commands

Implemented knowledge commands:
- list: Show all topics
- show: Display specific knowledge
- add/edit: Manage content
- search: Find relevant information
- history: View Git history
- consolidate: Run cleanup
- rules: Manage preferences

Beautiful table output with colors.

Tests: All commands, formatting, editor
```

---

#### Task 3.5: Auto-Learning from Research Results
**Status**: ⬜ Not Started  
**Estimated**: 45 minutes  
**Priority**: P2

**Description**:
Automatically extract and store valuable knowledge from successful research sessions.

**Acceptance Criteria**:
- Analyze research results for patterns
- Extract key information
- Store as knowledge entries
- Tag automatically
- Calculate confidence based on source
- Prompt user to review/approve

**Dependencies**: Task 3.2, Task 4.2

**Implementation**:
```go
type AutoLearner struct {
    km     *KnowledgeManager
    engine *ResearchEngine
}

func (al *AutoLearner) AnalyzeResult(result *ResearchResult) (*Knowledge, error)
func (al *AutoLearner) ExtractTopics(content string) []string
func (al *AutoLearner) CalculateConfidence(result *ResearchResult) float64
func (al *AutoLearner) ShouldStore(k *Knowledge) bool
```

**Analysis Strategy**:
1. Parse result markdown
2. Extract headers as topics
3. Identify code examples
4. Find URLs for sources
5. Calculate confidence:
   - Official docs: 0.9-1.0
   - GitHub repos: 0.7-0.9
   - Blog posts: 0.5-0.7
6. Prompt user if confidence < 0.8

**User Prompt**:
```
Found valuable information about "Swift Concurrency"

Would you like to save this to your knowledge base?
[Y]es / [N]o / [E]dit first / [A]lways for this topic
```

**Tests**:
- Topic extraction works
- Confidence calculation reasonable
- User prompts display correctly
- Knowledge stored properly
- Deduplication prevents duplicates

**Commit Template**:
```
[Knowledge] Add auto-learning from research results

Implemented auto-learning system:
- Analyzes research results for patterns
- Extracts topics and key information
- Calculates confidence scores
- Prompts user for approval
- Stores in knowledge base

Helps build knowledge over time.

Tests: Analysis, extraction, confidence, storage
```

---

### Phase 4: Prompt Management

#### Task 4.1: Prompt Loader
**Status**: ⬜ Not Started  
**Estimated**: 30 minutes  
**Priority**: P1

**Description**:
Implement prompt loading system that can read prompts from files, parse frontmatter, and perform template variable substitution.

**Acceptance Criteria**:
- Load prompts from `prompts/` directory
- Parse YAML frontmatter (name, description, version)
- Template variable substitution ({{query}}, {{mode}})
- Fall back to embedded default if file missing
- Cache loaded prompts in memory

**Dependencies**: Task 1.3

**Implementation**:
```go
type Prompt struct {
    Name        string
    Description string
    Version     string
    Template    string
}

type PromptLoader struct {
    promptsDir string
    cache      map[string]*Prompt
    mu         sync.RWMutex
}

func (l *PromptLoader) Load(name string) (*Prompt, error)
func (l *PromptLoader) Render(prompt *Prompt, vars map[string]string) string
func (l *PromptLoader) List() []string
```

**Tests**:
- Load valid prompt file
- Parse frontmatter correctly
- Variable substitution works
- Handles missing files gracefully
- Cache invalidation works

**Commit Template**:
```
[Prompts] Implement prompt loader and template system

Created PromptLoader with:
- File loading from prompts/ directory
- YAML frontmatter parsing
- Template variable substitution
- In-memory caching

Falls back to embedded default if custom not found.

Tests: Loading, parsing, substitution, caching
```

---

#### Task 4.2: Additional Prompt Templates
**Status**: ⬜ Not Started  
**Estimated**: 30 minutes  
**Priority**: P2

**Description**:
Create additional prompt templates for different research modes and use cases (deep-dive, compare, synthesis, etc.).

**Acceptance Criteria**:
- `prompts/deep-dive.md` - Comprehensive research
- `prompts/compare.md` - Compare multiple options
- `prompts/synthesis.md` - Synthesize from sources
- `prompts/quick.md` - Fast overview
- All have proper frontmatter and tested

**Dependencies**: Task 3.1

**Tests**:
- Each prompt produces appropriate output type
- Manual testing with gh copilot
- Different modes produce different results

**Commit Template**:
```
[Prompts] Add research mode prompt templates

Created specialized prompts:
- deep-dive.md - Comprehensive with examples
- compare.md - Side-by-side comparison
- synthesis.md - Multi-source synthesis
- quick.md - Fast overview

Tests: Manual validation with gh copilot
```

---

### Phase 5: Multi-Provider Authentication & Integration

#### Task 5.1: Provider Abstraction Layer
**Status**: ⬜ Not Started  
**Estimated**: 60 minutes  
**Priority**: P0

**Description**:
Design and implement provider abstraction layer supporting multiple AI backends (GitHub Copilot, OpenAI, Anthropic Claude) with unified interface.

**Acceptance Criteria**:
- Provider interface/trait with standard methods
- Factory pattern for provider instantiation  
- Adapter pattern for each provider SDK
- Configuration-driven provider selection
- Support for provider capabilities (streaming, function calling)
- Graceful provider fallback on error

**Dependencies**: Task 1.1, Task 3.1

**Research Context**:
Based on industry best practices:
- Use unified API gateway pattern (similar to Conduit, AWS Multi-Provider GenAI Gateway)
- Implement adapter pattern for each provider's SDK
- Support Model Context Protocol (MCP) for standardization
- Configuration via YAML/environment variables
- Allow runtime provider switching

**Provider Interface**:
```go
type AIProvider interface {
    Name() string
    Query(ctx context.Context, prompt string, opts QueryOptions) (*Response, error)
    IsAuthenticated() bool
    RequiresAuth() AuthInfo
    Capabilities() ProviderCapabilities
}

type ProviderCapabilities struct {
    Streaming      bool
    FunctionCall   bool
    MaxTokens      int
    SupportsImages bool
}

type AuthInfo struct {
    Type         string // "oauth", "apikey", "cli"
    IsConfigured bool
    HelpURL      string
    Instructions string
}

type ProviderFactory struct {
    providers map[string]AIProvider
}

func (f *ProviderFactory) Create(name string, config Config) (AIProvider, error)
func (f *ProviderFactory) Register(name string, provider AIProvider) error
func (f *ProviderFactory) List() []string
```

**Configuration**:
```yaml
# ~/.copilot-research/config.yaml
providers:
  primary: github-copilot
  fallback: openai
  
  github-copilot:
    enabled: true
    auth_type: cli  # Uses gh CLI
    timeout: 60s
    
  openai:
    enabled: true
    auth_type: apikey
    api_key_env: OPENAI_API_KEY
    model: gpt-4
    timeout: 30s
    
  anthropic:
    enabled: false
    auth_type: apikey
    api_key_env: ANTHROPIC_API_KEY
    model: claude-3-5-sonnet
    timeout: 30s
```

**Tests**:
- Factory creates providers correctly
- Interface methods work consistently
- Provider switching works
- Configuration parsing correct
- Fallback logic functions

**Commit Template**:
```
[Providers] Implement multi-provider abstraction layer

Created provider abstraction with:
- Unified AIProvider interface
- Factory pattern for instantiation
- Adapter pattern for each provider
- Configuration-driven selection
- Capability detection
- Fallback support

Supports GitHub Copilot, OpenAI, Anthropic.

Tests: Factory, interface, config, fallback
```

---

#### Task 5.2: GitHub Copilot Provider Implementation  
**Status**: ⬜ Not Started  
**Estimated**: 45 minutes  
**Priority**: P0

**Description**:
Implement GitHub Copilot provider that wraps `gh copilot` CLI, handles authentication via multiple methods, and provides excellent user onboarding.

**Acceptance Criteria**:
- Execute `gh copilot suggest` with prompt
- Support multiple authentication methods (priority order):
  1. `COPILOT_GITHUB_TOKEN` environment variable
  2. `GH_TOKEN` environment variable  
  3. Existing `gh` CLI authentication
  4. Interactive OAuth device flow
- Detect authentication status before queries
- Provide clear, actionable error messages
- Handle subscription/permission errors gracefully
- Beautiful onboarding flow for new users
- Timeout after reasonable duration
- Return clean markdown output

**Dependencies**: Task 5.1

**Authentication Research Context**:
Based on CLI best practices:
- Check credentials in priority order (COPILOT_GITHUB_TOKEN > GH_TOKEN > gh auth)
- Use OAuth device flow for interactive authentication
- Validate token has "Copilot Requests" permission
- Fine-grained PATs required (classic tokens rejected)
- Clear error messages with actionable solutions
- Help users understand subscription requirements

**Implementation**:
```go
type GitHubCopilotProvider struct {
    timeout    time.Duration
    authMethod string
    token      string
}

func (g *GitHubCopilotProvider) Name() string { return "github-copilot" }

func (g *GitHubCopilotProvider) Query(ctx context.Context, prompt string, opts QueryOptions) (*Response, error) {
    // Execute gh copilot suggest with prompt
    // Capture and parse output
    // Handle errors with helpful messages
}

func (g *GitHubCopilotProvider) IsAuthenticated() bool {
    // Check in priority order:
    // 1. COPILOT_GITHUB_TOKEN env var
    // 2. GH_TOKEN env var
    // 3. gh auth status
    // Return true if any valid
}

func (g *GitHubCopilotProvider) RequiresAuth() AuthInfo {
    if g.IsAuthenticated() {
        return AuthInfo{IsConfigured: true}
    }
    
    return AuthInfo{
        Type:         "oauth-device-flow",
        IsConfigured: false,
        HelpURL:      "https://github.com/features/copilot",
        Instructions: `GitHub Copilot authentication required.

Please authenticate using one of these methods:

1. GitHub CLI (recommended):
   gh auth login
   
2. Personal Access Token:
   export COPILOT_GITHUB_TOKEN=ghp_your_token_here
   
3. Interactive device flow:
   copilot-research auth login

Note: You need an active GitHub Copilot subscription.
Get one at https://github.com/features/copilot

Once authenticated, run your command again.`,
    }
}

func (g *GitHubCopilotProvider) Capabilities() ProviderCapabilities {
    return ProviderCapabilities{
        Streaming:      false,
        FunctionCall:   true,  // Via MCP
        MaxTokens:      8000,
        SupportsImages: false,
    }
}

// Authentication check with priority
func (g *GitHubCopilotProvider) detectAuth() (string, string) {
    // 1. Check COPILOT_GITHUB_TOKEN
    if token := os.Getenv("COPILOT_GITHUB_TOKEN"); token != "" {
        return "env:COPILOT_GITHUB_TOKEN", token
    }
    
    // 2. Check GH_TOKEN
    if token := os.Getenv("GH_TOKEN"); token != "" {
        return "env:GH_TOKEN", token
    }
    
    // 3. Check gh CLI authentication
    cmd := exec.Command("gh", "auth", "status")
    if cmd.Run() == nil {
        return "gh-cli", ""
    }
    
    return "none", ""
}
```

**Error Messages** (following UX best practices):
```
❌ Authentication Required

GitHub Copilot is not authenticated on this machine.

How to fix:
  1. Authenticate with GitHub CLI:
     $ gh auth login
     
  2. Or set an environment variable:
     $ export COPILOT_GITHUB_TOKEN=<your-token>
     
  3. Or use our interactive setup:
     $ copilot-research auth login

Need a subscription?
  Visit https://github.com/features/copilot to sign up.

---

❌ Subscription Required  

Your GitHub account is authenticated but doesn't have access to Copilot.

To use this tool, you need:
  • GitHub Copilot Individual ($10/month)
  • GitHub Copilot Business (via organization)
  • GitHub Copilot Enterprise (via organization)

Learn more: https://github.com/features/copilot

---

❌ Permission Error

Your token is valid but missing the "Copilot Requests" permission.

To fix:
  1. Create a new Personal Access Token at:
     https://github.com/settings/tokens/new
     
  2. Enable "Copilot Requests" permission
  
  3. Set the token:
     $ export COPILOT_GITHUB_TOKEN=<your-new-token>

Note: Classic tokens (ghp_*) are not supported for security reasons.
Use fine-grained PATs instead.
```

**Onboarding Flow**:
```
$ copilot-research "query"

👋 Welcome to Copilot Research!

To get started, we need to authenticate with GitHub Copilot.

You'll need:
  ✓ GitHub account
  ✓ Active Copilot subscription ($10/month)
  ✓ Terminal access

Choose authentication method:
  [1] GitHub CLI (recommended) - Uses existing gh authentication
  [2] Personal Access Token - Set COPILOT_GITHUB_TOKEN
  [3] Interactive OAuth - Browser-based device flow
  [Q] Quit

Selection: _
```

**Tests**:
- Mock `gh` command for testing
- Authentication priority order correct
- Error messages clear and actionable
- Onboarding flow guides users
- Timeout works correctly
- Parse output correctly
- Permission validation works

**Commit Template**:
```
[Providers] Implement GitHub Copilot provider with auth

Created GitHubCopilotProvider with:
- Multi-method authentication (PAT, gh CLI, OAuth)
- Priority-based credential detection
- Clear, actionable error messages
- Beautiful onboarding flow for new users
- Subscription and permission validation
- Configurable timeout

Follows CLI authentication best practices.

Tests: Auth methods, error messages, onboarding, execution
```

---

#### Task 5.3: OpenAI Provider Implementation
**Status**: ⬜ Not Started  
**Estimated**: 45 minutes  
**Priority**: P1

**Description**:
Implement OpenAI provider using official SDK with API key authentication and GPT-4 support.

**Acceptance Criteria**:
- Use OpenAI official Go SDK
- API key from environment variable or config
- Support multiple models (GPT-4, GPT-4-turbo, etc.)
- Handle rate limiting and errors
- Clear error messages for missing API key
- Streaming support (optional)

**Dependencies**: Task 5.1

**Implementation**:
```go
type OpenAIProvider struct {
    client  *openai.Client
    model   string
    timeout time.Duration
    apiKey  string
}

func (o *OpenAIProvider) Name() string { return "openai" }

func (o *OpenAIProvider) Query(ctx context.Context, prompt string, opts QueryOptions) (*Response, error) {
    // Use OpenAI SDK to query
}

func (o *OpenAIProvider) IsAuthenticated() bool {
    return o.apiKey != ""
}

func (o *OpenAIProvider) RequiresAuth() AuthInfo {
    if o.IsAuthenticated() {
        return AuthInfo{IsConfigured: true}
    }
    
    return AuthInfo{
        Type:         "apikey",
        IsConfigured: false,
        HelpURL:      "https://platform.openai.com/api-keys",
        Instructions: `OpenAI API key required.

Get your API key:
  1. Visit https://platform.openai.com/api-keys
  2. Create a new API key
  3. Set it in your environment:
     export OPENAI_API_KEY=sk-...
     
Or add to config:
  copilot-research config set providers.openai.api_key sk-...

Pricing: https://openai.com/pricing`,
    }
}

func (o *OpenAIProvider) Capabilities() ProviderCapabilities {
    return ProviderCapabilities{
        Streaming:      true,
        FunctionCall:   true,
        MaxTokens:      128000, // GPT-4-turbo
        SupportsImages: true,   // GPT-4-vision
    }
}
```

**Configuration**:
```yaml
providers:
  openai:
    api_key_env: OPENAI_API_KEY
    model: gpt-4-turbo-preview
    temperature: 0.7
    max_tokens: 4000
```

**Error Messages**:
```
❌ OpenAI API Key Missing

Set your OpenAI API key:
  $ export OPENAI_API_KEY=sk-...

Get a key at https://platform.openai.com/api-keys

---

❌ OpenAI Rate Limit Exceeded

You've hit OpenAI's rate limit. Please wait a moment.

Retry in: 30 seconds

Upgrade your plan at https://platform.openai.com/settings
```

**Tests**:
- API key detection works
- Query execution successful
- Rate limit handling
- Error messages helpful
- Model selection works

**Commit Template**:
```
[Providers] Add OpenAI provider implementation

Implemented OpenAIProvider with:
- Official OpenAI SDK integration
- API key authentication
- Multiple model support
- Rate limiting handling
- Clear error messages
- Streaming capability

Tests: Auth, querying, rate limits, errors
```

---

#### Task 5.4: Anthropic Claude Provider Implementation
**Status**: ⬜ Not Started  
**Estimated**: 45 minutes  
**Priority**: P2

**Description**:
Implement Anthropic Claude provider with API key authentication and Claude 3 support.

**Acceptance Criteria**:
- Use Anthropic official SDK
- API key from environment or config  
- Support Claude 3 models (Opus, Sonnet, Haiku)
- Handle rate limiting
- Clear error messages
- Support streaming (optional)

**Dependencies**: Task 5.1

**Implementation**:
```go
type AnthropicProvider struct {
    client  *anthropic.Client
    model   string
    timeout time.Duration
    apiKey  string
}

func (a *AnthropicProvider) Name() string { return "anthropic" }

func (a *AnthropicProvider) Query(ctx context.Context, prompt string, opts QueryOptions) (*Response, error) {
    // Use Anthropic SDK
}

func (a *AnthropicProvider) IsAuthenticated() bool {
    return a.apiKey != ""
}

func (a *AnthropicProvider) RequiresAuth() AuthInfo {
    if a.IsAuthenticated() {
        return AuthInfo{IsConfigured: true}
    }
    
    return AuthInfo{
        Type:         "apikey",
        IsConfigured: false,
        HelpURL:      "https://console.anthropic.com/",
        Instructions: `Anthropic API key required.

Get your API key:
  1. Visit https://console.anthropic.com/
  2. Create an account
  3. Generate API key
  4. Set in environment:
     export ANTHROPIC_API_KEY=sk-ant-...

Or add to config:
  copilot-research config set providers.anthropic.api_key sk-ant-...`,
    }
}

func (a *AnthropicProvider) Capabilities() ProviderCapabilities {
    return ProviderCapabilities{
        Streaming:      true,
        FunctionCall:   true,
        MaxTokens:      200000, // Claude 3
        SupportsImages: true,
    }
}
```

**Tests**:
- Authentication detection
- Query execution
- Error handling
- Model selection

**Commit Template**:
```
[Providers] Add Anthropic Claude provider

Implemented AnthropicProvider with:
- Official Anthropic SDK
- API key authentication  
- Claude 3 model support
- Streaming capability
- Clear error messages

Tests: Auth, querying, errors
```

---

#### Task 5.5: Provider Selection and Fallback Logic
**Status**: ⬜ Not Started  
**Estimated**: 30 minutes  
**Priority**: P1

**Description**:
Implement intelligent provider selection with automatic fallback on authentication or query failures.

**Acceptance Criteria**:
- Try primary provider first
- Fall back to secondary on auth failure
- Fall back on query errors (optional)
- Log provider switching
- User notification of fallback
- Configurable fallback behavior

**Dependencies**: Task 5.2, Task 5.3, Task 5.4

**Implementation**:
```go
type ProviderManager struct {
    factory  *ProviderFactory
    config   Config
    primary  string
    fallback string
}

func (pm *ProviderManager) Query(ctx context.Context, prompt string, opts QueryOptions) (*Response, error) {
    // Try primary provider
    provider, err := pm.getProvider(pm.primary)
    if err == nil && provider.IsAuthenticated() {
        resp, err := provider.Query(ctx, prompt, opts)
        if err == nil {
            return resp, nil
        }
        log.Printf("Primary provider %s failed: %v", pm.primary, err)
    }
    
    // Try fallback
    if pm.fallback != "" {
        log.Printf("Falling back to %s", pm.fallback)
        provider, err := pm.getProvider(pm.fallback)
        if err == nil && provider.IsAuthenticated() {
            resp, err := provider.Query(ctx, prompt, opts)
            if err == nil {
                fmt.Printf("ℹ️  Used %s (primary unavailable)\n", pm.fallback)
                return resp, nil
            }
        }
    }
    
    // All providers failed
    return nil, fmt.Errorf("all providers failed")
}

func (pm *ProviderManager) CheckAuthentication() ([]string, []string) {
    authenticated := []string{}
    unauthenticated := []string{}
    
    for _, name := range pm.factory.List() {
        provider, _ := pm.getProvider(name)
        if provider.IsAuthenticated() {
            authenticated = append(authenticated, name)
        } else {
            unauthenticated = append(unauthenticated, name)
        }
    }
    
    return authenticated, unauthenticated
}
```

**Configuration**:
```yaml
providers:
  primary: github-copilot
  fallback: openai
  auto_fallback: true  # Automatically try fallback on failure
  notify_fallback: true  # Show message when falling back
```

**User Messages**:
```
ℹ️  GitHub Copilot unavailable, using OpenAI GPT-4

⚠️  Primary provider failed, trying fallback...

❌ No authenticated providers available

Please authenticate at least one provider:
  • GitHub Copilot: gh auth login
  • OpenAI: export OPENAI_API_KEY=...
  • Anthropic: export ANTHROPIC_API_KEY=...
```

**Tests**:
- Primary provider used when available
- Fallback works on auth failure
- Fallback works on query failure
- User notifications display
- All providers failed handled

**Commit Template**:
```
[Providers] Add provider selection and fallback logic

Implemented intelligent provider management:
- Try primary provider first
- Automatic fallback on failure
- User notifications
- Configurable behavior
- Authentication status checking

Ensures queries succeed when any provider works.

Tests: Primary, fallback, notifications, config
```

---

#### Task 5.6: Auth Command for User Authentication
**Status**: ⬜ Not Started  
**Estimated**: 45 minutes  
**Priority**: P1

**Description**:
Add `auth` CLI command to help users authenticate providers interactively.

**Acceptance Criteria**:
- `auth status` - Show authentication status
- `auth login` - Interactive authentication
- `auth test` - Test provider connectivity
- Beautiful table output
- Guides users through setup
- Validates tokens

**Dependencies**: Task 5.5, Task 7.1

**Commands**:
```bash
copilot-research auth status
copilot-research auth login [provider]
copilot-research auth test [provider]
copilot-research auth logout [provider]
```

**Output Examples**:
```
$ copilot-research auth status

Authentication Status
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
Provider        Status          Method
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
github-copilot  ✅ Authenticated  gh CLI
openai          ❌ Not configured
anthropic       ❌ Not configured

Primary: github-copilot
Fallback: openai (not configured)

To authenticate OpenAI:
  export OPENAI_API_KEY=sk-...
  
Get a key: https://platform.openai.com/api-keys

---

$ copilot-research auth login

Choose a provider to authenticate:
  [1] GitHub Copilot (recommended)
  [2] OpenAI  
  [3] Anthropic Claude
  [Q] Quit

Selection: 1

Authenticating GitHub Copilot...

Opening browser for GitHub authentication...
✓ Authentication successful!

You're all set! Try:
  copilot-research "How do Swift actors work?"

---

$ copilot-research auth test github-copilot

Testing GitHub Copilot...
✓ Authentication valid
✓ Subscription active
✓ API accessible
✓ Test query successful

Provider ready to use!
```

**Tests**:
- Status shows correct information
- Login flow works
- Test validates connectivity
- Output formatted nicely

**Commit Template**:
```
[CLI] Add auth command for provider authentication

Implemented auth management:
- status: Show authentication status
- login: Interactive authentication
- test: Validate provider connectivity
- logout: Clear credentials

Beautiful table output guides users through setup.

Tests: Status, login, test, formatting
```

---

#### Task 5.7: Research Engine Core
**Status**: ⬜ Not Started  
**Estimated**: 45 minutes  
**Priority**: P0

**Description**:
Implement core research engine that coordinates prompt loading, Copilot querying, and result storage.

**Acceptance Criteria**:
- Orchestrates full research flow
- Loads appropriate prompt for mode
- Calls Copilot with rendered prompt
- Stores result in database
- Returns structured result
- Emits progress events for UI

**Dependencies**: Task 2.2, Task 3.1, Task 4.1

**Implementation**:
```go
type Engine struct {
    copilot *CopilotClient
    prompts *PromptLoader
    db      *DB
}

type ResearchOptions struct {
    Query      string
    Mode       string
    PromptName string
    NoStore    bool
}

type ResearchResult struct {
    Query     string
    Mode      string
    Content   string
    Duration  time.Duration
    SessionID int64
}

func (e *Engine) Research(ctx context.Context, opts ResearchOptions, progress chan<- string) (*ResearchResult, error)
```

**Progress events**:
- "Loading prompt..."
- "Querying GitHub Copilot..."
- "Processing results..."
- "Storing in database..."

**Tests**:
- Full research flow works
- Progress events emitted
- Context cancellation works
- Database storage optional (NoStore flag)
- Error handling at each step

**Commit Template**:
```
[Research] Implement core research engine

Created Engine that orchestrates:
- Prompt loading and rendering
- GitHub Copilot querying
- Result storage in database
- Progress event emission for UI

Supports context cancellation and optional storage.

Tests: Full flow, progress events, cancellation
```

---

### Phase 6: Beautiful UI (Bubble Tea)

#### Task 6.1: UI Components - Spinner and Progress
**Status**: ⬜ Not Started  
**Estimated**: 30 minutes  
**Priority**: P0

**Description**:
Create reusable Bubble Tea components for loading spinner and progress indicators using Charm bubbles.

**Acceptance Criteria**:
- Custom spinner with branded styling
- Progress bar component
- Status message display
- All use Lipgloss for styling
- Smooth animations

**Dependencies**: Task 1.1

**Implementation**:
```go
// internal/ui/spinner.go
type SpinnerModel struct {
    spinner  spinner.Model
    message  string
    styles   SpinnerStyles
}

// internal/ui/styles.go
var (
    TitleStyle = lipgloss.NewStyle().
        Bold(true).
        Foreground(lipgloss.Color("205"))
    
    SpinnerStyle = lipgloss.NewStyle().
        Foreground(lipgloss.Color("69"))
    
    ResultStyle = lipgloss.NewStyle().
        Border(lipgloss.RoundedBorder()).
        Padding(1, 2)
)
```

**Tests**:
- Spinner renders correctly
- Styles apply properly
- Can update message dynamically

**Commit Template**:
```
[UI] Create Bubble Tea spinner and progress components

Built reusable UI components:
- Animated spinner with custom styling
- Progress indicator
- Status message display

Uses Lipgloss for beautiful terminal styling.

Tests: Rendering, style application
```

---

#### Task 6.2: Main Research UI Model
**Status**: ⬜ Not Started  
**Estimated**: 60 minutes  
**Priority**: P0

**Description**:
Implement main Bubble Tea model for research UI that shows live progress, handles research completion, and displays results beautifully.

**Acceptance Criteria**:
- Shows spinner during research
- Updates status messages from progress channel
- Displays result with nice formatting
- Handles errors gracefully
- Supports quit/interrupt (Ctrl-C)
- Scrollable result view for long output

**Dependencies**: Task 4.2, Task 5.1

**States**:
- `researching` - Show spinner + progress
- `complete` - Show formatted result
- `error` - Show error message

**Implementation**:
```go
type ResearchModel struct {
    state    string
    query    string
    mode     string
    
    spinner  SpinnerModel
    status   string
    result   *ResearchResult
    err      error
    
    viewport viewport.Model
    ready    bool
}

func NewResearchModel(query, mode string) ResearchModel
func (m ResearchModel) Init() tea.Cmd
func (m ResearchModel) Update(msg tea.Msg) (tea.Model, tea.Cmd)
func (m ResearchModel) View() string
```

**Custom messages**:
```go
type progressMsg string
type completeMsg ResearchResult
type errorMsg error
```

**Tests**:
- UI updates on progress
- Result displays correctly
- Error handling works
- Viewport scrolling functional

**Commit Template**:
```
[UI] Implement main research UI model

Created Bubble Tea model with states:
- researching: Spinner + live progress
- complete: Formatted result display
- error: Error message

Supports scrolling for long results and graceful quit.

Tests: State transitions, rendering, error handling
```

---

### Phase 7: CLI Commands (Cobra)

#### Task 7.1: Root Command and Basic Structure
**Status**: ⬜ Not Started  
**Estimated**: 30 minutes  
**Priority**: P0

**Description**:
Set up Cobra root command with global flags, config loading, and help text.

**Acceptance Criteria**:
- `cmd/root.go` with root command
- Global flags: `--mode`, `--output`, `--prompt`, `--quiet`, `--json`
- Config file loading from `~/.copilot-research/config.yaml`
- Beautiful help text with examples
- Version command

**Dependencies**: Task 1.1, Task 1.2

**Global flags**:
```go
--mode string      Research mode (quick|deep|compare|synthesis)
--output string    Output file path
--prompt string    Prompt template to use
--quiet           Quiet mode (no UI, just output)
--json            Output as JSON
--no-store        Don't save to database
```

**Tests**:
- Flags parse correctly
- Config loads successfully
- Help text displays
- Version shows correctly

**Commit Template**:
```
[CLI] Set up Cobra root command structure

Created root command with:
- Global flags for mode, output, prompt, etc.
- Config file loading
- Help text and examples
- Version command

Tests: Flag parsing, config loading
```

---

#### Task 7.2: Main Research Command
**Status**: ⬜ Not Started  
**Estimated**: 45 minutes  
**Priority**: P0

**Description**:
Implement main research command that accepts query as argument or from stdin/file, runs research engine, and displays UI.

**Acceptance Criteria**:
- Accept query as: argument, --input file, or stdin
- Initialize engine with DB and prompts
- Run Bubble Tea UI (unless --quiet)
- Output result to stdout or file
- Handle all error cases gracefully
- Support JSON output format

**Dependencies**: Task 4.2, Task 5.2, Task 6.1

**Usage**:
```bash
copilot-research "query"
copilot-research --input file.txt
echo "query" | copilot-research
copilot-research "query" --output report.md
copilot-research "query" --json
copilot-research "query" --quiet
```

**Tests**:
- All input methods work
- UI launches correctly
- Output formats correct
- File output works
- Quiet mode bypasses UI

**Commit Template**:
```
[CLI] Implement main research command

Added research command supporting:
- Query from argument, file, or stdin
- Bubble Tea UI for progress
- Multiple output formats (markdown, JSON)
- Quiet mode for scripting
- File output option

Tests: Input methods, output formats, UI/quiet modes
```

---

#### Task 7.3: History Command
**Status**: ⬜ Not Started  
**Estimated**: 30 minutes  
**Priority**: P1

**Description**:
Implement history command to view past research sessions with search and filtering.

**Acceptance Criteria**:
- List recent sessions
- Search by query text
- Filter by mode
- Show session details
- Clear history option
- Pretty table output

**Dependencies**: Task 2.2, Task 6.1

**Subcommands**:
```bash
copilot-research history              # List recent
copilot-research history --search "Swift"
copilot-research history --mode deep
copilot-research history --id 123     # Show specific
copilot-research history --clear      # Clear all
```

**Tests**:
- Lists sessions correctly
- Search works
- Filtering works
- Clear requires confirmation

**Commit Template**:
```
[CLI] Add history command for viewing past research

Implemented history command with:
- List recent sessions
- Search by query
- Filter by mode
- View specific session
- Clear history option

Tests: List, search, filter, clear
```

---

#### Task 7.4: Stats Command
**Status**: ⬜ Not Started  
**Estimated**: 20 minutes  
**Priority**: P2

**Description**:
Show statistics about research usage, patterns, and database size.

**Acceptance Criteria**:
- Total sessions count
- Mode usage breakdown
- Database size
- Most common queries
- Nice chart/table formatting

**Dependencies**: Task 2.2, Task 6.1

**Output**:
```
Research Statistics
───────────────────────────────────

Total Sessions: 127
Database Size: 1.2 MB

Mode Usage:
  quick     82 (65%)
  deep      32 (25%)
  compare   13 (10%)

Top Queries:
  1. Swift 6 actors (23 times)
  2. iOS 26 APIs (15 times)
  3. SwiftUI best practices (12 times)
```

**Tests**:
- Calculates stats correctly
- Formats output nicely
- Handles empty database

**Commit Template**:
```
[CLI] Add stats command for usage analytics

Shows research statistics:
- Total sessions and database size
- Mode usage breakdown
- Most common queries

Beautiful table formatting with percentages.

Tests: Calculation, formatting, empty DB
```

---

#### Task 7.5: Config Command
**Status**: ⬜ Not Started  
**Estimated**: 25 minutes  
**Priority**: P2

**Description**:
Manage configuration settings via CLI.

**Acceptance Criteria**:
- Show current config
- Set config values
- Reset to defaults
- Validate config values
- Create config file if missing

**Dependencies**: Task 6.1

**Subcommands**:
```bash
copilot-research config                 # Show all
copilot-research config get key
copilot-research config set key value
copilot-research config reset
```

**Tests**:
- Get/set config works
- Validation catches bad values
- Reset works correctly

**Commit Template**:
```
[CLI] Add config command for settings management

Implemented config management:
- Show current configuration
- Get/set individual values
- Reset to defaults
- Value validation

Tests: Get, set, validation, reset
```

---

### Phase 8: Polish and Documentation

#### Task 8.1: Add Makefile Targets
**Status**: ⬜ Not Started  
**Estimated**: 15 minutes  
**Priority**: P1

**Description**:
Create comprehensive Makefile with all common development tasks.

**Acceptance Criteria**:
- build, test, install, clean targets
- fmt, lint targets
- run target for quick testing
- help target with descriptions

**Dependencies**: None

**Targets**:
```makefile
.PHONY: build test install clean fmt lint run help

build:
	go build -o copilot-research

test:
	go test ./... -v -cover

install:
	go install

clean:
	rm -f copilot-research
	rm -rf tmp/*

fmt:
	gofmt -s -w .

lint:
	golangci-lint run

run:
	go run main.go

help:
	@echo "Available targets:"
	@echo "  build   - Build binary"
	@echo "  test    - Run tests"
	@echo "  install - Install to GOPATH"
	@echo "  clean   - Remove build artifacts"
	@echo "  fmt     - Format code"
	@echo "  lint    - Run linter"
	@echo "  run     - Run directly"
```

**Commit Template**:
```
[Build] Add comprehensive Makefile

Added targets for:
- Building, testing, installing
- Code formatting and linting
- Cleanup and running
- Help documentation

Makes development workflow easier.
```

---

#### Task 8.2: GitHub Actions CI
**Status**: ⬜ Not Started  
**Estimated**: 30 minutes  
**Priority**: P1

**Description**:
Set up GitHub Actions for automated testing, building, and releases.

**Acceptance Criteria**:
- Test workflow (runs on PR and push)
- Build workflow (multi-platform)
- Release workflow (tag-based)
- Coverage reporting
- Linting checks

**Dependencies**: Task 7.1

**Workflows**:
- `.github/workflows/test.yml` - Run tests
- `.github/workflows/build.yml` - Build binaries
- `.github/workflows/release.yml` - Create releases

**Tests**:
- Workflows trigger correctly
- Tests run successfully
- Artifacts produced

**Commit Template**:
```
[CI] Add GitHub Actions workflows

Created workflows for:
- Running tests on PRs
- Building multi-platform binaries
- Creating releases from tags
- Code coverage reporting

Tests: Workflow syntax validation
```

---

#### Task 8.3: Usage Examples and Documentation
**Status**: ⬜ Not Started  
**Estimated**: 45 minutes  
**Priority**: P1

**Description**:
Write comprehensive documentation with examples, screenshots, and guides.

**Acceptance Criteria**:
- README.md updated with full usage
- docs/USAGE.md with detailed examples
- docs/PROMPTS.md explaining prompt system
- docs/DEVELOPMENT.md for contributors
- GIF/screenshots of UI in action

**Dependencies**: All previous tasks

**Documentation needed**:
- Quick start guide
- All command examples
- Prompt customization guide
- Configuration reference
- Contributing guide
- FAQ

**Commit Template**:
```
[Docs] Add comprehensive documentation and examples

Created detailed documentation:
- Complete usage guide
- Prompt customization examples
- Configuration reference
- Development setup guide
- FAQ and troubleshooting

Includes GIFs showing UI in action.
```

---

#### Task 8.4: Final Polish and Testing
**Status**: ⬜ Not Started  
**Estimated**: 60 minutes  
**Priority**: P1

**Description**:
Final round of testing, bug fixes, and polish before v1.0.0 release.

**Acceptance Criteria**:
- All tests passing
- No linter warnings
- Error messages are helpful
- Help text is clear
- Edge cases handled
- Performance acceptable
- Memory usage reasonable

**Dependencies**: All previous tasks

**Checklist**:
- [ ] All commands work correctly
- [ ] UI animations smooth
- [ ] Database operations fast
- [ ] Error handling comprehensive
- [ ] No memory leaks
- [ ] Cross-platform testing (macOS, Linux)
- [ ] README accurate
- [ ] Installation works

**Commit Template**:
```
[Polish] Final testing and bug fixes for v1.0.0

Final polish includes:
- Bug fixes from testing
- Improved error messages
- Performance optimizations
- Cross-platform testing
- Documentation corrections

Ready for v1.0.0 release.
```

---

## Task Summary

**Total Tasks**: 33

### By Phase:
- Phase 1 (Foundation): 3 tasks
- Phase 2 (Database): 2 tasks
- Phase 3 (Knowledge Management): 5 tasks
- Phase 4 (Prompts): 2 tasks
- Phase 5 (Multi-Provider Authentication & Integration): 7 tasks
- Phase 6 (UI): 2 tasks
- Phase 7 (CLI): 6 tasks
- Phase 8 (Polish): 4 tasks

### By Priority:
- P0 (Blocker): 16 tasks
- P1 (High): 13 tasks
- P2 (Medium): 4 tasks

**Estimated Total Time**: ~20 hours

---

## Progress Tracking

**Completed**: 6/33 (18%)  
**In Progress**: 0/33 (0%)  
**Not Started**: 27/33 (82%)

---

## Notes

- TDD approach for all tasks
- Commit after each task completion
- Update this plan with commit hashes
- Add learnings to docs/agents.md
- Test with real gh copilot integration frequently

---

**Last Updated**: 2025-11-17  
**Next Task**: Task 3.2 - Knowledge Manager Implementation

---

## Key Design Decisions

### Multi-Provider Architecture
- **Pattern**: Abstraction layer with factory + adapter patterns
- **Rationale**: Allows easy addition of new AI providers without changing core logic
- **Based on**: Industry best practices from LiteLLM, AISuite, AWS Multi-Provider Gateway
- **Authentication**: Priority-based credential checking with clear error messages
- **Fallback**: Automatic failover to secondary provider on authentication or query failures

### Knowledge Management
- **Storage**: Git-based versioning for audit trail and rollback
- **Format**: Markdown with YAML frontmatter for human readability  
- **Consolidation**: Automatic deduplication and content merging
- **Rules**: YAML-based user preferences for content filtering

### User Experience
- **Onboarding**: Friendly, guided authentication flow with clear instructions
- **Error Messages**: Specific, actionable, empathetic (following UX best practices)
- **Feedback**: Live progress updates during research with Bubble Tea
- **Documentation**: Extensive examples and troubleshooting guides
