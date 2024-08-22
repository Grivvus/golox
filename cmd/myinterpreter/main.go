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

	if command != "tokenize" && command != "parse"{
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

	scanner := NewScanner(fileContents)
	var tokens []Token
	token, err := scanner.NextToken()
	for token == nil || token.Token != EOF {
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
		} else {
			if token != nil {
				tokens = append(tokens, *token)
			}
		}
		token, err = scanner.NextToken()
	}
    tokens = append(tokens, *token)

    if command == "tokenize"{
	    for _, value := range tokens {
		    fmt.Println(value.String())
	    }
	} else if command == "parse"{
	    parser := NewParser(tokens)
        expressions := make([]Expr, 0, 1)
        printer := new(astPrinter)
        expr, err := parser.nextExpr()
        for err == nil {
            expressions = append(expressions, expr)
            expr, err = parser.nextExpr()
        }
        if err.Error() != "EOF"{
            fmt.Println(err)
            os.Exit(-1)
        }
        for _, v := range expressions{
            printer.print(v)
        }
	}

	os.Exit(scanner.ExitCode)
}
