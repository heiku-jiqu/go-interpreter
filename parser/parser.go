package parser

import (
	"strconv"

	"github.com/heiku-jiqu/go-interpreter/ast"
	"github.com/heiku-jiqu/go-interpreter/lexer"
	"github.com/heiku-jiqu/go-interpreter/token"
)

type Parser struct {
	l *lexer.Lexer

	curToken  token.Token
	peekToken token.Token
}

func NewParser(l *lexer.Lexer) *Parser {
	p := &Parser{l: l}

	p.NextToken()
	p.NextToken()

	return p
}

func (p *Parser) NextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) ParseProgram() *ast.Program {
	stmts := []ast.Statement{}
	if p.curToken.Type == token.LET {
		stmt := p.ParseLet()
		stmts = append(stmts, stmt)
	}
	prog := &ast.Program{Statements: stmts}
	return prog
}

func (p *Parser) ParseLet() *ast.LetStatement {
	letToken := p.curToken

	p.NextToken()
	identifier := ast.Identifier{
		Token: p.curToken,
		Value: p.curToken.Literal,
	}

	p.NextToken()
	if p.curToken.Type != token.ASSIGN {
		// error
	}

	p.NextToken()
	value := p.ParseExpression()

	s := &ast.LetStatement{
		Token: letToken,
		Name:  &identifier,
		Value: value,
	}
	return s
}

func (p *Parser) ParseExpression() ast.Expression {
	i, err := strconv.Atoi(p.curToken.Literal)
	if err != nil {
		// error
	}
	return &ast.IntLiteral{
		Token: p.curToken,
		Value: i,
	}
}
