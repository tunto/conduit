package conduit

import (
	"log"
	"net"
)

type conn struct {
	id int
	c  net.Conn
}

func NewConn(id int, c net.Conn) *conn {
	return &conn{id, c}
}

func (c *conn) Serve(p *Port) {
	defer c.c.Close()

	for {
		buf := make([]byte, 4096)

		n, err := c.c.Read(buf)
		if err != nil {
			log.Printf("[conn:%d] read error: %s", c.id, err)
			p.Close(c.id)
			return
		}

		pack := DataPacket(buf[:n])
		pack.ID = c.id

		if err := p.Send(pack); err != nil {
			log.Printf("[conn:%d] port send error: %s", c.id, err)
			p.Close(c.id)
			return
		}
	}
}

func (c *conn) Write(p Packet) error {
	_, err := c.c.Write(p.Data)
	return err
}
