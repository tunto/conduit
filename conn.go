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

func (c *conn) Serve(f Forwarder) {
	defer c.c.Close()

	for {
		buf := make([]byte, 4096)

		n, err := c.c.Read(buf)
		if err != nil {
			log.Printf("[conn:%d] read error: %s", c.id, err)
			f(Packet{Action: CloseConnection, ID: c.id})
			return
		}

		pack := DataPacket(buf[:n])
		pack.ID = c.id

		if err := f(pack); err != nil {
			log.Printf("[conn:%d] port send error: %s", c.id, err)
			return
		}
	}
}

func (c *conn) Write(p Packet) error {
	_, err := c.c.Write(p.Data)
	return err
}
