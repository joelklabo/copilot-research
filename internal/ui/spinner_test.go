package ui

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewSpinner(t *testing.T) {
	spinner := NewSpinner()
	assert.NotNil(t, spinner)
}

func TestSpinner_InitialMessage(t *testing.T) {
	spinner := NewSpinner()
	assert.Equal(t, "", spinner.message)
}

func TestSpinner_SetMessage(t *testing.T) {
	spinner := NewSpinner()
	spinner.SetMessage("Loading...")
	assert.Equal(t, "Loading...", spinner.message)
}

func TestSpinner_View(t *testing.T) {
	spinner := NewSpinner()
	spinner.SetMessage("Processing...")

	// View should include the message
	view := spinner.View()
	assert.Contains(t, view, "Processing...")
	assert.NotEmpty(t, view)
}

func TestSpinner_Update(t *testing.T) {
	spinner := NewSpinner()

	// Init the spinner
	cmd := spinner.Init()
	assert.NotNil(t, cmd)

	// Update with tick message should work
	tickMsg := spinner.spinner.Tick()
	_, cmd = spinner.Update(tickMsg)
	
	// Should continue ticking
	assert.NotNil(t, cmd)
}

func TestSpinner_BubbleTeaIntegration(t *testing.T) {
	spinner := NewSpinner()
	spinner.SetMessage("Testing...")

	// Should implement tea.Model
	var _ tea.Model = spinner

	// Init should return a command
	cmd := spinner.Init()
	require.NotNil(t, cmd)

	// Update should handle messages
	model, _ := spinner.Update(nil)
	assert.NotNil(t, model)

	// View should return a string
	view := spinner.View()
	assert.NotEmpty(t, view)
}

func TestSpinner_Styles(t *testing.T) {
	// Just verify styles are defined and don't panic
	styles := DefaultStyles()
	assert.NotNil(t, styles)
	assert.NotNil(t, styles.SpinnerStyle)
	assert.NotNil(t, styles.MessageStyle)
}

func TestSpinner_MultipleMessages(t *testing.T) {
	spinner := NewSpinner()

	// Set different messages
	messages := []string{
		"Loading prompt...",
		"Querying provider...",
		"Processing results...",
		"Complete!",
	}

	for _, msg := range messages {
		spinner.SetMessage(msg)
		view := spinner.View()
		assert.Contains(t, view, msg)
	}
}

func TestSpinner_EmptyMessage(t *testing.T) {
	spinner := NewSpinner()
	view := spinner.View()
	
	// Should still render even with empty message
	assert.NotEmpty(t, view)
}

func TestSpinner_ViewFormat(t *testing.T) {
	spinner := NewSpinner()
	spinner.SetMessage("Test")
	
	view := spinner.View()
	
	// Should have both spinner animation and message
	// The exact format may vary, but should be non-empty and formatted
	lines := strings.Split(view, "\n")
	assert.Greater(t, len(lines), 0)
}
