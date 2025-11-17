# Copilot Research - Detailed Usage Guide

This document provides comprehensive usage instructions and examples for the `copilot-research` CLI tool.

## Table of Contents
- [Basic Research](#basic-research)
- [Research Modes](#research-modes)
- [Input Sources](#input-sources)
- [Output Options](#output-options)
- [Authentication & Providers](#authentication--providers)
  - [Checking Status](#checking-status)
  - [Logging In](#logging-in)
  - [Testing Connectivity](#testing-connectivity)
  - [Logging Out](#logging-out)
- [History & Learning](#history--learning)
- [Knowledge Management](#knowledge-management)
- [Statistics](#statistics)
- [Configuration Management](#configuration-management)

---

## Basic Research

The core functionality of `copilot-research` is to conduct AI-powered research. You can provide your research query directly as an argument.

```bash
copilot-research "What are Swift 6 actors?"
```

## Research Modes

The tool supports different research modes to tailor the AI's response to your needs. You can specify a mode using the `--mode` or `-m` flag.

- `--mode quick` / `-m quick`: Provides a fast overview of the topic (default mode).
  ```bash
  copilot-research "Explain quantum computing" --mode quick
  ```
- `--mode deep` / `-m deep`: Conducts a deep dive into the topic, often including more examples and detailed explanations.
  ```bash
  copilot-research "iOS 26 new APIs" --mode deep
  ```
- `--mode compare` / `-m compare`: Compares multiple approaches or technologies.
  ```bash
  copilot-research "Compare React and Vue" --mode compare
  ```
- `--mode synthesis` / `-m synthesis`: Synthesizes information from multiple sources into a coherent narrative.
  ```bash
  copilot-research "Synthesize recent advancements in AI ethics" --mode synthesis
  ```

## Input Sources

You can provide your research query in several ways:

### As a command-line argument
```bash
copilot-research "Explain the observer pattern in Go"
```

### From a file
Use the `--input` or `-i` flag to specify a file containing your query.
```bash
# query.txt contains: "How does the Go garbage collector work?"
copilot-research --input query.txt
```

### From standard input (stdin)
Pipe your query directly to the tool. This is useful for scripting.
```bash
echo "What is the Elm architecture?" | copilot-research
```

## Output Options

You can control how the research results are outputted.

### Save to a file
Use the `--output` or `-o` flag to save the result to a specified file.
```bash
copilot-research "SwiftUI lifecycle" --output swiftui_lifecycle.md
```

### JSON format
Use the `--json` flag to get the output in JSON format, useful for programmatic consumption.
```bash
copilot-research "Rust ownership model" --json
```

### Quiet mode
Use the `--quiet` or `-q` flag to suppress the interactive UI and only print the final result to stdout. This is ideal for scripting.
```bash
copilot-research "Kubernetes deployments" --quiet
```

## Authentication & Providers

`copilot-research` supports multiple AI providers. The `auth` command helps you manage their authentication status.

### Checking Status
The `auth status` command displays the authentication status for all configured AI providers. It also shows which provider is set as primary and fallback, and provides instructions for unauthenticated providers.

```bash
copilot-research auth status
```

Example Output:
```
Authentication Status
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
Provider        Status          Method
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
github-copilot  ✅ Authenticated  CLI Tool
openai          ❌ Not Configured apikey
anthropic       ❌ Not Configured apikey

Primary: github-copilot
Fallback: openai

Authentication Required

To authenticate openai:
OpenAI API key required.

Get your API key:
  1. Visit https://platform.openai.com/api-keys
  2. Create a new API key
  3. Set it in your environment:
     export OPENAI_API_KEY=sk-...
     
Or add to config:
  copilot-research config set providers.openai.api_key sk-...

Pricing: https://openai.com/pricing
```

### Logging In
The `auth login` command guides you through the interactive authentication process for a specified AI provider. If no provider is specified, it will prompt you to choose one.

```bash
# Interactive login for GitHub Copilot
copilot-research auth login github-copilot

# Interactive login (will prompt for provider choice)
copilot-research auth login
```

### Testing Connectivity
The `auth test` command verifies the connectivity and authentication status for a specified AI provider. If no provider is specified, it will test all configured providers.

```bash
# Test a specific provider
copilot-research auth test openai

# Test all configured providers
copilot-research auth test
```

### Logging Out
The `auth logout` command clears the stored authentication credentials for a specified AI provider. If no provider is specified, it will clear credentials for all providers.

```bash
# Log out from GitHub Copilot
copilot-research auth logout github-copilot

# Log out from all providers
copilot-research auth logout
```

## History & Learning

The tool keeps a history of your research sessions and can learn from them.

### View Research History
The `history` command lists your past research sessions.
```bash
copilot-research history
```

### Search History
You can search your history by query text.
```bash
copilot-research history --search "Swift"
```

### Filter History
Filter your history by research mode.
```bash
copilot-research history --mode deep
```

### Show Specific Session
Display the full details of a specific research session by its ID.
```bash
copilot-research history --id 123
```

### Clear History
Clear all your research history. This action requires confirmation.
```bash
copilot-research history --clear
```

## Knowledge Management

The tool includes a knowledge management system to store and retrieve learned information.

### List Knowledge Topics
```bash
copilot-research knowledge list
```

### Show Specific Knowledge
```bash
copilot-research knowledge show swift-concurrency
```

### Add New Knowledge
```bash
copilot-research knowledge add "new-topic"
```

### Edit Knowledge
This will open the knowledge entry in your default `$EDITOR`.
```bash
copilot-research knowledge edit swift-concurrency
```

### Search Knowledge
```bash
copilot-research knowledge search "actor isolation"
```

### View Knowledge History
Show the Git history for a specific knowledge topic.
```bash
copilot-research knowledge history swift-concurrency
```

### Consolidate Knowledge Base
Run a consolidation pass to deduplicate and clean up your knowledge base.
```bash
copilot-research knowledge consolidate
```

### Manage Knowledge Rules
The rule system allows you to define preferences, exclusions, and content filtering.
```bash
# List all rules
copilot-research knowledge rules list

# Add an exclusion rule
copilot-research knowledge rules add --type exclude --pattern "Model View Controller|MVC" --reason "Using MV architecture instead"

# Remove a rule by ID
copilot-research knowledge rules remove <rule-id>
```

## Statistics

The `stats` command provides analytics about your research usage.

```bash
copilot-research stats
```

Example Output:
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

## Configuration Management

The `config` command allows you to manage application settings directly from the CLI.

### Show Current Configuration
```bash
# Show all configuration settings
copilot-research config

# Get a specific configuration value
copilot-research config get providers.openai.model
```

### Set Configuration Values
```bash
# Set a specific configuration value
copilot-research config set providers.openai.model gpt-4o
```

### Reset Configuration
Reset all configuration settings to their default values. This action requires confirmation.
```bash
copilot-research config reset
```
