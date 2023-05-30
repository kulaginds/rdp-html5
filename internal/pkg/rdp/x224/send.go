package x224

import (
	"bytes"
	"encoding/binary"
	"io"
)

type Data struct {
	LI       uint8
	DTROA    uint8
	NREOT    uint8
	UserData []byte
}

func (pdu *Data) Serialize() []byte {
	buf := new(bytes.Buffer)

	binary.Write(buf, binary.LittleEndian, pdu.LI)
	binary.Write(buf, binary.LittleEndian, pdu.DTROA)
	binary.Write(buf, binary.LittleEndian, pdu.NREOT)

	buf.Write(pdu.UserData)

	return buf.Bytes()
}

func (pdu *Data) Deserialize(wire io.Reader) error {
	err := binary.Read(wire, binary.LittleEndian, &pdu.LI)
	if err != nil {
		return err
	}

	if pdu.LI != fixedPartLen {
		return ErrWrongDataLength
	}

	err = binary.Read(wire, binary.LittleEndian, &pdu.DTROA)
	if err != nil {
		return err
	}

	err = binary.Read(wire, binary.LittleEndian, &pdu.NREOT)
	if err != nil {
		return err
	}

	return nil
}

const fixedPartLen uint8 = 0x02

func (p *Protocol) Send(userData []byte) error {
	req := Data{
		LI:       fixedPartLen,
		DTROA:    0xF0, // message type TPDU_DATA
		NREOT:    0x80, // EOT flag is up, which indicates that current TPDU is the last data unit of a complete TPDU sequence
		UserData: userData,
	}

	return p.tpktConn.Send(req.Serialize())
}
