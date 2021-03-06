package conduit

import (
	"errors"
	"fmt"
	"log"
	"net"
	"sync"
)

type Port struct {
	id int

	ln  net.Listener
	tun *Tunnel

	mu          sync.RWMutex
	connections map[int]*conn

	cancel chan struct{}
}

func OpenPort(port int, t *Tunnel) (*Port, error) {
	addr := fmt.Sprintf("localhost:%d", port)

	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}

	p := &Port{
		id: port,

		ln:  ln,
		tun: t,

		connections: make(map[int]*conn),
	}

	go p.listen()

	return p, nil
}

func (p *Port) Forward(pack Packet) error {
	p.mu.RLock()
	defer p.mu.RUnlock()

	conn, ok := p.connections[pack.ID]
	if !ok {
		return errors.New("[port] unknown destination id")
	}
	return conn.Write(pack)
}

func (p *Port) Send(pack Packet) error {
	if pack.Action == CloseConnection {
		p.mu.Lock()
		delete(p.connections, pack.ID)
		p.mu.Unlock()
	}

	pack.Port = p.id
	return p.tun.Send(pack)
}

func (p *Port) Stop() {
	p.ln.Close()
	p.mu.Lock()
	defer p.mu.Unlock()
	for _, conn := range p.connections {
		conn.c.Close()
	}
}

func (p *Port) listen() {
	var connID int
	for {
		conn, err := p.ln.Accept()
		if err != nil {
			log.Println(err)
			return
		}

		c := NewConn(connID, conn)

		p.mu.Lock()
		p.connections[connID] = c
		p.mu.Unlock()

		pack := Packet{Action: NewConnection, ID: connID}
		p.Send(pack)

		go c.Serve(p.Send)
		connID++
	}
}
