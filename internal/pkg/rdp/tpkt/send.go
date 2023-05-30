package tpkt

import (
	"bytes"
	"encoding/binary"
)

func (p *Protocol) Send(pduData []byte) error {
	buf := bytes.NewBuffer(make([]byte, 0, headerLen+len(pduData)))
	dataLen := uint16(headerLen + len(pduData))

	buf.Write([]byte{
		0x03, // TPKT version number
		0x00, // reserved for further study
	})

	// TPKT length
	_ = binary.Write(buf, binary.BigEndian, dataLen)

	buf.Write(pduData)

	if _, err := p.conn.Write(buf.Bytes()); err != nil {
		return err
	}

	return nil
}
