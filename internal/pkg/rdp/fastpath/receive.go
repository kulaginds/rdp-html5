package fastpath

import (
	"encoding/binary"
	"errors"
	"io"
)

type UpdateCode uint8

const (
	UpdateCodeOrders       UpdateCode = 0x0
	UpdateCodeBitmap       UpdateCode = 0x1
	UpdateCodePalette      UpdateCode = 0x2
	UpdateCodeSynchronize  UpdateCode = 0x3
	UpdateCodeSurfCMDs     UpdateCode = 0x4
	UpdateCodePTRNull      UpdateCode = 0x5
	UpdateCodePTRDefault   UpdateCode = 0x6
	UpdateCodePTRPosition  UpdateCode = 0x8
	UpdateCodeColor        UpdateCode = 0x9
	UpdateCodeCached       UpdateCode = 0xa
	UpdateCodePointer      UpdateCode = 0xb
	UpdateCodeLargePointer UpdateCode = 0xc
)

type Fragment uint8

const (
	FragmentSingle Fragment = 0x0
	FragmentLast   Fragment = 0x1
	FragmentFirst  Fragment = 0x2
	FragmentNext   Fragment = 0x3
)

type Compression uint8

const (
	CompressionUsed Compression = 0x2
)

type Update struct {
	UpdateCode       UpdateCode
	fragmentation    Fragment
	compression      Compression
	compressionFlags uint8
	size             uint16

	paletteUpdateData         *paletteUpdateData
	bitmapUpdateData          *bitmapUpdateData
	pointerPositionUpdateData *pointerPositionUpdateData
	colorPointerUpdateData    *colorPointerUpdateData
}

func (u *Update) Deserialize(wire io.Reader) error {
	var err error

	var updateHeader uint8
	err = binary.Read(wire, binary.LittleEndian, &updateHeader)
	if err != nil {
		return err
	}

	u.UpdateCode = UpdateCode(updateHeader & 0xf)
	u.fragmentation = Fragment((updateHeader >> 4) & 0x3)
	u.compression = Compression((updateHeader >> 6) & 0x3)

	if u.compression&CompressionUsed == CompressionUsed {
		err = binary.Read(wire, binary.LittleEndian, &u.compressionFlags)
		if err != nil {
			return err
		}
	}

	err = binary.Read(wire, binary.LittleEndian, &u.size)
	if err != nil {
		return err
	}

	//switch u.UpdateCode {
	//case UpdateCodePalette:
	//	u.paletteUpdateData = &paletteUpdateData{}
	//
	//	return u.paletteUpdateData.Deserialize(wire)
	//case UpdateCodeBitmap:
	//	u.bitmapUpdateData = &bitmapUpdateData{}
	//
	//	return u.bitmapUpdateData.Deserialize(wire)
	//case UpdateCodeSynchronize: // do nothing
	//case UpdateCodePTRPosition:
	//	u.pointerPositionUpdateData = &pointerPositionUpdateData{}
	//
	//	return u.pointerPositionUpdateData.Deserialize(wire)
	//case UpdateCodePTRNull: // do nothing
	//case UpdateCodePTRDefault: // do nothing
	//case UpdateCodeCached:
	//	u.colorPointerUpdateData = &colorPointerUpdateData{}
	//
	//	return u.colorPointerUpdateData.Deserialize(wire)
	//}

	d := make([]byte, u.size)
	_, err = wire.Read(d)
	if err != nil {
		return err
	}

	return nil
}

type UpdatePDUAction uint8

const (
	UpdatePDUActionFastPath UpdatePDUAction = 0x0
	UpdatePDUActionX224     UpdatePDUAction = 0x3
)

type UpdatePDUFlag uint8

const (
	UpdatePDUFlagSecureChecksum UpdatePDUFlag = 0x1
	UpdatePDUFlagEncrypted      UpdatePDUFlag = 0x2
)

type UpdatePDU struct {
	fpOutputHeader uint8

	Action UpdatePDUAction
	Flags  UpdatePDUFlag
	Data   []byte
}

func NewUpdatePDU(fpOutputHeader uint8) *UpdatePDU {
	return &UpdatePDU{
		fpOutputHeader: fpOutputHeader,
	}
}

var ErrUnexpectedX224 = errors.New("unexpected x224")

func (pdu *UpdatePDU) Deserialize(wire io.Reader) error {
	var err error

	pdu.Action = UpdatePDUAction(pdu.fpOutputHeader & 0x3)
	pdu.Flags = UpdatePDUFlag((pdu.fpOutputHeader >> 6) & 0x3)

	if pdu.Action == UpdatePDUActionX224 {
		return ErrUnexpectedX224
	}

	if pdu.Flags&UpdatePDUFlagSecureChecksum == UpdatePDUFlagSecureChecksum {
		return errors.New("checksum not supported")
	}

	if pdu.Flags&UpdatePDUFlagEncrypted == UpdatePDUFlagEncrypted {
		return errors.New("encryption not supported")
	}

	var (
		length           uint16
		length1, length2 uint8
	)

	err = binary.Read(wire, binary.LittleEndian, &length1)
	if err != nil {
		return err
	}

	length = uint16(length1)

	if length1&0x80 == 0x80 {
		err = binary.Read(wire, binary.LittleEndian, &length2)
		if err != nil {
			return err
		}

		length1 -= 0x80

		length = binary.BigEndian.Uint16([]byte{length1, length2})
	}

	if length > 0x4000 {
		return errors.New("too big packet")
	}

	pdu.Data = make([]byte, length)
	_, err = wire.Read(pdu.Data)
	if err != nil {
		return err
	}

	return nil
}

func (i *impl) Receive(fpOutputHeader uint8) (*UpdatePDU, error) {
	pdu := NewUpdatePDU(fpOutputHeader)
	if err := pdu.Deserialize(i.conn); err != nil {
		return nil, err
	}

	return pdu, nil
}
