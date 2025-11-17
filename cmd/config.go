package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage application configuration",
	Long: `The config command allows you to view, set, and reset application
configuration settings.

Examples:
  copilot-research config show
  copilot-research config set providers.openai.model gpt-4o
  copilot-research config reset`,
}

func init() {
	RootCmd.AddCommand(configCmd)
	configCmd.AddCommand(ConfigShowCmd)
	configCmd.AddCommand(configSetCmd)
	configCmd.AddCommand(configResetCmd)
}

var ConfigShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show current application configuration",
	Long:  `Displays the current application configuration in YAML format.`,
	Run: func(cmd *cobra.Command, args []string) {
		if AppConfig == nil {
			fmt.Fprintln(os.Stderr, "Error: Configuration not loaded.")
			os.Exit(1)
		}

		data, err := yaml.Marshal(AppConfig)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error marshalling config: %v\n", err)
			os.Exit(1)
		}
		fmt.Println(string(data))
	},
}

var configSetCmd = &cobra.Command{
	Use:   "set <key> <value>",
	Short: "Set a configuration value",
	Long: `Sets a specific configuration value. Nested keys can be specified
using dot notation (e.g., providers.openai.model).

Example:
  copilot-research config set providers.openai.model gpt-4o`,
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		key := args[0]
		value := args[1]

		// TODO: Implement setting nested values and saving
		fmt.Printf("Setting %s to %s (not yet implemented)\n", key, value)
	},
}

var configResetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Reset configuration to default values",
	Long:  `Resets the application configuration to its default settings.`,
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: Implement resetting to default and saving
		fmt.Println("Resetting config to defaults (not yet implemented)")
	},
}
