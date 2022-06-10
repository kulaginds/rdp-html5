package headers

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
)

func WrapMCSDomainPDU(tag uint8, data []byte) []byte {
	buf := &bytes.Buffer{}

	buf.WriteByte(tag << 2) // per-encoded DomainMCSPDU choice
	buf.Write(data)

	return buf.Bytes()
}

func UnwrapMCSDomainPDU(tag uint8, wire io.Reader) error {
	var choiceID uint8

	if err := binary.Read(wire, binary.BigEndian, &choiceID); err != nil {
		return err
	}

	choiceID >>= 2

	if choiceID != tag {
		return errors.New("unexpected MCS choice")
	}

	return nil
}

func WrapMCSSendData(data []byte) []byte {
	return WrapMCSDomainPDU(25, data)
}
