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
	return 1, &MyMockedAddr{}, nil
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

func TestSmoke(t *testing.T) {
	mockUDPConn := new(MyMockedUDPConn)
	server := NewServer(mockUDPConn)
	mockAddr := new(MyMockedAddr)

	client := server.GetOrMakeClient(mockAddr)
	server.Send("", client)
}
