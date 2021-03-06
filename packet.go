package conduit

const (
	PICNIC = iota

	Connect
	MakePort
	PortACK
	Data
	Disconnect

	NewConnection

	CloseConnection
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

type Forwarder func(Packet) error
