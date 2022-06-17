package tpkt

import "encoding/binary"

type Protocol uint8

func (a Protocol) IsFastpath() bool {
	return a&0x3 == 0
}

func (a Protocol) IsX224() bool {
	return a == 3
}

func (p *protocol) ReceiveProtocol() (Protocol, error) {
	var action Protocol

	return action, binary.Read(p.conn, binary.BigEndian, &action)
}
