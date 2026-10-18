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
	token.LPAREN:   CALL,
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
	p.registerPrefix(token.TRUE, p.parseBoolean)
	p.registerPrefix(token.FALSE, p.parseBoolean)
	p.registerPrefix(token.LPAREN, p.parseGroupedExpression)
	p.registerPrefix(token.IF, p.parseIfExpression)
	p.registerPrefix(token.FUNCTION, p.parseFunctionLiteral)

	p.infixParseFns = make(map[token.TokenType]infixParseFn)
	p.registerInfix(token.PLUS, p.parseInfixExpression)
	p.registerInfix(token.MINUS, p.parseInfixExpression)
	p.registerInfix(token.SLASH, p.parseInfixExpression)
	p.registerInfix(token.ASTERISK, p.parseInfixExpression)
	p.registerInfix(token.EQ, p.parseInfixExpression)
	p.registerInfix(token.NOT_EQ, p.parseInfixExpression)
	p.registerInfix(token.LT, p.parseInfixExpression)
	p.registerInfix(token.GT, p.parseInfixExpression)
	p.registerInfix(token.LPAREN, p.parseCallExpression)

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

// TODO: invariant should be everytime parseX runs, p.curToken is an UNPARSED token!
// similarly, try to refactor to LL(1) without the need for peekToken.
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

func (p *Parser) parseBoolean() ast.Expression {
	val := true
	if p.curToken.Type == token.FALSE {
		val = false
	}
	return &ast.Boolean{
		Token: p.curToken,
		Value: val,
	}
}
func (p *Parser) parseIfExpression() ast.Expression {
	ifToken := p.curToken

	if !p.expectPeek(token.LPAREN) {
		return nil
	}
	p.NextToken()

	exp := p.parseExpression(LOWEST)

	if !p.expectPeek(token.RPAREN) {
		return nil
	}
	if !p.expectPeek(token.LBRACE) {
		return nil
	}
	consequenceBlockStatement := p.parseBlockStatement()
	ifExp := &ast.IfExpression{
		Token:       ifToken,
		Condition:   exp,
		Consequence: consequenceBlockStatement,
	}

	if p.peekTokenIs(token.ELSE) { // don't use expectPeek to prevent adding parse error
		p.NextToken()
		if !p.expectPeek(token.LBRACE) {
			return nil
		}
		ifExp.Alternative = p.parseBlockStatement()
	}

	return ifExp
}
func (p *Parser) parseFunctionLiteral() ast.Expression {
	fl := &ast.FunctionLiteral{
		Token: p.curToken,
	}
	if !p.expectPeek(token.LPAREN) {
		return nil
	}
	fl.Parameters = p.parseFunctionParameters()
	if !p.expectPeek(token.LBRACE) {
		return nil
	}
	fl.Body = p.parseBlockStatement()
	return fl
}
func (p *Parser) parseFunctionParameters() []*ast.Identifier {
	params := []*ast.Identifier{}
	if p.peekTokenIs(token.RPAREN) {
		p.NextToken()
		return params
	}
	p.NextToken()
	ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	params = append(params, ident)
	for p.peekTokenIs(token.COMMA) {
		p.NextToken() // on comma
		p.NextToken() // on identifier
		ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
		params = append(params, ident)
	}
	if !p.expectPeek(token.RPAREN) {
		return nil
	}
	return params
}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	bs := &ast.BlockStatement{
		Token:      p.curToken,
		Statements: []ast.Statement{},
	}
	p.NextToken()
	for !p.curTokenIs(token.RBRACE) && !p.curTokenIs(token.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			bs.Statements = append(bs.Statements, stmt)
		}
		p.NextToken()
	}
	return bs
}
func (p *Parser) parseCallExpression(function ast.Expression) ast.Expression {
	ce := &ast.CallExpression{
		Token:    p.curToken,
		Function: function,
	}
	ce.Arguments = p.parseFunctionArguments()
	return ce
}
func (p *Parser) parseFunctionArguments() []ast.Expression {
	var args []ast.Expression
	if p.peekTokenIs(token.RPAREN) {
		p.NextToken()
		return args // no args, return empty
	}

	p.NextToken() // skip curr token which is LPAREN
	// first args
	exp := p.parseExpression(LOWEST)
	args = append(args, exp)

	// subsequent comma delimited args (if any)
	for p.peekTokenIs(token.COMMA) {
		p.NextToken() // skip over curr token arg
		p.NextToken() // skip over comma
		exp := p.parseExpression(LOWEST)
		args = append(args, exp)
	}
	// end of args, should be RPAREN
	if !p.expectPeek(token.RPAREN) {
		return nil
	}
	return args
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

func (p *Parser) parseGroupedExpression() ast.Expression {
	p.NextToken() // consume LPAREN
	exp := p.parseExpression(LOWEST)
	if !p.expectPeek(token.RPAREN) {
		return nil
	}
	return exp
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
