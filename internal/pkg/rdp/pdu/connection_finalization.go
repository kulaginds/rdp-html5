package pdu

import (
	"bytes"
	"encoding/binary"
	"io"
)

type MessageType uint16

const (
	MessageTypeSync MessageType = 1
)

type SynchronizePDUData struct {
	MessageType MessageType
}

func NewSynchronize(shareID uint32, userId uint16) *Data {
	return &Data{
		ShareDataHeader: *newShareDataHeader(shareID, userId, TypeData, Type2Synchronize),
		SynchronizePDUData: &SynchronizePDUData{
			MessageType: MessageTypeSync,
		},
	}
}

const ServerChannelID uint16 = 1002

func (pdu *SynchronizePDUData) Serialize() []byte {
	buf := new(bytes.Buffer)

	_ = binary.Write(buf, binary.LittleEndian, uint16(pdu.MessageType))
	_ = binary.Write(buf, binary.LittleEndian, ServerChannelID) // targetUser

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

func NewControl(shareID uint32, userId uint16, action ControlAction) *Data {
	return &Data{
		ShareDataHeader: *newShareDataHeader(shareID, userId, TypeData, Type2Control),
		ControlPDUData: &ControlPDUData{
			Action: action,
		},
	}
}

func (pdu *ControlPDUData) Serialize() []byte {
	buf := new(bytes.Buffer)

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

func NewFontList(shareID uint32, userId uint16) *Data {
	return &Data{
		ShareDataHeader: *newShareDataHeader(shareID, userId, TypeData, Type2Fontlist),
		FontListPDUData: &FontListPDUData{},
	}
}

func (pdu *FontListPDUData) Serialize() []byte {
	buf := new(bytes.Buffer)

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
