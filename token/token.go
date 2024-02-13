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
	ASSIGN   TokenType = "="
	PLUS     TokenType = "+"
	MINUS    TokenType = "-"
	SLASH    TokenType = "/"
	ASTERISK TokenType = "*"
	BANG     TokenType = "!"
	GT       TokenType = ">"
	LT       TokenType = "<"
	EQ       TokenType = "=="
	NOT_EQ   TokenType = "!="

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
	IF       TokenType = "IF"
	ELSE     TokenType = "ELSE"
	RETURN   TokenType = "RETURN"
	TRUE     TokenType = "TRUE"
	FALSE    TokenType = "FALSE"
)

var keywords = map[string]TokenType{
	"let":    LET,
	"fn":     FUNCTION,
	"if":     IF,
	"else":   ELSE,
	"return": RETURN,
	"true":   TRUE,
	"false":  FALSE,
}

func LookupIdent(ident string) TokenType {
	tok, ok := keywords[ident]
	if ok {
		return tok
	}
	return IDENT
}
