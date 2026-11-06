package main

import (
	"bytes"
	"fmt"
	"golox/internal/astprinter"
	"golox/internal/interpreter"
	"golox/internal/parser"
	"golox/internal/scanner"
	"os"
)

func main() {
	exitCode := 0
	exitMessage := "interpreted without errors"
	defer func() {
		fmt.Fprintln(os.Stderr, exitMessage)
		os.Exit(exitCode)
	}()
	if len(os.Args) < 3 {
		exitMessage = "Usage: ./your_program.sh tokenize <filename>"
		exitCode = 1
		return
	}

	command := os.Args[1]

	if command != "tokenize" && command != "parse" && command != "evaluate" && command != "run" {
		exitMessage = fmt.Sprintf("Unknown command: %s", command)
		exitCode = 1
		return
	}

	filename := os.Args[2]
	fileContents, err := os.ReadFile(filename)
	if err != nil {
		exitMessage = fmt.Sprintf("Error reading file: %v\n", err)
		exitCode = 1
		return
	}

	scanner := scanner.New(bytes.Runes(fileContents))
	tokens, err := scanner.Scan()
	if err != nil {
		exitCode = 70
		exitMessage = err.Error()
		return
	}

	if command == "tokenize" {
		for _, value := range tokens {
			fmt.Println(value.String())
		}
		exitCode = scanner.ExitCode
	} else if command == "parse" {
		if scanner.ExitCode != 0 {
			exitCode = scanner.ExitCode
		}
		parser := parser.New(tokens)
		parser.ParseExprs()
		printer := astprinter.New()
		for _, v := range parser.Errs {
			fmt.Fprintln(os.Stderr, v)
		}
		for _, v := range parser.Exprs {
			fmt.Println(v.Print(printer))
		}
		exitCode = parser.ExitCode
	} else if command == "evaluate" {
		parser := parser.New(tokens)
		exprs := parser.ParseExprs()
		interp := interpreter.New(parser)
		var res []any
		for _, expr := range exprs {
			res = append(res, expr.Accept(interp))
		}
		for _, v := range res {
			if v == nil {
				fmt.Println("nil")
			} else {
				fmt.Println(v)
			}
		}
	} else if command == "run" {
		parser := parser.New(tokens)
		interp := interpreter.New(parser)
		resolver := interpreter.NewResolver(interp)
		stmts := parser.ParseStmts()
		resolver.ResolveStmts(stmts)
		for _, stmt := range stmts {
			interp.Execute(stmt)
		}
	}
}
