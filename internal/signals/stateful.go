package signals

import (
	"sync"
)

type StatefulSignal[T any] struct {
	Signal[func(T) bool]
	currentValue T
	lock         sync.RWMutex
}

func (s *StatefulSignal[T]) CurrentValue() T {
	s.lock.RLock()
	defer s.lock.RUnlock()

	return s.currentValue
}

func (s *StatefulSignal[T]) Notify(callback func(oldValue T) T) {
	s.lock.Lock()
	defer s.lock.Unlock()

	oldValue := s.currentValue
	s.currentValue = callback(oldValue)
	s.Signal.Notify(s.currentValue)
}

func (s *StatefulSignal[T]) On(handler func(T) bool) *Subscription {
	s.lock.RLock()
	defer s.lock.RUnlock()

	handler(s.currentValue)
	return s.Signal.On(handler)
}

func NewStatefulSignal[T any](initialValue T) *StatefulSignal[T] {
	return &StatefulSignal[T]{
		currentValue: initialValue,
		lock:         sync.RWMutex{},
		Signal:       NewSignal[func(T) bool](),
	}
}
