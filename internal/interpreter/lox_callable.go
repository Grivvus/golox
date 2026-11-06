package interpreter

type LoxCallable interface {
	arity() int
	call(i Interpreter, args []any) any
}
