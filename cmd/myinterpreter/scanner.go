package main

import (
	"errors"
	"fmt"
	"strconv"
)

var reservedWords map[string]TokenType

type Scanner struct {
	Source       []rune
	ExitCode     int
	CurrentIndex int
	CurrentLine  uint
}

func NewScanner(source []rune) *Scanner {
	scanner := new(Scanner)
	scanner.Source = source
	scanner.CurrentIndex = 0
	scanner.CurrentLine = 1
	scanner.ExitCode = 0
	reservedWords = *fillMap()

	return scanner
}

func (s *Scanner) NextToken() (*Token, error) {
	if s.CurrentIndex < len(s.Source) && s.Source[s.CurrentIndex] == '\n' {
		s.CurrentIndex++
		s.CurrentLine++
	}

	if s.CurrentIndex < (len(s.Source)-1) && s.Source[s.CurrentIndex] == '/' && s.Source[s.CurrentIndex+1] == '/' {
		for s.CurrentIndex < len(s.Source) && s.Source[s.CurrentIndex] != '\n' {
			s.CurrentIndex++
		}
		s.CurrentLine++
		s.CurrentIndex++
	}

	if s.CurrentIndex >= len(s.Source) {
		return NewToken("", EOF, nil, s.CurrentLine), nil
	}

	char := s.Source[s.CurrentIndex]
	s.CurrentIndex++

	switch char {
	case '(':
		return NewToken("(", LEFT_PAREN, nil, s.CurrentLine), nil
	case ')':
		return NewToken(")", RIGHT_PAREN, nil, s.CurrentLine), nil
	case '{':
		return NewToken("{", LEFT_BRACE, nil, s.CurrentLine), nil
	case '}':
		return NewToken("}", RIGHT_BRACE, nil, s.CurrentLine), nil
	case '*':
		return NewToken("*", STAR, nil, s.CurrentLine), nil
	case '.':
		return NewToken(".", DOT, nil, s.CurrentLine), nil
	case ',':
		return NewToken(",", COMMA, nil, s.CurrentLine), nil
	case '+':
		return NewToken("+", PLUS, nil, s.CurrentLine), nil
	case '-':
		return NewToken("-", MINUS, nil, s.CurrentLine), nil
	case '/':
		return NewToken("/", SLASH, nil, s.CurrentLine), nil
	case ';':
		return NewToken(";", SEMICOLON, nil, s.CurrentLine), nil
	case '=':
		if s.CurrentIndex < len(s.Source) && s.Source[s.CurrentIndex] == '=' {
			s.CurrentIndex++
			return NewToken("==", EQUAL_EQUAL, nil, s.CurrentLine), nil
		} else {
			return NewToken("=", EQUAL, nil, s.CurrentLine), nil
		}
	case '!':
		if s.CurrentIndex < len(s.Source) && s.Source[s.CurrentIndex] == '=' {
			s.CurrentIndex++
			return NewToken("!=", BANG_EQUAL, nil, s.CurrentLine), nil
		} else {
			return NewToken("!", BANG, nil, s.CurrentLine), nil
		}
	case '<':
		if s.CurrentIndex < len(s.Source) && s.Source[s.CurrentIndex] == '=' {
			s.CurrentIndex++
			return NewToken("<=", LESS_EQUAL, nil, s.CurrentLine), nil
		} else {
			return NewToken("<", LESS, nil, s.CurrentLine), nil
		}
	case '>':
		if s.CurrentIndex < len(s.Source) && s.Source[s.CurrentIndex] == '=' {
			s.CurrentIndex++
			return NewToken(">=", GREATER_EQUAL, nil, s.CurrentLine), nil
		} else {
			return NewToken(">", GREATER, nil, s.CurrentLine), nil
		}
	case '"':
		var res_str []rune
		for s.CurrentIndex < len(s.Source) && s.Source[s.CurrentIndex] != '"'{
			res_str = append(res_str, rune(s.Source[s.CurrentIndex]))
			s.CurrentIndex++
		}

		if s.CurrentIndex < len(s.Source) && s.Source[s.CurrentIndex] == '"' {
			s.CurrentIndex++
			return NewToken("\""+string(res_str)+"\"", STRING, string(res_str), s.CurrentLine), nil
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
		if isDigit(char) {
			var digits []rune
			var isFloat bool = false

			digits = append(digits, char)
			for s.CurrentIndex < len(s.Source) && (isDigit(s.Source[s.CurrentIndex]) || (s.CurrentIndex < len(s.Source)-1 && s.Source[s.CurrentIndex] == '.' && isDigit(s.Source[s.CurrentIndex+1]))) {
				if s.Source[s.CurrentIndex] == '.' && isFloat == false {
					isFloat = true
				} else if s.Source[s.CurrentIndex] == '.' {
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
			parsed, err := strconv.ParseFloat(string(digits), 64)
			if err != nil {
				return nil, errors.New("Can't parse NUMBER")
			}
			return NewToken(numLiteral, NUMBER, parsed, s.CurrentLine), nil
		} else if isAlpha(char) {
			var identifier []rune
			identifier = append(identifier, char)

			for s.CurrentIndex < len(s.Source) && isAlphaOrDigit(s.Source[s.CurrentIndex]) {
				identifier = append(identifier, s.Source[s.CurrentIndex])
				s.CurrentIndex++
			}

			reserved := reservedWords[string(identifier)]
			if reserved == EOF {
				return NewToken(string(identifier), IDENTIFIER, nil, s.CurrentLine), nil
			}
			return NewToken(string(identifier), reserved, nil, s.CurrentLine), nil

		}
		s.ExitCode = 65
		return nil, errors.New(fmt.Sprintf("[line %v] Error: Unexpected character: %v", s.CurrentLine, string(char)))
	}
}

func isDigit(char rune) bool {
	if char >= '0' && char <= '9' {
		return true
	}
	return false
}

func isAlpha(char rune) bool {
	if char == '_' || (char >= 'A' && char <= 'Z') || (char >= 'a' && char <= 'z') {
		return true
	}
	return false
}

func isAlphaOrDigit(char rune) bool {
	return isDigit(char) || isAlpha(char)
}
