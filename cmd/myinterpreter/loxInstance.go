package main

import "fmt"

type LoxInstance struct {
	cls *LoxClass
}

func NewLoxInstance(cls *LoxClass) *LoxInstance {
	return &LoxInstance{
		cls: cls,
	}
}

func (instance *LoxInstance) String() string {
	return fmt.Sprintf("%v instance", instance.cls.String())
}
