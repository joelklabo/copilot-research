# Bubble Tea Patterns & Best Practices

## The Elm Architecture in Go

Bubble Tea implements The Elm Architecture (TEA):
- **Model**: Your application state
- **Update**: How state changes in response to messages
- **View**: How to render state as a string

## Essential Patterns

### 1. Async Operations with Messages

```go
// Custom message type
type researchCompleteMsg struct {
    result string
    err    error
}

// Return function that sends message when done
func doResearch(query string) tea.Cmd {
    return func() tea.Msg {
        result, err := performResearch(query)
        return researchCompleteMsg{result: result, err: err}
    }
}

// In Update():
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case researchCompleteMsg:
        if msg.err != nil {
            m.err = msg.err
            return m, tea.Quit
        }
        m.result = msg.result
        m.done = true
        return m, tea.Quit
    }
    return m, nil
}
```

### 2. Spinner During Loading

```go
import "github.com/charmbracelet/bubbles/spinner"

type model struct {
    spinner spinner.Model
    loading bool
}

func initialModel() model {
    s := spinner.New()
    s.Spinner = spinner.Dot
    s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
    return model{
        spinner: s,
        loading: true,
    }
}

func (m model) Init() tea.Cmd {
    return m.spinner.Tick
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case spinner.TickMsg:
        var cmd tea.Cmd
        m.spinner, cmd = m.spinner.Update(msg)
        return m, cmd
    }
    return m, nil
}

func (m model) View() string {
    if m.loading {
        return fmt.Sprintf("\n %s Researching...\n\n", m.spinner.View())
    }
    return m.result
}
```

### 3. Progress Bar

```go
import "github.com/charmbracelet/bubbles/progress"

type model struct {
    progress progress.Model
    percent  float64
}

type progressMsg float64

func (m model) Init() tea.Cmd {
    return tea.Batch(
        m.progress.Init(),
        simulateProgress(),
    )
}

func simulateProgress() tea.Cmd {
    return tea.Tick(time.Millisecond*100, func(t time.Time) tea.Msg {
        return progressMsg(0.05) // increment by 5%
    })
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case progressMsg:
        m.percent += float64(msg)
        if m.percent >= 1.0 {
            m.percent = 1.0
            return m, tea.Quit
        }
        cmd := m.progress.SetPercent(m.percent)
        return m, tea.Batch(cmd, simulateProgress())
    }
    return m, nil
}

func (m model) View() string {
    return "\n" + m.progress.View() + "\n"
}
```

### 4. Batch Commands

```go
// Start multiple things at once
func (m model) Init() tea.Cmd {
    return tea.Batch(
        m.spinner.Tick,
        doResearch(m.query),
        checkCache(m.query),
    )
}
```

### 5. Keyboard Handling

```go
import "github.com/charmbracelet/bubbletea"

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.String() {
        case "q", "ctrl+c":
            return m, tea.Quit
        case "enter":
            return m, submitQuery(m.input)
        case "up":
            m.cursor--
            if m.cursor < 0 {
                m.cursor = 0
            }
        case "down":
            m.cursor++
            if m.cursor >= len(m.items) {
                m.cursor = len(m.items) - 1
            }
        }
    }
    return m, nil
}
```

## Styling with Lipgloss

```go
import "github.com/charmbracelet/lipgloss"

var (
    titleStyle = lipgloss.NewStyle().
        Bold(true).
        Foreground(lipgloss.Color("#FAFAFA")).
        Background(lipgloss.Color("#7D56F4")).
        Padding(0, 1)

    errorStyle = lipgloss.NewStyle().
        Foreground(lipgloss.Color("#FF0000")).
        Bold(true)

    successStyle = lipgloss.NewStyle().
        Foreground(lipgloss.Color("#00FF00"))

    infoStyle = lipgloss.NewStyle().
        Foreground(lipgloss.Color("#00FFFF"))
)

func (m model) View() string {
    if m.err != nil {
        return errorStyle.Render("Error: " + m.err.Error())
    }
    
    title := titleStyle.Render("Research Results")
    body := lipgloss.NewStyle().
        Width(80).
        Border(lipgloss.RoundedBorder()).
        Render(m.result)
    
    return lipgloss.JoinVertical(lipgloss.Left, title, "", body)
}
```

## Testing Bubble Tea Apps

```go
import (
    "testing"
    tea "github.com/charmbracelet/bubbletea"
)

func TestModelUpdate(t *testing.T) {
    m := initialModel()
    
    // Simulate message
    msg := researchCompleteMsg{result: "test result"}
    newModel, cmd := m.Update(msg)
    
    // Assert state changed
    if !newModel.(model).done {
        t.Error("Expected done to be true")
    }
    
    // Assert command returned
    if cmd == nil {
        t.Error("Expected quit command")
    }
}
```

## Performance Tips

1. **Debounce rapid updates**: Don't re-render on every tiny change
2. **Limit view size**: Especially for large result sets
3. **Use viewport for scrolling**: From `bubbles/viewport`
4. **Cache expensive renders**: Store computed strings
5. **Profile with pprof**: Find bottlenecks in Update/View

## Common Pitfalls

❌ **Blocking in Update()**: Never do slow work in Update
✅ **Use tea.Cmd**: Return commands for async work

❌ **Mutating state directly**: Go values are passed by value
✅ **Return new model**: `return m, cmd` with modified m

❌ **Ignoring terminal size**: UI breaks on small terminals
✅ **Handle tea.WindowSizeMsg**: Adapt to terminal dimensions

❌ **Complex View logic**: Slows down rendering
✅ **Keep View simple**: Pre-compute in Update when possible

## Resources

- [Bubble Tea Tutorial](https://github.com/charmbracelet/bubbletea/tree/master/tutorials)
- [Bubbles Components](https://github.com/charmbracelet/bubbles)
- [Lipgloss Styling](https://github.com/charmbracelet/lipgloss)
- [Example Apps](https://github.com/charmbracelet/bubbletea/tree/master/examples)
