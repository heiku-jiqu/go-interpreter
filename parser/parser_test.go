package parser

import (
	"testing"

	"github.com/heiku-jiqu/go-interpreter/ast"
	"github.com/heiku-jiqu/go-interpreter/lexer"
	"github.com/heiku-jiqu/go-interpreter/token"
)

func TestLetStatement(t *testing.T) {
	input := "let x = 5;"
	l := lexer.NewLexer(input)
	p := NewParser(l)
	program := p.ParseProgram()

	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}
	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statement. got %d statement(s)", len(program.Statements))
	}

	expected := &ast.LetStatement{
		Token: token.Token{Type: token.LET, Literal: "let"},
		Name: &ast.Identifier{
			Token: token.Token{Type: token.IDENT, Literal: "x"},
			Value: "x",
		},
		Value: &ast.IntLiteral{
			Token: token.Token{Type: token.INT, Literal: "5"},
			Value: 5,
		},
	}
	letStmt, ok := program.Statements[0].(*ast.LetStatement)
	if !ok {
		t.Errorf("parsed statement is not a LetStatement")
	}

	if *letStmt.Name != *expected.Name {
		t.Errorf("invalid LetStatement.Name, got %+v, expected %+v", letStmt.Name, expected.Name)
	}
	letStmtValue, ok := letStmt.Value.(*ast.IntLiteral)
	if !ok {
		t.Errorf("parsed RHS of LetStatement is not IntLiteral")
	}
	if *letStmtValue != *expected.Value.(*ast.IntLiteral) {
		t.Errorf("invalid LetStatement.Value, got %+v, expected %+v", letStmt.Value, letStmt.Value)
	}
}
