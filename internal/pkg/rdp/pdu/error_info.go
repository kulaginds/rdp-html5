package pdu

import (
	"encoding/binary"
	"io"
)

type ErrorInfoPDUData struct {
	ErrorInfo uint32
}

func (pdu *ErrorInfoPDUData) Deserialize(wire io.Reader) error {
	return binary.Read(wire, binary.LittleEndian, &pdu.ErrorInfo)
}
