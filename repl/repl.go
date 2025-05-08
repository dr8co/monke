package repl

import (
	"bufio"
	"fmt"
	"github.com/dr8co/monke/lexer"
	"github.com/dr8co/monke/token"
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

		for tok := l.NextToken(); tok.Type != token.EOF; tok = l.NextToken() {
			_, errs := fmt.Fprintf(out, "%+v\n", tok)
			if errs != nil {
				panic(errs)
			}
		}
	}
}
