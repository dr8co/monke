// Monke is a simple programming language written in Go.
package main

import (
	"flag"
	"fmt"
	"os"
	"os/user"
	"path/filepath"

	"github.com/dr8co/monke/evaluator"
	"github.com/dr8co/monke/lexer"
	"github.com/dr8co/monke/object"
	"github.com/dr8co/monke/parser"
	"github.com/dr8co/monke/repl"
)

const VERSION = "0.1.0"

func main() {
	// Define command-line flags
	noColor := flag.Bool("no-color", false, "Disable syntax highlighting and colored output")
	fileFlag := flag.String("file", "", "Execute a Monkey script file")
	evalFlag := flag.String("eval", "", "Evaluate a Monkey expression and print the result")
	debugFlag := flag.Bool("debug", false, "Enable debug mode with more verbose output")
	versionFlag := flag.Bool("version", false, "Show version information")

	// Define short flag aliases
	flag.BoolVar(noColor, "n", false, "Disable syntax highlighting and colored output")
	flag.StringVar(fileFlag, "f", "", "Execute a Monkey script file")
	flag.StringVar(evalFlag, "e", "", "Evaluate a Monkey expression and print the result")
	flag.BoolVar(debugFlag, "d", false, "Enable debug mode with more verbose output")
	flag.BoolVar(versionFlag, "v", false, "Show version information")

	// Parse command-line flags
	flag.Parse()

	// Show version information if requested
	if *versionFlag {
		fmt.Printf("Monkey Programming Language v%s\n", VERSION)
		return
	}

	// Get current user
	usr, err := user.Current()
	if err != nil {
		panic(err)
	}

	// Create options struct for REPL
	options := repl.Options{
		NoColor: *noColor,
		Debug:   *debugFlag,
	}

	// Execute a file if specified
	if *fileFlag != "" {
		executeFile(*fileFlag, *debugFlag)
		return
	}

	// Evaluate an expression if specified
	if *evalFlag != "" {
		evaluateExpression(*evalFlag)
		return
	}

	// Start the REPL
	repl.Start(usr.Username, options)
}

// executeFile reads and executes a Monkey script file
func executeFile(filename string, debug bool) {
	cleaned := filepath.Clean(filename)
	absolute, err := filepath.Abs(cleaned)
	if err != nil {
		fmt.Printf("Error getting absolute path: %s\n", err)
		os.Exit(1)
	}
	fmt.Printf("Executing file: %s\n", absolute)

	// Read the file
	//nolint:gosec // We're not reading user input here
	content, err := os.ReadFile(absolute)
	if err != nil {
		fmt.Printf("Error reading file: %s\n", err)
		os.Exit(1)
	}

	// Create environment
	env := object.NewEnvironment()

	// Parse and evaluate the file
	l := lexer.New(string(content))
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) != 0 {
		printParserErrors(p.Errors())
		os.Exit(1)
	}

	evaluated := evaluator.Eval(program, env)

	// Print the result if in debug mode
	if debug && evaluated != nil {
		fmt.Println(evaluated.Inspect())
	}
}

// evaluateExpression evaluates a single Monkey expression
func evaluateExpression(expr string) {
	// Create environment
	env := object.NewEnvironment()

	// Parse and evaluate the expression
	l := lexer.New(expr)
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) != 0 {
		printParserErrors(p.Errors())
		os.Exit(1)
	}

	evaluated := evaluator.Eval(program, env)

	// Print the result
	if evaluated != nil {
		fmt.Println(evaluated.Inspect())
	}
}

// printParserErrors prints parser errors to stderr
func printParserErrors(errors []string) {
	_, err := fmt.Fprintln(os.Stderr, "Parser errors:")
	if err != nil {
		panic(err)
	}
	for _, msg := range errors {
		_, err := fmt.Fprintln(os.Stderr, "\t"+msg)
		if err != nil {
			panic(err)
		}
	}
}
