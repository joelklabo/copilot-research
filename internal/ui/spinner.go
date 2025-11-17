package ui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

// SpinnerModel is a Bubble Tea model for showing a loading spinner
type SpinnerModel struct {
	spinner spinner.Model
	message string
	styles  Styles
}

// NewSpinner creates a new spinner model
func NewSpinner() *SpinnerModel {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = DefaultStyles().SpinnerStyle

	return &SpinnerModel{
		spinner: s,
		message: "",
		styles:  DefaultStyles(),
	}
}

// SetMessage sets the message to display next to the spinner
func (m *SpinnerModel) SetMessage(msg string) {
	m.message = msg
}

// Init initializes the spinner
func (m *SpinnerModel) Init() tea.Cmd {
	return m.spinner.Tick
}

// Update handles messages
func (m *SpinnerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.spinner, cmd = m.spinner.Update(msg)
	return m, cmd
}

// View renders the spinner
func (m *SpinnerModel) View() string {
	if m.message == "" {
		return m.spinner.View()
	}
	
	return fmt.Sprintf("%s %s", 
		m.spinner.View(),
		m.styles.MessageStyle.Render(m.message),
	)
}
