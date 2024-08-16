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
    NUMBER
    IDENTIFIER
)

func (tt TokenType) String() string{
return [...]string{
        "EOF", "LEFT_PAREN", "RIGHT_PAREN", "LEFT_BRACE", "RIGHT_BRACE",
        "STAR", "DOT", "COMMA", "PLUS", "MINUS", "SLASH", "SEMICOLON",
        "EQUAL", "EQUAL_EQUAL", "BANG", "BANG_EQUAL",
        "LESS", "LESS_EQUAL", "GREATER", "GREATER_EQUAL",
        "STRING", "NUMBER", "IDENTIFIER",
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
    var literalStr string
    
    if t.Literal == nil{
        literalStr = "null"
    } else {
        switch v := t.Literal.(type) {
        case float64:
            if v == float64(int(v)){
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
