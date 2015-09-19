package main

import (
	"encoding/gob"
	"log"
	"net"

	"go.tun.to/conduit"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		log.Fatal("Connection error", err)
	}
	defer conn.Close()

	gc := &gobConn{
		conn,
		gob.NewEncoder(conn),
		gob.NewDecoder(conn),
	}

	state := gc.connector
	for {
		if state = state(); state == nil {
			return
		}
	}
}

type readerFunc func() readerFunc

type gobConn struct {
	net.Conn
	*gob.Encoder
	*gob.Decoder
}

func (gc *gobConn) connector() readerFunc {
	log.Println("waiting for a connection")

	var p conduit.Packet
	if err := gc.Decode(&p); err != nil {
		log.Fatal(err)
	}

	if p.Action != conduit.Connect {
		panic("protocol error: unexpected action")
	}

	conn, err := net.Dial("tcp", "localhost:8000")
	if err != nil {
		log.Println(err)
		return gc.connector
	}

	go gc.reader(conn)
	return gc.writer(conn)
}

func (gc *gobConn) writer(conn net.Conn) readerFunc {
	defer conn.Close()

	for {
		var p conduit.Packet

		// TODO: this blocks the whole connection thread... dont do that.
		if err := gc.Decode(&p); err != nil {
			log.Fatal(err)
		}

		switch p.Action {
		case conduit.Data:
			_, err := conn.Write(p.Data)
			if err != nil {
				log.Println(err)
				return gc.connector
			}
		case conduit.Disconnect:
			return gc.connector
		default:
			panic("protocol error: unexpected action")
		}
	}

	return gc.connector
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
			return
		}
	}
}
