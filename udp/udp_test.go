package udp

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockedUDPConn struct {
	mock.Mock
	net.UDPConn
}

// ReadFrom implements conn.
func (m *MockedUDPConn) ReadFromUDP(p []byte) (n int, addr *net.UDPAddr, err error) {
	args := m.Called(p)
	return args.Int(0), nil, args.Error(2)
}

// WriteTo implements conn.
func (m *MockedUDPConn) WriteToUDP(p []byte, addr *net.UDPAddr) (n int, err error) {
	args := m.Called(p, addr)
	return args.Int(0), args.Error(1)
}

type MockedAddr struct {
	mock.Mock
	net.UDPAddr
}

// // Network implements net.Addr.
// func (m *MockedAddr) Network() string {
// 	return m.Network()
// }

// // String implements net.Addr.
// // Subtle: this method shadows the method (Mock).String of MockedAddr.Mock.
// func (m *MockedAddr) String() string {
// 	return m.String()
// }

func TestSmoke(t *testing.T) {
	_, err := EmptyServer()
	assert.NoError(t, err)
}
