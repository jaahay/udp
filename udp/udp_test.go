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

type Message string

const (
	exit Message = "exit"
)

func TestSmoke(t *testing.T) {
	server, err := EmptyServer()
	assert.NoError(t, err)

	serverAddr, err := net.ResolveUDPAddr("udp", ":1053")
	assert.NoError(t, err)

	conn, err := net.DialUDP("udp", nil, serverAddr)
	assert.NoError(t, err)
	defer conn.Close()
	conn.Write([]byte(exit))

	done := make(chan bool)
	go func() {
		received := make([]byte, 1024)
		_, err = conn.Read(received)
		assert.NoError(t, err)
		done <- true
	}()
	<-done

	server.Close()
}
