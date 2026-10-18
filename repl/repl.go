package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/heiku-jiqu/go-interpreter/lexer"
	"github.com/heiku-jiqu/go-interpreter/parser"
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

		lex := lexer.New(input)
		p := parser.New(lex)

		prog := p.ParseProgram()
		if len(p.Errors()) != 0 {
			printParserErrors(out, p.Errors())
			continue
		}

		io.WriteString(out, prog.String())
		io.WriteString(out, "\n")
	}
}

func printParserErrors(out io.Writer, errors []string) {
	io.WriteString(out, " parser errors:\n")
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
