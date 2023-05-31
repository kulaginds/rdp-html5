package fastpath

import (
	"bytes"
	"encoding/binary"
	"io"
)

type InputEventPDU struct {
	action    uint8
	numEvents uint8
	flags     uint8
	eventData []byte
}

func NewInputEventPDU(eventData []byte) *InputEventPDU {
	return &InputEventPDU{
		numEvents: 1,
		eventData: eventData,
	}
}

func (pdu *InputEventPDU) Serialize() []byte {
	buf := new(bytes.Buffer)

	fpInputHeader := pdu.action&0x3 | ((pdu.numEvents & 0xf) << 2) | ((pdu.flags & 0x3) << 6)
	length := 1 + len(pdu.eventData) // without length bytes

	binary.Write(buf, binary.LittleEndian, fpInputHeader)
	pdu.SerializeLength(length, buf)
	buf.Write(pdu.eventData)

	return buf.Bytes()
}

func (pdu *InputEventPDU) SerializeLength(value int, w io.Writer) error {
	if value > 0x7f {
		value += 2 // 2 bytes length

		return binary.Write(w, binary.BigEndian, value|0x8000)
	}

	value += 1 // 1 byte length
	if _, err := w.Write([]byte{uint8(value)}); err != nil {
		return err
	}

	return nil
}

func (i *Protocol) Send(pdu *InputEventPDU) error {
	data := pdu.Serialize()

	_, err := i.conn.Write(data)
	if err != nil {
		return err
	}

	return nil
}
