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
