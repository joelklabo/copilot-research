package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/joelklabo/copilot-research/internal/research"
)

// States for the research UI
const (
	stateResearching = "researching"
	stateComplete    = "complete"
	stateError       = "error"
)

// ResearchModel is the main Bubble Tea model for research operations
type ResearchModel struct {
	state    string
	query    string
	mode     string
	
	spinner  *SpinnerModel
	status   string
	result   *research.ResearchResult
	err      error
	
	viewport viewport.Model
	ready    bool
	styles   Styles
}

// ProgressMsg is sent when research progress updates
type ProgressMsg string

// CompleteMsg is sent when research completes
type CompleteMsg struct {
	Result *research.ResearchResult
}

// ErrorMsg is sent when research errors
type ErrorMsg struct {
	Err error
}

// NewResearchModel creates a new research model
func NewResearchModel(query, mode string) ResearchModel {
	spinner := NewSpinner()
	
	return ResearchModel{
		state:   stateResearching,
		query:   query,
		mode:    mode,
		spinner: spinner,
		status:  "",
		styles:  DefaultStyles(),
	}
}

// Init initializes the model
func (m ResearchModel) Init() tea.Cmd {
	return m.spinner.Init()
}

// Update handles messages
func (m ResearchModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return m, tea.Quit
		case tea.KeyRunes:
			// Allow 'q' to quit when complete or errored
			if (m.state == stateComplete || m.state == stateError) && len(msg.Runes) > 0 && msg.Runes[0] == 'q' {
				return m, tea.Quit
			}
		}
		
		// Pass key events to viewport when in complete state
		if m.state == stateComplete && m.ready {
			var cmd tea.Cmd
			m.viewport, cmd = m.viewport.Update(msg)
			return m, cmd
		}

	case tea.WindowSizeMsg:
		// Initialize viewport with window size
		if m.state == stateComplete && !m.ready {
			m.viewport = viewport.New(msg.Width, msg.Height-10) // Leave space for header/footer
			m.viewport.SetContent(m.formatResult())
			m.ready = true
		}
		return m, nil

	case ProgressMsg:
		m.status = string(msg)
		m.spinner.SetMessage(m.status)
		return m, nil

	case CompleteMsg:
		m.state = stateComplete
		m.result = msg.Result
		m.ready = false // Reset viewport ready state
		return m, nil

	case ErrorMsg:
		m.state = stateError
		m.err = msg.Err
		return m, nil
	}

	// Update spinner in researching state
	if m.state == stateResearching {
		var cmd tea.Cmd
		spinnerModel, cmd := m.spinner.Update(msg)
		m.spinner = spinnerModel.(*SpinnerModel)
		return m, cmd
	}

	return m, nil
}

// View renders the model
func (m ResearchModel) View() string {
	switch m.state {
	case stateResearching:
		return m.viewResearching()
	case stateComplete:
		return m.viewComplete()
	case stateError:
		return m.viewError()
	default:
		return ""
	}
}

// viewResearching renders the researching state
func (m ResearchModel) viewResearching() string {
	var b strings.Builder
	
	b.WriteString(m.styles.TitleStyle.Render("🔍 Researching"))
	b.WriteString("\n\n")
	b.WriteString(m.styles.MessageStyle.Render(fmt.Sprintf("Query: %s", m.query)))
	b.WriteString("\n")
	b.WriteString(m.styles.MessageStyle.Render(fmt.Sprintf("Mode: %s", m.mode)))
	b.WriteString("\n\n")
	
	// Show spinner with status
	if m.status != "" {
		m.spinner.SetMessage(m.status)
	}
	b.WriteString(m.spinner.View())
	
	b.WriteString("\n\n")
	b.WriteString("Press Ctrl+C to cancel")
	
	return b.String()
}

// viewComplete renders the complete state
func (m ResearchModel) viewComplete() string {
	var b strings.Builder
	
	b.WriteString(m.styles.SuccessStyle.Render("✓ Research Complete"))
	b.WriteString("\n\n")
	b.WriteString(m.styles.MessageStyle.Render(fmt.Sprintf("Query: %s", m.query)))
	b.WriteString("\n")
	b.WriteString(m.styles.MessageStyle.Render(fmt.Sprintf("Mode: %s | Duration: %v", m.mode, m.result.Duration)))
	b.WriteString("\n\n")
	
	if m.ready {
		b.WriteString(m.viewport.View())
		b.WriteString("\n\n")
		b.WriteString("↑/↓: Scroll • q: Quit")
	} else {
		// Before viewport is ready, show result directly
		b.WriteString(m.styles.ResultStyle.Render(m.result.Content))
		b.WriteString("\n\n")
		b.WriteString("Press q to quit")
	}
	
	return b.String()
}

// viewError renders the error state
func (m ResearchModel) viewError() string {
	var b strings.Builder
	
	b.WriteString(m.styles.ErrorStyle.Render("✗ Error"))
	b.WriteString("\n\n")
	b.WriteString(m.styles.MessageStyle.Render(fmt.Sprintf("Query: %s", m.query)))
	b.WriteString("\n\n")
	b.WriteString(fmt.Sprintf("Error: %v", m.err))
	b.WriteString("\n\n")
	b.WriteString("Press q to quit")
	
	return b.String()
}

// formatResult formats the research result for display
func (m ResearchModel) formatResult() string {
	if m.result == nil {
		return ""
	}
	return m.result.Content
}
