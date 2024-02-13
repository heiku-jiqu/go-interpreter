package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/heiku-jiqu/go-interpreter/lexer"
	"github.com/heiku-jiqu/go-interpreter/token"
)

const PROMPT = ">> "

type Repl struct{}

func (r *Repl) Start(in io.Reader, out io.Writer) error {
	scanner := bufio.NewScanner(in)

	for {
		fmt.Print(PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			err := scanner.Err()
			if err != nil {
				return err
			}
			return nil
		}
		input := scanner.Text()

		lex := lexer.NewLexer(input)

		for tok := lex.NextToken(); tok.Type != token.EOF && tok.Type != token.ILLEGAL; tok = lex.NextToken() {
			fmt.Fprintf(out, "%+v\n", tok)
		}
	}
}
