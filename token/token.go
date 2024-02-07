package token

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

const (
	ILLEGAL TokenType = "ILLEGAL"
	EOF     TokenType = "EOF"

	// Operators
	ASSIGN TokenType = "="
	PLUS   TokenType = "+"

	// Identifiers and Literals
	IDENT TokenType = "IDENT"
	INT   TokenType = "INT"

	// Delimiters
	COMMA     TokenType = ","
	SEMICOLON TokenType = ";"

	LBRACE TokenType = "{"
	RBRACE TokenType = "}"
	LPAREN TokenType = "("
	RPAREN TokenType = ")"

	// Keywords
	FUNCTION TokenType = "FUNCTION"
	LET      TokenType = "LET"
)
