package main

import (
	"fmt"

	"github.com/jaahay/udp/udp"
)

func main() {

	fmt.Println("Hi")
	_, err := udp.EmptyServer()
	if err != nil {
		panic("failed to create a new empty server")
	}
	fmt.Println("Bye")
}
