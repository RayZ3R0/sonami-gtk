package signals

import (
	"reflect"
	"sync"

	"github.com/google/uuid"
)

const (
	Continue    = false
	Unsubscribe = true
)

type Subscription uuid.UUID

type Signal[T any] struct {
	mutex    sync.Mutex
	handlers map[*Subscription]T
}

func (b *Signal[T]) addHandler(handler T) *Subscription {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	sub := Subscription(uuid.New())
	b.handlers[&sub] = handler
	return &sub
}

func (b *Signal[T]) removeHandler(sub *Subscription) {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	delete(b.handlers, sub)
}

func (b *Signal[T]) Notify(args ...any) {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	for sub, handler := range b.handlers {
		reflectArgs := make([]reflect.Value, len(args))
		for i, arg := range args {
			reflectArgs[i] = reflect.ValueOf(arg)
		}
		result := reflect.ValueOf(handler).Call(reflectArgs)
		if len(result) > 0 && result[0].CanConvert(reflect.TypeFor[bool]()) && result[0].Bool() {
			b.removeHandler(sub)
		}
	}
}

func (b *Signal[T]) On(handler T) *Subscription {
	return b.addHandler(handler)
}

func (b *Signal[T]) Unsubscribe(sub *Subscription) {
	b.removeHandler(sub)
}

func NewSignal[T any]() Signal[T] {
	return Signal[T]{
		handlers: make(map[*Subscription]T),
	}
}
