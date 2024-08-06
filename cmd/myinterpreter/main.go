package main

import (
	"fmt"
	"os"
)

func main() {
	// You can use print statements as follows for debugging, they"ll be visible when running tests.
	fmt.Fprintln(os.Stderr, "Logs from your program will appear here!")

	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "Usage: ./your_program.sh tokenize <filename>")
		os.Exit(1)
	}

	command := os.Args[1]

	if command != "tokenize" {
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
		os.Exit(1)
	}

	// Uncomment this block to pass the first stage
	//
	filename := os.Args[2]
	fileContents, err := os.ReadFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}
    var line uint = 1
    var tokens []Token
    var isLexicalError bool
    for _, literal := range fileContents{
        if literal == '\n'{
            line++
        } else {
            token, err := newToken(string(literal), line)
            if err == nil{
                tokens = append(tokens, *token)
            } else {
                isLexicalError = true
                fmt.Fprintln(os.Stderr, err)
            }
        }
    }
    for _, token := range tokens{
        fmt.Println(token.toStr())
    }
    token, err := newToken("EOF", 0)
    if err == nil {
        fmt.Println(token.toStr())
    }
    if isLexicalError{
        os.Exit(65)
    }
}

type Token struct{
    lexeme string
    tokenName string
    literal any
    line uint
}

func newToken(lexeme string, line uint) (*Token, error){
    t := new(Token)
    t.line = line
    t.literal = nil
    switch lexeme{
    case "(": t.lexeme = "("; t.tokenName = "LEFT_PAREN"
    case ")": t.lexeme = ")"; t.tokenName = "RIGHT_PAREN"
    case "{": t.lexeme = "{"; t.tokenName = "LEFT_BRACE"
    case "}": t.lexeme = "}"; t.tokenName = "RIGHT_BRACE"
    case "*": t.lexeme = "*"; t.tokenName = "STAR"
    case ".": t.lexeme = "."; t.tokenName = "DOT"
    case ",": t.lexeme = ","; t.tokenName = "COMMA"
    case "+": t.lexeme = "+"; t.tokenName = "PLUS"
    case "-": t.lexeme = "-"; t.tokenName = "MINUS"
    case "/": t.lexeme = "/"; t.tokenName = "SLASH"
    case ";": t.lexeme = ";"; t.tokenName = "SEMICOLON"
    case "EOF": t.lexeme = ""; t.tokenName = "EOF"
    }
    if t.tokenName != ""{
        return t, nil
    } else {
        return nil, fmt.Errorf("[line %v] Error: Unexpected character: %v", t.line, lexeme)
    }
}

func (t Token)toStr() string{
    if (t.literal == nil){
        return t.tokenName + " " + t.lexeme + " " + "null";
    }
    return t.tokenName + " " + t.lexeme
}
