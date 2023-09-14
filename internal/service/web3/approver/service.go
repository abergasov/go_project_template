package approver

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"go_project_template/internal/logger"
	"log/slog"
	"math/big"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

// Service is a service that approves contracts to spend tokens
type Service struct {
	maxAllowed *big.Int
	log        logger.AppLogger
}

// InitService initializes the service
func InitService(appLog logger.AppLogger) *Service {
	two := big.NewInt(2)
	exponent := big.NewInt(256)
	power := new(big.Int).Exp(two, exponent, nil)
	result := new(big.Int).Sub(power, big.NewInt(1)) // (2 ** 256) -1
	return &Service{
		maxAllowed: result,
		log:        appLog,
	}
}

// ApproveContractUsageALL approves the contract to spend all the tokens
func (s *Service) ApproveContractUsageALL(web3Client *ethclient.Client, privateKey *ecdsa.PrivateKey, tokenAddress, holder, spender common.Address) (string, error) {
	log := s.log.With(
		slog.String("tokenAddress", tokenAddress.String()),
		slog.String("holder", holder.String()),
		slog.String("spender", spender.String()),
	)
	contract, err := NewErc20(tokenAddress, web3Client)
	if err != nil {
		return "", err
	}
	total, err := contract.Allowance(nil, holder, spender)
	if err != nil {
		return "", fmt.Errorf("unable to get allowance: %w", err)
	}
	log.Info("allowance", slog.String("allowance", total.String()))
	if total.Cmp(s.maxAllowed) == 0 { // they are equal
		log.Info("already approved", slog.String("allowance", total.String()))
		return "", nil
	}

	chainID, err := web3Client.ChainID(context.Background())
	if err != nil {
		return "", fmt.Errorf("unable to get chain id: %w", err)
	}
	log.Info("approving", slog.String("allowance", s.maxAllowed.String()), slog.Uint64("chainID", chainID.Uint64()))

	tx, err := contract.Approve(&bind.TransactOpts{
		From: holder,
		Signer: func(address common.Address, transaction *types.Transaction) (*types.Transaction, error) {
			return types.SignTx(transaction, types.LatestSignerForChainID(chainID), privateKey)
		},
	}, spender, s.maxAllowed)
	if err != nil {
		return "", fmt.Errorf("unable to approve: %w", err)
	}

	log.Info("tx sent", slog.String("txHash", tx.Hash().String()))
	return tx.Hash().String(), nil
}

func (s *Service) GetNativeTokenBalance(ctx context.Context, web3Client *ethclient.Client, address common.Address) (*big.Int, error) {
	val, err := web3Client.BalanceAt(ctx, address, nil)
	if err != nil {
		return nil, fmt.Errorf("unable to get balance: %w", err)
	}
	return val, nil
}

func (s *Service) WaitTransaction(ctx context.Context, web3Client *ethclient.Client, txHash string) error {
	if txHash == "" {
		return nil
	}
	ctx, cancel := context.WithTimeout(ctx, 5*time.Minute)
	defer cancel()
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			_, err := web3Client.TransactionReceipt(ctx, common.HexToHash(txHash))
			if err == nil {
				return nil
			}
			if !errors.Is(err, ethereum.NotFound) {
				return fmt.Errorf("unable to get transaction receipt: %w", err)
			}
			time.Sleep(5 * time.Second)
		}
	}
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
