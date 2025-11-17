package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"time" // Re-added for provider timeouts

	"github.com/joelklabo/copilot-research/internal/config" // Added
	"github.com/joelklabo/copilot-research/internal/provider" // Added
	"github.com/spf13/cobra"
)

var (
	CfgFile    string
	OutputFile string
	Quiet      bool
	JSONOutput bool
	Mode       string
	PromptName string
	NoStore    bool

	AppConfig *config.Config
	AppProviderManager *provider.ProviderManager
)

// RootCmd represents the base command
var RootCmd = &cobra.Command{
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
	return RootCmd.Execute()
}

func init() {

cobra.OnInitialize(InitConfig)

	// Global flags
	RootCmd.PersistentFlags().StringVar(&CfgFile, "config", "", "config file (default is $HOME/.copilot-research/config.yaml)")
	RootCmd.PersistentFlags().StringVarP(&OutputFile, "output", "o", "", "output file path")
	RootCmd.PersistentFlags().BoolVarP(&Quiet, "quiet", "q", false, "quiet mode (no UI, just output)")
	RootCmd.PersistentFlags().BoolVar(&JSONOutput, "json", false, "output as JSON")
	RootCmd.PersistentFlags().StringVarP(&Mode, "mode", "m", "quick", "research mode (quick|deep|compare|synthesis)")
	RootCmd.PersistentFlags().StringVarP(&PromptName, "prompt", "p", "default", "prompt template to use")
	RootCmd.PersistentFlags().BoolVar(&NoStore, "no-store", false, "don't save to database")
}

// InitConfig initializes the configuration
func InitConfig() {
	// Dummy usage to satisfy linter for time import
	var _ time.Duration

	// Determine config file path
	if CfgFile == "" {
	home, err := os.UserHomeDir()
				if err != nil {
				fmt.Fprintf(os.Stderr, "Error finding home directory: %v\n", err)
				os.Exit(1)
			}
		// Corrected typo: CfgFile instead of cfigFile
		CfgFile = filepath.Join(home, ".copilot-research", "config.yaml")
	}

	// Load config
	var err error
	AppConfig, err = config.LoadConfig(CfgFile)
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