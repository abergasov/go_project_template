package swapper_test

import (
	"context"
	"go_project_template/internal/logger"
	"go_project_template/internal/service/web3"
	"go_project_template/internal/service/web3/approver"
	"go_project_template/internal/service/web3/swapper"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestService_Swap1Inch(t *testing.T) {
	appLog := logger.NewAppSLogger("")
	erc20Approver := approver.InitService(appLog)
	service := swapper.NewService(appLog, erc20Approver)
	ethClient, privateKey, address := initTest(t)

	ctx, cancel := context.WithTimeout(context.Background(), 360*time.Second)
	defer cancel()

	err := service.Swap1Inch(
		ctx,
		ethClient,
		privateKey,
		address,
		1,
		web3.DAI,
		10,
		web3.USDC,
	)
	require.NoError(t, err)
}
