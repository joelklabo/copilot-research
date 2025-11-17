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

## Requirements

- Go 1.21+
- GitHub CLI (`gh`) installed and authenticated
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
