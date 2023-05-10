package swapper

import (
	"context"
	"fmt"
	"go_project_template/internal/logger"
	"go_project_template/internal/service/web3/approver"
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type Service struct {
	approver        *approver.Service
	log             logger.AppLogger
	stargateAddress map[string]common.Address
}

func NewService(appLog logger.AppLogger, erc20 *approver.Service) *Service {
	return &Service{
		log:      appLog,
		approver: erc20,
		stargateAddress: map[string]common.Address{
			"1": common.HexToAddress("0x150f94b44927f078737562f0fcf3c95c01cc2376"), // in eth mainnet
		},
	}
}

func (s *Service) suggestGas(ctx context.Context, web3Client *ethclient.Client, holder common.Address) (nonce uint64, gasPrice, gasTipCap *big.Int, err error) {
	var wg sync.WaitGroup
	var (
		errNonce     error
		errGasTipCap error
		errGasPrice  error
	)
	wg.Add(3)
	go func() {
		defer wg.Done()
		nonce, errNonce = web3Client.PendingNonceAt(context.Background(), holder)
	}()
	go func() {
		defer wg.Done()
		gasTipCap, errGasTipCap = web3Client.SuggestGasTipCap(ctx)
	}()
	go func() {
		defer wg.Done()
		gasPrice, errGasPrice = web3Client.SuggestGasPrice(ctx)
	}()
	wg.Wait()
	if errNonce != nil {
		return 0, nil, nil, fmt.Errorf("unable to get nonce: %w", errNonce)
	}
	if errGasTipCap != nil {
		return 0, nil, nil, fmt.Errorf("unable to get gas tip cap: %w", errGasTipCap)
	}
	if errGasPrice != nil {
		return 0, nil, nil, fmt.Errorf("unable to get gas price: %w", errGasPrice)
	}
	return nonce, gasPrice, gasTipCap, nil
}

func (s *Service) suggestedFeeAndTipWithNonce(ctx context.Context, web3Client *ethclient.Client, gasPrice *big.Int, boostPercent int, holder common.Address) (nonce uint64, gasFeeCap, gasTipCap *big.Int, err error) {
	var wg sync.WaitGroup
	var (
		errNonce     error
		errGasTipCap error
		errGasPrice  error
	)
	wg.Add(3)

	go func() {
		defer wg.Done()
		nonce, errNonce = web3Client.PendingNonceAt(context.Background(), holder)
	}()

	go func() {
		defer wg.Done()
		gasTipCap, errGasTipCap = web3Client.SuggestGasTipCap(ctx)
	}()

	go func() {
		defer wg.Done()
		if gasPrice != nil {
			return
		}
		gasPrice, errGasPrice = web3Client.SuggestGasPrice(ctx)
		if errGasPrice != nil {
			return
		}
		gasPrice = new(big.Int).Div(new(big.Int).Mul(big.NewInt(int64(boostPercent)+100), gasPrice), big.NewInt(100))
	}()

	wg.Wait()
	if errNonce != nil {
		return 0, nil, nil, fmt.Errorf("unable to get nonce: %w", errNonce)
	}
	if errGasTipCap != nil {
		return 0, nil, nil, fmt.Errorf("unable to get gas tip cap: %w", errGasTipCap)
	}
	if errGasPrice != nil {
		return 0, nil, nil, fmt.Errorf("unable to get gas price: %w", errGasPrice)
	}

	gasTipCap = new(big.Int).Div(new(big.Int).Mul(big.NewInt(int64(boostPercent)+100), gasTipCap), big.NewInt(100))
	gasFeeCap = new(big.Int).Add(gasTipCap, gasPrice)

	return nonce, gasFeeCap, gasTipCap, nil
}
