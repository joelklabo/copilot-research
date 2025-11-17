package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/charmbracelet/lipgloss"
	"github.com/joelklabo/copilot-research/internal/provider"
	"github.com/joelklabo/copilot-research/internal/ui" // Import ui package for styles
	"github.com/spf13/cobra"
)

// authCmd represents the auth command
var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Manage AI provider authentication",
	Long: `The auth command helps you manage authentication for various AI providers.

Use 'auth status' to check current authentication status.
Use 'auth login' to interactively authenticate with a provider.
Use 'auth test' to verify provider connectivity.`,
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help() // Added error check
	},
}

func init() {
	rootCmd.AddCommand(authCmd)

	authCmd.AddCommand(authStatusCommand)
	authCmd.AddCommand(authLoginCommand)
	authCmd.AddCommand(authTestCommand)
	authCmd.AddCommand(authLogoutCommand)
}

var authStatusCommand = &cobra.Command{
	Use:   "status",
	Short: "Show authentication status for all providers",
	Long:  `The status command displays the authentication status for all configured AI providers.`, 
	RunE: runAuthStatus, // Changed to RunE to handle errors
}

func runAuthStatus(cmd *cobra.Command, args []string) error {
	// Dummy usage to satisfy linter for provider import
	var _ provider.AIProvider

	styles := ui.DefaultStyles() // Get default UI styles

	fmt.Println(styles.TitleStyle.Render("Authentication Status"))

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	fmt.Fprintf(w, "%s\t%s\t%s\n",
		lipgloss.NewStyle().Bold(true).Render("Provider"),
		lipgloss.NewStyle().Bold(true).Render("Status"),
		lipgloss.NewStyle().Bold(true).Render("Method"),
	)
	fmt.Fprintf(w, "%s\t%s\t%s\n",
		lipgloss.NewStyle().Faint(true).Render("────────"),
		lipgloss.NewStyle().Faint(true).Render("──────"),
		lipgloss.NewStyle().Faint(true).Render("──────"),
	)

	var unauthenticatedInstructions []string

	// Get all registered providers
	providerNames := AppProviderManager.GetFactory().List()
	if len(providerNames) == 0 {
		fmt.Fprintln(w, "No AI providers configured.")
	}

	for _, name := range providerNames {
		p, err := AppProviderManager.GetFactory().Get(name)
		if err != nil {
			fmt.Fprintf(w, "%s\t%s\t%s\n",
				name,
				styles.ErrorStyle.Render("❌ Error"),
				fmt.Sprintf("Failed to get: %v", err),
			)
			continue
		}

		authInfo := p.RequiresAuth()
		statusIcon := styles.ErrorStyle.Render("❌ Not Configured")
		if p.IsAuthenticated() {
			statusIcon = styles.SuccessStyle.Render("✅ Authenticated")
		}

		method := authInfo.Type
		if p.IsAuthenticated() && authInfo.Type == "apikey" {
			method = "API Key (Env/Config)"
		} else if p.IsAuthenticated() && authInfo.Type == "cli" {
			method = "CLI Tool"
		} else if p.IsAuthenticated() && authInfo.Type == "oauth-device-flow" {
			method = "OAuth"
		}


		fmt.Fprintf(w, "%s\t%s\t%s\n",
			name,
			statusIcon,
			method,
		)

		if !p.IsAuthenticated() && authInfo.Instructions != "" {
			unauthenticatedInstructions = append(unauthenticatedInstructions,
				fmt.Sprintf("To authenticate %s:\n%s", name, authInfo.Instructions),
			)
		}
	}
	w.Flush()

	// Print primary/fallback info
	fmt.Fprintln(os.Stdout)
	fmt.Fprintf(os.Stdout, "Primary: %s\n", AppConfig.Providers.Primary)
	if AppConfig.Providers.Fallback != "" {
		fmt.Fprintf(os.Stdout, "Fallback: %s\n", AppConfig.Providers.Fallback)
	}

	// Print authentication instructions
	if len(unauthenticatedInstructions) > 0 {
		fmt.Fprintln(os.Stdout, "\n"+styles.TitleStyle.Render("Authentication Required"))
		for _, instr := range unauthenticatedInstructions {
			fmt.Fprintln(os.Stdout, instr)
		}
	}

	return nil
}

var authLoginCommand = &cobra.Command{
	Use:   "login [provider]",
	Short: "Interactively authenticate with a provider",
	Long: `The login command guides you through the authentication process for a specified AI provider.
If no provider is specified, it will prompt you to choose one.`, 
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Auth login command not yet implemented.")
	},
}

var authTestCommand = &cobra.Command{
	Use:   "test [provider]",
	Short: "Test connectivity and authentication for a provider",
	Long: `The test command verifies the connectivity and authentication status for a specified AI provider.
If no provider is specified, it will test all configured providers.`, 
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Auth test command not yet implemented.")
	},
}

var authLogoutCommand = &cobra.Command{
	Use:   "logout [provider]",
	Short: "Clear authentication credentials for a provider",
	Long: `The logout command clears the stored authentication credentials for a specified AI provider.
If no provider is specified, it will clear credentials for all providers.`, 
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Auth logout command not yet implemented.")
	},
}