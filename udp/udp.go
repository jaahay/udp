package udp

import (
	"fmt"
	"net"
)

type Server interface {
	GetOrMakeClient(addr net.Addr) client
	Send(s string, client client)
}

// net.PacketConn
type Conn interface {
	ReadFrom(p []byte) (n int, addr net.Addr, err error)
	WriteTo(p []byte, addr net.Addr) (n int, err error)

	// Close() error
	// LocalAddr() Addr
	// SetDeadline(t time.Time) error
	// SetReadDeadline(t time.Time) error
	// SetWriteDeadline(t time.Time) error
}

// // net.Addr
// type Addr interface {
// 	Network() string
// 	String() string
// }

type server struct {
	conn           Conn
	clients        map[net.Addr]client
	clientSessions map[client]clientSession
}

type client struct {
	addr net.Addr
}

type clientSession struct {
	clientConnection clientConnection
}

type clientConnection struct {
	addr net.Addr
}

func NewServer(conn Conn) Server {
	return server{
		conn,
		make(map[net.Addr]client),
		make(map[client]clientSession),
	}
}

func (server server) GetOrMakeClient(addr net.Addr) client {
	c, ok := server.clients[addr]
	if !ok {
		c := client{addr}
		server.clients[addr] = c
	}
	_, ok = server.clientSessions[c]
	if !ok {
		clientConnection := clientConnection{addr}
		udpSession := clientSession{clientConnection}
		server.clientSessions[c] = udpSession
		return c
	}
	return c
}

func (server server) Send(s string, client client) {
	session := server.clientSessions[client]
	server.conn.WriteTo([]byte(s), session.clientConnection.addr)
	fmt.Println(s)
}
