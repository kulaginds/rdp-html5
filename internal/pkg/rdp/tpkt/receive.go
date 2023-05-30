package tpkt

import (
	"bytes"
	"encoding/binary"
	"io"
)

func (p *Protocol) Receive() (io.Reader, error) {
	if p.fastpathEnabled {
		return p.receive(headerLen - 1)
	}

	return p.receive(headerLen)
}

func (p *Protocol) receive(len int) (io.Reader, error) {
	tpktPacket := make([]byte, len)

	if _, err := p.conn.Read(tpktPacket); err != nil {
		return nil, err
	}

	if len == 4 {
		tpktPacket = tpktPacket[2:4]
	} else {
		tpktPacket = tpktPacket[1:3]
	}

	dataLen := binary.BigEndian.Uint16(tpktPacket)
	dataLen -= uint16(len)

	data := make([]byte, dataLen)

	if _, err := p.conn.Read(data); err != nil {
		return nil, err
	}

	return bytes.NewBuffer(data), nil
}
