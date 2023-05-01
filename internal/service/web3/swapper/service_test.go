package swapper_test

import (
	"crypto/ecdsa"
	"go_project_template/internal/service/web3/swapper"
	"os"
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stretchr/testify/require"
)

func TestService_TransferETH(t *testing.T) {
	service := swapper.NewService()

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

	approveTxHash, err := service.TransferETH(
		client,
		privateKey,
		swapper.LZArbitrumChainId,
		crypto.PubkeyToAddress(*publicKeyECDSA),
		"0.001",
		0.5,
	)
	require.NoError(t, err)
	t.Log("bridgeTx:", approveTxHash)
}
