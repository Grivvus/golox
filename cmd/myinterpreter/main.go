package main

import (
	"fmt"
	"os"
)

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
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
    
    for _, literal := range fileContents{
        switch literal {
        case '(': fmt.Println("LEFT_PAREN ( null")
        case ')': fmt.Println("RIGHT_PAREN ) null")
        case '{': fmt.Println("LEFT_BRACE { null")
        case '}': fmt.Println("RIGHT_BRACE } null")
        case '*': fmt.Println("STAR * null")
        case '.': fmt.Println("DOT . null")
        case ',': fmt.Println("COMMA , null")
        case '+': fmt.Println("PLUS + null")
        case '-': fmt.Println("MINUS - null")
        case '/': fmt.Println("SLASH / null")
        case ';': fmt.Println("SEMICOLON ; null")
        }
    }
    fmt.Println("EOF  null")
}
