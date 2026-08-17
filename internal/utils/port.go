package utils

import stdnet "net"

// GetFreePort reserves an available TCP port for testing purposes
func GetFreePort() (int, error) {
	addr, err := stdnet.ResolveTCPAddr("tcp", "127.0.0.1:0")
	if err != nil {
		return 0, err
	}

	l, err := stdnet.ListenTCP("tcp", addr)
	if err != nil {
		return 0, err
	}
	defer l.Close() //nolint // closing listener is safe here
	return l.Addr().(*stdnet.TCPAddr).Port, nil
}
