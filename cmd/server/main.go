package main

import (
	"encoding/gob"
	"log"
	"net"
	"os"
	"time"

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

		// give handler a hopefully unique port to use
		go handle(conn)
	}
}

func handle(conn net.Conn) {
	defer conn.Close()
	gc := &gobConn{
		gob.NewEncoder(conn),
		gob.NewDecoder(conn),
	}
	if err := gc.Encode(conduit.Packet{Action: conduit.Connect}); err != nil {
		log.Println(err)
		return
	}

	portal(gc)
}

func portal(conn *gobConn) {
	ln, err := net.Listen("tcp", "localhost:9000")
	if err != nil {
		// cleanup
		log.Println(err)
		return
	}
	defer ln.Close()

	// TODO: exiting this loop would be nice
	for {
		t1 := time.Now()
		userConn, err := ln.Accept()
		if err != nil {
			// errors on the end user we dont care about. Just log.
			log.Println(err)
			continue
		}

		log.Println("new connection")
		go conn.reader(userConn)
		conn.writer(userConn)
		log.Println("sdfsdfsdf", time.Now().Sub(t1))
	}
}

func (gc *gobConn) writer(conn net.Conn) {
	defer conn.Close()

	for {
		var p conduit.Packet
		if err := gc.Decode(&p); err != nil {
			log.Fatal(err)
		}

		if p.Action != conduit.Data {
			return
		}

		_, err := conn.Write(p.Data)
		if err != nil {
			log.Println(err)
			return
		}
	}
}

func (gc *gobConn) reader(conn net.Conn) {
	defer conn.Close()
	for {
		buf := make([]byte, 4096)

		n, err := conn.Read(buf)
		if err != nil {
			log.Println(err)
			return
		}

		if err := gc.Encode(conduit.Packet{Action: conduit.Data, Data: buf[:n]}); err != nil {
			log.Println(err)
		}
	}
}

type gobConn struct {
	*gob.Encoder
	*gob.Decoder
}
