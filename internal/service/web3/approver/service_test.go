package approver_test

import (
	"crypto/ecdsa"
	"go_project_template/internal/service/web3/approver"
	"os"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stretchr/testify/require"
)

const (
	usdcAddress    = "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48" // in eth mainnet
	sampleContract = "0x8731d54E9D02c286767d56ac03e8037C07e01e98" // in eth mainnet
)

func TestService_ApproveContractUsage(t *testing.T) {
	service := approver.InitService()

	ethClientRPC := os.Getenv("ETH_CLIENT_RPC")
	if ethClientRPC == "" {
		t.Skip("ETH_CLIENT_RPC is not set")
	}
	client, err := ethclient.Dial(ethClientRPC)
	require.NoError(t, err)

	walletPK := os.Getenv("PRIVATE_KEY")
	if walletPK == "" {
		t.Skip("PRIVATE_KEY is not set")
	}
	privateKey, err := crypto.HexToECDSA(walletPK)
	require.NoError(t, err)

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	require.True(t, ok)

	approveTxHash, err := service.ApproveContractUsageALL(
		client,
		privateKey,
		common.HexToAddress(usdcAddress),
		crypto.PubkeyToAddress(*publicKeyECDSA),
		common.HexToAddress(sampleContract),
	)
	require.NoError(t, err)
	t.Log("approveTxHash:", approveTxHash)
}
