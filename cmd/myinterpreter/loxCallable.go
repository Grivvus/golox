package main

type LoxCallable interface {
	arity() int
	call(i Interpreter, args []any) any
}
