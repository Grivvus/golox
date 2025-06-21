package main

import (
	"fmt"
	"os"
)

type LoxInstance struct {
	cls    *LoxClass
	fields map[string]any
}

func NewLoxInstance(cls *LoxClass) *LoxInstance {
	return &LoxInstance{
		cls:    cls,
		fields: make(map[string]any, 0),
	}
}

func (instance *LoxInstance) String() string {
	return fmt.Sprintf("%v instance", instance.cls.String())
}

func (instance *LoxInstance) Get(name Token) any {
	value, ok := instance.fields[name.Lexeme]
	if ok {
		return value
	}
	method, ok := instance.cls.methods[name.Lexeme]
	if ok {
		return method
	}
	fmt.Fprintf(os.Stderr, "[line %v] Undefined property '%v'", name.line, name.Lexeme)
	os.Exit(70)
	return nil
}

func (instance *LoxInstance) Set(name Token, value any) {
	instance.fields[name.Lexeme] = value
}
