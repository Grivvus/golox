package main

import (
	"fmt"
	"os"
    "github.com/codecrafters-io/interpreter-starter-go/cmd/myinterpreter/lexer"
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

	filename := os.Args[2]
	fileContents, err := os.ReadFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}

    scanner := lexer.NewScanner(fileContents)
    var tokens []lexer.Token
    token, err := scanner.NextToken()
    for token == nil || token.Token != lexer.EOF{
        if err != nil{
            fmt.Fprintln(os.Stderr, err.Error())
        } else {
            tokens = append(tokens, *token)
        }
        token, err = scanner.NextToken()
    }
    tokens = append(tokens, *token)

    for _, value := range tokens{
        fmt.Println(value.String())
    }

    os.Exit(scanner.ExitCode)
}
