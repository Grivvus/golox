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

    if s.CurrentIndex < (len(s.Source) - 1) && s.Source[s.CurrentIndex] == '/' && s.Source[s.CurrentIndex + 1] == '/'{
        for s.CurrentIndex < len(s.Source) && s.Source[s.CurrentIndex] != '\n'{
            s.CurrentIndex++
        }
        s.CurrentLine++
        s.CurrentIndex++
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
    case '<':
        if s.CurrentIndex < len(s.Source) && s.Source[s.CurrentIndex] == '='{
            s.CurrentIndex++
            return NewToken("<=", LESS_EQUAL, nil), nil
        } else {
            return NewToken("<", LESS, nil), nil
        }
    case '>':
        if s.CurrentIndex < len(s.Source) && s.Source[s.CurrentIndex] == '='{
            s.CurrentIndex++
            return NewToken(">=", GREATER_EQUAL, nil), nil
        } else {
            return NewToken(">", GREATER, nil), nil
        }
    case '"':
        var res_str []rune;
        for s.CurrentIndex < len(s.Source) && s.Source[s.CurrentIndex] != '"' && s.Source[s.CurrentIndex] != '\n'{
            res_str = append(res_str, rune(s.Source[s.CurrentIndex]))
            s.CurrentIndex++
        }

        if s.CurrentIndex < len(s.Source) && s.Source[s.CurrentIndex] == '"'{
            s.CurrentIndex++
            return NewToken("\"" + string(res_str) + "\"", STRING, string(res_str)), nil
        } else {
            s.ExitCode = 65
            return nil, errors.New(fmt.Sprintf("[line %v] Error: Unterminated string.", s.CurrentLine))
        }
    case '\t':
        return nil, nil
    case ' ':
        return nil, nil
    case 10:
        // line feed ASCII code 10
        return nil, nil
    default:
        if isDigit(char){
            var digits []byte
            var isFloat bool = false
            
            digits = append(digits, char)
            for s.CurrentIndex < len(s.Source) && (isDigit(s.Source[s.CurrentIndex]) || (s.CurrentIndex < len(s.Source) - 1 && s.Source[s.CurrentIndex] == '.' && isDigit(s.Source[s.CurrentIndex + 1]))){
                if s.Source[s.CurrentIndex] == '.' && isFloat == false{
                    isFloat = true
                } else if s.Source[s.CurrentIndex] == '.'{
                    break
                }
                digits = append(digits, s.Source[s.CurrentIndex])
                s.CurrentIndex++
            }
            
            numLiteral := string(digits)
            if !isFloat {
                digits = append(digits, '.')
                digits = append(digits, '0')
            }
            return NewToken(numLiteral, NUMBER, string(digits)), nil
        }
        s.ExitCode = 65
        return nil, errors.New(fmt.Sprintf("[line %v] Error: Unexpected character: %v", s.CurrentLine, string(char)))
    }
}

func isDigit(char byte) bool {
    if char >= '0' && char <= '9'{
        return true
    }
    return false
}
