package parser

import (
	"fmt"
	"strconv"

	"github.com/heiku-jiqu/go-interpreter/ast"
	"github.com/heiku-jiqu/go-interpreter/lexer"
	"github.com/heiku-jiqu/go-interpreter/token"
)

type Parser struct {
	l *lexer.Lexer

	errors []string

	curToken  token.Token
	peekToken token.Token
}

func NewParser(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{},
	}

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
	for p.curToken.Type != token.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			stmts = append(stmts, stmt)
		}
		p.NextToken()
	}
	prog := &ast.Program{Statements: stmts}
	return prog
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.LET:
		return p.parseLetStatement()
	default:
		return nil
	}
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	s := &ast.LetStatement{
		Token: p.curToken,
	}
	p.NextToken()
	identifier := &ast.Identifier{
		Token: p.curToken,
		Value: p.curToken.Literal,
	}
	s.Name = identifier

	// assert next token is `=` before going to next token
	if !p.peekTokenIs(token.ASSIGN) {
		// error
		p.peekTokenTypeError(token.ASSIGN)
		return nil
	}
	p.NextToken()

	p.NextToken()
	value := p.parseExpression()
	s.Value = value

	for p.curToken.Type != token.SEMICOLON {
		p.NextToken()
	}

	return s
}

func (p *Parser) parseExpression() ast.Expression {
	i, err := strconv.Atoi(p.curToken.Literal)
	if err != nil {
		// error
	}
	return &ast.IntLiteral{
		Token: p.curToken,
		Value: i,
	}
}

func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

// Appends error mismatch TokenType message in p.errors
func (p *Parser) peekTokenTypeError(expectedType token.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, but got %s instead.", expectedType, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

func (p *Parser) Errors() []string {
	return p.errors
}
