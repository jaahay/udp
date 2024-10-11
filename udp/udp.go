package udp

import (
	"fmt"
	"net"
	"sync"
)

// net.PacketConn
type conn interface {
	ReadFromUDP(p []byte) (n int, addr *net.UDPAddr, err error)
	ReadFrom(p []byte) (n int, addr net.Addr, err error)
	WriteToUDP(p []byte, addr *net.UDPAddr) (n int, err error)
	Close() error
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
	clientConnection clientConnection
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
	Send(s string, clientId int) error
	newClientId() int
	Wait()
}

type server struct {
	id             int
	conn           conn
	clients        map[*net.UDPAddr]client
	clientSessions map[int]clientSession
	nextId         int
	starting       *sync.Mutex
	*sync.WaitGroup
}

func EmptyServer() (Server, error) {
	wg := &sync.WaitGroup{}
	wg.Add(1)

	available := &sync.Mutex{}
	available.Lock()

	addr, err := net.ResolveUDPAddr("udp", ":1053")
	if err != nil {
		panic("could not resolve udp addr")
	}
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		panic("could not dial udp")
	}
	server := &server{
		0,
		conn,
		make(map[*net.UDPAddr]client),
		make(map[int]clientSession),
		0,
		available,
		wg,
	}

	fmt.Println("Server starting")
	go server.start(available)
	fmt.Println("Server has started")

	// fmt.Println("UDP writing \"exit\"")
	// server.conn.WriteToUDP([]byte("exit"), addr)
	// fmt.Println("UDP has written \"exit\"")
	return server, nil
}

func (server server) start(available *sync.Mutex) error {
	defer server.conn.Close()
	defer server.Done()

	available.Unlock()
	fmt.Println("Server finished starting")
main:
	for {
		buffer := make([]byte, 1024)
		fmt.Println("Before:" + string(buffer))
		// number_of_bytes_read...
		_, addr, err := server.conn.ReadFromUDP(buffer)
		fmt.Println("After:" + string(buffer))
		fmt.Println("message received")
		if err != nil {
			panic(err)
		}
		msg := string(buffer)
		fmt.Println("server received: \"" + msg + "\"")

		switch msg {
		case "exit":
			break main
		default:
			client := server.GetOrMakeClient(addr)
			server.Send(string(buffer), client.Id())
		}
	}
	return nil
}

func NewServer(id int, conn conn, clients map[*net.UDPAddr]client, clientSessions map[int]clientSession, nextId int) Server {
	return &server{
		id, conn, clients, clientSessions, nextId, &sync.Mutex{}, &sync.WaitGroup{},
	}
}

func (server server) newClientId() int {
	server.nextId = server.nextId + 1
	return server.nextId
}
func (server server) GetOrMakeClient(addr *net.UDPAddr) Client {
	c, ok := server.clients[addr]
	if !ok {
		c = client{server.newClientId(), addr}
		server.clients[addr] = c
		conn, err := net.DialUDP("udp", nil, addr)
		if err != nil {
			panic(err)
		}
		// defer conn.Close()
		clientConnection := clientConnection{c.id, 0, addr, conn}
		clientSession := clientSession{clientConnection.id, 0, clientConnection}
		server.clientSessions[c.id] = clientSession
	}
	return c
}

func (server server) Send(s string, clientId int) error {
	session := server.clientSessions[clientId]
	_, err := server.conn.WriteToUDP([]byte(s), session.clientConnection.addr)
	if err != nil {
		panic(err)
	}
	return nil
}
