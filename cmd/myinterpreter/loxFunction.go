package main

import (

)

type LoxFunction struct {
    declaration Function
    closure State
    isInitialiser bool
}

func NewLoxFunction(declaration Function, closure State, isInitialiser bool) *LoxFunction{
    lf := new(LoxFunction)
    lf.declaration = declaration
    lf.closure = closure
    lf.isInitialiser = isInitialiser
    return lf
}

func (lf *LoxFunction) arity() int{
    return len(lf.declaration.arguments)
}

func (lf *LoxFunction) call(i Interpreter, args []any) any {
    funState := NewState(&lf.closure)
    for i, arg := range args {
        funState.define(lf.declaration.arguments[i].Lexeme, arg)
    }
    i.executeBlock(lf.declaration.body, funState)
    return 0;
}
