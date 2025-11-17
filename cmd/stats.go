package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/tabwriter"

	"github.com/joelklabo/copilot-research/internal/db"
	"github.com/joelklabo/copilot-research/internal/ui"
	"github.com/spf13/cobra"
)

// statsCmd represents the stats command
var statsCmd = &cobra.Command{
	Use:   "stats",
	Short: "Show research statistics",
	Long: `Display statistics about your research usage, patterns, and database size.

Examples:
  copilot-research stats`,
	RunE: func(cmd *cobra.Command, args []string) error {
		home, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("failed to get home directory: %w", err)
		}
		dbPath := filepath.Join(home, ".copilot-research", "research.db")
		database, err := db.NewSQLiteDB(dbPath)
		if err != nil {
			return fmt.Errorf("failed to open database: %w", err)
		}
		return _runStats(database, dbPath)
	},
}

func _runStats(database db.DB, dbPath string) error {
	styles := ui.DefaultStyles()

	defer database.Close()

	// Get total sessions
	totalSessions, err := database.GetTotalSessions()
	if err != nil {
		return fmt.Errorf("failed to get total sessions: %w", err)
	}

	// Get mode stats
	modeStats, err := database.GetModeStats()
	if err != nil {
		return fmt.Errorf("failed to get mode statistics: %w", err)
	}

	// Get database size (simple approximation)
	dbFileInfo, err := os.Stat(dbPath)
	dbSize := "N/A"
	if err == nil {
		dbSize = formatBytes(dbFileInfo.Size())
	}

	// Print stats
	fmt.Println(styles.TitleStyle.Render("Research Statistics"))
	fmt.Println(strings.Repeat("━", 80))
	fmt.Println()

	fmt.Printf("%s %d\n", styles.HeaderStyle.Render("Total Sessions:"), totalSessions)
	fmt.Printf("%s %s\n", styles.HeaderStyle.Render("Database Size:"), dbSize)

	if len(modeStats) > 0 {
		fmt.Println()
		fmt.Println(styles.HeaderStyle.Render("Mode Usage:"))
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
		for mode, count := range modeStats {
			percentage := 0.0
			if totalSessions > 0 {
				percentage = float64(count) / float64(totalSessions) * 100
			}
			fmt.Fprintf(w, "  %s\t%d (%.0f%%)\n", mode, count, percentage)
		}
		w.Flush()
	}

	// Get top queries
	topQueries, err := database.GetTopQueries(5) // Limit to top 5
	if err != nil {
		return fmt.Errorf("failed to get top queries: %w", err)
	}

	if len(topQueries) > 0 {
		fmt.Println()
		fmt.Println(styles.HeaderStyle.Render("Top Queries:"))
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
		for i, qc := range topQueries {
			fmt.Fprintf(w, "  %d. %s (%d times)\n", i+1, qc.Query, qc.Count)
		}
		w.Flush()
	}

	return nil
}

// formatBytes converts bytes to a human-readable format
func formatBytes(b int64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(b)/float64(div), "KMGTPE"[exp])
}

func init() {
	RootCmd.AddCommand(statsCmd)
	statsCmd.RunE = func(cmd *cobra.Command, args []string) error {
		home, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("failed to get home directory: %w", err)
		}
		dbPath := filepath.Join(home, ".copilot-research", "research.db")
		database, err := db.NewSQLiteDB(dbPath)
		if err != nil {
			return fmt.Errorf("failed to open database: %w", err)
		}
		return _runStats(database, dbPath)
	}
}
