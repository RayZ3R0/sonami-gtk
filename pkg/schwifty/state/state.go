package state

import (
	"github.com/google/uuid"
)

type StateCallback[T any] func(newValue T)

type State[T any] struct {
	value     T
	callbacks map[string]StateCallback[T]
	stateful  bool
}

func (s *State[T]) AddCallback(callback StateCallback[T]) string {
	id := uuid.NewString()
	s.callbacks[id] = callback
	if s.stateful {
		callback(s.value)
	}
	return id
}

func (s *State[T]) RemoveCallback(id string) {
	delete(s.callbacks, id)
}

func (s *State[T]) SetValue(value T) {
	s.value = value
	for _, callback := range s.callbacks {
		callback(value)
	}
}

func (s *State[T]) Value() T {
	return s.value
}

func New[T any](value T) *State[T] {
	return &State[T]{
		value:     value,
		callbacks: make(map[string]StateCallback[T]),
	}
}

func NewStateful[T any](value T) *State[T] {
	return &State[T]{
		value:     value,
		callbacks: make(map[string]StateCallback[T]),
		stateful:  true,
	}
}
