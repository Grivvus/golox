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
	if initializer, exist := cls.methods["init"]; exist {
		return initializer.arity()
	}
	return 0
}

func (cls *LoxClass) call(i Interpreter, args []any) any {
	instance := NewLoxInstance(cls)
	initializer, exist := cls.methods["init"]
	if exist {
		initializer.bind(instance).call(i, args)
	}
	return instance
}

func (cls *LoxClass) String() string {
	return cls.name
}
