package udp

import (
	"net"
	"testing"

	"github.com/stretchr/testify/mock"
)

type MyMockedUDPConn struct {
	mock.Mock
}

func (o *MyMockedUDPConn) ReadFrom(p []byte) (n int, addr net.Addr, err error) {
	return
}

func (o *MyMockedUDPConn) WriteTo(p []byte, addr net.Addr) (n int, err error) {
	return
}

type MyMockedAddr struct {
	mock.Mock
}

func (o *MyMockedAddr) Network() string {
	return ""
}

func (o *MyMockedAddr) String() string {
	return ""
}

func TestNewServer(t *testing.T) {
	mockUDPConn := new(MyMockedUDPConn)
	server := NewServer(mockUDPConn)
	mockAddr := new(MyMockedAddr)
	server.GetOrMakeClient(mockAddr)
}
