package utils

import (
	"errors"
	"fmt"
	stdmath "math"
	"sync/atomic"
)

// MustSafeIntToUint32 invokes SafeIntToUint32 and panics if an error is returned.
func MustSafeIntToUint32(i int) uint32 {
	v, err := SafeIntToUint32(i)
	if err != nil {
		panic(err)
	}
	return v
}

// SafeIntToUint32 converts an int to a uint32, returning an error if the value
// is negative or exceeds the maximum uint32 value.
func SafeIntToUint32(i int) (uint32, error) {
	if i < 0 {
		return 0, fmt.Errorf("cannot convert negative int %d to uint32", i)
	}
	if i > stdmath.MaxUint32 {
		return 0, fmt.Errorf("int %d exceeds maximum uint32 value", i)
	}
	return uint32(i), nil
}

// MustSafeUint64ToInt64 invokes SafeUint64ToInt64 and panics if an error is returned.
func MustSafeUint64ToInt64(u uint64) int64 {
	v, err := SafeUint64ToInt64(u)
	if err != nil {
		panic(err)
	}
	return v
}

// SafeUint64ToInt64 converts a uint64 to an int64, returning an error if the
// value is too large to fit in an int64.
func SafeUint64ToInt64(u uint64) (int64, error) {
	if u > stdmath.MaxInt64 {
		return 0, fmt.Errorf("uint64 value %d exceeds maximum int64 value", u)
	}
	return int64(u), nil
}

// SafeAddUint64Ptr adds the values to the counter pointer, and returns an error
// if the counter would overflow.
func SafeAddUint64Ptr(counter *uint64, values ...int) error {
	var totalLen int
	for i := range values {
		if values[i] < 0 {
			return errors.New("value is negative")
		}

		// Check for overflow before adding
		if totalLen > (int(^uint(0)>>1))-values[i] {
			return errors.New("integer overflow detected")
		}

		totalLen += values[i]
	}

	// Convert to uint64 for counter addition
	totalLenUint64 := uint64(totalLen) //nolint:gosec // int to uint64 conversion is safe after overflow check

	// Use compare-and-swap loop to safely check for uint64 overflow
	for {
		current := atomic.LoadUint64(counter)

		// Check if adding totalLen would overflow uint64
		if totalLenUint64 > stdmath.MaxUint64-current {
			return errors.New("uint64 counter overflow")
		}

		newValue := current + totalLenUint64

		// Attempt atomic compare-and-swap
		if atomic.CompareAndSwapUint64(counter, current, newValue) {
			return nil
		}
		// If CAS failed, another goroutine modified counter, retry
	}
}

// Clamp restricts a value to a range [lower, upper].
// If v < lower, returns lower.
// If v > upper, returns upper.
// Otherwise, returns v.
func Clamp(v, lower, upper float64) float64 {
	if v < lower {
		return lower
	}
	if v > upper {
		return upper
	}
	return v
}
