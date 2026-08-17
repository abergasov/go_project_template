package utils_test

import (
	"math"
	"testing"

	"go_project_template/internal/utils"

	"github.com/stretchr/testify/require"
)

// FuzzMathConversions tests safe integer conversion functions
func FuzzMathConversions(f *testing.F) {
	f.Add(int64(0), uint64(0))
	f.Add(int64(math.MaxInt32), uint64(math.MaxUint32))
	f.Add(int64(-1), uint64(math.MaxInt64))
	f.Add(int64(math.MaxInt64), uint64(math.MaxUint64))
	f.Add(int64(math.MinInt64), uint64(0))

	f.Fuzz(func(t *testing.T, i int64, u uint64) {
		// Test SafeIntToUint32
		val32, err32 := utils.SafeIntToUint32(int(i))
		if err32 == nil {
			require.False(t, int(i) < 0 || int(i) > math.MaxUint32, "SafeIntToUint32 should have failed for %d", i)
			require.Equal(t, uint32(i), val32, "SafeIntToUint32 returned wrong value") //nolint:gosec // safe: conversion guarded by SafeIntToUint32 success
		}

		// Test SafeUint64ToInt64
		val64, err64 := utils.SafeUint64ToInt64(u)
		if err64 == nil {
			require.False(t, u > math.MaxInt64, "SafeUint64ToInt64 should have failed for %d", u)
			require.Equal(t, int64(u), val64, "SafeUint64ToInt64 returned wrong value") //nolint:gosec // safe: conversion guarded by SafeUint64ToInt64 success
		}
	})
}

// FuzzSafeAddUint64Ptr tests atomic counter addition with overflow protection
func FuzzSafeAddUint64Ptr(f *testing.F) {
	f.Add(uint64(0), int(100), int(200))
	f.Add(uint64(math.MaxUint64-100), int(50), int(50))
	f.Add(uint64(math.MaxUint64-100), int(50), int(51))
	f.Add(uint64(0), int(-1), int(0))
	f.Add(uint64(1000), int(0), int(0))

	f.Fuzz(func(t *testing.T, initial uint64, v1, v2 int) {
		counter := initial
		err := utils.SafeAddUint64Ptr(&counter, v1, v2)

		if v1 < 0 || v2 < 0 {
			require.Error(t, err, "should fail for negative values")
			return
		}

		sum := int64(v1) + int64(v2)
		if sum < 0 || uint64(sum) > math.MaxUint64-initial {
			require.Error(t, err, "should fail for overflow")
		} else if err == nil {
			expected := initial + uint64(v1) + uint64(v2)
			require.Equal(t, expected, counter, "expected %d, got %d", expected, counter)
		}
	})
}
