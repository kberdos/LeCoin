package vswitch

import (
	"fmt"
	"lecoin/pkg/protocol"
	"net"
	"net/netip"
	"strconv"
	"sync"
)

type VSwitch struct {
	listener net.Listener
	clients  map[uint16]net.Conn
	mu       sync.Mutex
}

func NewVSwitch(port int) *VSwitch {
	if port < 0 || 1<<16 <= port {
		return nil
	}

	addr := netip.AddrPortFrom(netip.MustParseAddr("127.0.0.1"), uint16(port))

	l, err := net.ListenTCP("tcp4", net.TCPAddrFromAddrPort(addr))
	if err != nil {
		return nil
	}

	return &VSwitch{
		listener: l,
		clients:  make(map[uint16]net.Conn),
	}
}

func (vs *VSwitch) Run() {
	for {
		conn, _ := vs.listener.Accept()
		vs.mu.Lock()
		_, rportstr, err := net.SplitHostPort(conn.RemoteAddr().String())
		if err != nil {
			vs.mu.Unlock()
			continue
		}
		rport, _ := strconv.Atoi(rportstr) // won't error, too bad

		_, ok := vs.clients[uint16(rport)]
		if ok {
			// port already exists, no!
			vs.mu.Unlock()
			continue
		}

		vs.clients[uint16(rport)] = conn
		go handleConn(conn.(*net.TCPConn), uint16(rport), vs)
		vs.mu.Unlock()
	}
}

func handleConn(conn *net.TCPConn, port uint16, vs *VSwitch) {
	fmt.Printf("> Host at port %d has connected\n", conn.RemoteAddr().(*net.TCPAddr).AddrPort().Port())
	for {
		packet := protocol.ParsePacket(conn)
		if packet == nil {
			break
		}
		vs.mu.Lock()
		// check if rport exists
		switch packet.RcvrPort {
		case 0:
			// reps a broadcast
			for rport, rconn := range vs.clients {
				if rport != packet.SndrPort {
					Forward(rconn, packet)
				}
			}
		default:
			rconn, ok := vs.clients[packet.RcvrPort]
			if ok {
				Forward(rconn, packet)
			}
		}
		vs.mu.Unlock()
	}

	conn.SetLinger(0) // NOTE: this is so inbound port is actually reusable after closing
	conn.Close()

	vs.mu.Lock()
	delete(vs.clients, port)
	vs.mu.Unlock()
}

func Forward(conn net.Conn, p *protocol.Packet) {
	packetBytes := p.Marshal()
	conn.Write(packetBytes) // FIXME: too lazy to error check this
}
