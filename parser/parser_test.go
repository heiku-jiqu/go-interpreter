package parser

import (
	"testing"

	"github.com/heiku-jiqu/go-interpreter/ast"
	"github.com/heiku-jiqu/go-interpreter/lexer"
)

func TestLetStatement(t *testing.T) {
	input := `
let x = 5;
let test = 69;
let elephant = 420;
`
	l := lexer.NewLexer(input)
	p := NewParser(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}
	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain 3 statement. got %d statement(s)", len(program.Statements))
	}

	expecteds := []struct {
		name string
		val  int
	}{
		{"x", 5},
		{"test", 69},
		{"elephant", 420},
	}

	for i, s := range program.Statements {
		exp := expecteds[i]
		testLetStatement(t, s, exp.name, exp.val)
	}
}

func testLetStatement(t *testing.T, s ast.Statement, name string, expectedVal int) {
	letStmt, ok := s.(*ast.LetStatement)
	if !ok {
		t.Errorf("parsed statement is not a LetStatement")
	}

	if letStmt.Name.Value != name {
		t.Errorf("invalid LetStatement.Name.Value, got %s, expected %s", letStmt.Name.Value, name)
	}
	if letStmt.Name.TokenLiteral() != name {
		t.Errorf("invalid LetStatement.Name, got %s, expected %s", letStmt.Name, name)
	}
	letStmtValue, ok := letStmt.Value.(*ast.IntLiteral)
	if !ok {
		t.Errorf("parsed RHS of LetStatement is not IntLiteral")
	}

	if letStmtValue.Value != expectedVal {
		t.Errorf("invalid LetStatement.Value, got %+v, expected %+v", letStmtValue.Value, expectedVal)
	}
}

func checkParserErrors(t *testing.T, p *Parser) {
	t.Helper()
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}

	t.Errorf("parser has %d error(s)", len(errors))
	for _, msg := range errors {
		t.Errorf("parser error: %q", msg)
	}
	t.FailNow()
}
