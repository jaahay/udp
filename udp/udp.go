package udp

import (
	"fmt"
	"net"
	"sync"
)

// net.PacketConn
type conn interface {
	// Read([]byte) (int, error)
	// Write([]byte) (int, error)

	ReadFrom(p []byte) (n int, addr net.Addr, err error)
	WriteTo([]byte, net.Addr) (int, error)

	// ReadFromUDP(p []byte) (n int, addr *net.UDPAddr, err error)
	// WriteToUDP(p []byte, addr *net.UDPAddr) (n int, err error)

	Close() error
}

// // net.Addr
// type Addr interface {
// 	Network() string
// 	String() string
// }

type client struct {
	id   int
	addr *net.Addr
}

type Client interface {
	Id() int
}

type clientSession struct {
	id               int
	clientId         int
	clientConnection clientConnection
	*sync.WaitGroup
}

type clientConnection struct {
	id        int
	sessionId int
	addr      *net.Addr
}

func (client client) Id() int {
	return client.id
}

type server struct {
	id             int
	conn           conn
	clients        map[*net.Addr]client
	clientSessions map[int]clientSession
	nextId         int
	ready          *sync.Mutex
	*sync.WaitGroup
}

type Server interface {
	GetOrMakeClient(addr *net.Addr) Client
	Send(s string, clientId int) error
	newClientId() int
	Wait()
	Close()
}

func EmptyServer() (Server, error) {

	fmt.Println("Server starting")

	// addr, err := net.ResolveUDPAddr("udp", ":1053")
	// if err != nil {
	// 	panic("could not resolve udp addr")
	// }
	// conn, err := net.ListenUDP("udp", addr)
	// if err != nil {
	// 	panic("could not dial udp")
	// }

	conn, err := net.ListenPacket("udp", ":1053")
	if err != nil {
		panic("could not listen packet")
	}

	server := &server{
		0,
		conn,
		make(map[*net.Addr]client),
		make(map[int]clientSession),
		0,
		&sync.Mutex{},
		&sync.WaitGroup{},
	}

	// defer server.conn.Close()
	// defer server.Done()
	go func() {
	main:
		for {
			server.ready.Lock()
			buffer := make([]byte, 1024)
			// fmt.Println("Before:" + string(buffer))
			// number_of_bytes_read...
			n, addr, err := server.conn.ReadFrom(buffer)
			if err != nil {
				panic(err)
			}
			// fmt.Println("After:" + string(buffer))
			// fmt.Println("message received")
			msg := string(buffer[:n])
			fmt.Println("server handling: \"" + msg + "\"")

			switch msg {
			case "exit":
				client := server.GetOrMakeClient(&addr)
				server.Send(string(buffer), client.Id())
				fmt.Println("breaking main server")

				server.ready.Unlock()
				break main
			default:
				client := server.GetOrMakeClient(&addr)
				server.Send(string(buffer), client.Id())
				server.ready.Unlock()
			}
		}
	}()
	fmt.Println("Server has started")

	// fmt.Println("UDP writing \"exit\"")
	// server.conn.WriteToUDP([]byte("exit"), addr)
	// fmt.Println("UDP has written \"exit\"")
	return server, nil
}

func NewServer(id int, conn conn, clients map[*net.Addr]client, clientSessions map[int]clientSession, nextId int) Server {
	return &server{
		id, conn, clients, clientSessions, nextId, &sync.Mutex{}, &sync.WaitGroup{},
	}
}

func (server server) newClientId() int {
	server.nextId = server.nextId + 1
	return server.nextId
}

func (server server) GetOrMakeClient(addr *net.Addr) Client {
	c, ok := server.clients[addr]
	if !ok {
		c = client{server.newClientId(), addr}
		server.clients[addr] = c
		// defer conn.Close()
		clientConnection := clientConnection{c.id, 0, addr}
		clientSession := clientSession{clientConnection.id, 0, clientConnection, &sync.WaitGroup{}}
		server.clientSessions[c.id] = clientSession
	}
	return c
}

func (server server) Send(s string, clientId int) error {
	session := server.clientSessions[clientId]
	addr := session.clientConnection.addr
	session.Add(1)
	defer session.Done()

	fmt.Println("server sending msg: " + s)
	_, err := server.conn.WriteTo([]byte(s), *addr)
	if err != nil {
		panic(err)
	}
	fmt.Println("server sent: " + s)
	return nil
}

func (server server) Wait() {
	await := &sync.WaitGroup{}
	for _, clientSession := range server.clientSessions {
		await.Add(1)
		go func() {
			clientSession.Wait()
			await.Done()
		}()
	}
	await.Wait()
}

func (server server) Close() {
	server.Wait()
	server.conn.Close()
}
