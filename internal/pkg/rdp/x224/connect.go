package x224

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"
)

type ConnectionRequest struct {
	li           uint8
	CRCDT        uint8
	DSTREF       uint16
	SRCREF       uint16
	ClassOption  uint8
	VariablePart []byte // unsupported
	UserData     []byte
}

func (pdu *ConnectionRequest) Serialize() []byte {
	const x224FixedPartLen = 6 // without length indicator (LI)

	pdu.li = uint8(x224FixedPartLen + len(pdu.UserData))

	buf := new(bytes.Buffer)

	_ = binary.Write(buf, binary.LittleEndian, pdu.li)
	_ = binary.Write(buf, binary.LittleEndian, pdu.CRCDT)
	_ = binary.Write(buf, binary.LittleEndian, pdu.DSTREF)
	_ = binary.Write(buf, binary.LittleEndian, pdu.SRCREF)
	_ = binary.Write(buf, binary.LittleEndian, pdu.ClassOption)

	buf.Write(pdu.UserData)

	return buf.Bytes()
}

type ConnectionConfirm struct {
	LI          uint8
	CCCDT       uint8
	DSTREF      uint16
	SRCREF      uint16
	ClassOption uint8
}

func (pdu *ConnectionConfirm) Deserialize(wire io.Reader) error {
	const (
		fixedPartLen    uint8 = 0x06
		variablePartLen uint8 = 0x08
		packetLen             = fixedPartLen + variablePartLen
	)

	err := binary.Read(wire, binary.LittleEndian, &pdu.LI)
	if err != nil {
		return err
	}

	if pdu.LI != packetLen {
		return ErrSmallConnectionConfirmLength
	}

	err = binary.Read(wire, binary.LittleEndian, &pdu.CCCDT)
	if err != nil {
		return err
	}

	if pdu.CCCDT&0xf0 != 0xd0 { // connection confirm code
		return ErrWrongConnectionConfirmCode
	}

	err = binary.Read(wire, binary.LittleEndian, &pdu.DSTREF)
	if err != nil {
		return err
	}

	err = binary.Read(wire, binary.LittleEndian, &pdu.SRCREF)
	if err != nil {
		return err
	}

	err = binary.Read(wire, binary.LittleEndian, &pdu.ClassOption)
	if err != nil {
		return err
	}

	return nil
}

func (p *Protocol) Connect(userData []byte) (io.Reader, error) {
	var (
		wire io.Reader
		err  error
	)

	req := ConnectionRequest{
		CRCDT:        0xE0, // TPDU_CONNECTION_REQUEST
		DSTREF:       0,
		SRCREF:       0,
		ClassOption:  0,
		VariablePart: nil,
		UserData:     userData,
	}

	log.Println("X224: Client Connection Request")

	if err = p.tpktConn.Send(req.Serialize()); err != nil {
		return nil, fmt.Errorf("client connection request: %w", err)
	}

	log.Println("X224: Server Connection Confirm")

	wire, err = p.tpktConn.Receive()
	if err != nil {
		return nil, fmt.Errorf("recieve connection response: %w", err)
	}

	var resp ConnectionConfirm
	if err = resp.Deserialize(wire); err != nil {
		return nil, fmt.Errorf("server connection confirm: %w", err)
	}

	return wire, nil
}
