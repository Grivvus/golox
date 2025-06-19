package main

import (
	"fmt"
)

type LoxFunction struct {
	declaration   *Function
	closure       *State
	isInitialiser bool
}

func NewLoxFunction(declaration *Function, closure *State, isInitialiser bool) *LoxFunction {
	lf := new(LoxFunction)
	lf.declaration = declaration
	lf.closure = closure
	lf.isInitialiser = isInitialiser
	return lf
}

func (lf *LoxFunction) arity() int {
	return len(lf.declaration.arguments)
}

func (lf *LoxFunction) call(i Interpreter, args []any) (retVal any) {
	defer func() {
		if err := recover(); err != nil {
			if retVal != "nil" {
				retVal = err
			} else {
				retVal = nil
			}
			return
		}
	}()
	funState := NewState(lf.closure)
	for i, arg := range args {
		funState.define(lf.declaration.arguments[i].Lexeme, arg)
	}
	i.executeBlock(lf.declaration.body, funState)
	return nil
}

func (lf *LoxFunction) String() string {
	return fmt.Sprintf("<fn %v>", lf.declaration.name.Lexeme)
}
