package approver

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"sync"

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

func (s *Service) GetNativeTokenBalance(ctx context.Context, web3Client *ethclient.Client, address common.Address) (*big.Int, error) {
	val, err := web3Client.BalanceAt(ctx, address, nil)
	if err != nil {
		return nil, fmt.Errorf("unable to get balance: %w", err)
	}
	return val, nil
}

func (s *Service) GetERC20TokenBalance(ctx context.Context, web3Client *ethclient.Client, tokenAddress, address common.Address) (*big.Int, error) {
	contract, err := NewErc20(tokenAddress, web3Client)
	if err != nil {
		return nil, err
	}
	val, err := contract.BalanceOf(&bind.CallOpts{
		Context: ctx,
	}, address)
	if err != nil {
		return nil, fmt.Errorf("unable to get balance: %w", err)
	}
	return val, nil
}

func (s *Service) GetContractData(ctx context.Context, web3Client *ethclient.Client, tokenAddress common.Address) (ticker string, decimal uint8, err error) {
	contract, err := NewErc20(tokenAddress, web3Client)
	if err != nil {
		return "", 0, err
	}
	var wg sync.WaitGroup
	wg.Add(2)
	var (
		errTicker  error
		errDecimal error
	)
	go func() {
		defer wg.Done()
		ticker, errTicker = contract.Symbol(&bind.CallOpts{
			Context: ctx,
		})
	}()
	go func() {
		defer wg.Done()
		decimal, errDecimal = contract.Decimals(&bind.CallOpts{
			Context: ctx,
		})
	}()
	wg.Wait()
	if errTicker != nil {
		return "", 0, fmt.Errorf("unable to get ticker: %w", errTicker)
	}
	if errDecimal != nil {
		return "", 0, fmt.Errorf("unable to get decimal: %w", errDecimal)
	}
	return ticker, decimal, nil
}
