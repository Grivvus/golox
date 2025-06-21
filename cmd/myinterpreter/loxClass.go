package main

type LoxClass struct {
	name    string
	methods map[string]*LoxFunction
}

func NewLoxClass(name string, methods map[string]*LoxFunction) *LoxClass {
	return &LoxClass{
		name:    name,
		methods: methods,
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
