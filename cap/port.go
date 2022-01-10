package cap

import (
	"net"
)

type FreePortFinder struct{}

func (fpf FreePortFinder) FindPort() (uint, error) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return 0, err
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, err
	}
	defer l.Close()
	return uint(l.Addr().(*net.TCPAddr).Port), nil
}
