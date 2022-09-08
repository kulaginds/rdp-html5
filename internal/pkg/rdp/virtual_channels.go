package rdp

import (
	"bytes"
	"encoding/binary"
	"io"
)

type ChannelFlag uint32

const (
	// ChannelFlagFirst CHANNEL_FLAG_FIRST
	ChannelFlagFirst ChannelFlag = 0x00000001

	// ChannelFlagLast CHANNEL_FLAG_LAST
	ChannelFlagLast ChannelFlag = 0x00000002

	// ChannelFlagShowProtocol CHANNEL_FLAG_SHOW_PROTOCOL
	ChannelFlagShowProtocol ChannelFlag = 0x00000010

	// ChannelFlagSuspend CHANNEL_FLAG_SUSPEND
	ChannelFlagSuspend ChannelFlag = 0x00000020

	// ChannelFlagResume CHANNEL_FLAG_RESUME
	ChannelFlagResume ChannelFlag = 0x00000040

	// ChannelFlagShadowPersistent CHANNEL_FLAG_SHADOW_PERSISTENT
	ChannelFlagShadowPersistent ChannelFlag = 0x00000080

	// ChannelFlagCompressed CHANNEL_PACKET_COMPRESSED
	ChannelFlagCompressed ChannelFlag = 0x00200000

	// ChannelFlagAtFront CHANNEL_PACKET_AT_FRONT
	ChannelFlagAtFront ChannelFlag = 0x00400000

	// ChannelFlagFlushed CHANNEL_PACKET_FLUSHED
	ChannelFlagFlushed ChannelFlag = 0x00800000
)

// ChannelPDUHeader CHANNEL_PDU_HEADER
type ChannelPDUHeader struct {
	Flags ChannelFlag
}

func (h *ChannelPDUHeader) Serialize() []byte {
	buf := &bytes.Buffer{}

	binary.Write(buf, binary.LittleEndian, uint32(8))
	binary.Write(buf, binary.LittleEndian, uint32(h.Flags))

	return buf.Bytes()
}

func (h *ChannelPDUHeader) Deserialize(wire io.Reader) error {
	var (
		err    error
		length uint32
	)

	err = binary.Read(wire, binary.LittleEndian, &length)
	if err != nil {
		return err
	}

	err = binary.Read(wire, binary.LittleEndian, &h.Flags)
	if err != nil {
		return err
	}

	return nil
}
