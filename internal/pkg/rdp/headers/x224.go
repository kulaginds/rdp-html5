package headers

import (
	"bytes"
	"fmt"
	"io"
)

const (
	x224LI           = 1
	x224FixedPartLen = 6 // without length indicator (LI)
)

func WrapX224ConnectionRequestPDU(data []byte) []byte {
	pduLen := uint8(x224FixedPartLen + len(data))

	buf := bytes.NewBuffer(make([]byte, 0, pduLen))

	buf.Write([]byte{
		pduLen,     // LI
		0xE0,       // CR + CDT: TPDU_CONNECTION_REQUEST
		0x00, 0x00, // DST-REF
		0x00, 0x00, // SRC-REF
		0x00, // CLASS OPTION
	})

	// Variable part
	buf.Write(data)

	return buf.Bytes()
}

func WrapX224DataPDU(data []byte) []byte {
	buf := bytes.NewBuffer(make([]byte, 0, 3+len(data)))

	buf.Write([]byte{
		0x02, // packet size
		0xF0, // message type TPDU_DATA
		0x80, // EOT flag is up, which indicates that current TPDU is the last data unit of a complete TPDU sequence
	})
	buf.Write(data)

	return buf.Bytes()
}

func UnwrapX224ConnectionConfirmPDU(wire io.Reader, dataLen int) ([]byte, error) {
	x224Data := make([]byte, dataLen)

	_, err := wire.Read(x224Data)
	if err != nil {
		return nil, fmt.Errorf("read X224 data: %w", err)
	}

	x224Data = x224Data[x224LI+x224FixedPartLen:] // offset to rdpNegData

	return x224Data, nil
}

func UnwrapX224DataPDU(wire io.Reader, dataLen int) ([]byte, error) {
	const fixedPartLen = 2

	x224Data := make([]byte, dataLen)

	_, err := wire.Read(x224Data)
	if err != nil {
		return nil, fmt.Errorf("read X224 data: %w", err)
	}

	x224Data = x224Data[x224LI+fixedPartLen:]

	return x224Data, nil
}
