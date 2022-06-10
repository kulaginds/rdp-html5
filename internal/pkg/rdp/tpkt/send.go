package tpkt

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"log"
)

func (p *protocol) Send(pduData []byte) error {
	buf := bytes.NewBuffer(make([]byte, 0, headerLen+len(pduData)))
	dataLen := uint16(headerLen + len(pduData))

	buf.Write([]byte{
		0x03, // TPKT version number
		0x00, // reserved for further study
	})

	// TPKT length
	_ = binary.Write(buf, binary.BigEndian, dataLen)

	buf.Write(pduData)

	log.Printf("TPKT: Send: %s\n", hex.EncodeToString(buf.Bytes()))

	if _, err := p.conn.Write(buf.Bytes()); err != nil {
		return err
	}

	return nil
}
