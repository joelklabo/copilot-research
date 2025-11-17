package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var (
	cfgFile    string
	outputFile string
	quiet      bool
	jsonOutput bool
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
}

func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag
		return
	}

	// Find home directory
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error finding home directory: %v\n", err)
		os.Exit(1)
	}

	// Set default config file location
	cfgFile = filepath.Join(home, ".copilot-research", "config.yaml")
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
