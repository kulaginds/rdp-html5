package pdu

import (
	"bytes"
	"encoding/binary"
	"io"
)

type CapabilitySetType uint16

const (
	// CapabilitySetTypeGeneral CAPSTYPE_GENERAL
	CapabilitySetTypeGeneral CapabilitySetType = 0x0001

	// CapabilitySetTypeBitmap CAPSTYPE_BITMAP
	CapabilitySetTypeBitmap CapabilitySetType = 0x0002

	// CapabilitySetTypeOrder CAPSTYPE_ORDER
	CapabilitySetTypeOrder CapabilitySetType = 0x0003

	// CapabilitySetTypeBitmapCache CAPSTYPE_BITMAPCACHE
	CapabilitySetTypeBitmapCache CapabilitySetType = 0x0004

	// CapabilitySetTypeControl CAPSTYPE_CONTROL
	CapabilitySetTypeControl CapabilitySetType = 0x0005

	// CapabilitySetTypeActivation CAPSTYPE_ACTIVATION
	CapabilitySetTypeActivation CapabilitySetType = 0x0007

	// CapabilitySetTypePointer CAPSTYPE_POINTER
	CapabilitySetTypePointer CapabilitySetType = 0x0008

	// CapabilitySetTypeShare CAPSTYPE_SHARE
	CapabilitySetTypeShare CapabilitySetType = 0x0009

	// CapabilitySetTypeColorCache CAPSTYPE_COLORCACHE
	CapabilitySetTypeColorCache CapabilitySetType = 0x000A

	// CapabilitySetTypeSound CAPSTYPE_SOUND
	CapabilitySetTypeSound CapabilitySetType = 0x000C

	// CapabilitySetTypeInput CAPSTYPE_INPUT
	CapabilitySetTypeInput CapabilitySetType = 0x000D

	// CapabilitySetTypeFont CAPSTYPE_FONT
	CapabilitySetTypeFont CapabilitySetType = 0x000E

	// CapabilitySetTypeBrush CAPSTYPE_BRUSH
	CapabilitySetTypeBrush CapabilitySetType = 0x000F

	// CapabilitySetTypeGlyphCache CAPSTYPE_GLYPHCACHE
	CapabilitySetTypeGlyphCache CapabilitySetType = 0x0010

	// CapabilitySetTypeOffscreenBitmapCache CAPSTYPE_OFFSCREENCACHE
	CapabilitySetTypeOffscreenBitmapCache CapabilitySetType = 0x0011

	// CapabilitySetTypeBitmapCacheHostSupport CAPSTYPE_BITMAPCACHE_HOSTSUPPORT
	CapabilitySetTypeBitmapCacheHostSupport CapabilitySetType = 0x0012

	// CapabilitySetTypeBitmapCacheRev2 CAPSTYPE_BITMAPCACHE_REV2
	CapabilitySetTypeBitmapCacheRev2 CapabilitySetType = 0x0013

	// CapabilitySetTypeVirtualChannel CAPSTYPE_VIRTUALCHANNEL
	CapabilitySetTypeVirtualChannel CapabilitySetType = 0x0014

	// CapabilitySetTypeDrawNineGridCache CAPSTYPE_DRAWNINEGRIDCACHE
	CapabilitySetTypeDrawNineGridCache CapabilitySetType = 0x0015

	// CapabilitySetTypeDrawGDIPlus CAPSTYPE_DRAWGDIPLUS
	CapabilitySetTypeDrawGDIPlus CapabilitySetType = 0x0016

	// CapabilitySetTypeRail CAPSTYPE_RAIL
	CapabilitySetTypeRail CapabilitySetType = 0x0017

	// CapabilitySetTypeWindow CAPSTYPE_WINDOW
	CapabilitySetTypeWindow CapabilitySetType = 0x0018

	// CapabilitySetTypeCompDesk CAPSETTYPE_COMPDESK
	CapabilitySetTypeCompDesk CapabilitySetType = 0x0019

	// CapabilitySetTypeMultifragmentUpdate CAPSETTYPE_MULTIFRAGMENTUPDATE
	CapabilitySetTypeMultifragmentUpdate CapabilitySetType = 0x001A

	// CapabilitySetTypeLargePointer CAPSETTYPE_LARGE_POINTER
	CapabilitySetTypeLargePointer CapabilitySetType = 0x001B

	// CapabilitySetTypeSurfaceCommands CAPSETTYPE_SURFACE_COMMANDS
	CapabilitySetTypeSurfaceCommands CapabilitySetType = 0x001C

	// CapabilitySetTypeBitmapCodecs CAPSETTYPE_BITMAP_CODECS
	CapabilitySetTypeBitmapCodecs CapabilitySetType = 0x001D

	// CapabilitySetTypeFrameAcknowledge CAPSSETTYPE_FRAME_ACKNOWLEDGE
	CapabilitySetTypeFrameAcknowledge CapabilitySetType = 0x001E
)

type CapabilitySet struct {
	CapabilitySetType                   CapabilitySetType
	GeneralCapabilitySet                *GeneralCapabilitySet
	BitmapCapabilitySet                 *BitmapCapabilitySet
	OrderCapabilitySet                  *OrderCapabilitySet
	BitmapCacheCapabilitySetRev1        *BitmapCacheCapabilitySetRev1
	BitmapCacheCapabilitySetRev2        *BitmapCacheCapabilitySetRev2
	ColorCacheCapabilitySet             *ColorCacheCapabilitySet
	PointerCapabilitySet                *PointerCapabilitySet
	InputCapabilitySet                  *InputCapabilitySet
	BrushCapabilitySet                  *BrushCapabilitySet
	GlyphCacheCapabilitySet             *GlyphCacheCapabilitySet
	OffscreenBitmapCacheCapabilitySet   *OffscreenBitmapCacheCapabilitySet
	VirtualChannelCapabilitySet         *VirtualChannelCapabilitySet
	DrawNineGridCacheCapabilitySet      *DrawNineGridCacheCapabilitySet
	DrawGDIPlusCapabilitySet            *DrawGDIPlusCapabilitySet
	SoundCapabilitySet                  *SoundCapabilitySet
	BitmapCacheHostSupportCapabilitySet *BitmapCacheHostSupportCapabilitySet
	ControlCapabilitySet                *ControlCapabilitySet
	WindowActivationCapabilitySet       *WindowActivationCapabilitySet
	ShareCapabilitySet                  *ShareCapabilitySet
	FontCapabilitySet                   *FontCapabilitySet
	MultifragmentUpdateCapabilitySet    *MultifragmentUpdateCapabilitySet
	LargePointerCapabilitySet           *LargePointerCapabilitySet
	DesktopCompositionCapabilitySet     *DesktopCompositionCapabilitySet
	SurfaceCommandsCapabilitySet        *SurfaceCommandsCapabilitySet
	BitmapCodecsCapabilitySet           *BitmapCodecsCapabilitySet
	RailCapabilitySet                   *RailCapabilitySet
	WindowListCapabilitySet             *WindowListCapabilitySet
}

func (set *CapabilitySet) Serialize() []byte {
	var data []byte

	switch set.CapabilitySetType {
	case CapabilitySetTypeGeneral:
		data = set.GeneralCapabilitySet.Serialize()
	case CapabilitySetTypeBitmap:
		data = set.BitmapCapabilitySet.Serialize()
	case CapabilitySetTypeOrder:
		data = set.OrderCapabilitySet.Serialize()
	case CapabilitySetTypeBitmapCache:
		data = set.BitmapCacheCapabilitySetRev1.Serialize()
	case CapabilitySetTypeBitmapCacheRev2:
		data = set.BitmapCacheCapabilitySetRev2.Serialize()
	case CapabilitySetTypeColorCache:
		data = set.ColorCacheCapabilitySet.Serialize()
	case CapabilitySetTypeActivation:
		data = set.WindowActivationCapabilitySet.Serialize()
	case CapabilitySetTypeControl:
		data = set.ControlCapabilitySet.Serialize()
	case CapabilitySetTypePointer:
		data = set.PointerCapabilitySet.Serialize()
	case CapabilitySetTypeInput:
		data = set.InputCapabilitySet.Serialize()
	case CapabilitySetTypeBrush:
		data = set.BrushCapabilitySet.Serialize()
	case CapabilitySetTypeGlyphCache:
		data = set.GlyphCacheCapabilitySet.Serialize()
	case CapabilitySetTypeOffscreenBitmapCache:
		data = set.OffscreenBitmapCacheCapabilitySet.Serialize()
	case CapabilitySetTypeVirtualChannel:
		data = set.VirtualChannelCapabilitySet.Serialize()
	case CapabilitySetTypeSound:
		data = set.SoundCapabilitySet.Serialize()
	case CapabilitySetTypeShare:
		data = set.ShareCapabilitySet.Serialize()
	case CapabilitySetTypeFont:
		data = set.FontCapabilitySet.Serialize()
	case CapabilitySetTypeDrawNineGridCache:
		data = set.DrawNineGridCacheCapabilitySet.Serialize()
	case CapabilitySetTypeDrawGDIPlus:
		data = set.DrawGDIPlusCapabilitySet.Serialize()
	case CapabilitySetTypeMultifragmentUpdate:
		data = set.MultifragmentUpdateCapabilitySet.Serialize()
	case CapabilitySetTypeRail:
		data = set.RailCapabilitySet.Serialize()
	case CapabilitySetTypeWindow:
		data = set.WindowListCapabilitySet.Serialize()
	}

	buf := new(bytes.Buffer)

	lengthCapability := uint16(4 + len(data))

	_ = binary.Write(buf, binary.LittleEndian, set.CapabilitySetType)
	_ = binary.Write(buf, binary.LittleEndian, lengthCapability)
	buf.Write(data)

	return buf.Bytes()
}

func (set *CapabilitySet) DeserializeQuick(wire io.Reader) error {
	var (
		lengthCapability uint16
		err              error
	)

	err = binary.Read(wire, binary.LittleEndian, &set.CapabilitySetType)
	if err != nil {
		return err
	}

	err = binary.Read(wire, binary.LittleEndian, &lengthCapability)
	if err != nil {
		return err
	}

	data := make([]byte, lengthCapability-4)
	if _, err = wire.Read(data); err != nil {
		return err
	}

	return nil
}

func (set *CapabilitySet) Deserialize(wire io.Reader) error {
	var (
		lengthCapability uint16
		err              error
	)

	err = binary.Read(wire, binary.LittleEndian, &set.CapabilitySetType)
	if err != nil {
		return err
	}

	err = binary.Read(wire, binary.LittleEndian, &lengthCapability)
	if err != nil {
		return err
	}

	switch set.CapabilitySetType {
	case CapabilitySetTypeGeneral:
		set.GeneralCapabilitySet = &GeneralCapabilitySet{}

		return set.GeneralCapabilitySet.Deserialize(wire)
	case CapabilitySetTypeBitmap:
		set.BitmapCapabilitySet = &BitmapCapabilitySet{}

		return set.BitmapCapabilitySet.Deserialize(wire)
	case CapabilitySetTypeOrder:
		set.OrderCapabilitySet = &OrderCapabilitySet{}

		return set.OrderCapabilitySet.Deserialize(wire)
	case CapabilitySetTypeBitmapCache:
		set.BitmapCacheCapabilitySetRev1 = &BitmapCacheCapabilitySetRev1{}

		return set.BitmapCacheCapabilitySetRev1.Deserialize(wire)
	case CapabilitySetTypeBitmapCacheRev2:
		set.BitmapCacheCapabilitySetRev2 = &BitmapCacheCapabilitySetRev2{}

		return set.BitmapCacheCapabilitySetRev2.Deserialize(wire)
	case CapabilitySetTypeColorCache:
		set.ColorCacheCapabilitySet = &ColorCacheCapabilitySet{}

		return set.ColorCacheCapabilitySet.Deserialize(wire)
	case CapabilitySetTypePointer:
		set.PointerCapabilitySet = &PointerCapabilitySet{
			lengthCapability: lengthCapability - 4,
		}

		return set.PointerCapabilitySet.Deserialize(wire)
	case CapabilitySetTypeInput:
		set.InputCapabilitySet = &InputCapabilitySet{}

		return set.InputCapabilitySet.Deserialize(wire)
	case CapabilitySetTypeBrush:
		set.BrushCapabilitySet = &BrushCapabilitySet{}

		return set.BrushCapabilitySet.Deserialize(wire)
	case CapabilitySetTypeGlyphCache:
		set.GlyphCacheCapabilitySet = &GlyphCacheCapabilitySet{}

		return set.GlyphCacheCapabilitySet.Deserialize(wire)
	case CapabilitySetTypeOffscreenBitmapCache:
		set.OffscreenBitmapCacheCapabilitySet = &OffscreenBitmapCacheCapabilitySet{}

		return set.OffscreenBitmapCacheCapabilitySet.Deserialize(wire)
	case CapabilitySetTypeVirtualChannel:
		set.VirtualChannelCapabilitySet = &VirtualChannelCapabilitySet{}

		return set.VirtualChannelCapabilitySet.Deserialize(wire)
	case CapabilitySetTypeDrawNineGridCache:
		set.DrawNineGridCacheCapabilitySet = &DrawNineGridCacheCapabilitySet{}

		return set.DrawNineGridCacheCapabilitySet.Deserialize(wire)
	case CapabilitySetTypeDrawGDIPlus:
		set.DrawGDIPlusCapabilitySet = &DrawGDIPlusCapabilitySet{}

		return set.DrawGDIPlusCapabilitySet.Deserialize(wire)
	case CapabilitySetTypeSound:
		set.SoundCapabilitySet = &SoundCapabilitySet{}

		return set.SoundCapabilitySet.Deserialize(wire)
	}

	data := make([]byte, lengthCapability-4)
	if _, err = wire.Read(data); err != nil {
		return err
	}

	return nil
}

type ServerDemandActive struct {
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

func (pdu *ServerDemandActive) Deserialize(wire io.Reader) error {
	err := pdu.ShareControlHeader.Deserialize(wire)
	if err != nil {
		return err
	}

	err = binary.Read(wire, binary.LittleEndian, &pdu.ShareID)
	if err != nil {
		return err
	}

	err = binary.Read(wire, binary.LittleEndian, &pdu.LengthSourceDescriptor)
	if err != nil {
		return err
	}

	err = binary.Read(wire, binary.LittleEndian, &pdu.LengthCombinedCapabilities)
	if err != nil {
		return err
	}

	pdu.SourceDescriptor = make([]byte, pdu.LengthSourceDescriptor)

	_, err = wire.Read(pdu.SourceDescriptor)
	if err != nil {
		return err
	}

	err = binary.Read(wire, binary.LittleEndian, &pdu.NumberCapabilities)
	if err != nil {
		return err
	}

	err = binary.Read(wire, binary.LittleEndian, &pdu.pad2Octets)
	if err != nil {
		return err
	}

	pdu.CapabilitySets = make([]CapabilitySet, 0, pdu.NumberCapabilities)

	for i := uint16(0); i < pdu.NumberCapabilities; i++ {
		var capabilitySet CapabilitySet

		if err = capabilitySet.DeserializeQuick(wire); err != nil {
			return err
		}

		pdu.CapabilitySets = append(pdu.CapabilitySets, capabilitySet)
	}

	err = binary.Read(wire, binary.LittleEndian, &pdu.SessionId)
	if err != nil {
		return err
	}

	return nil
}

type ClientConfirmActive struct {
	ShareControlHeader ShareControlHeader
	ShareID            uint32
	SourceDescriptor   []byte
	CapabilitySets     []CapabilitySet
}

func NewClientConfirmActive(shareID uint32, userId, desktopWidth, desktopHeight uint16, withRemoteApp bool) *ClientConfirmActive {
	pdu := ClientConfirmActive{
		ShareControlHeader: ShareControlHeader{
			PDUType:   TypeConfirmActive,
			PDUSource: userId,
		},
		ShareID: shareID,
		SourceDescriptor: []byte{
			'w', 'e', 'b', '-', 'r', 'd', 'p', '-', 's', 'o', 'l', 'u', 't', 'i', 'o', 'n',
		},
		CapabilitySets: []CapabilitySet{
			NewGeneralCapabilitySet(),
			NewBitmapCapabilitySet(desktopWidth, desktopHeight),
			NewOrderCapabilitySet(),
			NewBitmapCacheCapabilitySetRev1(),
			NewPointerCapabilitySet(),
			NewInputCapabilitySet(),
			NewBrushCapabilitySet(),
			NewGlyphCacheCapabilitySet(),
			NewOffscreenBitmapCacheCapabilitySet(),
			NewVirtualChannelCapabilitySet(),
			NewSoundCapabilitySet(),
			NewMultifragmentUpdateCapabilitySet(),
		},
	}

	if withRemoteApp {
		pdu.CapabilitySets = append(pdu.CapabilitySets, NewRailCapabilitySet(), NewWindowListCapabilitySet())
	}

	return &pdu
}

func (pdu *ClientConfirmActive) Serialize() []byte {
	capBuf := bytes.Buffer{}

	for _, set := range pdu.CapabilitySets {
		capBuf.Write(set.Serialize())
	}

	lengthSourceDescriptor := uint16(len(pdu.SourceDescriptor))
	lengthCombinedCapabilities := uint16(4 + capBuf.Len())

	pdu.ShareControlHeader.PDUType = TypeConfirmActive
	pdu.ShareControlHeader.TotalLength = 6 + 4 + 2 + 2 + 2 + lengthSourceDescriptor + lengthCombinedCapabilities

	buf := new(bytes.Buffer)

	buf.Write(pdu.ShareControlHeader.Serialize())
	_ = binary.Write(buf, binary.LittleEndian, pdu.ShareID)
	_ = binary.Write(buf, binary.LittleEndian, uint16(0x03EA)) // originatorID
	_ = binary.Write(buf, binary.LittleEndian, lengthSourceDescriptor)
	_ = binary.Write(buf, binary.LittleEndian, lengthCombinedCapabilities)

	buf.Write(pdu.SourceDescriptor)

	_ = binary.Write(buf, binary.LittleEndian, uint16(len(pdu.CapabilitySets)))
	_ = binary.Write(buf, binary.LittleEndian, uint16(0)) // padding

	buf.Write(capBuf.Bytes())

	return buf.Bytes()
}

func (pdu *ClientConfirmActive) Deserialize(wire io.Reader) error {
	var err error

	err = pdu.ShareControlHeader.Deserialize(wire)
	if err != nil {
		return err
	}

	err = binary.Read(wire, binary.LittleEndian, &pdu.ShareID)
	if err != nil {
		return err
	}

	var originatorID uint16
	err = binary.Read(wire, binary.LittleEndian, &originatorID)
	if err != nil {
		return err
	}

	var lengthSourceDescriptor uint16
	err = binary.Read(wire, binary.LittleEndian, &lengthSourceDescriptor)
	if err != nil {
		return err
	}

	var lengthCombinedCapabilities uint16
	err = binary.Read(wire, binary.LittleEndian, &lengthCombinedCapabilities)
	if err != nil {
		return err
	}

	pdu.SourceDescriptor = make([]byte, lengthSourceDescriptor)
	_, err = wire.Read(pdu.SourceDescriptor)
	if err != nil {
		return err
	}

	var numberCapabilities uint16
	err = binary.Read(wire, binary.LittleEndian, &numberCapabilities)
	if err != nil {
		return err
	}

	var padding uint16
	err = binary.Read(wire, binary.LittleEndian, &padding)
	if err != nil {
		return err
	}

	pdu.CapabilitySets = make([]CapabilitySet, numberCapabilities)

	for i := range pdu.CapabilitySets {
		err = pdu.CapabilitySets[i].Deserialize(wire)
		if err != nil {
			return err
		}
	}

	return nil
}
