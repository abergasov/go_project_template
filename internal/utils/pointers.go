package utils

func ToPointer[T any](value T) *T {
	return &value
}

func FromPointer[T any](value *T) T {
	var result T
	if value == nil {
		return result
	}
	v := *value
	return v
}
