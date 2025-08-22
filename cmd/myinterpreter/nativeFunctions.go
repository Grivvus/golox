package main

import (
	"fmt"
	"time"
)

type nativeFnStringImpl struct{}

func (s nativeFnStringImpl) String() string {
	return "<native fn>"
}

type LoxTime struct {
	nativeFnStringImpl
}

func (t LoxTime) call(i Interpreter, args []any) any {
	return float64(time.Now().Unix())
}

func (t LoxTime) arity() int {
	return 0
}

type Floor struct {
	nativeFnStringImpl
}

func (f Floor) call(i Interpreter, args []any) any {
	switch arg := args[0].(type) {
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

type Str struct {
	nativeFnStringImpl
}

func (s Str) call(i Interpreter, args []any) any {
	return fmt.Sprintf("%v", args[0])
}

func (s Str) arity() int {
	return 1
}

type Len struct {
	nativeFnStringImpl
}

func (l Len) call(i Interpreter, args []any) any {
	switch arr := args[0].(type) {
	case []any:
		return float64(len(arr))
	default:
		i.error(i.parser.getCurrent(), "Only arrays have len")
	}
	panic("unreachable")
}

func (l Len) arity() int {
	return 1
}

type PrintLine struct {
	nativeFnStringImpl
}

func (p PrintLine) call(i Interpreter, args []any) any {
	fmt.Println(args[0])
	return nil
}

func (p PrintLine) arity() int {
	return 1
}
