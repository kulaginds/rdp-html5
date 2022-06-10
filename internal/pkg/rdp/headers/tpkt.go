package headers

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

const (
	tpktHeaderLen = 4
)

func WrapTPKT(data []byte) []byte {
	buf := bytes.NewBuffer(make([]byte, 0, tpktHeaderLen+len(data)))
	dataLen := uint16(tpktHeaderLen + len(data))

	buf.Write([]byte{
		3, // TPKT version number
		0, // reserved for further study
	})

	// TPKT length
	_ = binary.Write(buf, binary.BigEndian, dataLen)

	buf.Write(data)

	return buf.Bytes()
}

// UnwrapTPKT unwrap TPKT response header and return data length.
func UnwrapTPKT(wire io.Reader) (int, error) {
	tpktHeaderData := make([]byte, tpktHeaderLen)

	_, err := wire.Read(tpktHeaderData)
	if err != nil {
		return 0, fmt.Errorf("read TPKT header: %w", err)
	}

	dataLen := binary.BigEndian.Uint16(tpktHeaderData[2:4])
	dataLen -= tpktHeaderLen

	return int(dataLen), nil
}
