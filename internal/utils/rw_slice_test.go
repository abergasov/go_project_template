package utils_test

import (
	"go_project_template/internal/utils"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRWSlice_Add(t *testing.T) {
	rwSlice := utils.NewRWSlice[int]()

	// Add elements to the slice
	rwSlice.Add(1)
	rwSlice.Add(2)
	rwSlice.Add(3)

	require.True(t, equalSlices([]int{1, 2, 3}, rwSlice.LoadAll()))
}

func TestRWSlice_LoadAndErase(t *testing.T) {
	rwSlice := utils.NewRWSlice[int]()
	rwSlice.Add(1)
	rwSlice.Add(2)
	rwSlice.Add(3)

	require.True(t, equalSlices([]int{1, 2, 3}, rwSlice.LoadAndErase()))
	require.Equal(t, 0, rwSlice.Len())
}

func TestRWSlice_LoadAll(t *testing.T) {
	rwSlice := utils.NewRWSlice[string]()
	rwSlice.Add("hello")
	rwSlice.Add("world")

	require.True(t, equalSlices([]string{"hello", "world"}, rwSlice.LoadAll()))
}

func TestRWSlice_Len(t *testing.T) {
	rwSlice := utils.NewRWSlice[float64]()
	rwSlice.Add(1.1)
	rwSlice.Add(2.2)
	rwSlice.Add(3.3)
	require.Equal(t, 3, rwSlice.Len())
}

func TestRWSlice_Erase(t *testing.T) {
	rwSlice := utils.NewRWSlice[int]()
	rwSlice.Add(10)
	rwSlice.Add(20)

	rwSlice.Erase()

	require.Equal(t, 0, rwSlice.Len())
	require.Len(t, rwSlice.LoadAll(), 0)
}

func TestRWSlice_Concurrency(t *testing.T) {
	rwSlice := utils.NewRWSlice[int]()
	var wg sync.WaitGroup

	// Concurrent writes
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			rwSlice.Add(i)
		}(i)
	}

	// Concurrent reads
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = rwSlice.LoadAll()
		}()
	}

	wg.Wait()

	require.Equal(t, 100, rwSlice.Len())
}

func equalSlices[V comparable](a, b []V) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
