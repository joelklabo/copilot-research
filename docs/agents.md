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

---

*Keep updated as you discover patterns and solve problems.*
