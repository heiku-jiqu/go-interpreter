package lexer

import "github.com/heiku-jiqu/go-interpreter/token"

type Lexer struct {
	sourceCode string
}

func NewLexer(s string) *Lexer {
	return &Lexer{s}
}

func (l *Lexer) NextToken() token.Token {
	return token.Token{Type: token.LET, Literal: "let"}
}
