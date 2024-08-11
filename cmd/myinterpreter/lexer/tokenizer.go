package lexer

import "fmt"


type TokenType int

const (
    EOF TokenType = iota
    LEFT_PAREN
    RIGHT_PAREN
    LEFT_BRACE
    RIGHT_BRACE
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
)

func (tt TokenType) String() string{
return [...]string{
        "EOF", "LEFT_PAREN", "RIGHT_PAREN", "LEFT_BRACE", "RIGHT_BRACE",
        "STAR", "DOT", "COMMA", "PLUS", "MINUS", "SLASH", "SEMICOLON",
        "EQUAL", "EQUAL_EQUAL", "BANG", "BANG_EQUAL",
        "LESS", "LESS_EQUAL", "GREATER", "GREATER_EQUAL",
        "STRING",
    }[tt]
}

type Token struct{
    Lexeme string
    Token TokenType
    Literal any
}

func NewToken(lexeme string, token TokenType, literal any) *Token{
    t := new(Token)
    t.Lexeme = lexeme
    t.Token = token
    t.Literal = literal

    return t
}

func (t *Token) String() string{
    var s string
    if t.Token == STRING{
        s = fmt.Sprintf("%s \"%s\" ", t.Token, t.Lexeme)
        s += t.Lexeme
    } else {
        s = fmt.Sprintf("%s %s ", t.Token, t.Lexeme)
    }
    if t.Literal == nil{
        s += "null"
    }

    return s
}
