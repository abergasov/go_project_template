package web3

import (
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type Coin string

const (
	ETH  Coin = "ETH"
	USDC Coin = "USDC"
	DAI  Coin = "DAI"
)

var (
	ErrUnknownCoin  = errors.New("unknown coin")
	ErrUnknownChain = errors.New("unknown chain")
)

var Tokens = map[Coin]map[uint64]common.Address{
	ETH: {
		1: common.HexToAddress("0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE"),
	},
	USDC: {
		1: common.HexToAddress("0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48"),
	},
	DAI: {
		1: common.HexToAddress("0x6B175474E89094C44Da98b954EedeAC495271d0F"),
	},
}

func GetTokenAddress(chainID *big.Int, coin Coin) (common.Address, error) {
	if _, ok := Tokens[coin]; !ok {
		return common.Address{}, ErrUnknownCoin
	}
	if _, ok := Tokens[coin][chainID.Uint64()]; !ok {
		return common.Address{}, ErrUnknownChain
	}
	return Tokens[coin][chainID.Uint64()], nil
}
