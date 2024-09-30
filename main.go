package main

import (
	"net"

	"github.com/jaahay/udp/udp"
)

func main() {

	server := udp.EmptyServer()

	clientAddr, err := net.ResolveUDPAddr("udp", "10.0.0.1:2000")
	if err != nil {
		panic("could not resolve udp addr")
	}
	conn, err := net.DialUDP("udp", nil, clientAddr)
	if err != nil {
		panic("could not dial udp")
	}
	defer conn.Close()

	server.GetOrMakeClient(clientAddr)

	// for {
	// 	buffer := make([]byte, 1028)
	// 	// number_of_bytes_read...
	// 	_, addr, err := conn.ReadFromUDP(buffer)
	// 	if err != nil {
	// 		continue
	// 	}
	// 	client := server.GetOrMakeClient(addr)
	// 	server.Send(string(buffer), client.Id())
	// }
}
