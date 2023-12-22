package utils_test

import (
	"github.com/stretchr/testify/require"
	"go_project_template/internal/utils"
	"sort"
	"strings"
	"testing"
)

func TestMaps(t *testing.T) {
	m := map[int]string{228: "t", 1488: "e", 666: "s", 777: "t"}
	keys := utils.Keys(m)
	values := utils.Values(m)
	require.Equal(t, len(keys), len(m))
	require.Equal(t, len(values), len(m))
	sort.Ints(keys)
	sort.Strings(values)
	require.Equal(t, keys, []int{228, 666, 777, 1488})
	require.Equal(t, strings.Join(values, ""), "estt")
}
