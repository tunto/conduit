package main

import (
	"log"
	"net"
	"os"

	"go.tun.to/conduit"
)

func main() {
	ln, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer ln.Close()

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println(err)
			continue
		}

		t := conduit.NewTunnel(conn)
		go t.Open()
	}
}
