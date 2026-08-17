package utils_test

import (
	"fmt"
	stdnet "net"
	"testing"

	"go_project_template/internal/utils"

	"github.com/stretchr/testify/require"
)

func TestGetFreePortFunctions(t *testing.T) {
	port, err := utils.GetFreePort()
	require.NoError(t, err)
	require.Greater(t, port, 0)
	ensurePortBindable(t, port)
}

func ensurePortBindable(t *testing.T, port int) {
	t.Helper()

	listener, err := stdnet.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", port))
	require.NoError(t, err)

	t.Cleanup(func() {
		require.NoError(t, listener.Close())
	})
}
