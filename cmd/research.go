package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/joelklabo/copilot-research/internal/db"
	"github.com/joelklabo/copilot-research/internal/prompts"
	"github.com/joelklabo/copilot-research/internal/provider"
	"github.com/joelklabo/copilot-research/internal/research"
	"github.com/joelklabo/copilot-research/internal/ui"
	"github.com/spf13/cobra"
)

var (
	inputFile string
)

// researchCmd represents the research command
var researchCmd = &cobra.Command{
	Use:   "research [query]",
	Short: "Conduct AI-powered research",
	Long: `Conduct AI-powered research using GitHub Copilot and other AI providers.

The query can be provided as:
  - Command argument: copilot-research "How do Swift actors work?"
  - Input file: copilot-research --input query.txt
  - Standard input: echo "query" | copilot-research

Examples:
  copilot-research "How do Swift actors work?"
  copilot-research "Compare React and Vue" --mode compare
  copilot-research --input query.txt --output report.md
  echo "Explain Swift concurrency" | copilot-research --quiet`,
	RunE: runResearch,
}

func init() {
	RootCmd.AddCommand(researchCmd)
	
	// Command-specific flags
	researchCmd.Flags().StringVarP(&inputFile, "input", "i", "", "input file containing query")
}

func runResearch(cmd *cobra.Command, args []string) error {
	// Get query from args, file, or stdin
	query, err := determineQuery(args, inputFile)
	if err != nil {
		return fmt.Errorf("failed to get query: %w", err)
	}
	
	if query == "" {
		return fmt.Errorf("no query provided")
	}
	
	// Validate mode
	if err := validateMode(Mode); err != nil {
		return err
	}
	
	// Initialize database
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}
	
dbPath := filepath.Join(home, ".copilot-research", "research.db")
	if err := os.MkdirAll(filepath.Dir(dbPath), 0755); err != nil {
		return fmt.Errorf("failed to create database directory: %w", err)
	}
	
database, err := db.NewSQLiteDB(dbPath)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}
	defer database.Close()
	
	// Initialize prompt loader
	promptsDir := filepath.Join("prompts")
	loader := prompts.NewPromptLoader(promptsDir)
	
	// Initialize provider
	factory := provider.NewProviderFactory()
	ghProvider := provider.NewGitHubCopilotProvider(60 * time.Second)
	if err := factory.Register("github-copilot", ghProvider); err != nil {
		return fmt.Errorf("failed to register provider: %w", err)
	}
	
	// Updated call to NewProviderManager
	// Use AppConfig.Providers.AutoFallback and AppConfig.Providers.NotifyFallback
	providerMgr := provider.NewProviderManager(
		factory,
		AppConfig.Providers.Primary,
		AppConfig.Providers.Fallback,
		AppConfig.Providers.AutoFallback,
		AppConfig.Providers.NotifyFallback,
	)
	
	// Check authentication
	// This check should be done by the providerMgr, not a specific provider
	// For now, keep it for ghProvider as it's the only one registered here
	if !ghProvider.IsAuthenticated() {
		authInfo := ghProvider.RequiresAuth()
		return fmt.Errorf("authentication required:\n\n%s", authInfo.Instructions)
	}
	
	// Initialize research engine
	engine := research.NewEngine(database, loader, providerMgr)
	
	// Run research
	if Quiet {
		return runQuietResearch(engine, query)
	}
	
	return runInteractiveResearch(engine, query)
}

func runQuietResearch(engine *research.Engine, query string) error {
	ctx := context.Background()
	progress := make(chan string, 10)
	
	// Drain progress channel
	go func() {
		for range progress {
		}
	}()
	
	opts := research.ResearchOptions{
		Query:      query,
		Mode:       Mode,
		PromptName: PromptName,
		NoStore:    NoStore,
	}
	
	result, err := engine.Research(ctx, opts, progress)
	close(progress)
	
	if err != nil {
		return fmt.Errorf("research failed: %w", err)
	}
	
	// Format output
	format := "markdown"
	if JSONOutput {
		format = "json"
	}
	
	output := formatOutput(result.Content, format)
	
	// Write output
	if err := writeOutput(OutputFile, output); err != nil {
		return fmt.Errorf("failed to write output: %w", err)
	}
	
	return nil
}

func runInteractiveResearch(engine *research.Engine, query string) error {
	// Create UI model
	model := ui.NewResearchModel(query, Mode)
	
	// Create Bubble Tea program
	p := tea.NewProgram(model)
	
	// Start research in background
	go func() {
		ctx := context.Background()
		progress := make(chan string, 10)
		
		// Send progress updates to UI
		go func() {
			for msg := range progress {
				p.Send(ui.ProgressMsg(msg))
			}
		}()
		
		opts := research.ResearchOptions{
			Query:      query,
			Mode:       Mode,
			PromptName: PromptName,
			NoStore:    NoStore,
		}
		
		result, err := engine.Research(ctx, opts, progress)
		close(progress)
		
		if err != nil {
			p.Send(ui.ErrorMsg{Err: err})
			return
		}
		
		p.Send(ui.CompleteMsg{Result: result})
	}()
	
	// Run UI
	if _, err := p.Run(); err != nil {
		return fmt.Errorf("UI error: %w", err)
	}
	
	return nil
}

func determineQuery(args []string, inputFile string) (string, error) {
	// Priority: args > input file > stdin
	if len(args) > 0 {
		return getQueryFromArgs(args)
	}
	
	if inputFile != "" {
		return getQueryFromFile(inputFile)
	}
	
	// Check if stdin has data
	stat, err := os.Stdin.Stat()
	if err == nil && (stat.Mode()&os.ModeCharDevice) == 0 {
		return getQueryFromStdin()
	}
	
	return "", fmt.Errorf("no query provided")
}

func getQueryFromArgs(args []string) (string, error) {
	if len(args) == 0 {
		return "", fmt.Errorf("no arguments provided")
	}
	return strings.Join(args, " "), nil
}

func getQueryFromFile(filename string) (string, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}
	return strings.TrimSpace(string(data)), nil
}

func getQueryFromStdin() (string, error) {
	data, err := io.ReadAll(os.Stdin)
	if err != nil {
		return "", fmt.Errorf("failed to read stdin: %w", err)
	}
	return strings.TrimSpace(string(data)), nil
}

func formatOutput(content string, format string) string {
	switch format {
	case "json":
		output := map[string]interface{}{
			"content": content,
			"format":  "markdown",
		}
		data, _ := json.MarshalIndent(output, "", "  ")
		return string(data)
	default:
		return content
	}
}

func writeOutput(filename string, content string) error {
	if filename == "" {
		// Write to stdout
		fmt.Println(content)
		return nil
	}
	
	// Write to file
	if err := os.WriteFile(filename, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}
	
	return nil
}

func validateMode(mode string) error {
	if mode == "" {
		return nil // Will default to "quick"
	}
	
	validModes := map[string]bool{
		"quick":     true,
		"deep":      true,
		"compare":   true,
		"synthesis": true,
	}
	
	if !validModes[mode] {
		return fmt.Errorf("invalid mode: %s (valid modes: quick, deep, compare, synthesis)", mode)
	}
	
	return nil
}
