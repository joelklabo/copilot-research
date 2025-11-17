# Copilot Research

> Beautiful CLI tool for AI-powered research with learning capabilities

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

## Overview

Copilot Research is a fast, beautiful command-line tool that makes AI-powered research effortless. Built with the Charm ecosystem (Bubble Tea + Lipgloss), it provides real-time feedback, learns from your queries, and produces high-quality research reports.

## Features

- 🎨 **Beautiful UI** - Polished terminal interface with live progress
- ⚡ **Fast** - Single binary, instant startup
- 🧠 **Smart** - Learns from past research sessions
- 🔄 **Flexible** - Multiple research modes and customizable prompts
- 📝 **Output** - Clean markdown reports
- 🔌 **Scriptable** - Unix-friendly, pipeable output

## Quick Start

```bash
# Basic usage
copilot-research "Swift 6 actors"

# From a file
copilot-research --input research.txt

# Deep dive mode
copilot-research "iOS 26 new APIs" --mode deep

# See history
copilot-research history
```

## Installation

### Homebrew
```bash
brew install joelklabo/tap/copilot-research
```

### From source
```bash
go install github.com/joelklabo/copilot-research@latest
```

### Download binary
Download from [releases](https://github.com/joelklabo/copilot-research/releases)

## Usage

### Basic Research
```bash
copilot-research "What are Swift 6 actors?"
```

### Research Modes
- `--mode quick` - Fast overview (default)
- `--mode deep` - Deep dive with examples
- `--mode compare` - Compare multiple approaches
- `--mode synthesis` - Synthesize from multiple sources

### Input Sources
```bash
# String
copilot-research "topic"

# File
copilot-research --input file.txt

# Stdin
echo "topic" | copilot-research
```

### Output Options
```bash
# Save to file
copilot-research "topic" --output report.md

# JSON format
copilot-research "topic" --json

# Quiet mode (no UI)
copilot-research "topic" --quiet
```

### History & Learning
```bash
# View history
copilot-research history

# Search history
copilot-research history --search "Swift"

# Stats
copilot-research stats

# Clear history
copilot-research history --clear
```

## Configuration

Config file: `~/.copilot-research/config.yaml`

```yaml
# Default mode
default_mode: quick

# Database path
db_path: ~/.copilot-research/research.db

# Prompt directory
prompt_dir: ~/.copilot-research/prompts

# Active prompt
active_prompt: default

# Output preferences
output:
  format: markdown
  color: true
```

## Custom Prompts

Create custom prompts in `~/.copilot-research/prompts/`:

```bash
# List available prompts
copilot-research prompts list

# Use specific prompt
copilot-research "topic" --prompt claude

# Set default prompt
copilot-research config set active_prompt claude
```

## Architecture

```
copilot-research
├── cmd/              # CLI commands
├── internal/
│   ├── research/     # Research engine
│   ├── ui/           # Bubble Tea UI
│   ├── db/           # SQLite storage
│   └── prompts/      # Prompt management
├── prompts/          # Default prompts
└── docs/             # Documentation
```

## Development

```bash
# Clone
git clone https://github.com/joelklabo/copilot-research
cd copilot-research

# Install dependencies
go mod download

# Run tests
go test ./...

# Build
go build -o copilot-research

# Run
./copilot-research "test query"
```

## AI Provider System

Copilot Research uses a plugin-based provider system that supports multiple AI backends. Providers can be easily added without modifying core code.

### Supported Providers

Currently implemented:
- **GitHub Copilot** - Via `gh copilot` CLI (default)
- **OpenAI** - Coming soon
- **Anthropic Claude** - Coming soon

### Provider Configuration

Configure providers in `~/.copilot-research/config.yaml`:

```yaml
providers:
  primary: github-copilot
  fallback: openai
  
  github-copilot:
    enabled: true
    auth_type: cli
    timeout: 60s
    
  openai:
    enabled: true
    auth_type: apikey
    api_key_env: OPENAI_API_KEY
    model: gpt-4
    timeout: 30s
```

### Authentication

Each provider supports multiple authentication methods in priority order:

**GitHub Copilot:**
1. `COPILOT_GITHUB_TOKEN` environment variable
2. `GH_TOKEN` environment variable
3. `gh` CLI authentication (`gh auth login`)

**OpenAI:**
1. `OPENAI_API_KEY` environment variable
2. Configuration file: `copilot-research config set openai.api_key sk-...`

**Anthropic:**
1. `ANTHROPIC_API_KEY` environment variable
2. Configuration file: `copilot-research config set anthropic.api_key sk-ant-...`

### Provider Fallback

The system automatically falls back to secondary providers if the primary fails:

```bash
# Uses github-copilot, falls back to openai if unavailable
copilot-research "topic"

# Force specific provider
copilot-research "topic" --provider openai

# Check authentication status
copilot-research providers status
```

### Implementing New Providers

**For AI Agents:** See [Provider Implementation Guide](docs/provider-implementation-guide.md) for complete instructions on implementing new providers. This guide is designed to be read and followed by AI agents (Claude, Copilot, etc.) without human intervention.

**Quick Summary:**
1. Implement the `AIProvider` interface in `internal/provider/`
2. Handle authentication with clear error messages
3. Respect context cancellation and timeouts
4. Return responses in standardized format
5. Write comprehensive tests
6. Register in factory

Example provider interface:

```go
type AIProvider interface {
    Name() string
    Query(ctx context.Context, prompt string, opts QueryOptions) (*Response, error)
    IsAuthenticated() bool
    RequiresAuth() AuthInfo
    Capabilities() ProviderCapabilities
}
```

See `internal/provider/github_copilot.go` for a complete reference implementation.

## Requirements

- Go 1.21+
- **At least one AI provider:**
  - GitHub CLI (`gh`) with Copilot subscription
  - OpenAI API key
  - Anthropic API key
- SQLite3

## Contributing

Contributions welcome! Please see [CONTRIBUTING.md](CONTRIBUTING.md)

## License

MIT License - see [LICENSE](LICENSE)

## Acknowledgments

Built with:
- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - TUI framework
- [Lipgloss](https://github.com/charmbracelet/lipgloss) - Terminal styling
- [Charm](https://charm.sh/) - Excellent CLI tools

## Roadmap

- [ ] Interactive mode
- [ ] Multi-hop research
- [ ] Team collaboration (shared DB)
- [ ] Web UI for browsing
- [ ] Plugin system
- [ ] Export formats (PDF, HTML)

---

**Made with ❤️ by Joel Klabo**
