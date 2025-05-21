// Package repl implements the Read-Eval-Print Loop for the Monke programming language.
//
// The REPL provides an interactive interface for users to enter Monke code,
// have it evaluated, and see the results immediately. It uses the Charm libraries
// (Bubbletea, Bubbles, and Lipgloss) to create a modern, user-friendly terminal
// interface with features like syntax highlighting and command history.
//
// Key features:
//   - Interactive command input and execution
//   - Command history tracking
//   - Styled output with different colors for results and errors
//   - Persistent environment across commands
//
// The main entry point is the Start function, which initializes and runs the REPL
// with the given username.
package repl

import (
	"fmt"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/dr8co/monke/evaluator"
	"github.com/dr8co/monke/lexer"
	"github.com/dr8co/monke/object"
	"github.com/dr8co/monke/parser"
	"github.com/dr8co/monke/token"
	"strings"
	"time"
)

const PROMPT = ">> "

// Options contains configuration options for the REPL
type Options struct {
	NoColor bool // Disable syntax highlighting and colored output
	Debug   bool // Enable debug mode with more verbose output
}

// Start initializes and runs the REPL with the given username and options.
// It creates a new bubbletea program with an initial model and runs it.
// The username is displayed in the welcome message of the REPL.
// If an error occurs while running the program, it is printed to the console.
func Start(username string, options Options) {
	// Start the bubbletea program
	p := tea.NewProgram(initialModel(username, options))
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

	// Error styles
	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF5F87"))

	parseErrorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF5F87")).
			Bold(true)

	runtimeErrorStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FF8700")).
				Bold(true)

	errorTipStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFAF00"))

	historyStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#767676"))

	// Syntax highlighting styles
	keywordStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF79C6")).
			Bold(true)

	identifierStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#F8F8F2"))

	literalStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#F1FA8C"))

	operatorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF5555"))

	delimiterStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#BD93F9"))

	stringStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#50FA7B"))
)

// ErrorType represents the type of error that occurred
type ErrorType int

const (
	NoError ErrorType = iota
	ParseError
	RuntimeError
)

// Custom messages for async evaluation
type evalResultMsg struct {
	output    string
	isError   bool
	errorType ErrorType
	elapsed   time.Duration
}

// The model represents the state of the application
type model struct {
	textInput    textinput.Model
	history      []historyEntry
	env          *object.Environment
	username     string
	evaluating   bool
	currentInput string
	spinner      spinner.Model
	options      Options
}

// historyEntry represents a single entry in the REPL history
type historyEntry struct {
	input          string
	output         string
	isError        bool
	errorType      ErrorType
	evaluationTime time.Duration // Time taken to evaluate
}

// initialModel creates a new model with default values
func initialModel(username string, options Options) model {
	ti := textinput.New()
	ti.Placeholder = "Enter Monkey code"
	ti.Focus()
	ti.Width = 80
	ti.Prompt = promptStyle.Render(PROMPT)

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF79C6"))

	return model{
		textInput:  ti,
		history:    []historyEntry{},
		env:        object.NewEnvironment(),
		username:   username,
		evaluating: false,
		spinner:    s,
		options:    options,
	}
}

// Init is the first function that will be called
func (m model) Init() tea.Cmd {
	return tea.Batch(textinput.Blink, spinner.Tick)
}

// evalCmd is a command that evaluates Monkey code asynchronously
func evalCmd(input string, env *object.Environment, debug bool) tea.Cmd {
	return func() tea.Msg {
		start := time.Now()

		// Process the input
		l := lexer.New(input)
		p := parser.New(l)

		// Debug: Print tokenization time
		var tokenizeTime time.Duration
		if debug {
			tokenizeStart := time.Now()
			program := p.ParseProgram()
			tokenizeTime = time.Since(tokenizeStart)

			var output string
			var isError bool
			var errorType = NoError

			if len(p.Errors()) != 0 {
				isError = true
				errorType = ParseError
				output = formatParseErrors(p.Errors())

				if debug {
					fmt.Printf("DEBUG: Parse errors: %v\n", p.Errors())
				}
			} else {
				// Debug: Print evaluation time
				evalStart := time.Now()
				evaluated := evaluator.Eval(program, env)
				evalTime := time.Since(evalStart)

				if debug {
					fmt.Printf("DEBUG: Tokenize time: %v\n", tokenizeTime)
					fmt.Printf("DEBUG: Eval time: %v\n", evalTime)
				}

				if evaluated != nil {
					// Check if the result is an error object
					if evaluated.Type() == object.ERROR_OBJ {
						isError = true
						errorType = RuntimeError
						output = formatRuntimeError(evaluated.Inspect())

						if debug {
							fmt.Printf("DEBUG: Runtime error: %s\n", evaluated.Inspect())
						}
					} else {
						output = evaluated.Inspect()

						if debug {
							fmt.Printf("DEBUG: Result type: %s\n", evaluated.Type())
						}
					}
				} else {
					output = "nil"
				}
			}

			elapsed := time.Since(start)

			if debug {
				fmt.Printf("DEBUG: Total execution time: %v\n", elapsed)
			}

			return evalResultMsg{
				output:    output,
				isError:   isError,
				errorType: errorType,
				elapsed:   elapsed,
			}
		} else {
			// Non-debug path (original code)
			program := p.ParseProgram()

			var output string
			var isError bool
			var errorType = NoError

			if len(p.Errors()) != 0 {
				isError = true
				errorType = ParseError
				output = formatParseErrors(p.Errors())
			} else {
				evaluated := evaluator.Eval(program, env)
				if evaluated != nil {
					// Check if the result is an error object
					if evaluated.Type() == object.ERROR_OBJ {
						isError = true
						errorType = RuntimeError
						output = formatRuntimeError(evaluated.Inspect())
					} else {
						output = evaluated.Inspect()
					}
				} else {
					output = "nil"
				}
			}

			elapsed := time.Since(start)

			return evalResultMsg{
				output:    output,
				isError:   isError,
				errorType: errorType,
				elapsed:   elapsed,
			}
		}
	}
}

// Update handles all the updates to our model
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case spinner.TickMsg:
		if m.evaluating {
			m.spinner, cmd = m.spinner.Update(msg)
			return m, cmd
		}

	case evalResultMsg:
		// Evaluation completed
		m.evaluating = false

		// Add to history
		m.history = append(m.history, historyEntry{
			input:          m.currentInput,
			output:         msg.output,
			isError:        msg.isError,
			errorType:      msg.errorType,
			evaluationTime: msg.elapsed,
		})

		m.currentInput = ""
		return m, nil

	case tea.KeyMsg:
		// If we're evaluating, ignore key presses except for Ctrl+C
		if m.evaluating && msg.Type != tea.KeyCtrlC {
			return m, nil
		}

		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc, tea.KeyCtrlD:
			return m, tea.Quit
		case tea.KeyEnter:
			input := m.textInput.Value()
			if input == "" {
				return m, nil
			}

			// Start evaluation in the background
			m.evaluating = true
			m.currentInput = input
			m.textInput.SetValue("")

			return m, evalCmd(input, m.env, m.options.Debug)
		}
	}

	// Only update the text input if we're not evaluating
	if !m.evaluating {
		m.textInput, cmd = m.textInput.Update(msg)
	}

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
		s.WriteString(m.highlightCode(entry.input))
		s.WriteString("\n")

		if entry.isError {
			// Use different styles based on the error type
			switch entry.errorType {
			case ParseError:
				// Split the output to separate the error message from the tips
				parts := strings.Split(entry.output, "\nTips:")
				if len(parts) > 1 {
					s.WriteString(parseErrorStyle.Render(parts[0]))
					s.WriteString("\n")
					s.WriteString(errorTipStyle.Render("Tips:" + parts[1]))
				} else {
					s.WriteString(parseErrorStyle.Render(entry.output))
				}
			case RuntimeError:
				// Split the output to separate the error message from the tips
				parts := strings.Split(entry.output, "\nTips:")
				if len(parts) > 1 {
					s.WriteString(runtimeErrorStyle.Render(parts[0]))
					s.WriteString("\n")
					s.WriteString(errorTipStyle.Render("Tips:" + parts[1]))
				} else {
					s.WriteString(runtimeErrorStyle.Render(entry.output))
				}
			default:
				s.WriteString(errorStyle.Render(entry.output))
			}
		} else {
			s.WriteString(resultStyle.Render(entry.output))
		}

		// Show evaluation time if it took more than 10 ms
		if entry.evaluationTime > 10*time.Millisecond {
			timeStr := fmt.Sprintf(" (%.2fs)", entry.evaluationTime.Seconds())
			s.WriteString(historyStyle.Render(timeStr))
		}

		s.WriteString("\n\n")
	}

	// Current evaluation
	if m.evaluating {
		s.WriteString(promptStyle.Render(PROMPT))
		s.WriteString(m.highlightCode(m.currentInput))
		s.WriteString("\n")
		s.WriteString(m.spinner.View())
		s.WriteString(" Evaluating...")
		s.WriteString("\n\n")
	}

	// Input
	if !m.evaluating {
		s.WriteString(m.textInput.View())
		s.WriteString("\n")
	}

	// Help text
	s.WriteString(historyStyle.Render("\nPress Esc or Ctrl+C/D to exit"))

	return s.String()
}

// formatParseErrors formats parser errors into a string with improved readability
func formatParseErrors(errors []string) string {
	var s strings.Builder
	s.WriteString("Parser Errors:\n")

	for i, msg := range errors {
		s.WriteString(fmt.Sprintf("  %d. %s\n", i+1, msg))
	}

	s.WriteString("\nTips:\n")
	s.WriteString("  • Check for missing parentheses, braces, or semicolons\n")
	s.WriteString("  • Verify that all expressions are properly terminated\n")
	s.WriteString("  • Ensure variable names are valid identifiers\n")

	return s.String()
}

// formatRuntimeError formats a runtime error into a string with improved readability
func formatRuntimeError(errorMsg string) string {
	var s strings.Builder
	s.WriteString("Runtime Error:\n")
	s.WriteString("  " + errorMsg + "\n")

	s.WriteString("\nTips:\n")

	// Add specific tips based on common error patterns
	if strings.Contains(errorMsg, "identifier not found") {
		s.WriteString("  • Check if the variable is defined before use\n")
		s.WriteString("  • Verify the variable name is spelled correctly\n")
		s.WriteString("  • Make sure the variable is in scope\n")
	} else if strings.Contains(errorMsg, "wrong number of arguments") {
		s.WriteString("  • Check the function call has the correct number of arguments\n")
		s.WriteString("  • Verify the function definition matches its usage\n")
	} else if strings.Contains(errorMsg, "type mismatch") {
		s.WriteString("  • Ensure operands are of compatible types\n")
		s.WriteString("  • Check if you need to convert types before operation\n")
	} else if strings.Contains(errorMsg, "index") {
		s.WriteString("  • Verify array indices are within bounds\n")
		s.WriteString("  • Ensure you're indexing an array or hash\n")
	} else {
		s.WriteString("  • Review your code logic\n")
		s.WriteString("  • Check for type mismatches or undefined variables\n")
		s.WriteString("  • Consider breaking complex expressions into simpler steps\n")
	}

	return s.String()
}

// highlightCode applies syntax highlighting to Monkey code
func (m model) highlightCode(code string) string {
	// If NoColor option is enabled, return the code without highlighting
	if m.options.NoColor {
		return code
	}

	l := lexer.New(code)
	var s strings.Builder

	// Collect all tokens first
	var tokens []token.Token
	for {
		tok := l.NextToken()
		tokens = append(tokens, tok)
		if tok.Type == token.EOF {
			break
		}
	}

	// Simple approach: just highlight each token
	for _, tok := range tokens {
		if tok.Type == token.EOF {
			continue
		}

		// Apply the appropriate style based on token type
		switch tok.Type {
		case token.FUNCTION, token.LET, token.TRUE, token.FALSE, token.IF, token.ELSE, token.RETURN:
			s.WriteString(keywordStyle.Render(tok.Literal))
		case token.IDENT:
			s.WriteString(identifierStyle.Render(tok.Literal))
		case token.INT:
			s.WriteString(literalStyle.Render(tok.Literal))
		case token.STRING:
			s.WriteString(stringStyle.Render("\"" + tok.Literal + "\""))
		case token.ASSIGN, token.PLUS, token.MINUS, token.BANG, token.ASTERISK, token.SLASH,
			token.LT, token.GT, token.EQ, token.NOT_EQ:
			s.WriteString(operatorStyle.Render(tok.Literal))
		case token.COMMA, token.COLON, token.SEMICOLON, token.LPAREN, token.RPAREN,
			token.LBRACE, token.RBRACE, token.LBRACKET, token.RBRACKET:
			s.WriteString(delimiterStyle.Render(tok.Literal))
		default:
			s.WriteString(tok.Literal)
		}

		// Add a space after each token for readability
		// This is a simplification - in a real syntax highlighter, we'd preserve the original spacing
		if tok.Type != token.SEMICOLON && tok.Type != token.COMMA &&
			tok.Type != token.LPAREN && tok.Type != token.RPAREN &&
			tok.Type != token.LBRACE && tok.Type != token.RBRACE &&
			tok.Type != token.LBRACKET && tok.Type != token.RBRACKET {
			s.WriteString(" ")
		}
	}

	return s.String()
}
