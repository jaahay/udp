package main

import "net"

func main() {
	addr, err := net.ResolveUDPAddr("udp", "10.0.0.1:2000")
	if err != nil {
		return
	}
	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		return
	}
	defer conn.Close()

	for {
		buffer := make([]byte, 1028)
		bytes, addr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			continue
		}
		udpServer := UDPServer{}
		udpServer.Send(buffer[:bytes], addr)
	}
}
