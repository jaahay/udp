package udp

import (
	"fmt"
)

type Config struct {
}

// net.PacketConn
type Conn interface {
	ReadFrom(p []byte) (n int, addr Addr, err error)
	WriteTo(p []byte, addr Addr) (n int, err error)

	// Close() error
	// LocalAddr() Addr
	// SetDeadline(t time.Time) error
	// SetReadDeadline(t time.Time) error
	// SetWriteDeadline(t time.Time) error
}

// net.Addr
type Addr interface {
	Network() string
	String() string
}

type Server struct {
	conn           Conn
	clients        map[Addr]Client
	clientSessions map[Client]ClientSession
}

type Client struct {
	addr Addr
}

type ClientSession struct {
	clientConnection ClientConnection
	// session          func()
}

type ClientConnection struct {
	addr Addr
}

func NewServer(conn Conn) Server {
	return Server{
		conn,
		make(map[Addr]Client),
		make(map[Client]ClientSession),
	}
}

func (server *Server) GetOrMakeClient(addr Addr) Client {
	// todo: authentication, then;
	client, ok := server.clients[addr]
	if !ok {
		client := Client{addr}
		server.clients[addr] = client
	}
	_, ok = server.clientSessions[client]
	if !ok {
		clientConnection := ClientConnection{addr}
		// todo: client-based sessions; not Addr-based
		udpSession := ClientSession{
			clientConnection: clientConnection,
			// session:          go func (),
		}
		server.clientSessions[client] = udpSession
		return client
	}
	return client
}

func (server *Server) Send(s string, client Client) {
	session := server.clientSessions[client]
	server.conn.WriteTo([]byte(s), session.clientConnection.addr)
	fmt.Println(s)
}
