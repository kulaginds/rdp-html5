package tpkt

import (
	"bytes"
	"encoding/binary"
	"io"
)

func (p *Protocol) Receive() (io.Reader, error) {
	tpktPacket := make([]byte, headerLen)

	if _, err := p.conn.Read(tpktPacket); err != nil {
		return nil, err
	}

	tpktPacket = tpktPacket[2:4]

	dataLen := binary.BigEndian.Uint16(tpktPacket)
	dataLen -= uint16(headerLen)

	data := make([]byte, dataLen)

	if _, err := p.conn.Read(data); err != nil {
		return nil, err
	}

	return bytes.NewBuffer(data), nil
}
