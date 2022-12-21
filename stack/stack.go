package stack

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

// Peek vieww the top item on the stack
func (s *Stack[T]) Peek() T {
	if s.length == 0 {
		var t T
		return t
	}
	return s.top.value
}

// Pop the top item of the stack and return it
func (s *Stack[T]) Pop() T {
	if s.length == 0 {
		var t T
		return t
	}

	n := s.top
	s.top = n.prev
	s.length--
	return n.value
}

// Push a value onto the top of the stack
func (s *Stack[T]) Push(value T) {
	n := &node[T]{value, s.top}
	s.top = n
	s.length++
}
