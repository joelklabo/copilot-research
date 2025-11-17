package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"time" // Added for provider timeouts

	"github.com/joelklabo/copilot-research/internal/config" // Added
	"github.com/joelklabo/copilot-research/internal/provider" // Added
	"github.com/spf13/cobra"
)

var (
	cfgFile    string
	outputFile string
	quiet      bool
	jsonOutput bool
	mode       string
	promptName string
	noStore    bool

	AppConfig *config.Config // Added global config
	AppProviderManager *provider.ProviderManager // Added global provider manager
)

// rootCmd represents the base command
var rootCmd = &cobra.Command{
	Use:   "copilot-research",
	Short: "Beautiful CLI for AI-powered research",
	Long: `Copilot Research is a command-line tool that helps you conduct 
research using AI, with knowledge management and beautiful terminal UI.

Example usage:
  copilot-research "How do Swift actors work?"
  copilot-research knowledge list
  copilot-research history`,
	Version: "1.0.0",
}

// Execute runs the root command
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	// Global flags
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.copilot-research/config.yaml)")
	rootCmd.PersistentFlags().StringVarP(&outputFile, "output", "o", "", "output file path")
	rootCmd.PersistentFlags().BoolVarP(&quiet, "quiet", "q", false, "quiet mode (no UI, just output)")
	rootCmd.PersistentFlags().BoolVar(&jsonOutput, "json", false, "output as JSON")
	rootCmd.PersistentFlags().StringVarP(&mode, "mode", "m", "quick", "research mode (quick|deep|compare|synthesis)")
	rootCmd.PersistentFlags().StringVarP(&promptName, "prompt", "p", "default", "prompt template to use")
	rootCmd.PersistentFlags().BoolVar(&noStore, "no-store", false, "don't save to database")
}

func initConfig() {
	// Determine config file path
	if cfgFile == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error finding home directory: %v\n", err)
			os.Exit(1)
		}
		cfgFile = filepath.Join(home, ".copilot-research", "config.yaml")
	}

	// Load config
	var err error
	AppConfig, err = config.LoadConfig(cfgFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
		os.Exit(1)
	}

	// Initialize ProviderManager
	factory := provider.NewProviderFactory()

	// Register GitHub Copilot provider
	ghConfig := AppConfig.Providers.GitHubCopilot
	if ghConfig.Enabled {
		ghProvider := provider.NewGitHubCopilotProvider(ghConfig.Timeout)
		if err := factory.Register("github-copilot", ghProvider); err != nil {
			fmt.Fprintf(os.Stderr, "Error registering GitHub Copilot provider: %v\n", err)
			os.Exit(1)
		}
	}

	// Register OpenAI provider
	openaiConfig := AppConfig.Providers.OpenAI
	if openaiConfig.Enabled {
		// Corrected call to NewOpenAIProvider
		openaiProvider := provider.NewOpenAIProvider(
			openaiConfig.Model,
			openaiConfig.Timeout,
		)
		if err := factory.Register("openai", openaiProvider); err != nil {
			fmt.Fprintf(os.Stderr, "Error registering OpenAI provider: %v\n", err)
			os.Exit(1)
		}
	}

	// Register Anthropic provider
	anthropicConfig := AppConfig.Providers.Anthropic
	if anthropicConfig.Enabled {
		// NewAnthropicProvider does not exist yet, this will cause a compile error
		// I will implement this next.
		anthropicProvider := provider.NewAnthropicProvider(
			anthropicConfig.Model,
			anthropicConfig.Timeout,
			anthropicConfig.APIKeyEnv,
		)
		if err := factory.Register("anthropic", anthropicProvider); err != nil {
			fmt.Fprintf(os.Stderr, "Error registering Anthropic provider: %v\n", err)
			os.Exit(1)
		}
	}

	AppProviderManager = provider.NewProviderManager(
		factory,
		AppConfig.Providers.Primary,
		AppConfig.Providers.Fallback,
		AppConfig.Providers.AutoFallback,
		AppConfig.Providers.NotifyFallback,
	)
}

// GetKnowledgeDir returns the knowledge base directory
func GetKnowledgeDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error finding home directory: %v\n", err)
		os.Exit(1)
	}
	return filepath.Join(home, ".copilot-research", "knowledge")
}
