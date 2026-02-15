package g

//go:fix inline
func Ptr[T any](value T) *T {
	return new(value)
}
