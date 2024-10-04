package udp

import (
	"fmt"
	"net"
)

// net.PacketConn
type conn interface {
	ReadFromUDP(p []byte) (n int, addr *net.UDPAddr, err error)
	WriteToUDP(p []byte, addr *net.UDPAddr) (n int, err error)
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
	addr *net.UDPAddr
}

type clientSession struct {
	id               int
	clientId         int
	clientConnection *clientConnection
}

type clientConnection struct {
	id        int
	sessionId int
	addr      *net.UDPAddr
	conn      *net.UDPConn
}

func (client client) Id() int {
	return client.id
}

type Server interface {
	GetOrMakeClient(addr *net.UDPAddr) Client
	Send(s string, clientId int) <-chan error
	newClientId() int
}

type server struct {
	id             int
	conn           conn
	clients        map[*net.UDPAddr]*client
	clientSessions map[int]*clientSession
	nextId         int
}

func EmptyServer() (Server, error) {
	addr, err := net.ResolveUDPAddr("udp", "10.0.0.1:2000")
	if err != nil {
		panic("could not resolve udp addr")
	}
	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		panic("could not dial udp")
	}

	defer conn.Close()

	server := &server{
		0,
		conn,
		make(map[*net.UDPAddr]*client),
		make(map[int]*clientSession),
		0,
	}

main:
	for {
		buffer := make([]byte, 1028)
		// number_of_bytes_read...
		_, addr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println(err)
		}
		msg := string(buffer)

		switch msg {
		case "exit":
			break main
		default:
			client := server.GetOrMakeClient(addr)
			server.Send(string(buffer), client.Id())
		}
	}

	return server, nil
}

func NewServer(id int, conn conn, clients map[*net.UDPAddr]*client, clientSessions map[int]*clientSession, nextId int) Server {
	return &server{
		id, conn, clients, clientSessions, nextId,
	}
}

func (server server) newClientId() int {
	server.nextId = server.nextId + 1
	return server.nextId
}
func (server server) GetOrMakeClient(addr *net.UDPAddr) Client {
	c, ok := server.clients[addr]
	if !ok {
		server.clients[addr] = &client{server.newClientId(), addr}
		conn, err := net.DialUDP("udp", nil, addr)
		if err != nil {
			panic(err)
		}
		defer conn.Close()
		clientConnection := &clientConnection{c.id, 0, addr, conn}
		udpSession := &clientSession{clientConnection.id, 0, clientConnection}
		server.clientSessions[c.id] = udpSession
	}
	return c
}

func (server server) Send(s string, clientId int) <-chan error {
	session := server.clientSessions[clientId]
	_, err := server.conn.WriteToUDP([]byte(s), session.clientConnection.addr)
	if err != nil {
		panic(err)
	}
	return nil
}
