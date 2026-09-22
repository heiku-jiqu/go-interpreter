package parser

import (
	"fmt"
	"strconv"

	"github.com/heiku-jiqu/go-interpreter/ast"
	"github.com/heiku-jiqu/go-interpreter/lexer"
	"github.com/heiku-jiqu/go-interpreter/token"
)

const (
	_ int = iota
	LOWEST
	EQUALS
	LESSGREATER
	SUM
	PRODUCT
	PREFIX
	CALL
)

var precedences = map[token.TokenType]int{
	token.EQ:       EQUALS,
	token.NOT_EQ:   EQUALS,
	token.LT:       LESSGREATER,
	token.GT:       LESSGREATER,
	token.PLUS:     SUM,
	token.MINUS:    SUM,
	token.SLASH:    PRODUCT,
	token.ASTERISK: PRODUCT,
}

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

type Parser struct {
	l *lexer.Lexer

	errors []string

	curToken  token.Token
	peekToken token.Token

	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{},
	}

	p.NextToken()
	p.NextToken()

	p.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	p.registerPrefix(token.IDENT, p.parseIdentifier)
	p.registerPrefix(token.INT, p.parseIntegerLiteral)
	p.registerPrefix(token.BANG, p.parsePrefixExpression)
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)

	p.infixParseFns = make(map[token.TokenType]infixParseFn)
	p.registerInfix(token.PLUS, p.parseInfixExpression)
	p.registerInfix(token.MINUS, p.parseInfixExpression)
	p.registerInfix(token.SLASH, p.parseInfixExpression)
	p.registerInfix(token.ASTERISK, p.parseInfixExpression)
	p.registerInfix(token.EQ, p.parseInfixExpression)
	p.registerInfix(token.NOT_EQ, p.parseInfixExpression)
	p.registerInfix(token.LT, p.parseInfixExpression)
	p.registerInfix(token.GT, p.parseInfixExpression)

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
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
	}
}
func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	returnStmt := &ast.ReturnStatement{
		Token: p.curToken,
	}
	p.NextToken()

	for !p.curTokenIs(token.SEMICOLON) {
		p.NextToken()
	}
	return returnStmt
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	s := &ast.LetStatement{
		Token: p.curToken,
	}

	if !p.expectPeek(token.IDENT) {
		// error
		p.peekTokenTypeError(token.IDENT)
		return nil
	}
	identifier := &ast.Identifier{
		Token: p.curToken,
		Value: p.curToken.Literal,
	}
	s.Name = identifier

	// assert next token is `=` before going to next token
	if !p.expectPeek(token.ASSIGN) {
		// error
		p.peekTokenTypeError(token.ASSIGN)
		return nil
	}

	p.NextToken()
	value := p.parseExpression(LOWEST)
	s.Value = value

	for !p.curTokenIs(token.SEMICOLON) {
		p.NextToken()
	}

	return s
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.curToken}
	stmt.Expression = p.parseExpression(LOWEST)
	if p.peekTokenIs(token.SEMICOLON) {
		p.NextToken()
	}
	return stmt
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefixFn := p.prefixParseFns[p.curToken.Type]
	if prefixFn == nil {
		p.noPrefixParseFnError(p.curToken.Type)
		return nil
	}
	leftExp := prefixFn()
	// for p.peekPrecedence() > precedence && !p.peekTokenIs(token.SEMICOLON) {
	// 	infixParseFn := p.infixParseFns[p.peekToken.Type]
	// 	if infixParseFn == nil {
	// 		return leftExp
	// 	}
	// 	p.NextToken()
	// 	leftExp = infixParseFn(leftExp)
	// }
	for {
		infixParseFn := p.infixParseFns[p.peekToken.Type]
		if p.peekPrecedence() <= precedence ||
			infixParseFn == nil {
			// early return when
			// 	- precedence <= next precedence (i.e. build left expression first)
			// 	- no infix parse fn (i.e. token.EOF, etc)
			break
		}
		p.NextToken()
		leftExp = infixParseFn(leftExp)
	}
	return leftExp
}
func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{
		Token: p.curToken,
		Value: p.curToken.Literal,
	}
}
func (p *Parser) parseIntegerLiteral() ast.Expression {
	val, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}
	return &ast.IntegerLiteral{
		Token: p.curToken,
		Value: val,
	}
}
func (p *Parser) parsePrefixExpression() ast.Expression {
	prefixTok := p.curToken
	p.NextToken()
	right := p.parseExpression(PREFIX)
	return &ast.PrefixExpression{
		Token:    prefixTok,
		Operator: prefixTok.Literal,
		Right:    right,
	}
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	infixTok := p.curToken
	infixPrecedence := p.curPrecedence()
	p.NextToken()
	right := p.parseExpression(infixPrecedence)
	return &ast.InfixExpression{
		Token:    infixTok,
		Left:     left,
		Operator: infixTok.Literal,
		Right:    right,
	}

}
func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.curToken.Type == t
}
func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}
func (p *Parser) curPrecedence() int {
	if p, ok := precedences[p.curToken.Type]; ok {
		return p
	}
	return LOWEST
}
func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}
	return LOWEST
}

// expectPeek checks whether next token's type is `t` and only advances nexttoken if so while returning true.
// else it returns false while keeping current token.
//
// Nearly all parsers have this method.
func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.NextToken()
		return true
	} else {
		p.peekTokenTypeError(t)
		return false
	}
}

// Appends error mismatch TokenType message in p.errors
func (p *Parser) peekTokenTypeError(expectedType token.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, but got %s instead.", expectedType, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) registerPrefix(t token.TokenType, f prefixParseFn) {
	p.prefixParseFns[t] = f
}

func (p *Parser) registerInfix(t token.TokenType, f infixParseFn) {
	p.infixParseFns[t] = f
}

func (p *Parser) noPrefixParseFnError(t token.TokenType) {
	msg := fmt.Sprintf("no prefix parse function for %s found", t)
	p.errors = append(p.errors, msg)
}
