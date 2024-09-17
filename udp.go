package main

import (
	"fmt"
	"net"
)

type Server interface {
}
type server struct {
	conn           net.PacketConn
	clients        map[net.Addr]Client
	clientSessions map[Client]ClientSession
}

type Client struct {
	addr net.Addr
}

type ClientSession struct {
	clientConnection ClientConnection
	// session          func()
}

type ClientConnection struct {
	addr net.Addr
}

func NewServer(conn *net.UDPConn) Server {
	return server{
		conn,
		make(map[net.Addr]Client),
		make(map[Client]ClientSession),
	}
}

func (server *server) GetOrMakeClient(addr net.Addr) Client {
	// todo: authentication, then;
	client, ok := server.clients[addr]
	if !ok {
		client := Client{addr}
		server.clients[addr] = client
	}
	_, ok = server.clientSessions[client]
	if !ok {
		clientConnection := ClientConnection{addr}
		// todo: client-based sessions; not net.Addr-based
		udpSession := ClientSession{
			clientConnection: clientConnection,
			// session:          go func (),
		}
		server.clientSessions[client] = udpSession
		return client
	}
	return client
}

func (server *server) Send(s string, client Client) {
	session := server.clientSessions[client]
	server.conn.WriteTo([]byte(s), session.clientConnection.addr)
	fmt.Println(s)
}
