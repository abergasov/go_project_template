package utils_test

import (
	"go_project_template/internal/utils"
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWeiFromETHString(t *testing.T) {
	table := map[string]string{
		"0.02": "20000000000000000",
		"0.1":  "100000000000000000",
	}
	for eth, wei := range table {
		res := utils.WeiFromETHString(eth).String()
		require.Equal(t, wei, res)
	}
}

func TestGWeiFromETHString(t *testing.T) {
	table := map[string]string{
		"0.02": "20000000",
		"0.1":  "100000000",
	}
	for eth, gwei := range table {
		res := utils.GWeiFromETHString(eth).String()
		require.Equal(t, gwei, res)
	}
}

func TestETHFromWei(t *testing.T) {
	a, ok := big.NewInt(0).SetString("20000000000000000", 10)
	require.True(t, ok)
	b, ok := big.NewInt(0).SetString("100000000000000000", 10)
	require.True(t, ok)
	table := map[*big.Int]string{
		a: "0.02",
		b: "0.1",
	}
	for wei, eth := range table {
		res := utils.ETHFromWei(wei)
		require.Equal(t, eth, res)
	}
}

func TestETHFromGWei(t *testing.T) {
	a, ok := big.NewInt(0).SetString("20000000", 10)
	require.True(t, ok)
	b, ok := big.NewInt(0).SetString("100000000", 10)
	require.True(t, ok)
	table := map[*big.Int]string{
		a: "0.02",
		b: "0.1",
	}
	for wei, eth := range table {
		res := utils.ETHFromGWei(wei)
		require.Equal(t, eth, res)
	}
}

func TestCustomFromWei(t *testing.T) {
	table := map[string]string{
		"2815394107": "2815.394107",
	}
	for wei, eth := range table {
		a, ok := big.NewInt(0).SetString(wei, 10)
		require.True(t, ok)
		res := utils.CustomFromWei(a, 6)
		require.Equal(t, eth, res)
	}
}

func TestCustomToWei(t *testing.T) {
	table := map[float64]string{
		2815.394107: "2815394107",
	}
	for val, wei := range table {
		a, ok := big.NewInt(0).SetString(wei, 10)
		require.True(t, ok)
		res := utils.CustomToWei(val, 6)
		require.Equal(t, a, res)
	}
}
