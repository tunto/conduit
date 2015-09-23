package conduit

import (
	"errors"
	"log"
	"net"
	"time"
)

type ServiceManager struct {
	serviceList map[int]string
	connections map[int]*conn

	t *Tunnel

	portCh chan int
}

func NewServiceManager(t *Tunnel) *ServiceManager {
	sm := &ServiceManager{
		make(map[int]string),
		make(map[int]*conn),
		t,

		// keep one extra port before blocking just in case
		// there's an add timeout for retries
		make(chan int, 1),
	}

	go sm.run()

	return sm
}

func (sm *ServiceManager) run() {
	// TODO: close all service connections
	// defer c.c.Close()

	for {
		pack, err := sm.t.Read()

		if err != nil {
			log.Printf("[tunnel] connection error: %s", err)
			return
		}

		switch pack.Action {
		case PortACK:
			sm.portCh <- pack.Port
		case NewConnection:
			sm.addConnection(pack)
		default:
			conn, ok := sm.connections[pack.ID]
			if ok {
				conn.Write(pack)
			}
			/// forward to service
		}

	}
}

func (sm *ServiceManager) addConnection(pack Packet) {
	addr, ok := sm.serviceList[pack.Port]
	if !ok {
		log.Println("cant find service on that port")
		pack.Action = CloseConnection
		sm.Send(pack)
		return
	}

	_, ok = sm.connections[pack.ID]
	if ok {
		panic("protocol error: connection id in use. cannot add connection.")
	}

	conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Println(err)
		pack.Action = CloseConnection
		sm.Send(pack)
		return
	}

	c := NewConn(pack.ID, conn)
	sm.connections[pack.ID] = c

	port := pack.Port
	go c.Serve(func(pack Packet) error {
		pack.Port = port
		return sm.Send(pack)
	})

	// do we need to send something back?
}

func (sm *ServiceManager) AddService(addr string) (port int, err error) {
	sm.Send(Packet{Action: MakePort})

	t := time.NewTimer(10 * time.Second)
	defer t.Stop()

	select {
	case <-t.C:
		return 0, errors.New("[service manager] request timeout")
	case p, ok := <-sm.portCh:
		if !ok {
			return 0, errors.New("[service manager] manager closed")
		}
		port = p
	}

	sm.serviceList[port] = addr

	return port, nil
}

func (sm *ServiceManager) Send(pack Packet) error {
	return sm.t.Send(pack)
}
