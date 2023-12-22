package utils_test

import (
	"github.com/stretchr/testify/require"
	"go_project_template/internal/utils"
	"reflect"
	"slices"
	"testing"
)

func TestContainsInSlice(t *testing.T) {
	t.Parallel()
	t.Run("should return true if value is in slice", func(t *testing.T) {
		// given
		slice := []string{"a", "b", "c"}
		value := "b"

		// when
		result := utils.ContainsInSlice(slice, value)

		// then
		require.True(t, result)
	})
	t.Run("should return false if value is not in slice", func(t *testing.T) {
		// given
		slice := []string{"a", "b", "c"}
		value := "d"

		// when
		result := utils.ContainsInSlice(slice, value)

		// then
		require.False(t, result)
	})
}

func TestFindOutsideList(t *testing.T) {
	type tCase struct {
		slice []string
		list  []string
		res   []string
		name  string
	}
	table := []tCase{
		{
			list:  []string{"a", "b", "c"},
			slice: []string{},
			res:   []string{},
			name:  "should return empty slice if list is empty",
		},
		{
			name:  "should return diff",
			slice: []string{"a", "b", "c"},
			list:  []string{"a", "b"},
			res:   []string{"c"},
		},
		{
			name:  "should return empty slice all items are same",
			slice: []string{"a", "b", "c"},
			list:  []string{"a", "b", "c"},
			res:   []string{},
		},
		{
			name:  "should return empty list if items are less than list",
			slice: []string{"a", "b"},
			list:  []string{"a", "b", "c"},
			res:   []string{},
		},
	}
	for _, tc := range table {
		t.Run(tc.name, func(t *testing.T) {
			// given, when
			result := utils.FindOutsideList(tc.slice, tc.list)

			// then
			require.Equal(t, tc.res, result)
		})
	}
}

func TestFindGapsInBlockSlice(t *testing.T) {
	require.Empty(t, utils.FindGapsInBlockSlice([]uint64{}))
	require.Empty(t, utils.FindGapsInBlockSlice([]uint64{1, 2, 3, 4, 5, 6, 7, 8, 9}))
	require.Empty(t, utils.FindGapsInBlockSlice([]uint64{9, 3, 7, 4, 5, 6, 2, 8, 1}))
	require.Equal(t, []uint64{2, 3, 5}, utils.FindGapsInBlockSlice([]uint64{1, 4, 6, 7, 8, 9}))
	require.Equal(t, []uint64{2, 3, 5}, utils.FindGapsInBlockSlice([]uint64{9, 7, 4, 6, 8, 1}))
	require.Equal(t, []uint64{3}, utils.FindGapsInBlockSlice([]uint64{1, 7, 4, 5, 6, 2, 8, 9}))
	require.Equal(t, []uint64{3, 4}, utils.FindGapsInBlockSlice([]uint64{1, 7, 5, 6, 2, 8, 9}))
}

func TestChunkSlice(t *testing.T) {
	require.Equal(t, [][]uint64{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}}, utils.ChunkSlice([]uint64{1, 2, 3, 4, 5, 6, 7, 8, 9}, 3))
	require.Equal(t, [][]uint64{{1, 2, 3}, {4, 5, 6}, {7, 8}}, utils.ChunkSlice([]uint64{1, 2, 3, 4, 5, 6, 7, 8}, 3))
}

func TestUniqueSlice(t *testing.T) {
	require.Equal(t, []uint64{1, 2, 3, 4, 5, 6, 7, 8, 9}, utils.UniqueSlice([]uint64{1, 2, 3, 4, 5, 6, 7, 8, 9}))
	require.Equal(t, []uint64{1, 2, 3, 4, 5, 6, 7, 8}, utils.UniqueSlice([]uint64{1, 2, 3, 4, 5, 6, 7, 8, 8}))
	require.Equal(t, []uint64{1, 2, 3, 4, 5, 6, 7, 8}, utils.UniqueSlice([]uint64{1, 2, 3, 4, 3, 4, 5, 6, 7, 8, 8, 7}))
	require.Equal(t, []string{"a", "b", "c"}, utils.UniqueSlice([]string{"a", "b", "c", "a", "b", "c"}))
}

func TestMax(t *testing.T) {
	require.Equal(t, 9, utils.Max(2, 9))
	require.Equal(t, 4, utils.Max(4, 3))
}

type testPerson struct {
	ID   int
	Name string
}

func TestStringsFromObjectMap(t *testing.T) {
	// Example data
	personMap := map[int]testPerson{
		1: {ID: 1, Name: "John"},
		2: {ID: 2, Name: "Jane"},
		3: {ID: 3, Name: "Doe"},
	}

	// Extractor function for Person
	extractor := func(p testPerson) string {
		return p.Name
	}

	// Expected result
	expected := []string{"Doe", "Jane", "John"}

	// Test the function
	result := utils.StringsFromObjectMap(personMap, extractor)
	slices.Sort(result)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, but got %v", expected, result)
	}
}

func TestStringsFromObjectSlice(t *testing.T) {
	// Example data
	personSlice := []testPerson{
		{ID: 1, Name: "John"},
		{ID: 2, Name: "Jane"},
		{ID: 3, Name: "Doe"},
	}

	// Extractor function for Person
	extractor := func(p testPerson) string {
		return p.Name
	}

	// Expected result
	expected := []string{"John", "Jane", "Doe"}

	// Test the function
	result := utils.StringsFromObjectSlice(personSlice, extractor)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, but got %v", expected, result)
	}
}
