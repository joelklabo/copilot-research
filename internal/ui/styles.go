package ui

import "github.com/charmbracelet/lipgloss"

// Styles holds all the lipgloss styles used in the UI
type Styles struct {
	TitleStyle   lipgloss.Style
	SpinnerStyle lipgloss.Style
	MessageStyle lipgloss.Style
	ResultStyle  lipgloss.Style
	ErrorStyle   lipgloss.Style
	SuccessStyle lipgloss.Style
}

// DefaultStyles returns the default style configuration
func DefaultStyles() Styles {
	return Styles{
		TitleStyle: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("205")).
			MarginBottom(1),

		SpinnerStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("69")),

		MessageStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("241")).
			MarginLeft(2),

		ResultStyle: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("63")).
			Padding(1, 2).
			MarginTop(1),

		ErrorStyle: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("196")),

		SuccessStyle: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("42")),
	}
}
