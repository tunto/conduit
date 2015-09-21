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

	if err := tun.Send(conduit.Packet{Action: conduit.MakePort}); err != nil {
		panic(err)
	}

	if err := tun.Send(conduit.Packet{Action: conduit.MakePort}); err != nil {
		panic(err)
	}

	for {
		pack, err := tun.Read()
		if err != nil {
			panic(err)
		}
		fmt.Printf("%s\n", pack)
	}
}
