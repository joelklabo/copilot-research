# Go CLI Development Best Practices

## Project Structure

### Standard Layout
```
cmd/
  myapp/
    main.go          # Entry point
internal/
  app/
    app.go           # Core application logic
  config/
    config.go        # Configuration management
  ui/
    ui.go            # User interface (TUI)
pkg/                 # Public libraries (if any)
  something/
    something.go
```

### Why `internal/`?
- Cannot be imported by external packages
- Enforces encapsulation
- Makes public API explicit (only `pkg/` is public)

## Command Line Parsing with Cobra

### Basic Structure
```go
package cmd

import (
    "github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
    Use:   "copilot-research",
    Short: "AI-powered research tool",
    Long:  `Beautiful CLI for deep research using AI`,
    Run: func(cmd *cobra.Command, args []string) {
        // Main logic
    },
}

func Execute() error {
    return rootCmd.Execute()
}

func init() {
    rootCmd.Flags().StringP("prompt", "p", "default", "Prompt template to use")
    rootCmd.Flags().StringP("output", "o", "pretty", "Output format: pretty, json, md")
    rootCmd.Flags().BoolP("verbose", "v", false, "Verbose output")
}
```

### Subcommands
```go
var listCmd = &cobra.Command{
    Use:   "list",
    Short: "List available prompts",
    Run: func(cmd *cobra.Command, args []string) {
        // List prompts
    },
}

func init() {
    rootCmd.AddCommand(listCmd)
}
```

## Input Handling

### Multiple Input Sources
```go
func getInput(args []string, stdinAvailable bool) (string, error) {
    // 1. Command line arguments
    if len(args) > 0 {
        return strings.Join(args, " "), nil
    }
    
    // 2. File path
    if strings.HasPrefix(args[0], "@") {
        filepath := args[0][1:]
        data, err := os.ReadFile(filepath)
        if err != nil {
            return "", fmt.Errorf("read file: %w", err)
        }
        return string(data), nil
    }
    
    // 3. Stdin (piped input)
    if stdinAvailable {
        data, err := io.ReadAll(os.Stdin)
        if err != nil {
            return "", fmt.Errorf("read stdin: %w", err)
        }
        return string(data), nil
    }
    
    return "", errors.New("no input provided")
}

func isStdinAvailable() bool {
    stat, _ := os.Stdin.Stat()
    return (stat.Mode() & os.ModeCharDevice) == 0
}
```

### Usage Examples
```bash
# Direct argument
copilot-research "What is Swift 6?"

# From file
copilot-research @query.txt

# From pipe
echo "Research topic" | copilot-research

# From heredoc
copilot-research << EOF
Multi-line
research query
EOF
```

## Output Formatting

### Multiple Format Support
```go
type OutputFormat int

const (
    FormatPretty OutputFormat = iota
    FormatJSON
    FormatMarkdown
)

type Result struct {
    Query   string   `json:"query"`
    Answer  string   `json:"answer"`
    Sources []Source `json:"sources"`
}

func (r *Result) Format(format OutputFormat) string {
    switch format {
    case FormatJSON:
        data, _ := json.MarshalIndent(r, "", "  ")
        return string(data)
    case FormatMarkdown:
        return r.ToMarkdown()
    default:
        return r.ToPretty()
    }
}
```

## Configuration Management

### Config File Support
```go
type Config struct {
    DefaultPrompt string `yaml:"default_prompt"`
    DatabasePath  string `yaml:"database_path"`
    CacheExpiry   int    `yaml:"cache_expiry_days"`
}

func LoadConfig() (*Config, error) {
    // Check multiple locations
    paths := []string{
        "./.copilot-research.yaml",
        "~/.config/copilot-research/config.yaml",
        "/etc/copilot-research/config.yaml",
    }
    
    for _, path := range paths {
        if data, err := os.ReadFile(path); err == nil {
            var cfg Config
            if err := yaml.Unmarshal(data, &cfg); err != nil {
                return nil, err
            }
            return &cfg, nil
        }
    }
    
    return DefaultConfig(), nil
}
```

### Environment Variables
```go
func LoadFromEnv(cfg *Config) {
    if prompt := os.Getenv("COPILOT_RESEARCH_PROMPT"); prompt != "" {
        cfg.DefaultPrompt = prompt
    }
    if dbPath := os.Getenv("COPILOT_RESEARCH_DB"); dbPath != "" {
        cfg.DatabasePath = dbPath
    }
}
```

## Error Handling

### User-Friendly Errors
```go
// Bad
return fmt.Errorf("error: %v", err)

// Good
return fmt.Errorf("failed to load prompt '%s': %w", promptName, err)
```

### Exit Codes
```go
const (
    ExitSuccess = 0
    ExitError   = 1
    ExitUsage   = 2
)

func main() {
    if err := cmd.Execute(); err != nil {
        fmt.Fprintf(os.Stderr, "Error: %v\n", err)
        os.Exit(ExitError)
    }
}
```

## Testing

### Table-Driven Tests
```go
func TestParseQuery(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        want    Query
        wantErr bool
    }{
        {
            name:  "simple query",
            input: "test",
            want:  Query{Text: "test"},
        },
        {
            name:  "empty query",
            input: "",
            want:  Query{},
            wantErr: true,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := ParseQuery(tt.input)
            if (err != nil) != tt.wantErr {
                t.Errorf("ParseQuery() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if !reflect.DeepEqual(got, tt.want) {
                t.Errorf("ParseQuery() = %v, want %v", got, tt.want)
            }
        })
    }
}
```

### Integration Tests
```go
func TestCLIIntegration(t *testing.T) {
    // Build test binary
    binary := buildTestBinary(t)
    defer os.Remove(binary)
    
    tests := []struct {
        name       string
        args       []string
        stdin      string
        wantStdout string
        wantStderr string
        wantExit   int
    }{
        {
            name:       "help flag",
            args:       []string{"--help"},
            wantStdout: "Usage:",
            wantExit:   0,
        },
        {
            name:       "version flag",
            args:       []string{"--version"},
            wantStdout: "v1.0.0",
            wantExit:   0,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            cmd := exec.Command(binary, tt.args...)
            if tt.stdin != "" {
                cmd.Stdin = strings.NewReader(tt.stdin)
            }
            
            var stdout, stderr bytes.Buffer
            cmd.Stdout = &stdout
            cmd.Stderr = &stderr
            
            err := cmd.Run()
            exitCode := 0
            if err != nil {
                if exitErr, ok := err.(*exec.ExitError); ok {
                    exitCode = exitErr.ExitCode()
                }
            }
            
            if exitCode != tt.wantExit {
                t.Errorf("exit code = %d, want %d", exitCode, tt.wantExit)
            }
            
            if !strings.Contains(stdout.String(), tt.wantStdout) {
                t.Errorf("stdout = %q, want to contain %q", stdout.String(), tt.wantStdout)
            }
        })
    }
}

func buildTestBinary(t *testing.T) string {
    t.Helper()
    tmpfile, err := os.CreateTemp("", "test-*")
    if err != nil {
        t.Fatal(err)
    }
    tmpfile.Close()
    
    cmd := exec.Command("go", "build", "-o", tmpfile.Name(), ".")
    if err := cmd.Run(); err != nil {
        t.Fatal(err)
    }
    
    return tmpfile.Name()
}
```

### Mocking External Commands
```go
// For testing code that calls `gh copilot`
type Commander interface {
    Run(name string, args ...string) ([]byte, error)
}

type RealCommander struct{}

func (c RealCommander) Run(name string, args ...string) ([]byte, error) {
    return exec.Command(name, args...).CombinedOutput()
}

type MockCommander struct {
    Response []byte
    Err      error
}

func (c MockCommander) Run(name string, args ...string) ([]byte, error) {
    return c.Response, c.Err
}

// Usage in code
type ResearchEngine struct {
    cmd Commander
}

func NewResearchEngine() *ResearchEngine {
    return &ResearchEngine{cmd: RealCommander{}}
}

// In tests
func TestResearch(t *testing.T) {
    engine := &ResearchEngine{
        cmd: MockCommander{
            Response: []byte("test response"),
        },
    }
    // Test with mock
}
```

## Performance

### Benchmarking
```go
func BenchmarkResearch(b *testing.B) {
    engine := NewResearchEngine()
    query := "test query"
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _, err := engine.Research(query)
        if err != nil {
            b.Fatal(err)
        }
    }
}
```

### Profiling
```bash
# CPU profile
go test -cpuprofile=cpu.prof -bench=.
go tool pprof cpu.prof

# Memory profile
go test -memprofile=mem.prof -bench=.
go tool pprof mem.prof
```

## Build and Distribution

### Makefile
```makefile
.PHONY: build test install clean

VERSION := $(shell git describe --tags --always --dirty)
LDFLAGS := -ldflags "-X main.Version=$(VERSION)"

build:
	go build $(LDFLAGS) -o bin/copilot-research

test:
	go test -v ./...

test-coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

install:
	go install $(LDFLAGS)

clean:
	rm -rf bin/
	rm -f coverage.out

lint:
	golangci-lint run

release:
	goreleaser release --clean
```

### Cross-Platform Builds
```bash
# macOS (Apple Silicon)
GOOS=darwin GOARCH=arm64 go build -o dist/copilot-research-darwin-arm64

# macOS (Intel)
GOOS=darwin GOARCH=amd64 go build -o dist/copilot-research-darwin-amd64

# Linux
GOOS=linux GOARCH=amd64 go build -o dist/copilot-research-linux-amd64

# Windows
GOOS=windows GOARCH=amd64 go build -o dist/copilot-research-windows-amd64.exe
```

### Using goreleaser
```yaml
# .goreleaser.yaml
builds:
  - env:
      - CGO_ENABLED=0
    goos:
      - linux
      - darwin
      - windows
    goarch:
      - amd64
      - arm64
    ldflags:
      - -s -w
      - -X main.Version={{.Version}}

archives:
  - format: tar.gz
    name_template: "{{ .ProjectName }}_{{ .Os }}_{{ .Arch }}"

release:
  github:
    owner: your-username
    name: copilot-research
```

## Signal Handling

### Graceful Shutdown
```go
func main() {
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()
    
    // Handle signals
    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
    
    go func() {
        <-sigChan
        fmt.Println("\nShutting down gracefully...")
        cancel()
    }()
    
    // Run app with context
    if err := runApp(ctx); err != nil {
        if err == context.Canceled {
            os.Exit(ExitSuccess)
        }
        fmt.Fprintf(os.Stderr, "Error: %v\n", err)
        os.Exit(ExitError)
    }
}
```

## Logging

### Structured Logging
```go
import "log/slog"

func setupLogging(verbose bool) {
    level := slog.LevelInfo
    if verbose {
        level = slog.LevelDebug
    }
    
    handler := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
        Level: level,
    })
    
    logger := slog.New(handler)
    slog.SetDefault(logger)
}

// Usage
slog.Info("starting research", "query", query)
slog.Debug("cache hit", "key", cacheKey)
slog.Error("failed to connect", "error", err)
```

## Best Practices Summary

✅ **Do:**
- Accept input from args, stdin, and files
- Support multiple output formats (for piping)
- Provide verbose/debug modes
- Handle signals gracefully
- Use table-driven tests
- Write integration tests
- Profile performance
- Version your binary
- Document all flags in help text

❌ **Don't:**
- Assume terminal is interactive (check!)
- Print to stdout if it's for debugging (use stderr)
- Ignore errors (handle or propagate)
- Use global state
- Block indefinitely without timeout
- Forget to close resources
- Skip testing edge cases

## Resources

- [Cobra Documentation](https://github.com/spf13/cobra)
- [Go Testing Package](https://pkg.go.dev/testing)
- [Effective Go](https://go.dev/doc/effective_go)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
