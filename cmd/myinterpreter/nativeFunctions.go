package main

import (
	"time"
)

type loxTime struct{}

func newLoxTime() *loxTime {
	return new(loxTime)
}

func (t loxTime) call(i Interpreter, args []any) any {
	return float64(time.Now().Unix())
}

func (t loxTime) arity() int {
	return 0
}

func (t loxTime) String() string {
	return "<native fn>"
}
