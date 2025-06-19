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
	s.error("can't assign to variable that didnt exist")
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
		s.error(fmt.Sprintf("undefined variable '%v'", name))
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
	return env
}

func (s *State) accessAt(distance int, name string) any {
	return s.ancestor(distance).access(name)
}

func (s *State) assignAt(distance int, name string, value any) {
	s.ancestor(distance).assign(name, value)
}

func (s State) error(msg string) {
	fmt.Fprintln(os.Stderr, msg)
	os.Exit(70)
}
