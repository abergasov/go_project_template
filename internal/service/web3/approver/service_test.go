package approver_test

import (
	"context"
	"crypto/ecdsa"
	"go_project_template/internal/logger"
	"go_project_template/internal/service/web3"
	"go_project_template/internal/service/web3/approver"
	"go_project_template/internal/utils"
	"os"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stretchr/testify/require"
)

const (
	sampleContract = "0x8731d54E9D02c286767d56ac03e8037C07e01e98" // in eth mainnet
)

func TestService_ApproveContractUsage(t *testing.T) {
	appLog := logger.NewAppSLogger("")
	service := approver.InitService(appLog)
	ethClient, privateKey, accAddress := initTest(t)

	chainID, err := ethClient.ChainID(context.Background())
	require.NoError(t, err)
	tokenAddress, err := web3.GetTokenAddress(chainID, web3.USDC)
	require.NoError(t, err)

	approveTxHash, err := service.ApproveContractUsageALL(
		ethClient,
		privateKey,
		tokenAddress,
		accAddress,
		common.HexToAddress(sampleContract),
	)
	require.NoError(t, err)
	t.Log("approveTxHash:", approveTxHash)
}

func TestService_GetNativeTokenBalance(t *testing.T) {
	appLog := logger.NewAppSLogger("")
	service := approver.InitService(appLog)
	ethClient, _, accAddress := initTest(t)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	val, err := service.GetNativeTokenBalance(ctx, ethClient, accAddress)
	require.NoError(t, err)
	t.Log("val:", utils.ETHFromWei(val))
}

func TestService_GetERC20TokenBalance(t *testing.T) {
	appLog := logger.NewAppSLogger("")
	service := approver.InitService(appLog)
	ethClient, _, accAddress := initTest(t)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	chainID, err := ethClient.ChainID(context.Background())
	require.NoError(t, err)
	tokenAddress, err := web3.GetTokenAddress(chainID, web3.USDC)
	require.NoError(t, err)
	val, err := service.GetERC20TokenBalance(ctx, ethClient, tokenAddress, accAddress)
	require.NoError(t, err)
	t.Log("val:", utils.CustomFromWei(val, 6))
}

func TestService_GetContractData(t *testing.T) {
	appLog := logger.NewAppSLogger("")
	service := approver.InitService(appLog)
	ethClient, _, _ := initTest(t)
	chainID, err := ethClient.ChainID(context.Background())
	require.NoError(t, err)
	tokenAddress, err := web3.GetTokenAddress(chainID, web3.USDC)
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	ticker, decimal, err := service.GetContractData(ctx, ethClient, tokenAddress)
	require.NoError(t, err)
	require.Equal(t, "USDC", ticker)
	require.Equal(t, uint8(6), decimal)
}

func initTest(t *testing.T) (*ethclient.Client, *ecdsa.PrivateKey, common.Address) {
	ethClientRPC := os.Getenv("ETH_CLIENT_RPC")
	if ethClientRPC == "" {
		t.Skip("ETH_CLIENT_RPC is not set")
	}
	client, err := ethclient.Dial(ethClientRPC)
	require.NoError(t, err)
	t.Cleanup(func() {
		client.Close()

	})
	walletPK := os.Getenv("PRIVATE_KEY")
	if walletPK == "" {
		t.Skip("PRIVATE_KEY is not set")
	}
	privateKey, err := crypto.HexToECDSA(walletPK)
	require.NoError(t, err)

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	require.True(t, ok)
	return client, privateKey, crypto.PubkeyToAddress(*publicKeyECDSA)
}
