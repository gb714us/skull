package stack

import "errors"

var (
	ErrEmpty = errors.New("stack empty")
)

type Stack[T any] struct {
	arr []T
}

func New[T any]() *Stack[T] {
	return &Stack[T]{
		arr: make([]T, 0),
	}
}

func (s *Stack[T]) Push(t T) {
	s.arr = append(s.arr, t)
}

func (s *Stack[T]) Pop() (t T, err error) {
	if len(s.arr) == 0 {
		err = ErrEmpty
		return
	}

	t = s.arr[len(s.arr)-1]
	s.arr = s.arr[:len(s.arr)-1]
	return
}

func (s *Stack[T]) Size() int {
	return len(s.arr)
}

func (s *Stack[T]) Clear() {
	s.arr = s.arr[:0]
}
