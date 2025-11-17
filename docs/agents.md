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

### 2025-11-17: CLI Commands Implementation (Cobra + Lipgloss)

**Cobra Command Structure**:
- Root command with global flags (--quiet, --json, --output)
- Subcommands organized by feature (knowledge, auth, history, etc.)
- Each subcommand has its own flags and validation
- Use `cobra.ExactArgs(1)` for commands requiring specific argument count
- `RunE` instead of `Run` for proper error handling

**Lipgloss Styling Best Practices**:
```go
// Define styles once, reuse everywhere
var (
    titleStyle = lipgloss.NewStyle().
        Bold(true).
        Foreground(lipgloss.Color("205"))
    
    successStyle = lipgloss.NewStyle().
        Foreground(lipgloss.Color("42"))
)

// Use consistent color palette
// - 205: Pink/magenta for titles
// - 86: Green for headers
// - 240: Gray for info/metadata
// - 42: Bright green for success
// - 196: Red for errors
```

**Table Output with tabwriter**:
- Use `text/tabwriter` for aligned columns
- Set padding to 3 spaces for readability
- Flush writer before returning
- Style headers differently from data
```go
w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
fmt.Fprintf(w, "%s\t%s\n", headerStyle.Render("Name"), headerStyle.Render("Value"))
w.Flush()
```

**Time Formatting for CLI**:
- Use relative time for recent items ("2 hours ago", "just now")
- Switch to absolute dates for older items (>30 days)
- Helper function pattern:
```go
func formatTimeAgo(t time.Time) string {
    duration := time.Since(t)
    switch {
    case duration < time.Hour:
        return fmt.Sprintf("%d minutes ago", int(duration.Minutes()))
    case duration < 24*time.Hour:
        return fmt.Sprintf("%d hours ago", int(duration.Hours()))
    default:
        return t.Format("2006-01-02")
    }
}
```

**Editor Integration**:
- Respect $EDITOR environment variable
- Fall back to sensible default (nano, not vim - easier for new users)
- Use temp files for editing
- Clean up temp files after use
```go
editor := os.Getenv("EDITOR")
if editor == "" {
    editor = "nano"
}
cmd := exec.Command(editor, tempfile.Name())
cmd.Stdin = os.Stdin
cmd.Stdout = os.Stdout
cmd.Stderr = os.Stderr
cmd.Run()
```

**Empty State Messages**:
- Always provide helpful guidance when no data exists
- Tell users exactly how to add first item
- Example:
```
No knowledge entries found.

Add your first entry with:
  copilot-research knowledge add <topic>
```

**Flag Patterns**:
- Use long flags with descriptive names (--exclude, --reason)
- Provide short flags for common operations (-o for output)
- Set reasonable defaults
- Validate required combinations

**Common CLI Errors Fixed**:
1. `undefined: fmt` → Add import in test file
2. Editor doesn't open → Check $EDITOR is set or use fallback
3. Tables misaligned → Ensure consistent tab spacing in tabwriter
4. Color bleeding → Reset styles after use
5. Long topic names overflow → Truncate or wrap intelligently

### 2025-11-17: Prompt Loader & Template System

**Embed vs File System**:
- Initially tried `//go:embed` but path resolution is tricky
- Embedded files must be in same directory or subdirectories
- Solution: Load from filesystem with fallback to minimal default
- Benefit: Users can customize prompts without rebuilding

**Frontmatter Parsing Strategy**:
- Reuse logic from knowledge system (already tested)
- Split by lines, look for `---` delimiters
- YAML unmarshal for frontmatter
- Remaining content is template body
- Same pattern works for both knowledge and prompts

**Template Variable Substitution**:
- Simple string replacement works well
- Format: `{{variable_name}}`
- Use map[string]string for variables
- Apply all substitutions in order
```go
func Render(template string, vars map[string]string) string {
    result := template
    for key, value := range vars {
        placeholder := fmt.Sprintf("{{%s}}", key)
        result = strings.ReplaceAll(result, placeholder, value)
    }
    return result
}
```

**Caching Strategy**:
- Cache loaded prompts in memory (they rarely change)
- Use sync.RWMutex for thread-safe access
- Provide Reload() method to clear cache
- Check cache before hitting filesystem

**Prompt File Structure**:
```markdown
---
name: quick
description: Quick research prompt
version: 1.0.0
mode: quick
---

Your prompt content here...

Use {{query}} and {{mode}} variables.
```

**Multiple Prompt Modes**:
- `default.md` - Balanced, comprehensive research
- `quick.md` - Fast overviews (5 min read time)
- `deep-dive.md` - Exhaustive analysis with examples
- `compare.md` - Side-by-side comparison with matrix
- `synthesis.md` - Multi-source integration

Each optimized for different use cases and reading times.

**Testing Prompt Templates**:
- Validate frontmatter fields exist and are non-empty
- Check required template variables present ({{query}}, {{mode}})
- Test variable substitution works correctly
- Ensure no leftover {{}} after rendering
- Use table-driven tests for multiple prompts

**Prompt Design Principles**:
1. Clear instructions for AI assistant role
2. Explicit output format with sections
3. Emphasis on accuracy and citations
4. Examples of good responses
5. Mode-specific guidelines
6. Template variables for customization

**Common Errors Fixed**:
1. `pattern default.md: no matching files found` → Wrong embed path
2. Test fails with minimal default → Load from actual prompts directory
3. Variable not replaced → Typo in template variable name
4. Empty prompt content → Check file exists before parsing

### 2025-11-17: Provider Abstraction System

**Provider Architecture Pattern**:
- Interface-based design allows multiple AI backends
- Factory pattern for registration and retrieval
- Manager pattern for primary/fallback logic
- Thread-safe with sync.RWMutex for concurrent access

**AIProvider Interface Design**:
```go
type AIProvider interface {
    Name() string        // Unique identifier
    Query() (*Response, error)  // Main query method
    IsAuthenticated() bool      // Fast auth check
    RequiresAuth() AuthInfo     // User guidance
    Capabilities() ProviderCapabilities  // Feature detection
}
```

**Key Implementation Principles**:
1. **Context-aware**: All Query() methods must respect context cancellation
2. **Timeout handling**: Use context.WithTimeout for all external calls
3. **Authentication priority**: Check env vars → config files → CLI tools
4. **Error clarity**: Convert API errors to actionable user messages
5. **Standardized responses**: All providers return same Response format

**Authentication Patterns**:
- Fast checks: IsAuthenticated() should be < 1 second
- Cache auth status: Don't validate credentials on every call
- Multiple methods: Support env vars, config files, CLI tools
- Priority order: Most direct (env var) to most complex (OAuth)
- Clear instructions: RequiresAuth() tells users exactly what to do

**GitHub Copilot Provider Specifics**:
- Wraps `gh copilot suggest` CLI command
- Three auth methods: COPILOT_GITHUB_TOKEN > GH_TOKEN > gh CLI
- Uses `exec.CommandContext` for timeout support
- Parses markdown output from gh copilot
- Estimates token usage (4 chars per token)
- Handles subscription/permission errors gracefully

**Testing Strategies**:
- Mock providers for unit tests (see MockProvider in provider_test.go)
- Use `t.Skip()` when system auth prevents unauthenticated tests
- Test auth priority order with multiple env vars set
- Test context cancellation and timeout behavior
- Integration tests separate from unit tests (require real credentials)

**Provider Registration Pattern**:
```go
factory := NewProviderFactory()
provider := NewGitHubCopilotProvider(timeout)
factory.Register("github-copilot", provider)

// Manager handles fallback
manager := NewProviderManager(factory, "primary", "fallback")
resp, err := manager.Query(ctx, prompt, opts)
```

**Response Standardization**:
```go
type Response struct {
    Content    string                 // Clean text response
    Provider   string                 // Which provider answered
    Model      string                 // Model used
    TokensUsed TokenUsage            // Usage tracking
    Duration   time.Duration          // Performance metric
    Metadata   map[string]interface{} // Provider-specific data
}
```

**Common Provider Errors Fixed**:
1. Not checking context.Done() → Add `ctx.Err()` checks
2. Hanging on timeout → Use `context.WithTimeout`
3. Unclear auth errors → Parse and provide helpful messages
4. Token estimation wrong → Use provider data if available
5. Tests fail due to system auth → Use `t.Skip()` conditionally

**For Future AI Agents**:
- Read `docs/provider-implementation-guide.md` - comprehensive guide
- Follow pattern in `internal/provider/github_copilot.go`
- Implement all interface methods with proper error handling
- Write tests that adapt to local authentication state
- Document auth methods and API requirements
- Use standard Response format for consistency

---

*Keep updated as you discover patterns and solve problems.*
