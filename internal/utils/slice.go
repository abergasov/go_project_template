package utils

import (
	"sort"

	"golang.org/x/exp/constraints"
)

// ContainsInSlice checks if slice contains value
func ContainsInSlice[T comparable](slice []T, value T) bool {
	for _, entry := range slice {
		if entry == value {
			return true
		}
	}
	return false
}

// FindOutsideList finds elements from slice that are not in list
func FindOutsideList[T comparable](slice, list []T) []T {
	result := make([]T, 0, len(slice))
	for _, entry := range slice {
		if !ContainsInSlice(list, entry) {
			result = append(result, entry)
		}
	}
	return result
}

// FindGapsInBlockSlice finds gaps in block slice
func FindGapsInBlockSlice(blockList []uint64) []uint64 {
	if len(blockList) == 0 {
		return []uint64{}
	}
	// 1. sort slice
	sort.Slice(blockList, func(i, j int) bool {
		return blockList[i] < blockList[j]
	})
	firstBlock := blockList[0]
	lastBlock := blockList[len(blockList)-1]
	result := make([]uint64, 0, lastBlock-firstBlock)
	// 2. find gaps
	for i := 1; i < len(blockList); i++ {
		diff := blockList[i] - blockList[i-1]
		if diff > 1 {
			// there is a gap between the two numbers
			for j := 1; j < int(diff); j++ {
				gap := blockList[i-1] + uint64(j)
				result = append(result, gap)
			}
		}
	}
	// 3. return gaps
	return result
}

func ChunkSlice[T any](slice []T, chunkSize int) [][]T {
	chunks := make([][]T, 0, (len(slice)+chunkSize-1)/chunkSize)
	for i := 0; i < len(slice); i += chunkSize {
		end := i + chunkSize
		if end > len(slice) {
			end = len(slice)
		}
		chunks = append(chunks, slice[i:end])
	}
	return chunks
}

func UniqueSlice[T comparable](slice []T) []T {
	keys := make(map[T]struct{}, len(slice))
	list := make([]T, 0, len(slice))
	for _, entry := range slice {
		if _, value := keys[entry]; !value {
			keys[entry] = struct{}{}
			list = append(list, entry)
		}
	}
	return list
}

func Max[T constraints.Ordered](a, b T) T {
	if b > a {
		return b
	}
	return a
}

func StringsFromObjectMap[T any, K comparable](src map[K]T, extractor func(T) string) []string {
	result := make([]string, 0, len(src))
	for key := range src {
		result = append(result, extractor(src[key]))
	}
	return result
}

func StringsFromObjectSlice[T any](src []T, extractor func(T) string) []string {
	result := make([]string, 0, len(src))
	for key := range src {
		result = append(result, extractor(src[key]))
	}
	return result
}
