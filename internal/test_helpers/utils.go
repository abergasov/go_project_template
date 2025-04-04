package testhelpers

import (
	"github.com/stretchr/testify/require"
	"net"
	"testing"
)

// GetFreePort generate a free tcp port for testing
func GetFreePort(t *testing.T) int {
	t.Helper()
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	require.NoError(t, err)

	l, err := net.ListenTCP("tcp", addr)
	require.NoError(t, err)
	defer func() {
		require.NoError(t, l.Close())
	}()
	return l.Addr().(*net.TCPAddr).Port
}
