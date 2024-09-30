package udp

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockedUDPConn struct {
	mock.Mock
}

// ReadFrom implements conn.
func (m *MockedUDPConn) ReadFrom(p []byte) (n int, addr net.Addr, err error) {
	args := m.Called(p)
	return args.Int(0), nil, args.Error(2)
}

// WriteTo implements conn.
func (m *MockedUDPConn) WriteTo(p []byte, addr net.Addr) (n int, err error) {
	args := m.Called(p, addr)
	return args.Int(0), args.Error(1)
}

type MockedAddr struct {
	mock.Mock
	network string
	string  string
}

// Network implements net.Addr.
func (m *MockedAddr) Network() string {
	return m.network
}

// String implements net.Addr.
// Subtle: this method shadows the method (Mock).String of MockedAddr.Mock.
func (m *MockedAddr) String() string {
	return m.string
}

func TestSmoke(t *testing.T) {
	mockUDPConn := new(MockedUDPConn)
	server := NewServer(mockUDPConn)
	mockAddr := new(MockedAddr)

	client := server.GetOrMakeClient(mockAddr)
	assert.Equal(t, 0, client.Id())

	mockUDPConn.Mock.On("WriteTo", []byte("hello world"), mockAddr).Return(len([]byte("hello world")), nil)
	server.Send("hello world", client.Id())
	mockUDPConn.AssertCalled(t, "WriteTo", []byte("hello world"), mockAddr)
}
