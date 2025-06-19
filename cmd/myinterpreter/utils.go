package main

import (
	"errors"
)

type Stack[T any] struct {
	data   []T
	length int
}

func NewStack[T any]() *Stack[T] {
	return new(Stack[T]).Init()
}

// Sets the stack state to initial.
// Clears the stack and set Len to 0.
func (s *Stack[T]) Init() *Stack[T] {
	s.data = make([]T, 0)
	s.length = 0
	return s
}

func (s *Stack[T]) Push(elem T) {
	if s.length == len(s.data) {
		s.data = append(s.data, elem)
	} else {
		s.data[s.length] = elem
	}
	s.length++
}

func (s *Stack[T]) Pop() (T, error) {
	if s.IsEmpty() {
		var zeroValue T
		return zeroValue, errors.New("Can't Pop from empty stack")
	}
	s.length--
	return s.data[s.length], nil
}

func (s *Stack[T]) GetLast() T {
	return s.data[s.Len()-1]
}

func (s *Stack[T]) Len() int {
	return s.length
}

func (s *Stack[T]) IsEmpty() bool {
	return s.length == 0
}
