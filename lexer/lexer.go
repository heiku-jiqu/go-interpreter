package lexer

import "github.com/heiku-jiqu/go-interpreter/token"

type Lexer struct{}

func NewLexer(s string) *Lexer {
	return &Lexer{}
}

func (l *Lexer) NextToken() token.Token {
	return token.Token{Type: token.LET, Literal: "let"}
}
