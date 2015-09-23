package main

import (
	"fmt"

	"go.tun.to/conduit"
)

func main() {
	tun, err := conduit.OpenTunnel("localhost:8080")
	if err != nil {
		panic(err)
	}

	srvman := conduit.NewServiceManager(tun)

	port, err := srvman.AddService("localhost:8000")

	if err != nil {
		panic(err)
	}
	fmt.Println(port)

	// keep program alive
	<-make(chan struct{})
}
