package main

import (
	"fmt"
	"net"
)

type UDPServer struct {
	playerIds         map[net.Addr]PlayerId
	playerConnections map[PlayerId]PlayerConnection
	sessions          map[PlayerConnection]UDPSession
}

type PlayerId struct {
	addr net.Addr
}

type PlayerConnection struct {
}

type UDPSession struct {
}

func (udpServer *UDPServer) Send(b []byte, addr net.Addr) {
	session := udpServer.getOrMakeSession(addr)
	fmt.Println(addr.String)
}

func (udpServer UDPServer) getOrMakeSession(addr net.Addr) UDPSession {
	playerId, ok := udpServer.playerIds[addr]
	if !ok {
		playerId := PlayerId{addr}
		playerConnection := PlayerConnection{}
		udpSession := UDPSession{}
		udpServer.playerIds[addr] = playerId
		udpServer.playerConnections[playerId] = playerConnection
		udpServer.sessions[playerConnection] = udpSession
		return udpSession
	}
	return udpServer.sessions[udpServer.playerConnections[udpServer.playerIds[playerId.addr]]]
}
