package udp

import (
	"net"
)

// net.PacketConn
type conn interface {
	ReadFrom(p []byte) (n int, addr net.Addr, err error)
	WriteTo(p []byte, addr net.Addr) (n int, err error)
}

// // net.Addr
// type Addr interface {
// 	Network() string
// 	String() string
// }

type Client interface {
	Id() int
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
	conn      net.Conn
}

func (client client) Id() int {
	return client.id
}

type Server interface {
	GetOrMakeClient(addr net.Addr) Client
	Send(s string, clientId int)
	newClientId() int
}

type server struct {
	id             int
	conn           conn
	clients        map[net.Addr]client
	clientSessions map[int]clientSession
	nextId         int
}

func EmptyServer() Server {
	addr, err := net.ResolveUDPAddr("udp", "10.0.0.1:2000")
	if err != nil {
		panic("could not resolve udp addr")
	}
	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		panic("could not dial udp")
	}
	defer conn.Close()
	return server{
		0,
		conn,
		make(map[net.Addr]client),
		make(map[int]clientSession),
		0,
	}
}

func NewServer(id int, conn conn, clients map[net.Addr]client, clientSessions map[int]clientSession, nextId int) Server {
	return server{
		id, conn, clients, clientSessions, nextId,
	}
}

func (server server) newClientId() int {
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
		conn, err := net.DialUDP("udp", nil, addr.(*net.UDPAddr))
		if err != nil {
			panic("could not dial udp")
		}
		defer conn.Close()
		clientConnection := clientConnection{c.id, 0, addr, conn}
		udpSession := clientSession{clientConnection.id, 0, clientConnection}
		server.clientSessions[c.id] = udpSession
		return c
	}
	return c
}

func (server server) Send(s string, clientId int) {
	session := server.clientSessions[clientId]
	server.conn.WriteTo([]byte(s), session.clientConnection.addr)
}
