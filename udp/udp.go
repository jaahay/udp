package udp

import (
	"fmt"
	"net"
)

type Server interface {
	GetOrMakeClient(addr net.Addr) Client
	Send(s string, clientId int)
	NewClientId() int
}

type Client interface {
	Id() int
}

// net.PacketConn
type conn interface {
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
	id             int
	conn           conn
	clients        map[net.Addr]client
	clientSessions map[int]clientSession
	nextId         int
}

type client struct {
	id   int
	addr net.Addr
}

type clientSession struct {
	id               int
	clientId         int
	clientConnection clientConnection
}

type clientConnection struct {
	id        int
	sessionId int
	addr      net.Addr
}

func NewServer(conn conn) Server {
	return server{
		0,
		conn,
		make(map[net.Addr]client),
		make(map[int]clientSession),
		0,
	}
}

func (server server) NewClientId() int {
	server.nextId = server.nextId + 1
	return server.nextId
}
func (server server) GetOrMakeClient(addr net.Addr) Client {
	c, ok := server.clients[addr]
	if !ok {
		c := client{server.nextId, addr}
		server.clients[addr] = c
	}
	_, ok = server.clientSessions[c.id]
	if !ok {
		clientConnection := clientConnection{c.id, 0, addr}
		udpSession := clientSession{clientConnection.id, 0, clientConnection}
		server.clientSessions[c.id] = udpSession
		return c
	}
	return c
}

func (client client) Id() int {
	return client.id
}

func (server server) Send(s string, clientId int) {
	session := server.clientSessions[clientId]
	server.conn.WriteTo([]byte(s), session.clientConnection.addr)
	fmt.Println(s)
}
