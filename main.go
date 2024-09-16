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

	server := NewServer(conn)
	for {
		buffer := make([]byte, 1028)
		// number_of_bytes_read...
		_, addr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			continue
		}
		client := server.GetOrMakeClient(addr)
		server.Send(string(buffer), client)
	}
}
