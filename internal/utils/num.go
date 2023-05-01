package utils

import "math/big"

func BigFloatToInt(f *big.Float) *big.Int {
	i := new(big.Int)
	i, _ = f.Int(i)
	return i
}
