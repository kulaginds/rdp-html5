package fastpath

import (
	"bytes"
	"encoding/binary"
	"io"
)

type EventCode uint8

const (
	//EventCodeScanCode FASTPATH_INPUT_EVENT_SCANCODE
	EventCodeScanCode EventCode = 0

	//EventCodeMouse FASTPATH_INPUT_EVENT_MOUSE
	EventCodeMouse EventCode = 1

	//EventCodeMouseX FASTPATH_INPUT_EVENT_MOUSEX
	EventCodeMouseX EventCode = 2

	//EventCodeSync FASTPATH_INPUT_EVENT_SYNC
	EventCodeSync EventCode = 3

	//EventCodeUnicode FASTPATH_INPUT_EVENT_UNICODE
	EventCodeUnicode EventCode = 4

	//EventCodeQoETimestamp FASTPATH_INPUT_EVENT_QOE_TIMESTAMP
	EventCodeQoETimestamp EventCode = 6
)

type InputEvent struct {
	EventFlags           uint8
	EventCode            EventCode
	keyboardEvent        *keyboardEvent
	unicodeKeyboardEvent *unicodeKeyboardEvent
	mouseEvent           *mouseEvent
	extendedMouseEvent   *extendedMouseEvent
	qualityOfExperience  *qualityOfExperience
}

func (e *InputEvent) Serialize() []byte {
	buf := &bytes.Buffer{}

	// event flags in higher 5 bits
	// event code in lower 3 bits
	header := (e.EventFlags&0x1f)<<3 | uint8(e.EventCode)&0x7

	var data []byte

	switch e.EventCode {
	case EventCodeScanCode:
		data = e.keyboardEvent.Serialize()
	case EventCodeUnicode:
		data = e.unicodeKeyboardEvent.Serialize()
	case EventCodeMouse:
		data = e.mouseEvent.Serialize()
	case EventCodeMouseX:
		data = e.extendedMouseEvent.Serialize()
	case EventCodeSync: // do nothing
	case EventCodeQoETimestamp:
		data = e.qualityOfExperience.Serialize()
	}

	binary.Write(buf, binary.LittleEndian, header)
	buf.Write(data)

	return buf.Bytes()
}

type InputEventPDU struct {
	FPInputHeader uint8
	eventData     []byte
}

func NewInputEventPDU(eventData []byte) *InputEventPDU {
	return &InputEventPDU{
		FPInputHeader: 1 << 2, // action=FASTPATH, numEvents=1, no flags
		eventData:     eventData,
	}
}

func (pdu *InputEventPDU) Serialize() []byte {
	buf := &bytes.Buffer{}

	length := 1 + len(pdu.eventData) // without length bytes

	binary.Write(buf, binary.LittleEndian, pdu.FPInputHeader)
	pdu.SerializeLength(length, buf)
	buf.Write(pdu.eventData)

	return buf.Bytes()
}

func (pdu *InputEventPDU) SerializeLength(value int, w io.Writer) error {
	if value > 0x7f {
		value += 2 // 2 bytes length

		return binary.Write(w, binary.BigEndian, value|0x8000)
	}

	value += 1 // 1 byte length
	if _, err := w.Write([]byte{uint8(value)}); err != nil {
		return err
	}

	return nil
}

func (i *impl) Send(pdu *InputEventPDU) error {
	data := pdu.Serialize()

	_, err := i.conn.Write(data)
	if err != nil {
		return err
	}

	return nil
}
