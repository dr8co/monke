package repl

import (
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/dr8co/monke/evaluator"
	"github.com/dr8co/monke/lexer"
	"github.com/dr8co/monke/object"
	"github.com/dr8co/monke/parser"
	"strings"
)

const PROMPT = ">> "

func Start(username string) {
	// Start the bubbletea program
	p := tea.NewProgram(initialModel(username))
	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
	}
}

// Styling
var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#7D56F4")).
			Padding(0, 1)

	promptStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#7D56F4")).
			Bold(true)

	resultStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#04B575"))

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF5F87"))

	historyStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#767676"))
)

// The model represents the state of the application
type model struct {
	textInput textinput.Model
	history   []historyEntry
	env       *object.Environment
	username  string
}

// historyEntry represents a single entry in the REPL history
type historyEntry struct {
	input   string
	output  string
	isError bool
}

// initialModel creates a new model with default values
func initialModel(username string) model {
	ti := textinput.New()
	ti.Placeholder = "Enter Monkey code"
	ti.Focus()
	ti.Width = 80
	ti.Prompt = promptStyle.Render(PROMPT)

	return model{
		textInput: ti,
		history:   []historyEntry{},
		env:       object.NewEnvironment(),
		username:  username,
	}
}

// Init is the first function that will be called
func (m model) Init() tea.Cmd {
	return textinput.Blink
}

// Update handles all the updates to our model
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyEnter:
			input := m.textInput.Value()
			if input == "" {
				return m, nil
			}

			// Process the input
			l := lexer.New(input)
			p := parser.New(l)
			program := p.ParseProgram()

			var output string
			var isError bool

			if len(p.Errors()) != 0 {
				isError = true
				output = formatParseErrors(p.Errors())
			} else {
				evaluated := evaluator.Eval(program, m.env)
				if evaluated != nil {
					output = evaluated.Inspect()
				} else {
					output = "nil"
				}
			}

			// Add to history
			m.history = append(m.history, historyEntry{
				input:   input,
				output:  output,
				isError: isError,
			})

			// Clear the input
			m.textInput.SetValue("")
			return m, nil
		}
	}

	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

// View renders the current UI
func (m model) View() string {
	var s strings.Builder

	// Title
	s.WriteString(titleStyle.Render(" Monkey Programming Language REPL "))
	s.WriteString("\n")

	// Welcome message
	if m.username != "" {
		s.WriteString(fmt.Sprintf("\nHello %s! Feel free to type in commands\n", m.username))
	}
	s.WriteString("\n")

	// History
	for _, entry := range m.history {
		s.WriteString(promptStyle.Render(PROMPT))
		s.WriteString(entry.input)
		s.WriteString("\n")

		if entry.isError {
			s.WriteString(errorStyle.Render(entry.output))
		} else {
			s.WriteString(resultStyle.Render(entry.output))
		}
		s.WriteString("\n\n")
	}

	// Input
	s.WriteString(m.textInput.View())
	s.WriteString("\n")

	// Help text
	s.WriteString(historyStyle.Render("\nPress Esc or Ctrl+C to exit"))

	return s.String()
}

// formatParseErrors formats parser errors into a string
func formatParseErrors(errors []string) string {
	var s strings.Builder
	s.WriteString("parser errors:\n")
	for _, msg := range errors {
		s.WriteString("\t" + msg + "\n")
	}
	return s.String()
}
