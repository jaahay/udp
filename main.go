package main

import (
	"github.com/jaahay/udp/udp"
)

func main() {

	_, err := udp.EmptyServer()
	if err != nil {
		panic("failed to create a new empty server")
	}
}
