package utils

import (
	"sync/atomic"
)

type RoundRobinBalancer[T any] struct {
	values []T
	index  *atomic.Int64
}

func NewRoundRobinBalancer[T any](values []T) *RoundRobinBalancer[T] {
	return &RoundRobinBalancer[T]{
		values: values,
		index:  &atomic.Int64{},
	}
}

func (wb *RoundRobinBalancer[T]) Next() T {
	return wb.values[(wb.index.Add(1)-1)%int64(len(wb.values))]
}

func (wb *RoundRobinBalancer[T]) Values() []T {
	return wb.values
}
