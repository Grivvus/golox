package main

import "fmt"

type TokenType int

const (
	EOF TokenType = iota
	LEFT_PAREN
	RIGHT_PAREN
	LEFT_BRACE
	RIGHT_BRACE
	LEFT_SQUARE_BRACKET
	RIGHT_SQUARE_BRACKET
	STAR
	DOT
	COMMA
	PLUS
	MINUS
	SLASH
	SEMICOLON
	EQUAL
	EQUAL_EQUAL
	BANG
	BANG_EQUAL
	LESS
	LESS_EQUAL
	GREATER
	GREATER_EQUAL
	STRING
	NUMBER
	IDENTIFIER
	AND
	OR
	CLASS
	SUPER
	THIS
	IF
	ELSE
	TRUE
	FALSE
	FOR
	WHILE
	FUN
	RETURN
	NIL
	PRINT
	VAR
)

func fillMap() *map[string]TokenType {
	res := map[string]TokenType{
		"and":    AND,
		"or":     OR,
		"class":  CLASS,
		"super":  SUPER,
		"this":   THIS,
		"if":     IF,
		"else":   ELSE,
		"true":   TRUE,
		"false":  FALSE,
		"for":    FOR,
		"while":  WHILE,
		"fun":    FUN,
		"return": RETURN,
		"nil":    NIL,
		"print":  PRINT,
		"var":    VAR,
	}

	return &res
}

func (tt TokenType) String() string {
	return [...]string{
		"EOF", "LEFT_PAREN", "RIGHT_PAREN", "LEFT_BRACE", "RIGHT_BRACE",
		"LEFT_SQUARE_BRACKET", "RIGHT_SQUARE_BRACKET",
		"STAR", "DOT", "COMMA", "PLUS", "MINUS", "SLASH", "SEMICOLON",
		"EQUAL", "EQUAL_EQUAL", "BANG", "BANG_EQUAL",
		"LESS", "LESS_EQUAL", "GREATER", "GREATER_EQUAL",
		"STRING", "NUMBER", "IDENTIFIER", "AND", "OR",
		"CLASS", "SUPER", "THIS", "IF", "ELSE", "TRUE", "FALSE",
		"FOR", "WHILE", "FUN", "RETURN", "NIL", "PRINT", "VAR",
	}[tt]
}

type Token struct {
	Lexeme  string
	Token   TokenType
	Literal any
	Line    uint
}

func NewToken(lexeme string, token TokenType, literal any, line uint) *Token {
	t := new(Token)
	t.Lexeme = lexeme
	t.Token = token
	t.Literal = literal
	t.Line = line

	return t
}

func (t *Token) String() string {
	var literalStr string

	if t.Literal == nil {
		literalStr = "null"
	} else {
		switch v := t.Literal.(type) {
		case float64:
			if v == float64(int(v)) {
				literalStr = fmt.Sprintf("%.1f", v)
			} else {
				literalStr = fmt.Sprintf("%v", v)
			}
		default:
			literalStr = fmt.Sprintf("%v", v)
		}
	}
	s := fmt.Sprintf("%v %v %v", t.Token, t.Lexeme, literalStr)

	return s
}
