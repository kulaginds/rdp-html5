package x224

import (
	"encoding/binary"
	"io"
)

func (p *protocol) Receive() (io.Reader, error) {
	const fixedPartLen uint8 = 0x02

	wire, err := p.tpktConn.Receive()
	if err != nil {
		return nil, err
	}

	var li uint8
	binary.Read(wire, binary.LittleEndian, &li)

	if li != fixedPartLen {
		return nil, ErrWrongDataLength
	}

	packetData := make([]byte, li)
	wire.Read(packetData)

	return wire, nil
}
