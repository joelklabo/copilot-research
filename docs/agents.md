# Agent Context & Learnings

This document captures essential context, patterns, and solutions for AI agents working on this codebase.

## Project Overview

**Copilot Research** is a CLI tool for AI-powered research built with Go and the Charm ecosystem (Bubble Tea + Lipgloss).

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

- **Why Go?**: Fast, single binary, excellent CLI/TUI libraries, cross-platform.
- **Why Bubble Tea?**: Polished TUI, active community, Elm architecture for predictable state.
- **Why Separate Prompts?**: Better git history, easy A/B testing, swappable backends, version control friendly.

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
- Unit tests for business logic, integration tests for CLI commands, UI tests for Bubble Tea components, and performance benchmarks.

## Common Patterns

### Bubble Tea Model
Follows the Elm architecture: `Init()`, `Update(msg)`, `View()`. Long-running operations should run in goroutines and communicate via custom messages. Use `tea.Cmd` for async operations.

## Known Issues & Solutions

- **Bubble Tea blocks during research**: Run research in a goroutine and send progress messages via a channel to update the UI.
- **SQLite concurrent access**: Use WAL mode (`PRAGMA journal_mode=WAL`) for safe concurrent database access.

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

## Key Learnings for Agents

This section distills critical insights and actionable advice for AI agents working on this codebase.

### General Principles
- **Live Feedback**: Always provide visual feedback during long operations (spinners, progress bars). Never let the CLI appear frozen.
- **Database for Learning**: SQLite is used to accumulate research knowledge (query → results → citations → timestamp). Use WAL mode for concurrent access.
- **Prompt Management**: Prompts are separate `.md` files for versioning, A/B testing, and backend swapping.

### Go & Bubble Tea
- **Concurrency**: Run long operations in goroutines, communicate via custom message types.
- **UI Styling**: Use Lipgloss for consistent, beautiful styling.
- **CLI Structure**: Cobra for command structure and flag parsing.
- **Testing**: Use `exec.Command` for integration tests, mock external commands, and Bubble Tea's testing helpers for UI.

### Knowledge Management System
- **Version Control**: Knowledge is stored in `~/.copilot-research/knowledge/` with Git tracking for versioning and history.
- **Structure**: Markdown with YAML frontmatter for human-readability, metadata, and programmatic parsing.
- **Deduplication**: Implement a strategy (e.g., SHA-256 ID + similarity check) to prevent bloat and merge similar entries.
- **Thread Safety**: Use `sync.RWMutex` for cache access to ensure concurrent safety.
- **Filename Safety**: Sanitize topic names for safe use as filenames (e.g., replace `/` with `-`).

### Provider Abstraction System
- **Interface-based Design**: `AIProvider` interface allows multiple AI backends.
- **Context-Awareness**: All `Query()` methods must respect context cancellation and timeouts.
- **Authentication**: Prioritize env vars > config files > CLI tools. Provide clear instructions for authentication.
- **Error Clarity**: Convert API errors into actionable user messages.
- **Configurable Fallback**: `ProviderManager` handles primary/fallback logic, with options for auto-fallback and user notifications.

### Research Engine
- **Orchestration**: Coordinates prompt loading, provider querying, and storage.
- **Progress Events**: Send status updates via a channel for UI feedback.
- **Context Checking**: Regularly check `ctx.Err()` for graceful cancellation.
- **Error Handling**: Distinguish critical from non-critical errors (e.g., log DB save failures, return provider query failures).

### Bubble Tea UI Components
- **Init/Update/View**: Strictly follow this pattern for all components.
- **Centralized Styles**: Organize all Lipgloss styles in one place for consistency and theming.
- **Component Reusability**: Design components to focus on behavior, with styles handled separately.
- **State Machines**: Use state machines (e.g., `stateResearching`, `stateComplete`) to simplify complex UI flows and rendering.
- **Viewport**: Initialize `viewport.New()` on `tea.WindowSizeMsg` for correct sizing.

### CLI Command Implementation
- **Input Priority**: Handle input from args > file > stdin. Use `os.Stdin.Stat()` to check for piped input.
- **Interactive vs. Quiet**: Support both modes, with UI for interactive and silent execution for quiet.
- **Background Tasks**: Run long operations in goroutines, bridging progress channels between engine and UI.
- **Destructive Actions**: Always prompt for confirmation before executing destructive commands.
- **Output**: Format tables for readability (fixed-width, truncation) and support different output formats (e.g., JSON).

### OpenAI Provider Implementation
- **SDK Integration**: Use `github.com/sashabaranov/go-openai`.
- **Lazy Client Init**: Create client only if API key exists.
- **Query Options**: Allow request-level overrides for model, max tokens, etc.
- **Token Usage**: Extract token usage from API responses.
- **Error Classification**: Detect and provide helpful messages for rate limits, timeouts, and authentication errors.

---

*This document will be updated as new patterns and solutions are discovered.*