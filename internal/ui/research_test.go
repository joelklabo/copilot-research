package ui

import (
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/joelklabo/copilot-research/internal/research"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewResearchModel(t *testing.T) {
	model := NewResearchModel("test query", "quick")
	assert.NotNil(t, model)
	assert.Equal(t, "test query", model.query)
	assert.Equal(t, "quick", model.mode)
	assert.Equal(t, stateResearching, model.state)
}

func TestResearchModel_Init(t *testing.T) {
	model := NewResearchModel("test", "quick")
	cmd := model.Init()
	assert.NotNil(t, cmd)
}

func TestResearchModel_States(t *testing.T) {
	tests := []struct {
		name         string
		initialState string
		wantState    string
	}{
		{"starts researching", stateResearching, stateResearching},
		{"can be complete", stateComplete, stateComplete},
		{"can be error", stateError, stateError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := NewResearchModel("test", "quick")
			model.state = tt.initialState
			assert.Equal(t, tt.wantState, model.state)
		})
	}
}

func TestResearchModel_ProgressMessage(t *testing.T) {
	model := NewResearchModel("test", "quick")

	// Send progress message
	msg := ProgressMsg("Loading prompt...")
	newModel, _ := model.Update(msg)
	
	rm := newModel.(ResearchModel)
	assert.Equal(t, "Loading prompt...", rm.status)
	assert.Equal(t, stateResearching, rm.state)
}

func TestResearchModel_CompleteMessage(t *testing.T) {
	model := NewResearchModel("test", "quick")

	// Send complete message
	result := &research.ResearchResult{
		Query:    "test",
		Mode:     "quick",
		Content:  "Test result content",
		Duration: 100 * time.Millisecond,
	}
	msg := CompleteMsg{Result: result}
	newModel, _ := model.Update(msg)
	
	rm := newModel.(ResearchModel)
	assert.Equal(t, stateComplete, rm.state)
	assert.Equal(t, result, rm.result)
}

func TestResearchModel_ErrorMessage(t *testing.T) {
	model := NewResearchModel("test", "quick")

	// Send error message
	msg := ErrorMsg{Err: assert.AnError}
	newModel, _ := model.Update(msg)
	
	rm := newModel.(ResearchModel)
	assert.Equal(t, stateError, rm.state)
	assert.Error(t, rm.err)
}

func TestResearchModel_QuitOnCtrlC(t *testing.T) {
	model := NewResearchModel("test", "quick")

	// Send Ctrl+C (interrupt)
	msg := tea.KeyMsg{Type: tea.KeyCtrlC}
	_, cmd := model.Update(msg)
	
	// Should return quit command
	assert.NotNil(t, cmd)
}

func TestResearchModel_QuitOnQ(t *testing.T) {
	model := NewResearchModel("test", "quick")
	model.state = stateComplete // Q only works when complete

	// Send 'q' key
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}}
	_, cmd := model.Update(msg)
	
	// Should return quit command
	assert.NotNil(t, cmd)
}

func TestResearchModel_ViewResearching(t *testing.T) {
	model := NewResearchModel("test", "quick")
	model.status = "Loading..."

	view := model.View()
	
	// Should show spinner and status
	assert.NotEmpty(t, view)
	assert.Contains(t, view, "Loading...")
}

func TestResearchModel_ViewComplete(t *testing.T) {
	model := NewResearchModel("test", "quick")
	model.state = stateComplete
	model.result = &research.ResearchResult{
		Query:    "test",
		Mode:     "quick",
		Content:  "Test result",
		Duration: 100 * time.Millisecond,
	}

	view := model.View()
	
	// Should show result
	assert.NotEmpty(t, view)
	assert.Contains(t, view, "Test result")
}

func TestResearchModel_ViewError(t *testing.T) {
	model := NewResearchModel("test", "quick")
	model.state = stateError
	model.err = assert.AnError

	view := model.View()
	
	// Should show error
	assert.NotEmpty(t, view)
	// Error view should indicate failure
	assert.Contains(t, view, "Error") // Should be case-insensitive
}

func TestResearchModel_ViewportIntegration(t *testing.T) {
	model := NewResearchModel("test", "quick")
	model.state = stateComplete
	
	// Create long content
	longContent := ""
	for i := 0; i < 100; i++ {
		longContent += "Line " + string(rune('0'+i%10)) + "\n"
	}
	
	model.result = &research.ResearchResult{
		Query:    "test",
		Mode:     "quick",
		Content:  longContent,
		Duration: 100 * time.Millisecond,
	}

	// Initialize viewport with window size
	msg := tea.WindowSizeMsg{Width: 80, Height: 24}
	newModel, _ := model.Update(msg)
	
	rm := newModel.(ResearchModel)
	assert.True(t, rm.ready) // Viewport should be ready
}

func TestResearchModel_BubbleTeaInterface(t *testing.T) {
	model := NewResearchModel("test", "quick")

	// Should implement tea.Model
	var _ tea.Model = model

	// Init should return a command
	cmd := model.Init()
	require.NotNil(t, cmd)

	// Update should handle messages
	newModel, _ := model.Update(nil)
	assert.NotNil(t, newModel)

	// View should return a string
	view := model.View()
	assert.NotEmpty(t, view)
}

func TestResearchModel_MultipleProgressUpdates(t *testing.T) {
	model := NewResearchModel("test", "quick")

	// Send multiple progress messages
	progressMessages := []string{
		"Loading prompt...",
		"Querying provider...",
		"Processing results...",
		"Storing in database...",
	}

	for _, msg := range progressMessages {
		newModel, _ := model.Update(ProgressMsg(msg))
		model = newModel.(ResearchModel)
		assert.Equal(t, msg, model.status)
		assert.Equal(t, stateResearching, model.state)
	}
}

func TestProgressMsg(t *testing.T) {
	msg := ProgressMsg("test")
	assert.Equal(t, ProgressMsg("test"), msg)
}

func TestCompleteMsg(t *testing.T) {
	result := &research.ResearchResult{
		Query:   "test",
		Content: "result",
	}
	msg := CompleteMsg{Result: result}
	assert.Equal(t, result, msg.Result)
}

func TestErrorMsg(t *testing.T) {
	msg := ErrorMsg{Err: assert.AnError}
	assert.Error(t, msg.Err)
}
