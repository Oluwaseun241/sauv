package main

import (
	"fmt"
	"os"
	"sauv/cmd"
	"strings"

	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Styles
var (
	docStyle            = lipgloss.NewStyle().Margin(1, 2)
	focusedStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	blurredStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	cursorStyle         = focusedStyle
	noStyle             = lipgloss.NewStyle()
	helpStyle           = blurredStyle
	cursorModeHelpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))

	focusedButton = focusedStyle.Render("[ Submit ]")
	blurredButton = fmt.Sprintf("[ %s ]", blurredStyle.Render("Submit"))
)

type item struct {
	title, desc string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

type model struct {
	//dbConfigPath   string
	list           list.Model
	inputs         []textinput.Model
	oldDBUrl       string
	newDBUrl       string
	focusIndex     int
	isMigrating    bool
	migrationError error
	cursorMode     cursor.Mode

	selectedOption string
	showInputs     bool
}

func initialModel() model {
	items := []list.Item{
		item{title: "Backup", desc: "Perform a database backup"},
		item{title: "Backup and Migrate", desc: "Backup the database and migrate to a new database"},
	}

	l := list.New(items, list.NewDefaultDelegate(), 0, 0)
	l.Title = "Sauv! Choose an Option"

	m := model{
		list:   l,
		inputs: make([]textinput.Model, 2),
	}

	var t textinput.Model
	for i := range m.inputs {
		t = textinput.New()
		t.Cursor.Style = cursorStyle
		t.CharLimit = 32

		switch i {
		case 0:
			t.Placeholder = "Enter old database URL"
			t.Focus()
			t.PromptStyle = focusedStyle
			t.TextStyle = focusedStyle
		case 1:
			t.Placeholder = "Enter new database URL"
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

	// Handle key events
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "enter":

			if m.showInputs {
				if m.focusIndex == len(m.inputs) {
					oldDBUrl := m.inputs[0].Value()
					newDBUrl := m.inputs[1].Value()
					m.isMigrating = true

					// Call your backup/migration logic here
					if m.selectedOption == "Backup and Migrate" {
						err := cmd.PerformBackupAndMigration(oldDBUrl, newDBUrl)
						if err != nil {
							m.migrationError = err
						}
					} else {
						fmt.Println("\nStarting Backup...")
						// Call backup function here
					}

					m.isMigrating = false
					return m, tea.Quit
				}

				// Move focus between inputs
				if msg.String() == "tab" || msg.String() == "down" {
					m.focusIndex++
					if m.focusIndex > len(m.inputs) {
						m.focusIndex = 0
					}
				} else if msg.String() == "up" {
					m.focusIndex--
					if m.focusIndex < 0 {
						m.focusIndex = len(m.inputs)
					}
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
			} else {
				i, ok := m.list.SelectedItem().(item)
				if ok {
					m.selectedOption = i.title
					m.showInputs = true

					// Reset focus index and refocus the first input field
					m.focusIndex = 0
					m.inputs[0].Focus()
					m.inputs[0].PromptStyle = focusedStyle
					m.inputs[0].TextStyle = focusedStyle
				}
			}
		}

	// Handle window resize
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}

	// Update the list or inputs based on current view
	if m.showInputs {
		cmd := m.updateInputs(msg)
		return m, cmd
	} else {
		var cmd tea.Cmd
		m.list, cmd = m.list.Update(msg)
		return m, cmd
	}
}

func (m *model) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))

	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func (m model) View() string {
	if m.showInputs {
		var b strings.Builder

		if !m.isMigrating {
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
			fmt.Fprintf(&b, "\n\n%s\n\n", button)

		} else {
			if m.migrationError != nil {
				b.WriteString(fmt.Sprintf("Migration failed: %v\n", m.migrationError))
			} else {
				b.WriteString("Migration successful!\n")
			}
		}
		return b.String()
	} else {
		return docStyle.Render(m.list.View())
	}
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
