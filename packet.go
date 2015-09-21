package conduit

const (
	EOF = iota
	Connect
	MakePort
	PortACK
	Data
	Disconnect
)

var ConnectPacket = Packet{Action: Connect}

type Packet struct {
	Action, ID, Port int

	Data []byte
}

func DataPacket(b []byte) Packet {
	return Packet{
		Action: Data,
		Data:   b,
	}
}

func PortReply(port int) Packet {
	return Packet{
		Action: PortACK,
		Port:   port,
	}
}
