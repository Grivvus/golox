package main

import (
	"time"
)

type LoxTime struct{}

func NewLoxTime() *LoxTime {
	return new(LoxTime)
}

func (t LoxTime) call(i Interpreter, args []any) any {
	return float64(time.Now().Unix())
}

func (t LoxTime) arity() int {
	return 0
}

func (t LoxTime) String() string {
	return "<native fn>"
}

type Floor struct{}

func NewFloor() *Floor {
	return &Floor{}
}

func (f Floor) call(i Interpreter, args []any) any {
	if len(args) == 0 || len(args) > 1 {
		i.error(i.parser.getCurrent(), "Expect only 1 argument")
	}
	arg := args[0]
	switch arg := arg.(type) {
	case float64:
		return float64(int64(arg))
	default:
		i.error(i.parser.getCurrent(), "Argument should be a number")
	}
	panic("unreachable")
}

func (f Floor) arity() int {
	return 1
}

func (f Floor) String() string {
	return "<native fn>"
}
