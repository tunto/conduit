package conduit

const (
	EOF = iota
	Connect
	Data
	Disconnect
)

type Packet struct {
	Action int
	Data   []byte
}
