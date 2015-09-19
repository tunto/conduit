package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

var (
	portIN  = flag.Int("in", 8080, "input port to listen on")
	portOUT = flag.Int("out", 8000, "output port to proxy for")
)

func main() {
	flag.Parse()

	if *portIN < 1 || 65535 < *portIN {
		log.Println("input port out of range (1-65535)")
		os.Exit(1)
	}
	if *portOUT < 1 || 65535 < *portOUT {
		log.Println("output port out of range (1-65535)")
		os.Exit(1)
	}

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", *portIN))
	if err != nil {
		log.Fatal(err)
	}

	defer l.Close()

	for {
		if conn, err := l.Accept(); err == nil {
			go handler(conn, fmt.Sprintf(":%d", *portOUT))
		}
	}
}

func handler(conn net.Conn, addr string) {
	defer conn.Close()

	proxyConn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Println(err)
		return
	}

	go io.Copy(conn, proxyConn)
	io.Copy(proxyConn, conn)
}
