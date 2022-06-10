package tpkt

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"io"
	"log"
)

func (p *protocol) Receive() (io.Reader, error) {
	tpktPacket := make([]byte, headerLen)

	if _, err := p.conn.Read(tpktPacket); err != nil {
		return nil, err
	}

	dataLen := binary.BigEndian.Uint16(tpktPacket[2:4])
	dataLen -= headerLen

	data := make([]byte, dataLen)

	if _, err := p.conn.Read(data); err != nil {
		return nil, err
	}

	log.Printf("TPKT: Receive: %s\n", hex.EncodeToString(data))

	return bytes.NewBuffer(data), nil
}
