package main

import (
	"fmt"
	"os"
)

type State struct {
    enclosing *State
	values map[string]any
}

func NewState(enclosing *State) *State{
    s := new(State)
    if enclosing != nil {
        s.enclosing = enclosing
    }
    s.values = make(map[string]any)
    return s
}

func (s *State) assign(name string, value any){
    _, exist := s.values[name]
    if !exist {
        if s.enclosing != nil {
            s.enclosing.assign(name, value)
        } else {
            fmt.Fprintln(os.Stderr, "cannot assign unexisted variable")
            os.Exit(70)
        }
    } else {
        s.values[name] = value
    }
}

func (s *State) define(name string, value any) {
	s.values[name] = value
}

func (s *State) access(name string) any {
	value, exist := s.values[name]
	if !exist {
        if s.enclosing != nil {
            return s.enclosing.access(name)
        }
		fmt.Fprintln(os.Stderr, "undefined variable")
		os.Exit(70)
	}
	return value
}
