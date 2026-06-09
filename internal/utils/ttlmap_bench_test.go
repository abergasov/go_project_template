package utils_test

import (
	"strconv"
	"testing"
	"time"

	"go_project_template/internal/utils"
)

// newBenchMap returns a TTL map pre-populated with n string→int entries.
func newBenchMap(b *testing.B, n int) *utils.TTLMap[string, int] {
	b.Helper()
	m := utils.NewTTLMap[string, int](time.Hour, time.Hour)
	for i := range n {
		m.Put("key"+strconv.Itoa(i), i)
	}
	return m
}

// BenchmarkPut measures pure write throughput.
func BenchmarkPut(b *testing.B) {
	m := utils.NewTTLMap[string, int](time.Hour, time.Hour)

	b.ResetTimer()
	for i := range b.N {
		m.Put("key"+strconv.Itoa(i%10_000), i)
	}
}

// BenchmarkPutParallel measures concurrent write throughput.
func BenchmarkPutParallel(b *testing.B) {
	m := utils.NewTTLMap[string, int](time.Hour, time.Hour)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			m.Put("key"+strconv.Itoa(i%10_000), i)
			i++
		}
	})
}

// BenchmarkGet measures pure read throughput on a warm cache.
func BenchmarkGet(b *testing.B) {
	const n = 10_000
	m := newBenchMap(b, n)
	b.ResetTimer()
	for i := range b.N {
		m.Get("key" + strconv.Itoa(i%n))
	}
}

// BenchmarkGetParallel measures concurrent read throughput — this is where
// sharding and the RLock hot-path show the biggest gains.
func BenchmarkGetParallel(b *testing.B) {
	const n = 10_000
	m := newBenchMap(b, n)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			m.Get("key" + strconv.Itoa(i%n))
			i++
		}
	})
}

// BenchmarkMixed80R20W simulates a realistic read-heavy workload.
func BenchmarkMixed80R20W(b *testing.B) {
	const n = 10_000
	m := newBenchMap(b, n)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			if i%5 == 0 {
				m.Put("key"+strconv.Itoa(i%n), i)
			} else {
				m.Get("key" + strconv.Itoa(i%n))
			}
			i++
		}
	})
}
