package protocol

import (
	"encoding/binary"
	"io"
)

type Packet struct {
	SndrPort uint16
	RcvrPort uint16
	MsgType  MessageType
	MsgLen   uint32
	Msg      []byte
}

type MessageType uint16

// create a new Serializable struct (can be empty)
type MessageConstructor func() Serializable

func (p *Packet) Marshal() []byte {
	p.MsgLen = uint32(len(p.Msg)) // make sure this is facts
	buf := make([]byte, 10+len(p.Msg))

	binary.BigEndian.PutUint16(buf[:2], p.SndrPort)
	binary.BigEndian.PutUint16(buf[2:4], p.RcvrPort)
	binary.BigEndian.PutUint16(buf[4:6], uint16(p.MsgType))
	binary.BigEndian.PutUint32(buf[6:10], p.MsgLen)
	copy(buf[10:10+len(p.Msg)], p.Msg)

	return buf
}

func ParsePacket(rdr io.Reader) *Packet {
	p := Packet{}
	metadata := make([]byte, 10)
	n, err := io.ReadFull(rdr, metadata)
	if err != nil || n != 10 {
		return nil
	}

	p.SndrPort = binary.BigEndian.Uint16(metadata[:2])
	p.RcvrPort = binary.BigEndian.Uint16(metadata[2:4])
	p.MsgType = MessageType(binary.BigEndian.Uint16(metadata[4:6]))
	p.MsgLen = binary.BigEndian.Uint32(metadata[6:10])

	p.Msg = make([]byte, p.MsgLen)
	n, err = io.ReadFull(rdr, p.Msg)
	if err != nil || uint32(n) != p.MsgLen {
		return nil
	}

	return &p
}
