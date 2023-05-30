package tpkt

import "encoding/binary"

type ProtocolCode uint8

func (a ProtocolCode) IsFastpath() bool {
	return a&0x3 == 0
}

func (a ProtocolCode) IsX224() bool {
	return a == 3
}

func (p *Protocol) ReceiveProtocol() (ProtocolCode, error) {
	var action ProtocolCode

	return action, binary.Read(p.conn, binary.BigEndian, &action)
}
