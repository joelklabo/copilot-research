# Copilot Research - Development Guide

This document provides guidelines and instructions for developers who want to contribute to `copilot-research`.

## Table of Contents
- [Getting Started](#getting-started)
- [Project Structure](#project-structure)
- [Development Workflow](#development-workflow)
  - [Test-Driven Development (TDD)](#test-driven-development-tdd)
  - [Commit Messages](#commit-messages)
- [Building the Project](#building-the-project)
- [Running Tests](#running-tests)
- [Code Formatting and Linting](#code-formatting-and-linting)
- [Debugging](#debugging)
- [Adding New Features](#adding-new-features)
- [Troubleshooting](#troubleshooting)

---

## Getting Started

To set up your development environment:

1.  **Clone the repository:**
    ```bash
    git clone https://github.com/joelklabo/copilot-research
    cd copilot-research
    ```

2.  **Install Go:** Ensure you have Go 1.21+ installed. You can download it from [go.dev/dl](https://go.dev/dl/).

3.  **Install Dependencies:**
    ```bash
    go mod download
    ```

4.  **Install `golangci-lint`:** This tool is used for linting and static analysis.
    ```bash
    go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
    ```

## Project Structure

The project follows a standard Go project layout:

-   `cmd/`: Contains the main packages for CLI commands. Each subcommand typically has its own file (e.g., `cmd/research.go`, `cmd/auth.go`).
-   `internal/`: Contains private application code that should not be imported by external projects.
    -   `internal/config/`: Application configuration management.
    -   `internal/db/`: Database models and SQLite implementation.
    -   `internal/knowledge/`: Knowledge base management system.
    -   `internal/prompts/`: Prompt loading and templating.
    -   `internal/provider/`: AI provider abstraction and implementations (GitHub Copilot, OpenAI, Anthropic).
    -   `internal/research/`: Core research engine logic.
    -   `internal/ui/`: Bubble Tea UI components and styling.
-   `prompts/`: Default prompt templates (Markdown files with YAML frontmatter).
-   `docs/`: Project documentation (including this file).
-   `main.go`: The main entry point of the application.
-   `Makefile`: Contains common development tasks.

## Development Workflow

We follow a strict Test-Driven Development (TDD) approach and maintain clear commit messages.

### Test-Driven Development (TDD)

1.  **Understand the Problem**: Fully grasp the feature or bug you're addressing.
2.  **Write a Failing Test**: Create a test that clearly demonstrates the bug or the absence of the feature. Run it to confirm it fails.
3.  **Implement the Solution**: Write the minimum amount of code necessary to make the failing test pass.
4.  **Run Tests**: Execute all relevant tests to ensure everything passes and no regressions are introduced.
5.  **Refactor**: Improve the code's design, readability, and performance while ensuring tests still pass.
6.  **Commit**: Create a detailed commit message.
7.  **Push**: Push your changes to the remote repository.
8.  **Update Plan**: Mark the task as complete in `docs/plan.md` with the commit hash.

### Commit Messages

Follow the Conventional Commits specification. Each commit message should be structured as follows:

```
<type>(<scope>): <short description>

[optional body]

[optional footer(s)]
```

**Examples:**
-   `feat(cli): Add auth status command`
-   `fix(db): Handle concurrent SQLite access`
-   `docs(readme): Update usage section`

## Building the Project

To build the `copilot-research` binary:

```bash
make build
# or
go build -o copilot-research
```

## Running Tests

To run all tests in the project:

```bash
make test
# or
go test ./... -v -cover
```

To run tests for a specific package (e.g., `internal/provider`):

```bash
go test ./internal/provider/... -v
```

## Code Formatting and Linting

Maintain consistent code style and catch potential issues early:

```bash
# Format code
make fmt
# or
gofmt -s -w .

# Run linter
make lint
# or
golangci-lint run
```

## Debugging

-   Use `fmt.Println()` for quick debugging output.
-   Utilize your IDE's debugger (e.g., VS Code's Go extension).
-   For Bubble Tea UI components, `tea.WithAltScreen()` and `tea.WithLogFile()` can be helpful.

## Adding New Features

When adding new features, especially new AI providers:

-   **Follow existing patterns**: Mimic the structure and style of existing code (e.g., `internal/provider/github_copilot.go`).
-   **Implement interfaces**: Ensure your new components correctly implement the defined interfaces (e.g., `provider.AIProvider`).
-   **Write tests**: Always accompany new code with comprehensive unit and integration tests.
-   **Update documentation**: Ensure `README.md`, `docs/USAGE.md`, and `docs/PROMPTS.md` are updated as necessary. If it's a new provider, consider updating `docs/provider-implementation-guide.md`.

## Troubleshooting

-   **"imported and not used" errors**: If you encounter these during `go test`, it might be due to a linter treating warnings as errors. Ensure all imports are genuinely used. Sometimes, adding a dummy variable like `var _ time.Duration` can satisfy the linter for indirectly used imports.
-   **Build failures**: Check for missing dependencies (`go mod tidy`), syntax errors, or type mismatches.
-   **Test failures**: Analyze the test output carefully. Use verbose mode (`-v`) for more details.
-   **GitHub Actions failures**: Check the workflow logs on GitHub for specific error messages.
