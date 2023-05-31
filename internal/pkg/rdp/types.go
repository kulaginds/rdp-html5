package rdp

import (
	"bufio"
)

type ProtocolCode uint8

func (a ProtocolCode) IsFastpath() bool {
	return a&0x3 == 0
}

func (a ProtocolCode) IsX224() bool {
	return a == 3
}

func receiveProtocol(bufReader *bufio.Reader) (ProtocolCode, error) {
	action, err := bufReader.ReadByte()
	if err != nil {
		return 0, err
	}

	err = bufReader.UnreadByte()
	if err != nil {
		return 0, err
	}

	return ProtocolCode(action), nil
}
