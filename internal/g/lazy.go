package g

import "reflect"

func Lazy[T any](fn func() T) func() T {
	var result T
	return func() T {
		if reflect.TypeOf(result) == nil {
			result = fn()
		}
		return result
	}
}
