package textInput

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	focusedStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	noStyle       = lipgloss.NewStyle()
	focusedButton = focusedStyle.Render("[ Submit ]")
	blurredButton = lipgloss.NewStyle().Foreground(lipgloss.Color("")).Render("Submit")
)

type model struct {
	focusIndex int
	inputs     []textinput.Model
	cursorMode cursor.Mode
	submitted  bool
	values     map[string]string
}

func initialModel() model {
	m := model{
		inputs: make([]textinput.Model, 3),
		values: make(map[string]string),
	}

	var t textinput.Model
	for i := range m.inputs {
		t = textinput.New()
		t.Cursor.Style = focusedStyle

		switch i {
		case 0:
			t.Placeholder = "Old database url"
			t.Focus()
			t.PromptStyle = focusedStyle
			t.TextStyle = focusedStyle
		case 1:
			t.Placeholder = "New database url"
		case 2:
			t.Placeholder = "Backup Destination"
		}
		m.inputs[i] = t
	}

	return m
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit
		// Set focus to next input
		case "tab", "shift+tab", "enter", "up", "down":
			s := msg.String()

			if s == "enter" && m.focusIndex == len(m.inputs) {
				m.submitted = true
				for _, input := range m.inputs {
					m.values[input.Placeholder] = input.Value()
				}
				return m, tea.Quit
			}

			if s == "up" || s == "shift+tab" {
				m.focusIndex--
			} else {
				m.focusIndex++
			}

			if m.focusIndex > len(m.inputs) {
				m.focusIndex = 0
			} else if m.focusIndex < 0 {
				m.focusIndex = len(m.inputs)
			}

			cmds := make([]tea.Cmd, len(m.inputs))
			for i := 0; i <= len(m.inputs)-1; i++ {
				if i == m.focusIndex {
					cmds[i] = m.inputs[i].Focus()
					m.inputs[i].PromptStyle = focusedStyle
					m.inputs[i].TextStyle = focusedStyle
				} else {
					m.inputs[i].Blur()
					m.inputs[i].PromptStyle = noStyle
					m.inputs[i].TextStyle = noStyle
				}
			}

			return m, tea.Batch(cmds...)
		}
	}
	cmd := m.updateInputs(msg)
	return m, cmd
}

func (m *model) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))

	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func (m model) View() string {
	var b strings.Builder

	for i := range m.inputs {
		b.WriteString(m.inputs[i].View())
		if i < len(m.inputs)-1 {
			b.WriteRune('\n')
		}
	}

	button := blurredButton
	if m.focusIndex == len(m.inputs) {
		button = focusedButton
	}
	b.WriteString("\n\n" + button)
	return b.String()
}

func GetInputs() (map[string]string, error) {
	p := tea.NewProgram(initialModel())
	m, err := p.Run()
	if err != nil {
		return nil, err
	}

	if m, ok := m.(model); ok && m.submitted {
		return m.values, nil
	}
	return nil, fmt.Errorf("input submission canceled")
}
