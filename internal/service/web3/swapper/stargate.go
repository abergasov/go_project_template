package swapper

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"go_project_template/internal/utils"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

const (
	LZArbitrumChainId = 110
	LZOptimismChainId = 111
)

// TransferETH transfers ETH from one chain into another
func (s *Service) TransferETH(web3Client *ethclient.Client, privateKey *ecdsa.PrivateKey, targetChainID uint16, holder common.Address, ethAmount string, slippagePercent float64) (string, error) {
	chainID, err := web3Client.ChainID(context.Background())
	if err != nil {
		return "", fmt.Errorf("unable to get chain id: %w", err)
	}
	contractAddress, ok := s.stargateAddress[chainID.String()]
	if !ok {
		return "", fmt.Errorf("unable to find contract address for chain id %s", chainID.String())
	}
	contract, err := NewStargateRouter(contractAddress, web3Client)
	if err != nil {
		return "", fmt.Errorf("unable to create contract: %w", err)
	}

	sendAmount := utils.WeiFromETHString(ethAmount)
	sendAmountFloat := big.NewFloat(0).SetInt(sendAmount)
	minAmount := big.NewFloat(0).Mul(sendAmountFloat, big.NewFloat(100-slippagePercent))
	minAmount = minAmount.Quo(minAmount, big.NewFloat(100))
	ma := utils.BigFloatToInt(minAmount)

	sendValueAdd := utils.WeiFromETHString("0.00038411")
	sendValueAdd = sendValueAdd.Add(sendValueAdd, sendAmount)

	tx, err := contract.SwapETH(&bind.TransactOpts{
		From:  holder,
		Value: sendValueAdd,
		Signer: func(address common.Address, transaction *types.Transaction) (*types.Transaction, error) {
			return types.SignTx(transaction, types.LatestSignerForChainID(chainID), privateKey)
		},
	}, targetChainID, holder, holder.Bytes(), sendAmount, ma)
	if err != nil {
		return "", fmt.Errorf("unable to approve: %w", err)
	}
	return tx.Hash().String(), nil
}
