package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
)

type TokenType string

const (
	// single character token types
	PLUS      = "+"
	MINUS     = "-"
	MUL       = "*"
	FLOAT_DIV = "/"
	LPAREN    = "("
	RPAREN    = ")"
	LBRACE    = "{"
	RBRACE    = "}"
	SEMI      = ";"
	DOT       = "."
	COLON     = ":"
	COMMA     = ","
	ASSIGN    = ":="
	LESS      = "<"
	// block of reserved words
	PROGRAM     = "PROGRAM"
	INTEGER     = "INTEGER"
	REAL        = "REAL"
	INTEGER_DIV = "DIV"
	VAR         = "VAR"
	PROCEDURE   = "PROCEDURE"
	BEGIN       = "BEGIN"
	END         = "END"
	// misc
	IDENT         = "IDENT"
	INTEGER_CONST = "INTEGER_CONST"
	REAL_CONST    = "REAL_CONST"
	EOF           = "EOF"
	NONE          = 255
)

// Reserved keywords
var keywords = map[string]TokenType{
	"PROGRAM":   PROGRAM,
	"INTEGER":   INTEGER,
	"REAL":      REAL,
	"DIV":       INTEGER_DIV,
	"VAR":       VAR,
	"PROCEDURE": PROCEDURE,
	"BEGIN":     BEGIN,
	"END":       END,
}

type Token struct {
	Type   TokenType
	Value  interface{}
	Line   int
	Column int
}

type Lexer struct {
	Text        string
	CurrentChar byte
	Pos         int
	Line        int
	Column      int
}

func NewLexer(text string) *Lexer {
	l := Lexer{Text: text}
	l.Pos = 0
	l.CurrentChar = l.Text[l.Pos]
	l.Column = 1
	l.Line = 1
	return &l
}

func (l *Lexer) LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}

func (l *Lexer) error() {
	fmt.Printf("Unexpected character '%s'\n", string(l.CurrentChar))
	os.Exit(1)
}

func (l *Lexer) advance() {
	if l.CurrentChar == '\r' { // for Windows compatibility
		l.nextChar()
	}
	if l.CurrentChar == '\n' {
		l.Line++
		l.Column = 0
		l.nextChar()
	}
	l.nextChar()
}

func (l *Lexer) nextChar() {
	l.Pos++
	if l.Pos > len(l.Text)-1 {
		l.CurrentChar = NONE
	} else {
		l.CurrentChar = l.Text[l.Pos]
	}
}

func (l *Lexer) peek() byte {
	peekPos := l.Pos + 1
	if peekPos > len(l.Text)-1 {
		return NONE
	} else {
		return l.Text[peekPos]
	}
}

func (l *Lexer) isSpace(ch byte) bool {
	return ch == ' '
}

func (l *Lexer) isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func (l *Lexer) isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func (l *Lexer) skipWhiteSpace() {
	for l.CurrentChar != NONE && l.isSpace(l.CurrentChar) {
		l.advance()
	}
}

func (l *Lexer) number() Token {
	result := ""
	for l.CurrentChar != NONE && l.isDigit(l.CurrentChar) {
		result += string(l.CurrentChar)
		l.advance()
	}
	if l.CurrentChar == '.' && l.isDigit(l.peek()) {
		result += string(l.CurrentChar)
		l.advance() // eat '.'
		for l.CurrentChar != NONE && l.isDigit(l.CurrentChar) {
			result += string(l.CurrentChar)
			l.advance()
		}
		real, _ := strconv.ParseFloat(result, 64)
		return Token{Type: REAL_CONST, Value: real, Line: l.Line, Column: l.Column}
	}
	integer, _ := strconv.Atoi(result)
	return Token{Type: INTEGER_CONST, Value: integer, Line: l.Line, Column: l.Column}
}

func (l *Lexer) identifier() Token {
	result := ""
	for l.CurrentChar != NONE && l.isLetter(l.CurrentChar) {
		result += string(l.CurrentChar)
		l.advance()
	}
	return Token{Type: l.LookupIdent(result), Value: result, Line: l.Line, Column: l.Column}
}

func (l *Lexer) skipComments() {
	for l.CurrentChar != NONE && l.CurrentChar != '}' {
		l.advance()
	}
	if l.CurrentChar == NONE {
		fmt.Print("Unexpected End Of File.")
		os.Exit(1)
	}
	l.advance() // eat closing comment '}'
}

func (l *Lexer) GetNextToken() Token {
	for l.CurrentChar != NONE {
		if l.isSpace(l.CurrentChar) {
			l.skipWhiteSpace()
			continue
		}
		if l.CurrentChar == '{' {
			l.skipComments()
			continue
		}
		if l.isDigit(l.CurrentChar) {
			return l.number()
		}
		if l.isLetter(l.CurrentChar) {
			return l.identifier()
		}
		if l.CurrentChar == '+' {
			l.advance()
			return Token{Type: PLUS, Value: "+", Line: l.Line, Column: l.Column}
		}
		if l.CurrentChar == '-' {
			l.advance()
			return Token{Type: MINUS, Value: "-", Line: l.Line, Column: l.Column}
		}
		if l.CurrentChar == '*' {
			l.advance()
			return Token{Type: MUL, Value: "*", Line: l.Line, Column: l.Column}
		}
		if l.CurrentChar == '/' {
			l.advance()
			return Token{Type: FLOAT_DIV, Value: "/", Line: l.Line, Column: l.Column}
		}
		if l.CurrentChar == '(' {
			l.advance()
			return Token{Type: LPAREN, Value: "(", Line: l.Line, Column: l.Column}
		}
		if l.CurrentChar == ')' {
			l.advance()
			return Token{Type: RPAREN, Value: ")", Line: l.Line, Column: l.Column}
		}
		if l.CurrentChar == '.' {
			l.advance()
			return Token{Type: DOT, Value: ".", Line: l.Line, Column: l.Column}
		}
		if l.CurrentChar == ',' {
			l.advance()
			return Token{Type: COMMA, Value: ",", Line: l.Line, Column: l.Column}
		}
		if l.CurrentChar == ';' {
			l.advance()
			return Token{Type: SEMI, Value: ";", Line: l.Line, Column: l.Column}
		}
		if l.CurrentChar == ':' {
			if l.peek() == '=' {
				l.advance()
				l.advance()
				return Token{Type: ASSIGN, Value: ":=", Line: l.Line, Column: l.Column}
			}
			l.advance()
			return Token{Type: COLON, Value: ":", Line: l.Line, Column: l.Column}
		}
		if l.CurrentChar == '\r' || l.CurrentChar == '\n' {
			l.advance()
		} else {
			l.error()
		}
	}
	return Token{Type: EOF, Value: NONE}
}

func main() {
	//text := "program Main;procedure Alpha(a : integer; b : integer); var x : integer;"

	pwd, _ := os.Getwd()
	txt, _ := ioutil.ReadFile(pwd + "/test.txt")
	content := string(txt)

	lexer := NewLexer(content)
	token := lexer.GetNextToken()
	for token.Value != NONE {
		fmt.Printf("%-v\n", token)
		token = lexer.GetNextToken()
	}
}
