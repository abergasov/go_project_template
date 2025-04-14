package utils

import (
	"math/big"
)

func WeiFromETHString(eth string) *big.Int {
	ethFloat, _ := new(big.Float).SetString(eth)
	amount, _ := ethFloat.Mul(ethFloat, new(big.Float).SetInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil))).Int(nil)
	return amount
}

func GWeiFromETHString(eth string) *big.Int {
	amount, _, err := big.ParseFloat(eth, 10, 256, big.ToNearestEven)
	if err != nil {
		return nil
	}

	multiplier := new(big.Float).SetInt(big.NewInt(1_000_000_000))
	result := new(big.Float).Mul(amount, multiplier)

	gwei := new(big.Int)
	result.Int(gwei)
	return gwei
}

func ETHFromWei(wei *big.Int) string {
	return new(big.Float).Quo(new(big.Float).SetInt(wei), new(big.Float).SetInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil))).String()
}

func ETHFromGWei(gwei *big.Int) string {
	amount := new(big.Float).SetInt(gwei)
	divisor := new(big.Float).SetInt(big.NewInt(1_000_000_000))
	result := new(big.Float).Quo(amount, divisor)
	return result.Text('f', -1)
}

func CustomFromWei(wei *big.Int, decimals int) string {
	amount := new(big.Float).SetInt(wei)
	divisor := new(big.Float).SetFloat64(1)
	divisor.Quo(divisor, new(big.Float).SetFloat64(1).SetPrec(256).SetFloat64(float64(10)).SetPrec(256).SetInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(decimals)), nil)))
	dec := new(big.Float).Quo(amount, new(big.Float).SetInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(decimals)), nil)))
	return dec.Text('f', decimals)
}

func CustomToWei(amount float64, decimals int) *big.Int {
	amountFloat := new(big.Float).SetFloat64(amount)
	exp := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(decimals)), nil)
	expFloat := new(big.Float).SetInt(exp)
	result := new(big.Float).Mul(amountFloat, expFloat)
	wei := new(big.Int)
	result.Int(wei)
	return wei
}
