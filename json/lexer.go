package main

import (
	"fmt"
	"unicode/utf8"
)

func main() {
	// f, err := os.Open()
	json_text := `{"hi":123, "abc":"xyz", "arr":[1,2,3]}`
	l := New(json_text)
	for {
		t := l.NextToken()
		fmt.Printf("%v\n", t)
		if t.Type == ERROR || t.Type==EOF{
			break
		}
	}
}

type TokenType int

const (
	ERROR TokenType = iota
	LBRACE
	RBRACE
	LBRACKET // [
	RBRACKET // ]
	COMMA
	COLON

	TRUE
	FALSE
	NULL

	DIGIT // 0 - 9
	MINUS
	PLUS
	DOT      // .
	EXPONENT // E OR e

	DQUOTE // "
	ESCAPE //
	CHARS

	WHITESPACE // SPACE / LF / CR / TAB
	EOF
)

type Token struct {
	Type  TokenType
	Value string // the value from the input string source code
}

type Lexer struct {
	input    string // the input string source code
	len      int    // the input's len
	pos      int    // current position of lexer
	start    int    // starting position of current token
	inquotes bool   // whether lexer is currently in double quotes

	line int // for error messages
}

func New(input string) *Lexer {
	return &Lexer{
		input:    input,
		len:      utf8.RuneCountInString(input),
		pos:      0,
		start:    0,
		inquotes: false,

		line: 0,
	}
}

func (l *Lexer) NextToken() Token {
	// if l.inquotes {
	// 	return l.readChars()
	// }
	var t Token = Token{ERROR, "error"}
	if l.pos >= l.len {
		return Token{EOF, "EOF"}
	}
	r, err := l.peekRune()
	if err != nil {
		return t
	}
	switch r {
	case '{':
		t = Token{LBRACE, string(r)}
	case '}':
		t = Token{RBRACE, string(r)}
	case '[':
		t = Token{LBRACKET, string(r)}
	case ']':
		t = Token{RBRACKET, string(r)}
	case ':':
		t = Token{COLON, string(r)}
	case ',':
		t = Token{COMMA, string(r)}
	case '"':
		l.inquotes = !l.inquotes
		t = Token{DQUOTE, string(r)}
	case '-':
		t = Token{MINUS, string(r)}
	case '+':
		t = Token{PLUS, string(r)}
	case '.':
		t = Token{DOT, string(r)}
	case 'E', 'e':
		t = Token{EXPONENT, string(r)}
	case '\\':
		t = Token{ESCAPE, string(r)}
	case ' ', '\r', '\n':
		t = Token{WHITESPACE, string(r)}
	default:
		if l.inquotes {
			return l.readString()
		}
		if isDigit(r) {
			return l.readDigits()
		}

	}

	l.nextRune()
	return t
}

func (l *Lexer) nextRune() (rune, error) {
	r, size := utf8.DecodeRuneInString(l.input[l.pos:])
	if r == utf8.RuneError {
		return r, fmt.Errorf("Error decoding rune")
	}
	l.pos += size
	return r, nil
}
func (l *Lexer) peekRune() (rune, error) {
	r, _ := utf8.DecodeRuneInString(l.input[l.pos:])
	if r == utf8.RuneError {
		return r, fmt.Errorf("Error decoding rune")
	}
	return r, nil
}

func (l *Lexer) readDigits() Token {
	l.start = l.pos
	for {
		r, err := l.peekRune()
		if err != nil {
			return Token{ERROR, "decode rune error"}
		}
		if !isDigit(r) {
			break
		}
		l.pos += 1
	}
	t := Token{DIGIT, l.input[l.start:l.pos]}
	return t
}

func (l *Lexer) readString() Token {
	start := l.pos
	for {
		r, err := l.peekRune()
		if err != nil {
			return Token{ERROR, "not utf8"}
		}
		if r == '\\' || r == '"' {
			return Token{CHARS, l.input[start:l.pos]}
		}
		l.nextRune()
	}
}

func isDigit(ch rune) bool {
	return ch >= '0' && ch <= '9'
}

var keywords = map[string]TokenType{
	"null":  NULL,
	"true":  TRUE,
	"false": FALSE,
}

// returns keyword tokens else returns CHARS
func lookupKeywords(s string) Token {
	t, ok := keywords[s]
	if !ok {
		return Token{CHARS, s}
	}
	return Token{t, s}

}
