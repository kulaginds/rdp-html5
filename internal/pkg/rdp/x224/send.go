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

	binary.Write(buf, binary.BigEndian, pdu.LI)
	binary.Write(buf, binary.BigEndian, pdu.DTROA)
	binary.Write(buf, binary.BigEndian, pdu.NREOT)

	buf.Write(pdu.UserData)

	return buf.Bytes()
}

func (pdu *Data) Deserialize(wire io.Reader) error {
	err := binary.Read(wire, binary.BigEndian, &pdu.LI)
	if err != nil {
		return err
	}

	if pdu.LI != dataFixedPartLen {
		return ErrWrongDataLength
	}

	err = binary.Read(wire, binary.BigEndian, &pdu.DTROA)
	if err != nil {
		return err
	}

	err = binary.Read(wire, binary.BigEndian, &pdu.NREOT)
	if err != nil {
		return err
	}

	return nil
}

const dataFixedPartLen uint8 = 0x02

func (p *Protocol) Send(userData []byte) error {

	req := Data{
		LI:       dataFixedPartLen,
		DTROA:    0xF0, // message type TPDU_DATA
		NREOT:    0x80, // EOT flag is up, which indicates that current TPDU is the last data unit of a complete TPDU sequence
		UserData: userData,
	}

	return p.tpktConn.Send(req.Serialize())
}
