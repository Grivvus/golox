package main

import (
	"fmt"
	"os"
)

type State struct {
	enclosing *State
	values    map[string]any
}

func NewState(enclosing *State) *State {
	return &State{
		enclosing: enclosing,
		values:    make(map[string]any),
	}
}

func (s *State) assign(name string, value any) {
	if _, exist := s.values[name]; exist {
		s.values[name] = value
		return
	}
	if s.enclosing != nil {
		s.enclosing.assign(name, value)
		return
	}
	fmt.Fprintln(os.Stderr, "cannot assign unexisted variable")
	os.Exit(70)
}

func (s *State) define(name string, value any) {
	// fmt.Printf("DEFINE: %s = %v in state %p\n", name, value, s)
	s.values[name] = value
}

func (s *State) access(name string) any {
	value, exist := s.values[name]
	if !exist {
		if s.enclosing != nil {
			return s.enclosing.access(name)
		}
		fmt.Fprintf(os.Stderr, "undefined variable '%v'", name)
		os.Exit(70)
	}
	return value
}

func (s *State) ancestor(distance int) *State {
	env := s
	for range distance {
		if env.enclosing == nil {
			break
		}
		env = env.enclosing
	}
	// fmt.Printf("found ancestor with depth %v, addresss is %p\n", distance, env)
	return env
}

func (s *State) accessAt(distance int, name string) any {
	return s.ancestor(distance).access(name)
}

func (s *State) assignAt(distance int, name string, value any) {
	s.ancestor(distance).assign(name, value)
}
