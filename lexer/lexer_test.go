package lexer

import (
	"testing"

	"github.com/heiku-jiqu/go-interpreter/token"
)

func TestLexing(t *testing.T) {
	t.Run("test next token", func(t *testing.T) {
		input := "let five = 5;"
		lexer := NewLexer(input)

		expected := []token.Token{
			{Type: token.LET, Literal: "let"},
			{Type: token.IDENT, Literal: "five"},
			{Type: token.ASSIGN, Literal: "="},
			{Type: token.INT, Literal: "INT"},
			{Type: token.SEMICOLON, Literal: ";"},
		}
		for _, expectTok := range expected {
			tok := lexer.NextToken()
			if tok != expectTok {
				t.Fatalf("expected %+v but got %+v", expectTok, tok)
			}
		}
	})
}
