package utils

func FromPointer[T any](value *T) T {
	var result T
	if value == nil {
		return result
	}
	v := *value
	return v
}
