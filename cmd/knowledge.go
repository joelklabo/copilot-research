package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/joelklabo/copilot-research/internal/knowledge"
	"github.com/spf13/cobra"
)

var (
	excludePattern string
	excludeReason  string
)

// Styles
var (
	titleStyle   = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("205"))

	headerStyle  = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("86"))

	infoStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("240"))

	successStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("42"))
)

// knowledgeCmd represents the knowledge command
var knowledgeCmd = &cobra.Command{
	Use:   "knowledge",
	Short: "Manage knowledge base",
	Long:  `Commands for managing the research knowledge base with Git versioning.`,
}

// listCmd lists all knowledge topics
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all knowledge topics",
	Long:  `Display all topics in the knowledge base with their metadata.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		km, err := knowledge.NewKnowledgeManager(GetKnowledgeDir())
		if err != nil {
			return fmt.Errorf("failed to initialize knowledge manager: %w", err)
		}

		entries, err := km.List()
		if err != nil {
			return fmt.Errorf("failed to list knowledge: %w", err)
		}

		if len(entries) == 0 {
			fmt.Println("No knowledge entries found.")
			fmt.Println("\nAdd your first entry with:")
			fmt.Println("  copilot-research knowledge add <topic>")
			return nil
		}

		// Print header
		fmt.Println(titleStyle.Render(fmt.Sprintf("Knowledge Base (%d topics)", len(entries))))
		fmt.Println(strings.Repeat("━", 80))

		// Create table writer
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\n",
			headerStyle.Render("Topic"),
			headerStyle.Render("Version"),
			headerStyle.Render("Confidence"),
			headerStyle.Render("Updated"))

		for _, entry := range entries {
			// Format confidence as percentage
			confidence := fmt.Sprintf("%.0f%%", entry.Confidence*100)
			
			// Format time ago
			timeAgo := formatTimeAgo(entry.UpdatedAt)
			
			fmt.Fprintf(w, "%s\t%d\t%s\t%s\n",
				entry.Topic,
				entry.Version,
				confidence,
				infoStyle.Render(timeAgo))
		}

		w.Flush()
		return nil
	},
}

// showCmd displays a specific knowledge entry
var showCmd = &cobra.Command{
	Use:   "show <topic>",
	Short: "Display a knowledge entry",
	Long:  `Show the full content of a specific knowledge entry.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		topic := args[0]

		km, err := knowledge.NewKnowledgeManager(GetKnowledgeDir())
		if err != nil {
			return fmt.Errorf("failed to initialize knowledge manager: %w", err)
		}

		entry, err := km.Get(topic)
		if err != nil {
			return fmt.Errorf("knowledge not found: %s", topic)
		}

		// Print header
		fmt.Println(titleStyle.Render(fmt.Sprintf("%s (v%d)", entry.Topic, entry.Version)))
		fmt.Println(strings.Repeat("━", 80))
		fmt.Println()

		// Print metadata
		fmt.Printf("%s %s\n", headerStyle.Render("Confidence:"), fmt.Sprintf("%.0f%%", entry.Confidence*100))
		if len(entry.Tags) > 0 {
			fmt.Printf("%s %s\n", headerStyle.Render("Tags:"), strings.Join(entry.Tags, ", "))
		}
		if entry.Source != "" {
			fmt.Printf("%s %s\n", headerStyle.Render("Source:"), entry.Source)
		}
		fmt.Printf("%s %s\n", headerStyle.Render("Updated:"), formatTimeAgo(entry.UpdatedAt))
		fmt.Println()

		// Print content
		fmt.Println(entry.Content)
		return nil
	},
}

// addCmd adds new knowledge
var addCmd = &cobra.Command{
	Use:   "add <topic>",
	Short: "Add new knowledge entry",
	Long:  `Create a new knowledge entry by opening your $EDITOR.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		topic := args[0]

		km, err := knowledge.NewKnowledgeManager(GetKnowledgeDir())
		if err != nil {
			return fmt.Errorf("failed to initialize knowledge manager: %w", err)
		}

		// Check if already exists
		if _, err := km.Get(topic); err == nil {
			return fmt.Errorf("topic already exists: %s (use 'edit' to modify)", topic)
		}

		// Create template content
		template := fmt.Sprintf(`# %s

Write your knowledge content here in Markdown format.

## Key Points

- Point 1
- Point 2

## Examples

Add examples here...
`, topic)

		// Open editor
		content, err := openEditor(template)
		if err != nil {
			return fmt.Errorf("failed to open editor: %w", err)
		}

		if strings.TrimSpace(content) == strings.TrimSpace(template) {
			return fmt.Errorf("no changes made, aborting")
		}

		// Create knowledge entry
		k := &knowledge.Knowledge{
			Topic:      topic,
			Content:    content,
			Source:     "manual",
			Confidence: 0.8,
			Tags:       []string{},
		}

		if err := km.Add(k); err != nil {
			return fmt.Errorf("failed to add knowledge: %w", err)
		}

		fmt.Println(successStyle.Render("✓") + " Added: " + topic)
		return nil
	},
}

// editCmd edits existing knowledge
var editCmd = &cobra.Command{
	Use:   "edit <topic>",
	Short: "Edit existing knowledge entry",
	Long:  `Edit an existing knowledge entry in your $EDITOR.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		topic := args[0]

		km, err := knowledge.NewKnowledgeManager(GetKnowledgeDir())
		if err != nil {
			return fmt.Errorf("failed to initialize knowledge manager: %w", err)
		}

		// Get existing entry
		entry, err := km.Get(topic)
		if err != nil {
			return fmt.Errorf("knowledge not found: %s", topic)
		}

		// Open editor with existing content
		content, err := openEditor(entry.Content)
		if err != nil {
			return fmt.Errorf("failed to open editor: %w", err)
		}

		if content == entry.Content {
			return fmt.Errorf("no changes made, aborting")
		}

		// Update entry
		entry.Content = content
		if err := km.Update(topic, entry); err != nil {
			return fmt.Errorf("failed to update knowledge: %w", err)
		}

		fmt.Println(successStyle.Render("✓") + " Updated: " + topic)
		return nil
	},
}

// searchCmd searches knowledge
var searchCmd = &cobra.Command{
	Use:   "search <query>",
	Short: "Search knowledge base",
	Long:  `Search for knowledge entries by topic, content, or tags.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		query := args[0]

		km, err := knowledge.NewKnowledgeManager(GetKnowledgeDir())
		if err != nil {
			return fmt.Errorf("failed to initialize knowledge manager: %w", err)
		}

		results, err := km.Search(query)
		if err != nil {
			return fmt.Errorf("search failed: %w", err)
		}

		if len(results) == 0 {
			fmt.Printf("No results found for: %s\n", query)
			return nil
		}

		// Print results
		fmt.Println(titleStyle.Render(fmt.Sprintf("Search Results (%d found)", len(results))))
		fmt.Println(strings.Repeat("━", 80))

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
		fmt.Fprintf(w, "%s\t%s\t%s\n",
			headerStyle.Render("Topic"),
			headerStyle.Render("Confidence"),
			headerStyle.Render("Tags"))

		for _, entry := range results {
			confidence := fmt.Sprintf("%.0f%%", entry.Confidence*100)
			tags := strings.Join(entry.Tags, ", ")
			if tags == "" {
				tags = infoStyle.Render("(none)")
			}

			fmt.Fprintf(w, "%s\t%s\t%s\n",
				entry.Topic,
				confidence,
				tags)
		}

		w.Flush()
		return nil
	},
}

// historyCmd shows git history
var historyCmd = &cobra.Command{
	Use:   "history <topic>",
	Short: "Show git history for a topic",
	Long:  `Display the commit history for a knowledge entry.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		topic := args[0]

		km, err := knowledge.NewKnowledgeManager(GetKnowledgeDir())
		if err != nil {
			return fmt.Errorf("failed to initialize knowledge manager: %w", err)
		}

		commits, err := km.History(topic)
		if err != nil {
			return fmt.Errorf("failed to get history: %w", err)
		}

		if len(commits) == 0 {
			fmt.Printf("No history found for: %s\n", topic)
			return nil
		}

		// Print history
		fmt.Println(titleStyle.Render(fmt.Sprintf("History: %s", topic)))
		fmt.Println(strings.Repeat("━", 80))

		for _, commit := range commits {
			fmt.Printf("%s %s\n",
				headerStyle.Render(commit.Hash[:8]),
				commit.Message)
			fmt.Printf("  %s by %s\n",
				infoStyle.Render(formatTimeAgo(commit.Date)),
				infoStyle.Render(commit.Author))
			fmt.Println()
		}

		return nil
	},
}

// consolidateCmd runs consolidation
var consolidateCmd = &cobra.Command{
	Use:   "consolidate",
	Short: "Consolidate knowledge entries",
	Long:  `Run consolidation to merge and optimize knowledge entries.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		km, err := knowledge.NewKnowledgeManager(GetKnowledgeDir())
		if err != nil {
			return fmt.Errorf("failed to initialize knowledge manager: %w", err)
		}

		fmt.Println("Running consolidation...")
		if err := km.Consolidate(); err != nil {
			return fmt.Errorf("consolidation failed: %w", err)
		}

		fmt.Println(successStyle.Render("✓") + " Consolidation complete")
		return nil
	},
}

// statsCmd shows knowledge stats
var statsCmd = &cobra.Command{
	Use:   "stats",
	Short: "Show knowledge base statistics",
	Long:  `Display statistics about the knowledge base.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		km, err := knowledge.NewKnowledgeManager(GetKnowledgeDir())
		if err != nil {
			return fmt.Errorf("failed to initialize knowledge manager: %w", err)
		}

		entries, err := km.List()
		if err != nil {
			return fmt.Errorf("failed to list knowledge: %w", err)
		}

		// Calculate stats
		totalTopics := len(entries)
		totalTags := make(map[string]int)
		avgConfidence := 0.0

		for _, entry := range entries {
			avgConfidence += entry.Confidence
			for _, tag := range entry.Tags {
				totalTags[tag]++
			}
		}

		if totalTopics > 0 {
			avgConfidence /= float64(totalTopics)
		}

		// Print stats
		fmt.Println(titleStyle.Render("Knowledge Base Statistics"))
		fmt.Println(strings.Repeat("━", 80))
		fmt.Println()

		fmt.Printf("%s %d\n", headerStyle.Render("Total Topics:"), totalTopics)
		fmt.Printf("%s %.0f%%\n", headerStyle.Render("Average Confidence:"), avgConfidence*100)
		fmt.Printf("%s %d\n", headerStyle.Render("Unique Tags:"), len(totalTags))

		if len(totalTags) > 0 {
			fmt.Println()
			fmt.Println(headerStyle.Render("Top Tags:"))
			// Simple display of tags
			for tag, count := range totalTags {
				fmt.Printf("  %s: %d\n", tag, count)
			}
		}

		return nil
	},
}

// rulesCmd manages rules
var rulesCmd = &cobra.Command{
	Use:   "rules",
	Short: "Manage knowledge rules",
	Long:  `Commands for managing user preferences and content filtering rules.`,
}

// rulesListCmd lists all rules
var rulesListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all rules",
	Long:  `Display all configured rules.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		km, err := knowledge.NewKnowledgeManager(GetKnowledgeDir())
		if err != nil {
			return fmt.Errorf("failed to initialize knowledge manager: %w", err)
		}

		re, err := knowledge.NewRuleEngine(km)
		if err != nil {
			return fmt.Errorf("failed to initialize rule engine: %w", err)
		}

		rules := re.ListRules()

		if len(rules) == 0 {
			fmt.Println("No rules configured.")
			fmt.Println("\nAdd a rule with:")
			fmt.Println("  copilot-research knowledge rules add --exclude <pattern>")
			return nil
		}

		// Print rules
		fmt.Println(titleStyle.Render(fmt.Sprintf("Rules (%d configured)", len(rules))))
		fmt.Println(strings.Repeat("━", 80))

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\n",
			headerStyle.Render("Type"),
			headerStyle.Render("Pattern"),
			headerStyle.Render("Reason"),
			headerStyle.Render("ID"))

		for _, rule := range rules {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\n",
				rule.Type,
				rule.Pattern,
				rule.Reason,
				infoStyle.Render(rule.ID[:8]))
		}

		w.Flush()
		return nil
	},
}

// rulesAddCmd adds a rule
var rulesAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new rule",
	Long:  `Add a new content filtering rule.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if excludePattern == "" {
			return fmt.Errorf("--exclude pattern is required")
		}

		km, err := knowledge.NewKnowledgeManager(GetKnowledgeDir())
		if err != nil {
			return fmt.Errorf("failed to initialize knowledge manager: %w", err)
		}

		re, err := knowledge.NewRuleEngine(km)
		if err != nil {
			return fmt.Errorf("failed to initialize rule engine: %w", err)
		}

		rule := knowledge.Rule{
			Type:    "exclude",
			Pattern: excludePattern,
			Reason:  excludeReason,
		}

		if err := re.AddRule(rule); err != nil {
			return fmt.Errorf("failed to add rule: %w", err)
		}

		fmt.Println(successStyle.Render("✓") + " Rule added")
		return nil
	},
}

// rulesRemoveCmd removes a rule
var rulesRemoveCmd = &cobra.Command{
	Use:   "remove <id>",
	Short: "Remove a rule",
	Long:  `Remove a rule by its ID.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ruleID := args[0]

		km, err := knowledge.NewKnowledgeManager(GetKnowledgeDir())
		if err != nil {
			return fmt.Errorf("failed to initialize knowledge manager: %w", err)
		}

		re, err := knowledge.NewRuleEngine(km)
		if err != nil {
			return fmt.Errorf("failed to initialize rule engine: %w", err)
		}

		// Find rule by prefix match
		rules := re.ListRules()
		var matchID string
		for _, rule := range rules {
			if strings.HasPrefix(rule.ID, ruleID) {
				matchID = rule.ID
				break
			}
		}

		if matchID == "" {
			return fmt.Errorf("rule not found: %s", ruleID)
		}

		if err := re.RemoveRule(matchID); err != nil {
			return fmt.Errorf("failed to remove rule: %w", err)
		}

		fmt.Println(successStyle.Render("✓") + " Rule removed")
		return nil
	},
}

// Helper functions

func formatTimeAgo(t time.Time) string {
	duration := time.Since(t)

	switch {
	case duration < time.Minute:
		return "just now"
	case duration < time.Hour:
		mins := int(duration.Minutes())
		if mins == 1 {
			return "1 minute ago"
		}
		return fmt.Sprintf("%d minutes ago", mins)
	case duration < 24*time.Hour:
		hours := int(duration.Hours())
		if hours == 1 {
			return "1 hour ago"
		}
		return fmt.Sprintf("%d hours ago", hours)
	case duration < 30*24*time.Hour:
		days := int(duration.Hours() / 24)
		if days == 1 {
			return "1 day ago"
		}
		return fmt.Sprintf("%d days ago", days)
	default:
		return t.Format("2006-01-02")
	}
}

func openEditor(initialContent string) (string, error) {
	// Get editor from environment
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "nano" // fallback
	}

	// Create temp file
	tmpfile, err := os.CreateTemp("", "knowledge-*.md")
	if err != nil {
		return "", err
	}
	defer os.Remove(tmpfile.Name())

	// Write initial content
	if _, err := tmpfile.WriteString(initialContent); err != nil {
		return "", err
	}
	tmpfile.Close()

	// Open editor
	cmd := exec.Command(editor, tmpfile.Name())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return "", err
	}

	// Read edited content
	content, err := os.ReadFile(tmpfile.Name())
	if err != nil {
		return "", err
	}

	return string(content), nil
}

func init() {
	// Apply MarginBottom in init function
	titleStyle = titleStyle.MarginBottom(1)

	rootCmd.AddCommand(knowledgeCmd)

	// Add subcommands
	knowledgeCmd.AddCommand(listCmd)
	knowledgeCmd.AddCommand(showCmd)
	knowledgeCmd.AddCommand(addCmd)
	knowledgeCmd.AddCommand(editCmd)
	knowledgeCmd.AddCommand(searchCmd)
	knowledgeCmd.AddCommand(historyCmd)
	knowledgeCmd.AddCommand(consolidateCmd)
	knowledgeCmd.AddCommand(statsCmd)
	knowledgeCmd.AddCommand(rulesCmd)

	// Rules subcommands
	rulesCmd.AddCommand(rulesListCmd)
	rulesCmd.AddCommand(rulesAddCmd)
	rulesCmd.AddCommand(rulesRemoveCmd)

	// Flags for rules add
	rulesAddCmd.Flags().StringVar(&excludePattern, "exclude", "", "Pattern to exclude")
	rulesAddCmd.Flags().StringVar(&excludeReason, "reason", "", "Reason for the rule")
}
