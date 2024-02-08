package token

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
	// TODO: line number, column number and filename
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

var keywords = map[string]TokenType{
	"let": LET,
	"fn":  FUNCTION,
}

func LookupIdent(ident string) TokenType {
	tok, ok := keywords[ident]
	if ok {
		return tok
	}
	return IDENT
}
