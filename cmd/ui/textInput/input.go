package textInput

import (
	"fmt"
	"strings"

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
	inputs    []textinput.Model
	focusIdx  int
	submitted bool
	err       error
	values    map[string]string
	showAll   bool
}

func initialModel(showAll bool) model {
	labels := []string{"Old database URL", "New database URL", "Backup Destination"}
	inputs := make([]textinput.Model, 3)

	for i := 0; i < 3; i++ {
		t := textinput.New()
		t.Placeholder = labels[i]
		t.CharLimit = 150
		t.Width = 50
		t.Prompt = "> "

		if i == 0 {
			t.Focus()
			t.PromptStyle = focusedStyle
			t.TextStyle = focusedStyle
		}

		inputs[i] = t
	}

	return model{
		inputs:    inputs,
		focusIdx:  0,
		submitted: false,
		values:    make(map[string]string),
		showAll:   showAll,
	}
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

			if s == "enter" && m.focusIdx == len(m.inputs) {
				if err := m.validateInputs(); err != nil {
					m.err = err
					return m, nil
				}
				m.submitted = true
				for i, input := range m.inputs {
					switch i {
					case 0:
						m.values["Old database URL"] = input.Value()
					case 1:
						m.values["New database URL"] = input.Value()
					case 2:
						m.values["Backup Destination"] = input.Value()
					}
				}
				return m, tea.Quit
			}

			if s == "up" || s == "shift+tab" {
				m.focusIdx--
			} else {
				m.focusIdx++
			}

			if m.focusIdx > len(m.inputs) {
				m.focusIdx = 0
			} else if m.focusIdx < 0 {
				m.focusIdx = len(m.inputs)
			}

			cmds := make([]tea.Cmd, len(m.inputs))
			for i := 0; i <= len(m.inputs)-1; i++ {
				if i == m.focusIdx {
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
	var cmd tea.Cmd
	// Only update the focused input
	if m.focusIdx < len(m.inputs) {
		m.inputs[m.focusIdx], cmd = m.inputs[m.focusIdx].Update(msg)
	}
	return cmd
}

func (m model) View() string {
	var b strings.Builder

	b.WriteString("Enter database details\n\n")

	maxInputs := 3
	if m.showAll {
		maxInputs = 3
	}

	for i := 0; i < maxInputs; i++ {
		if !m.showAll && i == 1 {
			continue
		}

		b.WriteString(m.inputs[i].View())
		if i < maxInputs-1 {
			b.WriteRune('\n')
		}
	}

	if m.err != nil {
		b.WriteString("\n\n" + focusedStyle.Render(m.err.Error()))
	}

	button := blurredButton
	if m.focusIdx == len(m.inputs) {
		button = focusedButton
	}
	b.WriteString("\n\n" + button)
	return b.String()
}

func GetInputs(showAll bool) (map[string]string, error) {
	p := tea.NewProgram(initialModel(showAll))
	m, err := p.Run()
	if err != nil {
		return nil, err
	}

	if m, ok := m.(model); ok && m.submitted {
		return m.values, nil
	}
	return nil, fmt.Errorf("input submission canceled")
}

func (m model) validateInputs() error {
	oldDB := strings.TrimSpace(m.inputs[0].Value())
	newDB := strings.TrimSpace(m.inputs[1].Value())
	dest := strings.TrimSpace(m.inputs[2].Value())

	if oldDB == "" {
		return fmt.Errorf("old database URL is required")
	}

	if !strings.HasPrefix(oldDB, "postgresql://") {
		return fmt.Errorf("invalid old database URL format")
	}

	if newDB != "" && !strings.HasPrefix(newDB, "postgres://") {
		return fmt.Errorf("invalid new database URL format")
	}

	if dest == "" {
		return fmt.Errorf("backup destination is required")
	}

	return nil
}
