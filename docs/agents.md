# Agent Context & Learnings

This document captures context, patterns, and solutions for AI agents working on this codebase.

## Project Overview

**Copilot Research** is a beautiful CLI tool for AI-powered research built with Go and the Charm ecosystem (Bubble Tea + Lipgloss).

**Key principles:**
- Design first (beautiful > functional)
- Live feedback always (never appear hung)
- Simple by default, powerful when needed
- Unix-friendly (scriptable, pipeable)
- Learn from history

## Technology Stack

- **Language**: Go 1.21+
- **TUI Framework**: Bubble Tea (charmbracelet)
- **Styling**: Lipgloss (charmbracelet)
- **Spinners/Progress**: Charm ecosystem components
- **Database**: SQLite3
- **CLI Parsing**: Cobra
- **Testing**: Standard Go testing + testify

## Architecture Decisions

### Why Go?
- Fast compilation and execution
- Single binary deployment
- Excellent CLI/TUI libraries (Charm ecosystem)
- Cross-platform support
- Strong standard library

### Why Bubble Tea?
- Most polished TUI framework available
- Active development and community
- Excellent live updates and animations
- Professional appearance out of box
- Elm architecture (predictable state management)

### Why Separate Prompts?
- Better git history (see what changed)
- Easy A/B testing
- Swappable backends (Claude, GPT, etc.)
- Community contributions
- Version control friendly

## Directory Structure

```
copilot-research/
├── cmd/                    # CLI entry points
│   └── root.go            # Root command
├── internal/              # Private application code
│   ├── research/          # Research engine
│   │   ├── engine.go      # Core research logic
│   │   └── modes.go       # Research modes
│   ├── ui/                # Bubble Tea UI components
│   │   ├── spinner.go     # Loading states
│   │   ├── progress.go    # Progress indicators
│   │   └── output.go      # Result rendering
│   ├── db/                # Database layer
│   │   ├── sqlite.go      # SQLite implementation
│   │   └── models.go      # Data models
│   ├── prompts/           # Prompt management
│   │   ├── loader.go      # Load prompts from files
│   │   └── registry.go    # Available prompts
│   └── config/            # Configuration
│       └── config.go      # Config management
├── prompts/               # Prompt templates (versioned)
│   ├── default.md         # Default Copilot prompt
│   ├── claude.md          # Claude-optimized
│   └── deep-dive.md       # Deep research mode
├── docs/                  # Documentation
│   ├── agents.md          # This file
│   ├── plan.md            # Implementation plan
│   └── architecture.md    # Architecture docs
├── tmp/                   # Temporary files (gitignored)
├── go.mod                 # Go dependencies
├── go.sum                 # Dependency checksums
├── Makefile               # Build commands
└── README.md              # User documentation
```

## Development Workflow

### TDD Approach
1. Write failing test
2. Run test (prove it fails)
3. Implement minimum code to pass
4. Run test (prove it passes)
5. Refactor if needed
6. Commit with detailed message
7. Push to origin
8. Update plan.md

### Commit Messages
```
[Component] Brief description

Detailed explanation of what changed and why.
Mention any trade-offs or decisions made.

Tests: Added X, modified Y
```

### Testing Strategy
- Unit tests for all business logic
- Integration tests for CLI commands
- UI tests for Bubble Tea components
- Test with real `gh copilot` calls (optional flag)
- Benchmark performance-critical paths

## Common Patterns

### Bubble Tea Model Pattern
```go
type model struct {
    state       string
    spinner     spinner.Model
    progress    float64
    result      string
    err         error
}

func (m model) Init() tea.Cmd {
    return m.spinner.Tick
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        if msg.String() == "q" {
            return m, tea.Quit
        }
    case spinner.TickMsg:
        var cmd tea.Cmd
        m.spinner, cmd = m.spinner.Update(msg)
        return m, cmd
    case researchCompleteMsg:
        m.result = msg.result
        return m, tea.Quit
    }
    return m, nil
}

func (m model) View() string {
    if m.err != nil {
        return errorStyle.Render(m.err.Error())
    }
    if m.result != "" {
        return resultStyle.Render(m.result)
    }
    return fmt.Sprintf("%s Researching...\n", m.spinner.View())
}
```

## Known Issues & Solutions

### Issue: Bubble Tea blocks during research
**Problem**: Can't update UI while waiting for `gh copilot`
**Solution**: Run research in goroutine, send progress messages via channel

### Issue: SQLite concurrent access
**Problem**: Multiple goroutines accessing DB
**Solution**: Use WAL mode
```go
db.Exec("PRAGMA journal_mode=WAL")
```

## Useful Commands

### Development
```bash
go test ./... -v
go build -o copilot-research
./copilot-research "test"
go install
```

### Testing GitHub Copilot
```bash
gh copilot suggest "Research Swift 6 actors"
gh auth status
```

---

## Learning Log

### 2025-11-17: Initial Setup
- Created repository structure
- Chose Go + Bubble Tea for best TUI experience
- Decided on separate prompt files for flexibility

### 2025-11-17: CLI Research Tool Design
**Key insight**: When building research tools, the experience matters as much as the results.

**Live Feedback is Critical**:
- Users need constant visual feedback during long operations
- Bubble Tea's spinner + progress components solve this elegantly
- Never let CLI appear frozen (use async operations + UI updates)

**Database for Learning**:
- SQLite enables accumulating research knowledge over time
- Store: query → results → citations → timestamp
- Future queries can reference past research
- Use WAL mode for concurrent access: `PRAGMA journal_mode=WAL`

**Prompt Management**:
- Separate prompt files (`.md` format) for git-friendly versioning
- Easy to swap between AI backends (Copilot, Claude, GPT)
- Template variables for dynamic content injection
- Default prompt should be comprehensive but not overwhelming

### 2025-11-17: Go & Bubble Tea Learnings

**Bubble Tea Best Practices**:
- Run long operations in goroutines, communicate via messages
- Create custom message types for different events
- Use `tea.Cmd` for async operations
- Lipgloss for consistent, beautiful styling
- Spinner + progress components from Charm ecosystem

**Go CLI Patterns**:
- Cobra for command structure and flag parsing
- Accept input from: args, stdin, files (make it Unix-friendly)
- Output format options: pretty (default), JSON, markdown
- Handle SIGINT/SIGTERM gracefully
- Single binary deployment (no dependencies)

**Testing Go CLIs**:
- Use `exec.Command` for integration tests
- Table-driven tests for different input formats
- Mock external commands (`gh copilot`) for unit tests
- Benchmark long-running operations
- Test UI rendering with Bubble Tea's testing helpers

### 2025-11-17: Research Agent Methodology

**Multi-Query Synthesis Pattern**:
1. Break down complex topic into focused sub-queries
2. Execute parallel research on each sub-query
3. Deduplicate findings across sources
4. Synthesize into coherent narrative
5. Provide inline citations

**Quality Indicators**:
- Multiple source corroboration
- Recency of information (timestamp everything)
- Authority of source (official docs > blogs)
- Depth vs breadth trade-off based on query

**Iterative Refinement**:
- Initial broad sweep (understand landscape)
- Follow-up targeted queries (fill gaps)
- Cross-reference contradictions
- Update knowledge base with learnings

### 2025-11-17: Knowledge Management System

**Git-Based Knowledge Versioning**:
- Store all knowledge in `~/.copilot-research/knowledge/` with Git tracking
- Automatic commits for every knowledge change with descriptive messages
- View history of any topic: when info changed and why
- Rollback capability if needed
- Separate repos for app code vs knowledge (knowledge is user data)

**Knowledge Structure**:
```
~/.copilot-research/knowledge/
├── .git/                     # Auto-managed Git repo
├── topics/                   # Topic-specific knowledge
│   ├── swift-concurrency.md
│   ├── swiftui-patterns.md
│   └── ...
├── patterns/                 # Learned patterns
│   ├── common-errors.md
│   └── best-practices.md
├── rules/                    # User preferences
│   ├── preferences.yaml
│   └── exclusions.yaml
└── MANIFEST.yaml            # Central registry
```

**Markdown + YAML Frontmatter**:
- Human-readable and editable
- Version control friendly (good diffs)
- Easy to parse programmatically
- Standard format used by many tools
```yaml
---
topic: swift-concurrency
version: 3
confidence: 0.95
tags: [swift, concurrency, actors]
source: https://docs.swift.org/
created: 2025-11-17T12:00:00Z
updated: 2025-11-17T14:00:00Z
---

# Swift Concurrency
[Content...]
```

**Deduplication Strategy**:
- Generate SHA-256 ID from topic + content
- Compare new knowledge against existing
- If similar (>90% match), merge instead of duplicate
- Keep highest confidence version
- Preserve all unique information
- Update version number on merge

**Auto-Learning from Research**:
- Analyze successful research results
- Extract key topics and patterns
- Calculate confidence score based on source quality
- Prompt user for approval before storing
- Build knowledge base over time automatically

**Rule System for Preferences**:
```yaml
rules:
  - type: exclude
    pattern: "Model View Controller|MVC"
    reason: "Using MV architecture instead"
  
  - type: prefer
    pattern: "Swift Testing"
    over: "XCTest"
```

**Key Decisions**:
- User home dir (`~/`) not project dir (knowledge is global across projects)
- Git for versioning (proven, reliable, no custom format)
- Markdown for human readability (can edit manually)
- YAML frontmatter for metadata (standard, well-supported)
- Automatic consolidation to prevent bloat

**Testing Approach**:
- Test knowledge CRUD operations
- Test frontmatter parsing/serialization
- Test Git operations (commit, history, diff)
- Test deduplication algorithm
- Test rule matching and application

### 2025-11-17: KnowledgeManager Implementation Learnings

**Thread-Safe Concurrent Access**:
- Used `sync.RWMutex` for cache access (read-write lock)
- Multiple readers can access simultaneously
- Writers get exclusive lock
- Lock at method level, not individual operations
- Critical for concurrent CLI operations

**Git Command Execution in Go**:
- Use `exec.Command` for git operations
- Set `cmd.Dir` to set working directory
- Use `CombinedOutput()` to capture both stdout/stderr
- Always check for specific error messages (e.g., "did not match any files")
- Initialize git config (user.name, user.email) on repo creation

**Filename Safety**:
- Topics can contain `/` or spaces - must sanitize for filenames
- Replace `/` with `-` (swift/feature → swift-feature)
- Replace spaces with `_`
- Filter to alphanumeric + `-` + `_` + `.` only
- Use `strings.Map()` for character filtering

**Frontmatter Parsing**:
- Split by lines, not by string search (handles edge cases)
- Look for `---` as trimmed line content
- First `---` starts frontmatter, second ends it
- Trim body content (remove leading/trailing whitespace)
- YAML parsing is strict - test edge cases

**Deduplication Algorithm**:
- Word overlap metric for simplicity (can upgrade to embeddings)
- Use map for "toRemove" to avoid double-marking
- Threshold of 0.85 similarity works well
- Keep higher confidence or newer version
- Skip already-marked entries in nested loops

**Testing File Operations**:
- Use `t.TempDir()` - automatically cleaned up
- Git operations may fail in temp dirs (path length, permissions)
- Test with realistic content similarity (not just different words)
- For concurrent tests, use channels to synchronize goroutines

**Search Implementation**:
- Case-insensitive matching (strings.ToLower)
- Search across: topic, content, tags
- Helper function for tag matching
- Return all matches (let caller filter/rank)

**Common Errors Fixed**:
1. Undefined function → Missing import (strings, fmt)
2. Invalid filepath → Special character in topic name
3. Git command fails → Wrong working directory or bad permissions
4. Tests fail → Content not similar enough for dedup threshold
5. Empty search results → Case sensitivity or whitespace in parsed content

**Performance Considerations**:
- Cache all knowledge in memory for fast access
- Reload cache only on initialization
- Git operations are slow - minimize commits
- Use `commitAll` for batch operations
- Consider async git operations for large knowledge bases

### 2025-11-17: RuleEngine Implementation Learnings

**Mutex Deadlock Prevention**:
- NEVER call functions that acquire locks while holding a lock yourself
- Pattern that causes deadlock:
  ```go
  func RemoveRule() {
      re.mu.Lock()
      defer re.mu.Unlock()
      return re.save()  // save() tries to get RLock = DEADLOCK!
  }
  ```
- Solution: Unlock before calling save()
  ```go
  func RemoveRule() {
      re.mu.Lock()
      re.rules = newRules
      re.mu.Unlock()  // Unlock BEFORE calling save()
      return re.save()
  }
  ```
- Test symptom: Test hangs indefinitely with no output

**Git Operations in Tests**:
- Auto-committing on every save causes tests to hang
- Git commands in temp directories can fail unexpectedly
- Solution: Don't auto-commit in library code, let caller decide
- Document that manual commit is needed for Git tracking

**Regex Pattern Matching**:
- Use `regexp.Compile()` once, cache if performance matters
- `ReplaceAllString()` for simple replacements
- `MatchString()` for presence checks
- Case-insensitive: Use `(?i)` prefix in pattern

**Rule Application Order**:
- Order matters when applying multiple rules
- Exclude rules can affect what prefer rules can replace
- Solution: Apply rules in specific order or be smarter about removal
- For exclude: Replace matched text, not whole sentences
  - Bad: Remove entire sentence containing "MVVM"
  - Good: Replace just "MVVM" with ""

**YAML Persistence**:
- Use struct tags for YAML marshaling
- Empty slices marshal as `[]` not `null`
- Read file, unmarshal, modify, marshal, write
- Handle missing file gracefully (os.IsNotExist)

**Test-Driven Development Flow**:
1. Write test that calls undefined function → compile error ✓
2. Implement minimal function signature → test fails ✓
3. Implement logic → tests pass ✓
4. Refactor → tests still pass ✓
5. Commit with detailed message

**UUID Generation**:
- Use `github.com/google/uuid` for unique IDs
- `uuid.New().String()` for new UUID
- Store as string in structs (easier for YAML)

**Common Errors Fixed**:
1. Deadlock → Unlock before calling methods that lock
2. Tests hang → Remove blocking git operations
3. Unused import → Remove after refactoring
4. Rule application fails → Fix pattern matching logic

---

*Keep updated as you discover patterns and solve problems.*
