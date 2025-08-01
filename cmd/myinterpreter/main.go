package main

import (
	"bytes"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "Usage: ./your_program.sh tokenize <filename>")
		os.Exit(1)
	}

	command := os.Args[1]

	if command != "tokenize" && command != "parse" && command != "evaluate" && command != "run" {
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
		os.Exit(1)
	}

	filename := os.Args[2]
	fileContents, err := os.ReadFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}

	scanner := NewScanner(bytes.Runes(fileContents))
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

	if command == "tokenize" {
		for _, value := range tokens {
			fmt.Println(value.String())
		}
		os.Exit(scanner.ExitCode)
	} else if command == "parse" {
		if scanner.ExitCode != 0 {
			os.Exit(scanner.ExitCode)
		}
		parser := NewParser(tokens)
		parser.parseExprs()
		printer := NewPrinter()
		for _, v := range parser.errs {
			fmt.Fprintln(os.Stderr, v)
		}
		for _, v := range parser.exprs {
			fmt.Println(v.print(printer))
		}
		os.Exit(parser.exitCode)
	} else if command == "evaluate" {
		parser := NewParser(tokens)
		exprs := parser.parseExprs()
		interp := NewInterpreter(parser)
		var res []any
		for _, expr := range exprs {
			res = append(res, expr.accept(interp))
		}
		for _, v := range res {
			if v == nil {
				fmt.Println("nil")
			} else {
				fmt.Println(v)
			}
		}
	} else if command == "run" {
		parser := NewParser(tokens)
		interp := NewInterpreter(parser)
		resolver := NewResolver(interp)
		stmts := parser.parseStmts()
		resolver.resolveStmts(stmts)
		for _, stmt := range stmts {
			interp.execute(stmt)
		}
	}
}
