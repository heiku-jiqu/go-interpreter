package ast

import "github.com/heiku-jiqu/go-interpreter/token"

type Node interface {
	TokenLiteral() string
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}

type LetStatement struct {
	Token token.Token
	Name  *Identifier
	Value Expression
}

func (l *LetStatement) TokenLiteral() string { return l.Token.Literal }
func (l *LetStatement) statementNode()       {}

type Identifier struct {
	Token token.Token
	Value string
}

func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
func (i *Identifier) expressionNode()      {}

type IntLiteral struct {
	Token token.Token
	Value int
}

func (i *IntLiteral) TokenLiteral() string { return i.Token.Literal }
func (i *IntLiteral) expressionNode()      {}
