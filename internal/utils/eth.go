package utils

import (
	"math/big"

	"github.com/shopspring/decimal"
)

func WeiFromETHString(eth string) *big.Int {
	amount, _ := decimal.NewFromString(eth)
	mul := decimal.NewFromFloat(float64(10)).Pow(decimal.NewFromFloat(18))
	return amount.Mul(mul).BigInt()
}

func GWeiFromETHString(eth string) *big.Int {
	amount, _ := decimal.NewFromString(eth)
	mul := decimal.NewFromFloat(float64(10)).Pow(decimal.NewFromFloat(9))
	return amount.Mul(mul).BigInt()
}

func ETHFromWei(wei *big.Int) string {
	amount := decimal.NewFromBigInt(wei, 0)
	mul := decimal.NewFromFloat(float64(10)).Pow(decimal.NewFromFloat(18))
	return amount.Div(mul).String()
}

func ETHFromGWei(wei *big.Int) string {
	amount := decimal.NewFromBigInt(wei, 0)
	mul := decimal.NewFromFloat(float64(10)).Pow(decimal.NewFromFloat(9))
	return amount.Div(mul).String()
}

func CustomFromWei(wei *big.Int, decimals int) string {
	amount := decimal.NewFromBigInt(wei, 0)
	mul := decimal.NewFromFloat(float64(10)).Pow(decimal.NewFromFloat(float64(decimals)))
	return amount.Div(mul).String()
}

func CustomToWei(amount float64, decimals int) *big.Int {
	amountDecimal := decimal.NewFromFloat(amount)
	mul := decimal.NewFromFloat(float64(10)).Pow(decimal.NewFromFloat(float64(decimals)))
	return amountDecimal.Mul(mul).BigInt()
}
