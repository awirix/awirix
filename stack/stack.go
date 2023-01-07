package stack

import "github.com/awirix/awirix/option"

type (
	Stack[T any] struct {
		top    *node[T]
		length int
	}
	node[T any] struct {
		value T
		prev  *node[T]
	}
)

func New[T any]() *Stack[T] {
	return &Stack[T]{length: 0}
}

// Len returns the number of items on the stack
func (s *Stack[T]) Len() int {
	return s.length
}

// Peek view the top item on the stack
func (s *Stack[T]) Peek() *option.Option[T] {
	if s.length == 0 {
		return option.None[T]()
	}
	return option.Some(s.top.value)
}

// Pop the top item of the stack and return it
func (s *Stack[T]) Pop() *option.Option[T] {
	if s.length == 0 {
		return option.None[T]()
	}

	n := s.top
	s.top = n.prev
	s.length--
	return option.Some(n.value)
}

// Push a value onto the top of the stack
func (s *Stack[T]) Push(value T) {
	n := &node[T]{value, s.top}
	s.top = n
	s.length++
}
