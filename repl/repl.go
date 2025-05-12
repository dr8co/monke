package repl

import (
	"bufio"
	"fmt"
	"github.com/dr8co/monke/lexer"
	"github.com/dr8co/monke/parser"
	"io"
)

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	for {
		_, err := fmt.Fprintf(out, PROMPT)
		if err != nil {
			panic(err)
		}
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		l := lexer.New(line)
		p := parser.New(l)

		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			printParseErrors(out, p.Errors())
			continue
		}

		_, err = io.WriteString(out, program.String())
		if err != nil {
			panic(err)
		}
		_, err = io.WriteString(out, "\n")
		if err != nil {
			panic(err)
		}
	}
}

func printParseErrors(out io.Writer, errors []string) {
	_, err := io.WriteString(out, "parser errors:\n")
	if err != nil {
		panic(err)
	}

	for _, msg := range errors {
		_, err = io.WriteString(out, "\t"+msg+"\n")
		if err != nil {
			panic(err)
		}
	}
}
