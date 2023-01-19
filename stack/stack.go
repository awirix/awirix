package stack

import "github.com/samber/mo"

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
func (s *Stack[T]) Peek() mo.Option[T] {
	if s.length == 0 {
		return mo.None[T]()
	}
	return mo.Some(s.top.value)
}

// Pop the top item of the stack and return it
func (s *Stack[T]) Pop() mo.Option[T] {
	if s.length == 0 {
		return mo.None[T]()
	}

	n := s.top
	s.top = n.prev
	s.length--
	return mo.Some(n.value)
}

// Push a value onto the top of the stack
func (s *Stack[T]) Push(value T) {
	n := &node[T]{value, s.top}
	s.top = n
	s.length++
}
