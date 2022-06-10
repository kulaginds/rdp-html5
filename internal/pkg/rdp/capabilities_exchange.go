package rdp

import (
	"bytes"
	"encoding/binary"
	"io"
)

type CapabilitySetType uint16

type PDUType uint16

type CapabilitySet struct {
	CapabilitySetType CapabilitySetType
	LengthCapability  uint16
	CapabilityData    []byte
}

func (set *CapabilitySet) Serialize() []byte {
	buf := &bytes.Buffer{}

	binary.Write(buf, binary.LittleEndian, set.CapabilitySetType)
	binary.Write(buf, binary.LittleEndian, set.LengthCapability)

	buf.Write(set.CapabilityData)

	return buf.Bytes()
}

func (set *CapabilitySet) Deserialize(wire io.Reader) error {
	binary.Read(wire, binary.LittleEndian, &set.CapabilitySetType)
	binary.Read(wire, binary.LittleEndian, &set.LengthCapability)

	set.CapabilityData = make([]byte, set.LengthCapability)

	if _, err := wire.Read(set.CapabilityData); err != nil {
		return err
	}

	return nil
}

type ShareControlHeader struct {
	TotalLength uint16
	PDUType     PDUType
	PDUSource   uint16
}

func (header *ShareControlHeader) Serialize() []byte {
	buf := &bytes.Buffer{}

	binary.Write(buf, binary.LittleEndian, header.TotalLength)
	binary.Write(buf, binary.LittleEndian, header.PDUType)
	binary.Write(buf, binary.LittleEndian, header.PDUSource)

	return buf.Bytes()
}

func (header *ShareControlHeader) Deserialize(wire io.Reader) error {
	binary.Read(wire, binary.LittleEndian, &header.TotalLength)
	binary.Read(wire, binary.LittleEndian, &header.PDUType)
	binary.Read(wire, binary.LittleEndian, &header.PDUSource)

	return nil
}

type ServerDemandActivePDU struct {
	ShareControlHeader         ShareControlHeader
	ShareID                    uint32
	LengthSourceDescriptor     uint16
	LengthCombinedCapabilities uint16
	SourceDescriptor           []byte
	NumberCapabilities         uint16
	pad2Octets                 uint16
	CapabilitySets             []CapabilitySet
	SessionId                  uint32
}

func (pdu *ServerDemandActivePDU) Deserialize(wire io.Reader) error {
	err := pdu.ShareControlHeader.Deserialize(wire)
	if err != nil {
		return err
	}

	binary.Read(wire, binary.LittleEndian, &pdu.ShareID)
	binary.Read(wire, binary.LittleEndian, &pdu.LengthSourceDescriptor)
	binary.Read(wire, binary.LittleEndian, &pdu.LengthCombinedCapabilities)

	pdu.SourceDescriptor = make([]byte, pdu.LengthSourceDescriptor)

	_, err = wire.Read(pdu.SourceDescriptor)
	if err != nil {
		return err
	}

	binary.Read(wire, binary.LittleEndian, &pdu.NumberCapabilities)
	binary.Read(wire, binary.LittleEndian, &pdu.pad2Octets)

	pdu.CapabilitySets = make([]CapabilitySet, 0, pdu.NumberCapabilities)

	for i := uint16(0); i < pdu.NumberCapabilities; i++ {
		var capabilitySet CapabilitySet

		if err = capabilitySet.Deserialize(wire); err != nil {
			return err
		}

		pdu.CapabilitySets = append(pdu.CapabilitySets, capabilitySet)
	}

	binary.Read(wire, binary.LittleEndian, &pdu.SessionId)

	return nil
}

type ClientConfirmActivePDU struct {
	ShareControlHeader         ShareControlHeader
	ShareID                    uint32
	originatorID               uint16
	LengthSourceDescriptor     uint16
	LengthCombinedCapabilities uint16
	SourceDescriptor           []byte
	NumberCapabilities         uint16
	pad2Octets                 uint16
	CapabilitySets             []CapabilitySet
}

func NewClientConfirmActivePDU(shareID uint32, userId uint16) *ClientConfirmActivePDU {
	return &ClientConfirmActivePDU{
		ShareControlHeader: ShareControlHeader{
			TotalLength: 6 + 4 + 2 + 2 + 2 + 16 + 2 + 2,
			PDUType:     0x13, // TS_PROTOCOL_VERSION, PDUTYPE_CONFIRMACTIVEPDU
			PDUSource:   userId,
		},
		ShareID:                    shareID,
		originatorID:               0x03EA,
		LengthSourceDescriptor:     16,
		LengthCombinedCapabilities: 0,
		SourceDescriptor: []byte{
			'w', 'e', 'b', '-', 'r', 'd', 'p', '-', 's', 'o', 'l', 'u', 't', 'i', 'o', 'n',
		},
		NumberCapabilities: 0,
		pad2Octets:         0,
		CapabilitySets:     nil,
	}
}

func (pdu *ClientConfirmActivePDU) Serialize() []byte {
	buf := &bytes.Buffer{}

	buf.Write(pdu.ShareControlHeader.Serialize())
	binary.Write(buf, binary.LittleEndian, pdu.ShareID)
	binary.Write(buf, binary.LittleEndian, pdu.originatorID)
	binary.Write(buf, binary.LittleEndian, pdu.LengthSourceDescriptor)
	binary.Write(buf, binary.LittleEndian, pdu.LengthCombinedCapabilities)

	buf.Write(pdu.SourceDescriptor)

	binary.Write(buf, binary.LittleEndian, pdu.NumberCapabilities)
	binary.Write(buf, binary.LittleEndian, pdu.pad2Octets)

	for _, set := range pdu.CapabilitySets {
		buf.Write(set.Serialize())
	}

	return buf.Bytes()
}
