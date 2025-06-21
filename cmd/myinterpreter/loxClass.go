package main

type LoxClass struct {
	name string
}

func NewLoxClass(name string) *LoxClass {
	return &LoxClass{
		name: name,
	}
}

func (cls *LoxClass) arity() int {
	return 0
}

func (cls *LoxClass) call(i Interpreter, args []any) any {
	instance := NewLoxInstance(cls)
	return instance
}

func (cls *LoxClass) String() string {
	return cls.name
}
