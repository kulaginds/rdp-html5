package rdp

import (
	"encoding/binary"
	"errors"
	"io"

	"github.com/kulaginds/web-rdp-solution/internal/pkg/rdp/headers"
)

type LicensingBinaryBlob struct {
	BlobType uint16
	BlobLen  uint16
	BlobData []byte
}

func (b *LicensingBinaryBlob) Deserialize(wire io.Reader) error {
	binary.Read(wire, binary.LittleEndian, &b.BlobType)
	binary.Read(wire, binary.LittleEndian, &b.BlobLen)

	if b.BlobLen == 0 {
		return nil
	}

	b.BlobData = make([]byte, b.BlobLen)

	if _, err := wire.Read(b.BlobData); err != nil {
		return err
	}

	return nil
}

type LicensingErrorMessage struct {
	ErrorCode       uint32
	StateTransition uint32
	ErrorInfo       LicensingBinaryBlob
}

func (m *LicensingErrorMessage) Deserialize(wire io.Reader) error {
	binary.Read(wire, binary.LittleEndian, &m.ErrorCode)
	binary.Read(wire, binary.LittleEndian, &m.StateTransition)

	return m.ErrorInfo.Deserialize(wire)
}

type LicensingPreamble struct {
	MsgType uint8
	Flags   uint8
	MsgSize uint16
}

func (p *LicensingPreamble) Deserialize(wire io.Reader) error {
	binary.Read(wire, binary.LittleEndian, &p.MsgType)
	binary.Read(wire, binary.LittleEndian, &p.Flags)
	binary.Read(wire, binary.LittleEndian, &p.MsgSize)

	return nil
}

type ServerLicenseErrorPDU struct {
	Preamble           LicensingPreamble
	ValidClientMessage LicensingErrorMessage
}

func (pdu *ServerLicenseErrorPDU) Deserialize(wire io.Reader) error {
	securityFlag, err := headers.UnwrapSecurityFlag(wire)
	if err != nil {
		return err
	}

	if securityFlag != 0x0080 { // SEC_LICENSE_PKT
		return errors.New("bad license header")
	}

	err = pdu.Preamble.Deserialize(wire)
	if err != nil {
		return err
	}

	err = pdu.ValidClientMessage.Deserialize(wire)
	if err != nil {
		return err
	}

	return nil
}
