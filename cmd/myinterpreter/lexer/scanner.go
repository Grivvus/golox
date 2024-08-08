package lexer

import (
	"errors"
    "fmt"
)

type Scanner struct{
    Source []byte
    ExitCode int
    CurrentIndex int
    CurrentLine uint
}

func NewScanner(source []byte) *Scanner{
    scanner := new(Scanner)
    scanner.Source = source
    scanner.CurrentIndex = 0
    scanner.CurrentLine = 1
    scanner.ExitCode = 0

    return scanner
}

func (s *Scanner) NextToken() (*Token, error){
    if s.CurrentIndex < len(s.Source) && s.Source[s.CurrentIndex] == '\n'{
        s.CurrentIndex++
        s.CurrentLine++
    }

    if s.CurrentIndex >= len(s.Source){
        return NewToken("", EOF, nil), nil
    }

    char := s.Source[s.CurrentIndex]
    s.CurrentIndex++

    switch char{
    case '(':
        return NewToken("(", LEFT_PAREN, nil), nil
    case ')':
        return NewToken(")", RIGHT_PAREN, nil), nil
    case '{':
        return NewToken("{", LEFT_BRACE, nil), nil
    case '}':
        return NewToken("}", RIGHT_BRACE, nil), nil
    case '*':
        return NewToken("*", STAR, nil), nil
    case '.':
        return NewToken(".", DOT, nil), nil
    case ',':
        return NewToken(",", COMMA, nil), nil
    case '+':
        return NewToken("+", PLUS, nil), nil
    case '-':
        return NewToken("-", MINUS, nil), nil
    case '/':
        return NewToken("/", SLASH, nil), nil
    case ';':
        return NewToken(";", SEMICOLON, nil), nil
    case '=':
        if s.CurrentIndex < len(s.Source) && s.Source[s.CurrentIndex] == '='{
            s.CurrentIndex++
            return NewToken("==", EQUAL_EQUAL, nil), nil
        } else {
            return NewToken("=", EQUAL, nil), nil
        }
    case '!':
        if s.CurrentIndex < len(s.Source) && s.Source[s.CurrentIndex] == '='{
            s.CurrentIndex++
            return NewToken("!=", BANG_EQUAL, nil), nil
        } else {
            return NewToken("!", BANG, nil), nil
        }
    default:
        s.ExitCode = 65
        return nil, errors.New(fmt.Sprintf("[line %v] Error: Unexpected character: %v", s.CurrentLine, string(char)))
    }
}
