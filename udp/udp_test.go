package udp

import (
	"fmt"
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

func TestSmoke(t *testing.T) {
	fmt.Println("fmt Begin")
	server, err := EmptyServer()
	assert.NoError(t, err)
	fmt.Println("fmt Begun")
	serverAddr, err := net.ResolveUDPAddr("udp", ":1053")
	if err != nil {
		panic("could not resolve server udp addr")
	}

	conn, err := net.DialUDP("udp", nil, serverAddr)
	if err != nil {
		panic("could not dial server")
	}

	// clientAddr, err := net.ResolveUDPAddr("udp", ":1054")
	// if err != nil {
	// 	panic("could not resolve client udp addr")
	// }
	// conn, err := net.DialUDP("udp", nil, clientAddr)
	// if err != nil {
	// 	panic("could not dial udp")
	// }

	defer conn.Close()

	fmt.Println("fmt writing \"exit\"")
	conn.Write([]byte("exit"))

	received := make([]byte, 1024)
	_, err = conn.Read(received)
	if err != nil {
		panic("Read data failed:")
	}
	fmt.Println("client received: \"" + string(received) + "\"")

	fmt.Println("awaiting...")

	server.Wait()
}
