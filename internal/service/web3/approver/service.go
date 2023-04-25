package approver

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

// Service is a service that approves contracts to spend tokens
type Service struct {
	maxAllowed *big.Int
}

// InitService initializes the service
func InitService() *Service {
	two := big.NewInt(2)
	exponent := big.NewInt(256)
	power := new(big.Int).Exp(two, exponent, nil)
	result := new(big.Int).Sub(power, big.NewInt(1)) // (2 ** 256) -1
	return &Service{
		maxAllowed: result,
	}
}

// ApproveContractUsageALL approves the contract to spend all the tokens
func (s *Service) ApproveContractUsageALL(web3Client *ethclient.Client, privateKey *ecdsa.PrivateKey, tokenAddress, holder, spender common.Address) (string, error) {
	contract, err := NewErc20(tokenAddress, web3Client)
	if err != nil {
		return "", err
	}
	total, err := contract.Allowance(nil, holder, spender)
	if err != nil {
		return "", fmt.Errorf("unable to get allowance: %w", err)
	}
	if total.Cmp(s.maxAllowed) == 0 { // they are equal
		return "", nil
	}

	chainID, err := web3Client.ChainID(context.Background())
	if err != nil {
		return "", fmt.Errorf("unable to get chain id: %w", err)
	}

	tx, err := contract.Approve(&bind.TransactOpts{
		From: holder,
		Signer: func(address common.Address, transaction *types.Transaction) (*types.Transaction, error) {
			return types.SignTx(transaction, types.LatestSignerForChainID(chainID), privateKey)
		},
	}, spender, s.maxAllowed)
	if err != nil {
		return "", fmt.Errorf("unable to approve: %w", err)
	}

	return tx.Hash().String(), nil
}
