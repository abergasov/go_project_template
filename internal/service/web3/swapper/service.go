package swapper

import "github.com/ethereum/go-ethereum/common"

type Service struct {
	stargateAddress map[string]common.Address
}

func NewService() *Service {
	return &Service{
		stargateAddress: map[string]common.Address{
			"1": common.HexToAddress("0x150f94b44927f078737562f0fcf3c95c01cc2376"), // in eth mainnet
		},
	}
}
