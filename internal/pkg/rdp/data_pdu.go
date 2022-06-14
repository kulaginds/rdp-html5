package rdp

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

type PDUType uint16

const (
	// PDUTypeDemandActive PDUTYPE_DEMANDACTIVEPDU
	PDUTypeDemandActive PDUType = 0x11

	// PDUTypeConfirmActive PDUTYPE_CONFIRMACTIVEPDU
	PDUTypeConfirmActive PDUType = 0x13

	// PDUTypeDeactivateAll PDUTYPE_DEACTIVATEALLPDU
	PDUTypeDeactivateAll PDUType = 0x16

	// PDUTypeData PDUTYPE_DATAPDU
	PDUTypeData PDUType = 0x17
)

func (t PDUType) IsDemandActive() bool {
	return t == PDUTypeDemandActive
}

func (t PDUType) IsConfirmActive() bool {
	return t == PDUTypeConfirmActive
}

func (t PDUType) IsDeactivateAll() bool {
	return t == PDUTypeDeactivateAll
}

func (t PDUType) IsData() bool {
	return t == PDUTypeData
}

type ShareControlHeader struct {
	TotalLength uint16
	PDUType     PDUType
	PDUSource   uint16
}

func newShareControlHeader(pduType PDUType, pduSource uint16) *ShareControlHeader {
	return &ShareControlHeader{
		PDUType:   pduType,
		PDUSource: pduSource,
	}
}

func (header *ShareControlHeader) Serialize() []byte {
	buf := &bytes.Buffer{}

	binary.Write(buf, binary.LittleEndian, header.TotalLength)
	binary.Write(buf, binary.LittleEndian, uint16(header.PDUType))
	binary.Write(buf, binary.LittleEndian, header.PDUSource)

	return buf.Bytes()
}

func (header *ShareControlHeader) Deserialize(wire io.Reader) error {
	binary.Read(wire, binary.LittleEndian, &header.TotalLength)
	binary.Read(wire, binary.LittleEndian, &header.PDUType)
	binary.Read(wire, binary.LittleEndian, &header.PDUSource)

	return nil
}

type PDUType2 uint8

const (
	// PDUType2Update PDUTYPE2_UPDATE
	PDUType2Update PDUType2 = 0x02

	// PDUType2Control PDUTYPE2_CONTROL
	PDUType2Control PDUType2 = 0x14

	// PDUType2Pointer PDUTYPE2_POINTER
	PDUType2Pointer PDUType2 = 0x1B

	// PDUType2Input PDUTYPE2_INPUT
	PDUType2Input PDUType2 = 0x1C

	// PDUType2Synchronize PDUTYPE2_SYNCHRONIZE
	PDUType2Synchronize PDUType2 = 0x1F

	// PDUType2Fontlist PDUTYPE2_FONTLIST
	PDUType2Fontlist PDUType2 = 0x27

	// PDUType2ErrorInfo PDUTYPE2_SET_ERROR_INFO_PDU
	PDUType2ErrorInfo PDUType2 = 0x2f
)

func (t PDUType2) IsUpdate() bool {
	return t == PDUType2Update
}

func (t PDUType2) IsControl() bool {
	return t == PDUType2Control
}

func (t PDUType2) IsPointer() bool {
	return t == PDUType2Pointer
}

func (t PDUType2) IsInput() bool {
	return t == PDUType2Input
}

func (t PDUType2) IsSynchronize() bool {
	return t == PDUType2Synchronize
}

func (t PDUType2) IsFontlist() bool {
	return t == PDUType2Fontlist
}

func (t PDUType2) IsErrorInfo() bool {
	return t == PDUType2ErrorInfo
}

type ShareDataHeader struct {
	ShareControlHeader ShareControlHeader
	ShareID            uint32
	StreamID           uint8
	UncompressedLength uint16
	PDUType2           PDUType2
	CompressedType     uint8
	CompressedLength   uint16
}

func newShareDataHeader(shareID uint32, pduSource uint16, pduType PDUType, pduType2 PDUType2) *ShareDataHeader {
	return &ShareDataHeader{
		ShareControlHeader: *newShareControlHeader(pduType, pduSource),
		ShareID:            shareID,
		StreamID:           0x01, // STREAM_LOW
		PDUType2:           pduType2,
	}
}

func (header *ShareDataHeader) Serialize() []byte {
	buf := &bytes.Buffer{}

	buf.Write(header.ShareControlHeader.Serialize())
	binary.Write(buf, binary.LittleEndian, header.ShareID)
	binary.Write(buf, binary.LittleEndian, uint8(0)) // padding
	binary.Write(buf, binary.LittleEndian, header.StreamID)
	binary.Write(buf, binary.LittleEndian, header.UncompressedLength)
	binary.Write(buf, binary.LittleEndian, uint8(header.PDUType2))
	binary.Write(buf, binary.LittleEndian, header.CompressedType)
	binary.Write(buf, binary.LittleEndian, header.CompressedLength)

	return buf.Bytes()
}

func (header *ShareDataHeader) Deserialize(wire io.Reader) error {
	var (
		padding uint8
		err     error
	)

	if err = header.ShareControlHeader.Deserialize(wire); err != nil {
		return err
	}

	err = binary.Read(wire, binary.LittleEndian, &header.ShareID)
	if err != nil {
		return err
	}

	err = binary.Read(wire, binary.LittleEndian, &padding)
	if err != nil {
		return err
	}

	err = binary.Read(wire, binary.LittleEndian, &header.StreamID)
	if err != nil {
		return err
	}

	err = binary.Read(wire, binary.LittleEndian, &header.UncompressedLength)
	if err != nil {
		return err
	}

	err = binary.Read(wire, binary.LittleEndian, &header.PDUType2)
	if err != nil {
		return err
	}

	err = binary.Read(wire, binary.LittleEndian, &header.CompressedType)
	if err != nil {
		return err
	}

	err = binary.Read(wire, binary.LittleEndian, &header.CompressedLength)
	if err != nil {
		return err
	}

	return nil
}

type DataPDU struct {
	ShareDataHeader    ShareDataHeader
	SynchronizePDUData *SynchronizePDUData
	ControlPDUData     *ControlPDUData
	FontListPDUData    *FontListPDUData
	FontMapPDUData     *FontMapPDUData
	ErrorInfoPDUData   *ErrorInfoPDUData
}

func (pdu *DataPDU) Serialize() []byte {
	var data []byte

	switch {
	case pdu.ShareDataHeader.PDUType2.IsSynchronize():
		data = pdu.SynchronizePDUData.Serialize()
	case pdu.ShareDataHeader.PDUType2.IsControl():
		data = pdu.ControlPDUData.Serialize()
	case pdu.ShareDataHeader.PDUType2.IsFontlist():
		data = pdu.FontListPDUData.Serialize()
	}

	pdu.ShareDataHeader.ShareControlHeader.TotalLength = uint16(18 + len(data))
	pdu.ShareDataHeader.UncompressedLength = uint16(4 + len(data))

	buf := &bytes.Buffer{}

	buf.Write(pdu.ShareDataHeader.Serialize())
	buf.Write(data)

	return buf.Bytes()
}

func (pdu *DataPDU) Deserialize(wire io.Reader) error {
	var err error

	if err = pdu.ShareDataHeader.Deserialize(wire); err != nil {
		return err
	}

	switch {
	case pdu.ShareDataHeader.PDUType2.IsSynchronize():
		pdu.SynchronizePDUData = &SynchronizePDUData{}

		return pdu.SynchronizePDUData.Deserialize(wire)
	case pdu.ShareDataHeader.PDUType2.IsControl():
		pdu.ControlPDUData = &ControlPDUData{}

		return pdu.ControlPDUData.Deserialize(wire)
	case pdu.ShareDataHeader.PDUType2.IsErrorInfo():
		pdu.ErrorInfoPDUData = &ErrorInfoPDUData{}

		return pdu.ErrorInfoPDUData.Deserialize(wire)
	}

	return fmt.Errorf("unknown data pdu: %d", pdu.ShareDataHeader.PDUType2)
}
