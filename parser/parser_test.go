package parser

import (
	"testing"

	"github.com/heiku-jiqu/go-interpreter/ast"
	"github.com/heiku-jiqu/go-interpreter/lexer"
)

// func TestLetStatement(t *testing.T) {
// 	input := `
// let x = 5;
// let test = 69;
// let elephant = 420;
// `
// 	l := lexer.New(input)
// 	p := New(l)
// 	program := p.ParseProgram()
// 	checkParserErrors(t, p)

// 	if program == nil {
// 		t.Fatalf("ParseProgram() returned nil")
// 	}
// 	if len(program.Statements) != 3 {
// 		t.Fatalf("program.Statements does not contain 3 statement. got %d statement(s)", len(program.Statements))
// 	}

// 	expecteds := []struct {
// 		name string
// 		val  int
// 	}{
// 		{"x", 5},
// 		{"test", 69},
// 		{"elephant", 420},
// 	}

// 	for i, s := range program.Statements {
// 		exp := expecteds[i]
// 		testLetStatement(t, s, exp.name, exp.val)
// 	}
// }

// func testLetStatement(t *testing.T, s ast.Statement, name string, expectedVal int) {
// 	letStmt, ok := s.(*ast.LetStatement)
// 	if !ok {
// 		t.Errorf("parsed statement is not a LetStatement")
// 	}

// 	if letStmt.Name.Value != name {
// 		t.Errorf("invalid LetStatement.Name.Value, got %s, expected %s", letStmt.Name.Value, name)
// 	}
// 	if letStmt.Name.TokenLiteral() != name {
// 		t.Errorf("invalid LetStatement.Name, got %s, expected %s", letStmt.Name, name)
// 	}
// 	letStmtValue, ok := letStmt.Value.(*ast.IntLiteral)
// 	if !ok {
// 		t.Errorf("parsed RHS of LetStatement is not IntLiteral")
// 	}

// 	if letStmtValue.Value != expectedVal {
// 		t.Errorf("invalid LetStatement.Value, got %+v, expected %+v", letStmtValue.Value, expectedVal)
// 	}
// }
// func TestReturnStatements(t *testing.T) {
// 	input := `
//  return 5;
//  return 10;
//  return 993322;`
// 	l := lexer.New(input)
// 	p := New(l)
// 	program := p.ParseProgram()
// 	checkParserErrors(t, p)
// 	if len(program.Statements) != 3 {
// 		t.Fatalf("program.Statements does not contain 3 statements. got=%d", len(program.Statements))
// 	}
// 	for _, stmt := range program.Statements {
// 		returnStmt, ok := stmt.(*ast.ReturnStatement)
// 		if !ok {
// 			t.Errorf("stmt not *ast.returnStatement. got=%T", stmt)
// 			continue
// 		}
// 		if returnStmt.TokenLiteral() != "return" {
// 			t.Errorf("returnStmt.TokenLiteral not 'return', got %q",
// 				returnStmt.TokenLiteral())
// 		}
// 	}
// }

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

func TestIdentifierExpression(t *testing.T) {
	input := "foobar;"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("Program has not enough statements. got=%d", len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not *ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	ident, ok := stmt.Expression.(*ast.Identifier)
	if !ok {
		t.Fatalf("exp not *ast.Identifier. got=%T", stmt.Expression)
	}
	if ident.Value != "foobar" {
		t.Errorf("ident.Value not %s. got=%s", "foobar", ident.Value)
	}
	if ident.TokenLiteral() != "foobar" {
		t.Errorf("ident.TokenLiteral not %s. got=%s", "foobar", ident.TokenLiteral())
	}

}

func TestIntegerLiteralExpression(t *testing.T) {
	input := "5;"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)
	if len(program.Statements) != 1 {
		t.Fatalf("Program has not enough statements. got=%d", len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not *ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	literal, ok := stmt.Expression.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("exp not *ast.IntegerLiteral. got=%T", stmt.Expression)
	}
	if literal.Value != 5 {
		t.Errorf("ident.Value not %s. got=%d", "foobar", literal.Value)
	}
	if literal.TokenLiteral() != "5" {
		t.Errorf("ident.TokenLiteral not %s. got=%s", "foobar", literal.TokenLiteral())
	}
}
