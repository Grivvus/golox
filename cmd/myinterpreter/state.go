package main

import (
	"fmt"
	"os"
)

type State struct {
	values map[string]any
}

func NewState() *State{
    s := new(State)
    s.values = make(map[string]any)
    return s
}

func (s *State) define(name string, value any) {
	s.values[name] = value
}

func (s *State) access(name string) any {
	value, exist := s.values[name]
	if !exist {
		fmt.Fprintln(os.Stderr, "undefined name")
		os.Exit(65)
	}
	return value
}
