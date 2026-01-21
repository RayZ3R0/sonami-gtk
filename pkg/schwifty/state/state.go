package state

import (
	"sync"

	"github.com/google/uuid"
)

type StateCallback[T any] func(newValue T)

type State[T any] struct {
	value         T
	callbacks     map[string]StateCallback[T]
	callbacksLock sync.RWMutex
	stateful      bool
}

func (s *State[T]) AddCallback(callback StateCallback[T]) string {
	id := uuid.NewString()
	s.callbacksLock.Lock()
	s.callbacks[id] = callback
	s.callbacksLock.Unlock()
	if s.stateful {
		callback(s.value)
	}
	return id
}

func (s *State[T]) RemoveCallback(id string) {
	s.callbacksLock.Lock()
	delete(s.callbacks, id)
	s.callbacksLock.Unlock()
}

func (s *State[T]) SetValue(value T) {
	s.callbacksLock.RLock()
	for _, callback := range s.callbacks {
		callback(value)
	}
	s.callbacksLock.RUnlock()
	s.value = value
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
