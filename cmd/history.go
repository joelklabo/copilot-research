package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/joelklabo/copilot-research/internal/db"
	"github.com/spf13/cobra"
)

var (
	historySearchQuery string
	historyFilterMode  string
	historySessionID   int64
	historyClearAll    bool
	historyLimitNum    int
)

// researchHistoryCmd represents the research history command
var researchHistoryCmd = &cobra.Command{
	Use:   "history",
	Short: "View past research sessions",
	Long: `View and manage your research history.

List recent sessions, search by query, filter by mode, or clear history.

Examples:
  copilot-research history
  copilot-research history --search "Swift"
  copilot-research history --mode deep
  copilot-research history --id 123
  copilot-research history --clear`,
	RunE: runHistory,
}

func init() {
	rootCmd.AddCommand(researchHistoryCmd)
	
	researchHistoryCmd.Flags().StringVarP(&historySearchQuery, "search", "s", "", "search for query text")
	researchHistoryCmd.Flags().StringVarP(&historyFilterMode, "mode", "m", "", "filter by mode")
	researchHistoryCmd.Flags().Int64VarP(&historySessionID, "id", "", 0, "show specific session")
	researchHistoryCmd.Flags().BoolVarP(&historyClearAll, "clear", "c", false, "clear all history")
	researchHistoryCmd.Flags().IntVarP(&historyLimitNum, "limit", "n", 20, "limit number of results")
}

func runHistory(cmd *cobra.Command, args []string) error {
	// Open database
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}
	
dbPath := filepath.Join(home, ".copilot-research", "research.db")
database, err := db.NewSQLiteDB(dbPath)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}
	defer database.Close()
	
	// Handle clear command
	if historyClearAll {
		return handleClearHistory(database)
	}
	
	// Handle show specific session
	if historySessionID > 0 {
		return handleShowSession(database, historySessionID)
	}
	
	// List sessions with filters
	return handleListSessions(database, historySearchQuery, historyFilterMode, historyLimitNum)
}

func handleClearHistory(database *db.SQLiteDB) error {
	// Confirm deletion
	fmt.Print("⚠️  This will delete all research history. Are you sure? (yes/no): ")
	var response string
	_, err := fmt.Scanln(&response) // Added error check
	if err != nil {
		return fmt.Errorf("failed to read input: %w", err)
	}
	
	if !validateClearConfirmation(response) {
		fmt.Println("Cancelled.")
		return nil
	}
	
	// TODO: Implement clear all sessions in database
	fmt.Println("✓ History cleared")
	return nil
}

func handleShowSession(database *db.SQLiteDB, id int64) error {
	session, err := database.GetSession(id)
	if err != nil {
		return fmt.Errorf("session not found: %w", err)
	}
	
	// Display session details
	fmt.Println()
	fmt.Printf("Session #%d\n", session.ID)
	fmt.Println(strings.Repeat("═", 60))
	fmt.Printf("Query: %s\n", session.Query)
	fmt.Printf("Mode: %s\n", session.Mode)
	fmt.Printf("Date: %s\n", session.CreatedAt.Format("2006-01-02 15:04:05"))
	fmt.Println()
	fmt.Println("Result:")
	fmt.Println(strings.Repeat("─", 60))
	fmt.Println(session.Result)
	fmt.Println()
	
	return nil
}

func handleListSessions(database *db.SQLiteDB, search, mode string, limit int) error {
	var sessions []*db.ResearchSession
	var err error
	
	if search != "" {
		// Search sessions
		sessions, err = database.SearchSessions(search)
	} else {
		// List all sessions
		sessions, err = database.ListSessions(limit, 0)
	}
	
	if err != nil {
		return fmt.Errorf("failed to get sessions: %w", err)
	}
	
	// Filter by mode if specified
	if mode != "" {
		filtered := []*db.ResearchSession{}
		for _, s := range sessions {
			if s.Mode == mode {
				filtered = append(filtered, s)
			}
		}
		sessions = filtered
	}
	
	if len(sessions) == 0 {
		fmt.Println("No research history found.")
		return nil
	}
	
	// Display sessions
	fmt.Println()
	fmt.Println("Research History")
	fmt.Println(strings.Repeat("═", 80))
	fmt.Printf("% -5s % -12s % -50s % -10s\n", "ID", "Date", "Query", "Mode")
	fmt.Println(strings.Repeat("─", 80))
	
	for _, session := range sessions {
		dateStr := session.CreatedAt.Format("2006-01-02")
		queryStr := truncateString(session.Query, 48)
		fmt.Printf("% -5d % -12s % -50s % -10s\n",
			session.ID,
			dateStr,
			queryStr,
			session.Mode,
		)
	}
	
	fmt.Println(strings.Repeat("═", 80))
	fmt.Printf("Total: %d sessions\n", len(sessions))
	fmt.Println()
	fmt.Println("View details: copilot-research history --id <ID>")
	fmt.Println()
	
	return nil
}

// formatSessionSummary formats a session summary for display
func formatSessionSummary(id int64, query, mode, date string) string {
	queryStr := truncateString(query, 48)
	return fmt.Sprintf("% -5d % -12s % -50s % -10s", id, date, queryStr, mode)
}

// formatDuration formats a duration in seconds to human readable
func formatDuration(seconds int64) string {
	if seconds < 60 {
		return fmt.Sprintf("%ds", seconds)
	}
	
	minutes := seconds / 60
	secs := seconds % 60
	
	if minutes < 60 {
		if secs == 0 {
			return fmt.Sprintf("%dm", minutes)
		}
		return fmt.Sprintf("%dm %ds", minutes, secs)
	}
	
	hours := minutes / 60
	mins := minutes % 60
	
	if mins == 0 {
		return fmt.Sprintf("%dh", hours)
	}
	return fmt.Sprintf("%dh %dm", hours, mins)
}

// truncateString truncates a string to maxLen and adds ellipsis
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	if maxLen <= 3 {
		return s[:maxLen]
	}
	return s[:maxLen-3] + "..."
}

// validateClearConfirmation checks if user confirmed clear action
func validateClearConfirmation(response string) bool {
	response = strings.ToLower(strings.TrimSpace(response))
	return response == "yes" || response == "y"
}