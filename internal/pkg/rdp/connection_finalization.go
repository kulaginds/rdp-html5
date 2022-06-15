package rdp

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"

	"github.com/kulaginds/web-rdp-solution/internal/pkg/rdp/mcs"
)

type MessageType uint16

const (
	MessageTypeSync MessageType = 1
)

type SynchronizePDUData struct {
	MessageType MessageType
}

func NewSynchronizePDU(shareID uint32, userId uint16) *DataPDU {
	return &DataPDU{
		ShareDataHeader: *newShareDataHeader(shareID, userId, PDUTypeData, PDUType2Synchronize),
		SynchronizePDUData: &SynchronizePDUData{
			MessageType: MessageTypeSync,
		},
	}
}

func (pdu *SynchronizePDUData) Serialize() []byte {
	buf := &bytes.Buffer{}

	_ = binary.Write(buf, binary.LittleEndian, uint16(pdu.MessageType))
	_ = binary.Write(buf, binary.LittleEndian, uint16(mcs.ServerChannelID)) // targetUser

	return buf.Bytes()
}

func (pdu *SynchronizePDUData) Deserialize(wire io.Reader) error {
	var err error

	err = binary.Read(wire, binary.LittleEndian, &pdu.MessageType)
	if err != nil {
		return err
	}

	var targetUser uint16
	err = binary.Read(wire, binary.LittleEndian, &targetUser)
	if err != nil {
		return err
	}

	return nil
}

type ControlAction uint16

const (
	// ControlActionRequestControl CTRLACTION_REQUEST_CONTROL
	ControlActionRequestControl ControlAction = 0x0001

	// ControlActionGrantedControl CTRLACTION_GRANTED_CONTROL
	ControlActionGrantedControl ControlAction = 0x0002

	// ControlActionDetach CTRLACTION_DETACH
	ControlActionDetach ControlAction = 0x0003

	// ControlActionCooperate CTRLACTION_COOPERATE
	ControlActionCooperate ControlAction = 0x0004
)

type ControlPDUData struct {
	Action    ControlAction
	GrantID   uint16
	ControlID uint32
}

func NewControlPDU(shareID uint32, userId uint16, action ControlAction) *DataPDU {
	return &DataPDU{
		ShareDataHeader: *newShareDataHeader(shareID, userId, PDUTypeData, PDUType2Control),
		ControlPDUData: &ControlPDUData{
			Action: action,
		},
	}
}

func (pdu *ControlPDUData) Serialize() []byte {
	buf := &bytes.Buffer{}

	_ = binary.Write(buf, binary.LittleEndian, uint16(pdu.Action))
	_ = binary.Write(buf, binary.LittleEndian, pdu.GrantID)
	_ = binary.Write(buf, binary.LittleEndian, pdu.ControlID)

	return buf.Bytes()
}

func (pdu *ControlPDUData) Deserialize(wire io.Reader) error {
	var err error

	err = binary.Read(wire, binary.LittleEndian, &pdu.Action)
	if err != nil {
		return err
	}

	err = binary.Read(wire, binary.LittleEndian, &pdu.GrantID)
	if err != nil {
		return err
	}

	err = binary.Read(wire, binary.LittleEndian, &pdu.ControlID)
	if err != nil {
		return err
	}

	return nil
}

type FontListPDUData struct{}

func NewFontListPDU(shareID uint32, userId uint16) *DataPDU {
	return &DataPDU{
		ShareDataHeader: *newShareDataHeader(shareID, userId, PDUTypeData, PDUType2Fontlist),
	}
}

func (pdu *FontListPDUData) Serialize() []byte {
	buf := &bytes.Buffer{}

	_ = binary.Write(buf, binary.LittleEndian, uint16(0x0000)) // numberFonts
	_ = binary.Write(buf, binary.LittleEndian, uint16(0x0000)) // totalNumFonts
	_ = binary.Write(buf, binary.LittleEndian, uint16(0x0003)) // listFlags = FONTLIST_FIRST | FONTLIST_LAST
	_ = binary.Write(buf, binary.LittleEndian, uint16(0x0032)) // entrySize

	return buf.Bytes()
}

type FontMapPDUData struct{}

func (pdu *FontMapPDUData) Deserialize(wire io.Reader) error {
	var (
		numberEntries   uint16
		totalNumEntries uint16
		mapFlags        uint16
		entrySize       uint16
		err             error
	)

	err = binary.Read(wire, binary.LittleEndian, &numberEntries)
	if err != nil {
		return err
	}

	err = binary.Read(wire, binary.LittleEndian, &totalNumEntries)
	if err != nil {
		return err
	}

	err = binary.Read(wire, binary.LittleEndian, &mapFlags)
	if err != nil {
		return err
	}

	err = binary.Read(wire, binary.LittleEndian, &entrySize)
	if err != nil {
		return err
	}

	return nil
}

func (c *client) ConnectionFinalization() error {
	var err error

	synchronize := NewSynchronizePDU(c.shareID, c.mcsLayer.UserId())
	if err = c.mcsLayer.Send("global", synchronize.Serialize()); err != nil {
		return err
	}

	controlCooperate := NewControlPDU(c.shareID, c.mcsLayer.UserId(), ControlActionCooperate)
	if err = c.mcsLayer.Send("global", controlCooperate.Serialize()); err != nil {
		return err
	}

	controlRequestControl := NewControlPDU(c.shareID, c.mcsLayer.UserId(), ControlActionRequestControl)
	if err = c.mcsLayer.Send("global", controlRequestControl.Serialize()); err != nil {
		return err
	}

	fontList := NewFontListPDU(c.shareID, c.mcsLayer.UserId())

	err = c.mcsLayer.Send("global", fontList.Serialize())
	if err != nil {
		return err
	}

	var (
		serverSynchronizeReceived bool
		controlCooperateReceived  bool
		grantedControlReceived    bool
		fontMapReceived           bool

		dataPDU *DataPDU
		wire    io.Reader
	)

	for {
		if serverSynchronizeReceived && controlCooperateReceived && grantedControlReceived && fontMapReceived {
			return nil
		}

		_, wire, err = c.mcsLayer.Receive()
		if err != nil {
			return err
		}

		dataPDU = &DataPDU{}
		if err = dataPDU.Deserialize(wire); err != nil {
			return err
		}

		pduType2 := dataPDU.ShareDataHeader.PDUType2

		switch {
		case pduType2.IsSynchronize():
			serverSynchronizeReceived = true
		case pduType2.IsControl():
			if dataPDU.ControlPDUData.Action == ControlActionCooperate {
				controlCooperateReceived = true
			}

			if dataPDU.ControlPDUData.Action == ControlActionGrantedControl {
				grantedControlReceived = true
			}
		case pduType2.IsFontmap():
			fontMapReceived = true
		case pduType2.IsErrorInfo():
			return fmt.Errorf("server error info: %d", dataPDU.ErrorInfoPDUData.ErrorInfo)
		default:
			return fmt.Errorf("unknown server message with pduType2 = %d", pduType2)
		}
	}
}
