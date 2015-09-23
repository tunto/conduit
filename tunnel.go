package conduit

import (
	"encoding/gob"
	"log"
	"net"
)

type Tunnel struct {
	c net.Conn
	*gob.Encoder
	*gob.Decoder

	ports map[int]*Port
}

func OpenTunnel(addr string) (*Tunnel, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}

	return NewTunnel(conn), nil
}

func NewTunnel(conn net.Conn) *Tunnel {
	return &Tunnel{
		conn,
		gob.NewEncoder(conn),
		gob.NewDecoder(conn),
		make(map[int]*Port),
	}
}

// TODO: this is server specific and should be moved somewhere else.
func (t *Tunnel) Open() {
	defer t.Close()

	nextPort := 9000
	for {
		pack, err := t.Read()
		if err != nil {
			log.Printf("[tunnel] connection error: %s", err)
			return
		}

		switch pack.Action {
		case Data:
			p, ok := t.ports[pack.Port]
			if !ok {
				log.Println("[tunnel] bad port")
				continue
			}
			if err := p.Forward(pack); err != nil {
				log.Printf("[tunnel] port forward error: %s", err)
				continue
			}
		case MakePort:
			p, err := OpenPort(nextPort, t)
			if err != nil {
				log.Printf("[tunnel] port create error: %s", err)
				continue
			}

			if err := t.Send(PortReply(nextPort)); err != nil {
				log.Printf("[tunnel] send error: %s", err)
				continue
			}

			t.ports[nextPort] = p

			nextPort++
		}

	}
}

func (t *Tunnel) Read() (pack Packet, err error) {
	err = t.Decode(&pack)
	return
}

func (t *Tunnel) Send(pack Packet) error {
	if pack.Action == PICNIC {
		panic("tunnel: packet sent with unset action type")
	}

	return t.Encode(pack)
}

func (t Tunnel) Close() {
	t.c.Close()
	for _, p := range t.ports {
		p.Stop()
	}
}
