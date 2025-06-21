package main

type LoxClass struct {
	name       string
	superclass *LoxClass
	methods    map[string]*LoxFunction
}

func NewLoxClass(name string, superclass *LoxClass, methods map[string]*LoxFunction) *LoxClass {
	return &LoxClass{
		name:       name,
		superclass: superclass,
		methods:    methods,
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

func (cls *LoxClass) findMethod(name string) *LoxFunction {
	if method, exist := cls.methods[name]; exist {
		return method
	}
	if cls.superclass != nil {
		return cls.superclass.findMethod(name)
	}
	return nil
}

func (cls *LoxClass) String() string {
	return cls.name
}
