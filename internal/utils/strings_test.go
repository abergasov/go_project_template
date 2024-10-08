package utils_test

import (
	"go_project_template/internal/utils"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCleanString(t *testing.T) {
	table := map[string]string{
		"Hello, World!": "helloworld",
		"Hello, 123!":   "hello",
		"Hello, 123":    "hello",
		"Привет, 123":   "привет",
		"Это тест":      "этотест",
	}
	for src, res := range table {
		require.Equal(t, res, utils.CleanString(src))
	}
}

func TestGetFirstValidString(t *testing.T) {
	table := map[string][]string{
		"":  {},
		"a": {"a"},
		"b": {"b", "a"},
		"c": {"", "c", "b", "a"},
		"d": {"", " ", "  ", "d"},
	}
	for expected, params := range table {
		require.Equal(t, expected, utils.GetFirstValidString(params...))
	}
}
