package vsocket

import (
	"fmt"
	"net"
	"net/netip"

	"lecoin/pkg/protocol"
)

type (
	NetChan      chan protocol.Serializable
	ChanMap      map[protocol.MessageType][]NetChan
	ConstructMap map[protocol.MessageType]protocol.MessageConstructor
)

type VSocket struct {
	lport        uint16
	conn         net.Conn
	channels     ChanMap // maybe lock, but should be defined before run anyways
	constructors ConstructMap
}

func NewVSocket(lport uint16, rport uint16) (*VSocket, error) {
	laddr := netip.AddrPortFrom(netip.MustParseAddr("127.0.0.1"), lport)
	raddr := netip.AddrPortFrom(netip.MustParseAddr("127.0.0.1"), rport)

	conn, err := net.DialTCP("tcp4", net.TCPAddrFromAddrPort(laddr), net.TCPAddrFromAddrPort(raddr))
	if err != nil {
		return nil, err
	}

	return &VSocket{
		lport:        lport,
		conn:         conn,
		channels:     make(ChanMap),
		constructors: make(ConstructMap),
	}, nil
}

func (vskt *VSocket) RegisterChannel(msgType protocol.MessageType, constructor protocol.MessageConstructor) NetChan {
	// even if something is already registered, we allow overwriting
	chans, ok := vskt.channels[msgType]
	if !ok {
		vskt.channels[msgType] = make([]NetChan, 0)
		chans = vskt.channels[msgType]
	}
	ret := make(NetChan)
	vskt.channels[msgType] = append(chans, ret)
	vskt.constructors[msgType] = constructor
	return ret
}

func (vskt *VSocket) RegisterChannels(msgTypes []protocol.MessageType, constructors []protocol.MessageConstructor) map[protocol.MessageType]NetChan {
	if len(msgTypes) != len(constructors) {
		panic("invalid call to register chan")
	}
	clientHandlers := make(map[protocol.MessageType]NetChan)
	for i, msgType := range msgTypes {
		clientHandlers[msgType] = vskt.RegisterChannel(msgType, constructors[i])
	}

	return clientHandlers
}

func (vskt *VSocket) Run() {
	for {
		packet := protocol.ParsePacket(vskt.conn)
		if packet == nil {
			break
		}

		// determine message from socket protocols
		msgType := packet.MsgType
		outChans, ok := vskt.channels[msgType]
		if !ok {
			fmt.Printf("invalid message type %d\n", msgType)
			continue
		}
		constructor, ok := vskt.constructors[msgType]
		if !ok {
			fmt.Printf("message type %d does not have constructor\n", msgType)
			continue
		}
		// construct and deserialize struct
		obj := constructor()
		obj.Unmarshal(packet.Msg)

		for _, outChan := range outChans {
			go func() {
				outChan <- obj
			}()
		}
	}
}

func (vskt *VSocket) Send(msg protocol.Serializable, rport uint16) error {
	buf := msg.FullMarshal()
	packet := protocol.Packet{
		SndrPort: vskt.lport,
		RcvrPort: rport,
		MsgType:  msg.MsgType(),
		Msg:      buf,
	}
	pBytes := packet.Marshal()
	_, err := vskt.conn.Write(pBytes)
	return err
}

func (vskt *VSocket) Broadcast(msg protocol.Serializable) error {
	return vskt.Send(msg, 0)
}
