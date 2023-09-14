package swapper

import (
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"go_project_template/internal/service/web3"
	"go_project_template/internal/utils"
	"log/slog"
	"math/big"
	"net/http"

	"github.com/shopspring/decimal"

	"github.com/ethereum/go-ethereum/core/types"

	"github.com/ethereum/go-ethereum/common"

	"github.com/ethereum/go-ethereum/ethclient"
)

type inchChain string

type inchResponse struct {
	FromToken struct {
		Symbol        string   `json:"symbol"`
		Name          string   `json:"name"`
		Address       string   `json:"address"`
		Decimals      int      `json:"decimals"`
		LogoURI       string   `json:"logoURI"`
		Eip2612       bool     `json:"eip2612"`
		DomainVersion string   `json:"domainVersion"`
		Tags          []string `json:"tags"`
	} `json:"fromToken"`
	ToToken struct {
		Symbol   string   `json:"symbol"`
		Name     string   `json:"name"`
		Decimals int      `json:"decimals"`
		Address  string   `json:"address"`
		LogoURI  string   `json:"logoURI"`
		Eip2612  bool     `json:"eip2612"`
		Tags     []string `json:"tags"`
	} `json:"toToken"`
	ToTokenAmount   string `json:"toTokenAmount"`
	FromTokenAmount string `json:"fromTokenAmount"`
	Protocols       [][][]struct {
		Name             string `json:"name"`
		Part             int    `json:"part"`
		FromTokenAddress string `json:"fromTokenAddress"`
		ToTokenAddress   string `json:"toTokenAddress"`
	} `json:"protocols"`
	Tx struct {
		From     string `json:"from"`
		To       string `json:"to"`
		Data     string `json:"data"`
		Value    string `json:"value"`
		Gas      uint64 `json:"gas"`
		GasPrice string `json:"gasPrice"`
	} `json:"tx"`
	EstimatedGas int `json:"estimatedGas"`
}

const (
	inchVersion = 5
)

func (s *Service) Swap1Inch(
	ctx context.Context,
	web3Client *ethclient.Client,
	privateKey *ecdsa.PrivateKey,
	holder common.Address,
	slippagePercent float64,
	from web3.Coin,
	swapAmount float64,
	to web3.Coin,
) error {
	log := s.log.With(
		slog.String("from", string(from)),
		slog.String("to", string(to)),
		slog.Float64("amount", swapAmount),
	)
	log.Info("start 1inch swap")
	chainID, err := web3Client.ChainID(ctx)
	if err != nil {
		return fmt.Errorf("unable to get chain id: %w", err)
	}
	fromToken, err := web3.GetTokenAddress(chainID, from)
	if err != nil {
		return fmt.Errorf("unable to get token FROM address: %w", err)
	}

	_, decimalToken, err := s.approver.GetContractData(ctx, web3Client, fromToken)
	if err != nil {
		return fmt.Errorf("unable to get token decimals: %w", err)
	}

	toToken, err := web3.GetTokenAddress(chainID, to)
	if err != nil {
		return fmt.Errorf("unable to get token TO address: %w", err)
	}

	log.Info("approve inch usage")
	if err = s.approveInchUsage(ctx, web3Client, privateKey, from, chainID, holder); err != nil {
		return fmt.Errorf("unable to approve inch usage: %w", err)
	}

	amountToSwap := utils.CustomToWei(swapAmount, int(decimalToken))
	res, err := s.getInchData(ctx, chainID, slippagePercent, amountToSwap.String(), holder.String(), fromToken.String(), toToken.String())
	if err != nil {
		return fmt.Errorf("unable to get 1inch data: %w", err)
	}

	gasPrice, err := decimal.NewFromString(res.Tx.GasPrice)
	if err != nil {
		return fmt.Errorf("unable to get gas price: %w", err)
	}
	nonce, gasFeeCap, gasTipCap, err := s.suggestedFeeAndTipWithNonce(ctx, web3Client, gasPrice.BigInt(), 100, holder)
	if err != nil {
		return fmt.Errorf("unable to suggest gas: %w", err)
	}

	broadcastTx := types.NewTx(&types.DynamicFeeTx{
		ChainID:   chainID,
		To:        utils.ToPointer(common.HexToAddress(res.Tx.To)),
		Nonce:     nonce,
		GasTipCap: gasTipCap,
		GasFeeCap: gasFeeCap,
		Gas:       res.Tx.Gas,
		Value:     big.NewInt(0),
		Data:      common.FromHex(res.Tx.Data),
	})
	log.Info(
		"run transaction",
		slog.String("gas tip cap", gasTipCap.String()),
		slog.String("gas fee cap", gasFeeCap.String()),
		slog.Uint64("gas limit", res.Tx.Gas),
	)
	signedTx, err := types.SignTx(broadcastTx, types.LatestSignerForChainID(chainID), privateKey)
	if err != nil {
		return fmt.Errorf("unable to sign tx: %w", err)
	}
	if err = web3Client.SendTransaction(ctx, signedTx); err != nil {
		return fmt.Errorf("unable to send tx: %w", err)
	}
	return nil
}

func (s *Service) approveInchUsage(ctx context.Context, web3Client *ethclient.Client, privateKey *ecdsa.PrivateKey, from web3.Coin, chain *big.Int, holder common.Address) error {
	if from == web3.ETH {
		return nil
	}
	// need to approve
	spender, err := s.getInchSpender(ctx, chain)
	if err != nil {
		return fmt.Errorf("unable to get 1inch spender: %w", err)
	}
	fromToken, err := web3.GetTokenAddress(chain, from)
	if err != nil {
		return fmt.Errorf("unable to get token FROM address: %w", err)
	}
	txHash, err := s.approver.ApproveContractUsageALL(web3Client, privateKey, fromToken, holder, common.HexToAddress(spender))
	if err != nil {
		return fmt.Errorf("unable to approve: %w", err)
	}
	if err = s.approver.WaitTransaction(ctx, web3Client, txHash); err != nil {
		return fmt.Errorf("unable to wait approve tx mined: %w", err)
	}
	return nil
}

func (s *Service) getInchData(ctx context.Context, chain *big.Int, slippagePercent float64, amountToSwap, wallet, fromToken, toToken string) (*inchResponse, error) {
	url := fmt.Sprintf(
		"https://api.1inch.io/v%d.0/%s/swap?fromTokenAddress=%s&toTokenAddress=%s&amount=%s&fromAddress=%s&slippage=%.2f",
		inchVersion,
		chain.String(),
		fromToken,
		toToken,
		amountToSwap,
		wallet,
		slippagePercent,
	)
	resp, respCode, err := utils.Get(ctx, url)
	if err != nil {
		return nil, fmt.Errorf("unable to get 1inch data: %w", err)
	}
	if respCode != http.StatusOK {
		type inchErrorResponse struct {
			StatusCode  int    `json:"statusCode"`
			Error       string `json:"error"`
			Description string `json:"description"`
			Meta        []struct {
				Type  string `json:"type"`
				Value string `json:"value"`
			} `json:"meta"`
		}
		var inchErrResp inchErrorResponse
		if err = json.Unmarshal(resp, &inchErrResp); err != nil {
			return nil, fmt.Errorf("failed to unmarshal response body: %d, %w", respCode, err)
		}
		return nil, fmt.Errorf("unexpected status code: %d, %s, %s", respCode, inchErrResp.Error, inchErrResp.Description)
	}

	var inchResp inchResponse
	if err = json.Unmarshal(resp, &inchResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response body: %w", err)
	}
	return &inchResp, nil
}

func (s *Service) getInchSpender(ctx context.Context, chainID *big.Int) (string, error) {
	url := fmt.Sprintf("https://api.1inch.io/v%d.0/%s/approve/spender", inchVersion, chainID.String())
	resp, respCode, err := utils.Get(ctx, url)
	if err != nil {
		return "", fmt.Errorf("unable to get 1inch data: %w", err)
	}
	if respCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", respCode)
	}
	type inchSpenderResponse struct {
		Address string `json:"address"`
	}
	var inchResp inchSpenderResponse
	if err = json.Unmarshal(resp, &inchResp); err != nil {
		return "", fmt.Errorf("failed to unmarshal response body: %w", err)
	}
	return inchResp.Address, nil
}
