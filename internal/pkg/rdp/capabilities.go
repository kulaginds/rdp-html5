package rdp

import (
	"bytes"
	"encoding/binary"
	"io"
)

type GeneralCapabilitySet struct {
	OSMajorType           uint16
	OSMinorType           uint16
	ExtraFlags            uint16
	RefreshRectSupport    uint8
	SuppressOutputSupport uint8
}

func NewGeneralCapabilitySet() *CapabilitySet {
	return &CapabilitySet{
		CapabilitySetType: CapabilitySetTypeGeneral,
		GeneralCapabilitySet: &GeneralCapabilitySet{
			OSMajorType: 0x0008,                   // Chrome OS platform
			ExtraFlags:  0x0001 | 0x0004 | 0x0400, // required: FASTPATH_OUTPUT_SUPPORTED, LONG_CREDENTIALS_SUPPORTED, NO_BITMAP_COMPRESSION_HDR
		},
	}
}

func (s *GeneralCapabilitySet) Serialize() []byte {
	buf := new(bytes.Buffer)

	_ = binary.Write(buf, binary.LittleEndian, s.OSMajorType)
	_ = binary.Write(buf, binary.LittleEndian, s.OSMinorType)
	_ = binary.Write(buf, binary.LittleEndian, uint16(0x0200)) // protocolVersion
	_ = binary.Write(buf, binary.LittleEndian, uint16(0x0000)) // padding
	_ = binary.Write(buf, binary.LittleEndian, uint16(0x0000)) // compressionTypes
	_ = binary.Write(buf, binary.LittleEndian, s.ExtraFlags)
	_ = binary.Write(buf, binary.LittleEndian, uint16(0x0000)) // updateCapabilityFlag
	_ = binary.Write(buf, binary.LittleEndian, uint16(0x0000)) // remoteUnshareFlag
	_ = binary.Write(buf, binary.LittleEndian, uint16(0x0000)) // compressionLevel
	_ = binary.Write(buf, binary.LittleEndian, s.RefreshRectSupport)
	_ = binary.Write(buf, binary.LittleEndian, s.SuppressOutputSupport)

	return buf.Bytes()
}

func (s *GeneralCapabilitySet) Deserialize(wire io.Reader) error {
	var err error

	err = binary.Read(wire, binary.LittleEndian, &s.OSMajorType)
	if err != nil {
		return err
	}

	err = binary.Read(wire, binary.LittleEndian, &s.OSMinorType)
	if err != nil {
		return err
	}

	var protocolVersion uint16
	err = binary.Read(wire, binary.LittleEndian, &protocolVersion)
	if err != nil {
		return err
	}

	var padding uint16
	err = binary.Read(wire, binary.LittleEndian, &padding)
	if err != nil {
		return err
	}

	var compressionTypes uint16
	err = binary.Read(wire, binary.LittleEndian, &compressionTypes)
	if err != nil {
		return err
	}

	err = binary.Read(wire, binary.LittleEndian, &s.ExtraFlags)
	if err != nil {
		return err
	}

	var updateCapabilityFlag uint16
	err = binary.Read(wire, binary.LittleEndian, &updateCapabilityFlag)
	if err != nil {
		return err
	}

	var remoteUnshareFlag uint16
	err = binary.Read(wire, binary.LittleEndian, &remoteUnshareFlag)
	if err != nil {
		return err
	}

	var compressionLevel uint16
	err = binary.Read(wire, binary.LittleEndian, &compressionLevel)
	if err != nil {
		return err
	}

	err = binary.Read(wire, binary.LittleEndian, &s.RefreshRectSupport)
	if err != nil {
		return err
	}

	err = binary.Read(wire, binary.LittleEndian, &s.SuppressOutputSupport)
	if err != nil {
		return err
	}

	return nil
}

type BitmapCapabilitySet struct {
	PreferredBitsPerPixel uint16
	Receive1BitPerPixel   uint16
	Receive4BitsPerPixel  uint16
	Receive8BitsPerPixel  uint16
	DesktopWidth          uint16
	DesktopHeight         uint16
	DesktopResizeFlag     uint16
	DrawingFlags          uint8
}

func NewBitmapCapabilitySet(desktopWidth, desktopHeight uint16) *CapabilitySet {
	return &CapabilitySet{
		CapabilitySetType: CapabilitySetTypeBitmap,
		BitmapCapabilitySet: &BitmapCapabilitySet{
			PreferredBitsPerPixel: 0x0018, // HIGH_COLOR_24BPP
			Receive1BitPerPixel:   0x0001,
			Receive4BitsPerPixel:  0x0001,
			Receive8BitsPerPixel:  0x0001,
			DesktopWidth:          desktopWidth,
			DesktopHeight:         desktopHeight,
		},
	}
}

func (s *BitmapCapabilitySet) Serialize() []byte {
	buf := new(bytes.Buffer)

	_ = binary.Write(buf, binary.LittleEndian, s.PreferredBitsPerPixel)
	_ = binary.Write(buf, binary.LittleEndian, s.Receive1BitPerPixel)
	_ = binary.Write(buf, binary.LittleEndian, s.Receive4BitsPerPixel)
	_ = binary.Write(buf, binary.LittleEndian, s.Receive8BitsPerPixel)
	_ = binary.Write(buf, binary.LittleEndian, s.DesktopWidth)
	_ = binary.Write(buf, binary.LittleEndian, s.DesktopHeight)
	_ = binary.Write(buf, binary.LittleEndian, uint16(0)) // padding
	_ = binary.Write(buf, binary.LittleEndian, s.DesktopResizeFlag)
	_ = binary.Write(buf, binary.LittleEndian, uint16(0x0001)) // bitmapCompressionFlag
	_ = binary.Write(buf, binary.LittleEndian, uint8(0))       // highColorFlags
	_ = binary.Write(buf, binary.LittleEndian, s.DrawingFlags)
	_ = binary.Write(buf, binary.LittleEndian, uint16(0x0001)) // multipleRectangleSupport
	_ = binary.Write(buf, binary.LittleEndian, uint16(0))      // padding

	return buf.Bytes()
}

func (s *BitmapCapabilitySet) Deserialize(wire io.Reader) error {
	var err error

	err = binary.Read(wire, binary.LittleEndian, &s.PreferredBitsPerPixel)
	if err != nil {
		return err
	}

	err = binary.Read(wire, binary.LittleEndian, &s.Receive1BitPerPixel)
	if err != nil {
		return err
	}

	err = binary.Read(wire, binary.LittleEndian, &s.Receive4BitsPerPixel)
	if err != nil {
		return err
	}

	err = binary.Read(wire, binary.LittleEndian, &s.Receive8BitsPerPixel)
	if err != nil {
		return err
	}

	err = binary.Read(wire, binary.LittleEndian, &s.DesktopWidth)
	if err != nil {
		return err
	}

	err = binary.Read(wire, binary.LittleEndian, &s.DesktopHeight)
	if err != nil {
		return err
	}

	var padding uint16
	err = binary.Read(wire, binary.LittleEndian, &padding)
	if err != nil {
		return err
	}

	err = binary.Read(wire, binary.LittleEndian, &s.DesktopResizeFlag)
	if err != nil {
		return err
	}

	var bitmapCompressionFlag uint16
	err = binary.Read(wire, binary.LittleEndian, &bitmapCompressionFlag)
	if err != nil {
		return err
	}

	var highColorFlags uint8
	err = binary.Read(wire, binary.LittleEndian, &highColorFlags)
	if err != nil {
		return err
	}

	err = binary.Read(wire, binary.LittleEndian, &s.DrawingFlags)
	if err != nil {
		return err
	}

	var multipleRectangleSupport uint16
	err = binary.Read(wire, binary.LittleEndian, &multipleRectangleSupport)
	if err != nil {
		return err
	}

	err = binary.Read(wire, binary.LittleEndian, &padding)
	if err != nil {
		return err
	}

	return nil
}

type OrderCapabilitySet struct {
	OrderFlags          uint16
	OrderSupport        [32]byte
	textFlags           uint16
	OrderSupportExFlags uint16
	DesktopSaveSize     uint32
	textANSICodePage    uint16
}

func NewOrderCapabilitySet() *CapabilitySet {
	return &CapabilitySet{
		CapabilitySetType: CapabilitySetTypeOrder,
		OrderCapabilitySet: &OrderCapabilitySet{
			OrderFlags:      0x2 | 0x0008, // NEGOTIATEORDERSUPPORT, ZEROBOUNDSDELTASSUPPORT this flags must be set
			DesktopSaveSize: 480 * 480,
		},
	}
}

func (s *OrderCapabilitySet) Serialize() []byte {
	buf := new(bytes.Buffer)

	buf.Write(make([]byte, 16))                            // terminalDescriptor
	_ = binary.Write(buf, binary.LittleEndian, uint32(0))  // padding
	_ = binary.Write(buf, binary.LittleEndian, uint16(1))  // desktopSaveXGranularity
	_ = binary.Write(buf, binary.LittleEndian, uint16(20)) // desktopSaveYGranularity
	_ = binary.Write(buf, binary.LittleEndian, uint16(0))  // padding
	_ = binary.Write(buf, binary.LittleEndian, uint16(1))  // maximumOrderLevel = ORD_LEVEL_1_ORDERS
	_ = binary.Write(buf, binary.LittleEndian, uint16(0))  // numberFonts
	_ = binary.Write(buf, binary.LittleEndian, s.OrderFlags)
	_ = binary.Write(buf, binary.LittleEndian, s.OrderSupport)
	_ = binary.Write(buf, binary.LittleEndian, s.textFlags) // textFlags
	_ = binary.Write(buf, binary.LittleEndian, s.OrderSupportExFlags)
	_ = binary.Write(buf, binary.LittleEndian, uint32(0)) // padding
	_ = binary.Write(buf, binary.LittleEndian, s.DesktopSaveSize)
	_ = binary.Write(buf, binary.LittleEndian, uint32(0))          // padding
	_ = binary.Write(buf, binary.LittleEndian, s.textANSICodePage) // textANSICodePage
	_ = binary.Write(buf, binary.LittleEndian, uint16(0))          // padding

	return buf.Bytes()
}

func (s *OrderCapabilitySet) Deserialize(wire io.Reader) error {
	var (
		err                     error
		terminalDescriptor      [16]byte
		padding                 uint32
		desktopSaveXGranularity uint16
		desktopSaveYGranularity uint16
		padding2                uint16
		maximumOrderLevel       uint16
		numberFonts             uint16
		textFlags               uint16
		textANSICodePage        uint16
	)

	err = binary.Read(wire, binary.LittleEndian, &terminalDescriptor)
	if err != nil {
		return err
	}

	err = binary.Read(wire, binary.LittleEndian, &padding)
	if err != nil {
		return err
	}

	err = binary.Read(wire, binary.LittleEndian, &desktopSaveXGranularity)
	if err != nil {
		return err
	}

	err = binary.Read(wire, binary.LittleEndian, &desktopSaveYGranularity)
	if err != nil {
		return err
	}

	err = binary.Read(wire, binary.LittleEndian, &padding2)
	if err != nil {
		return err
	}

	err = binary.Read(wire, binary.LittleEndian, &maximumOrderLevel)
	if err != nil {
		return err
	}

	err = binary.Read(wire, binary.LittleEndian, &numberFonts)
	if err != nil {
		return err
	}

	err = binary.Read(wire, binary.LittleEndian, &s.OrderFlags)
	if err != nil {
		return err
	}

	err = binary.Read(wire, binary.LittleEndian, &s.OrderSupport)
	if err != nil {
		return err
	}

	err = binary.Read(wire, binary.LittleEndian, &textFlags)
	if err != nil {
		return err
	}

	err = binary.Read(wire, binary.LittleEndian, &s.OrderSupportExFlags)
	if err != nil {
		return err
	}

	err = binary.Read(wire, binary.LittleEndian, &padding)
	if err != nil {
		return err
	}

	err = binary.Read(wire, binary.LittleEndian, &s.DesktopSaveSize)
	if err != nil {
		return err
	}

	err = binary.Read(wire, binary.LittleEndian, &padding)
	if err != nil {
		return err
	}

	err = binary.Read(wire, binary.LittleEndian, &textANSICodePage)
	if err != nil {
		return err
	}

	err = binary.Read(wire, binary.LittleEndian, &padding2)
	if err != nil {
		return err
	}

	return nil
}

type BitmapCacheCapabilitySetRev1 struct {
	Cache0Entries         uint16
	Cache0MaximumCellSize uint16
	Cache1Entries         uint16
	Cache1MaximumCellSize uint16
	Cache2Entries         uint16
	Cache2MaximumCellSize uint16
}

func NewBitmapCacheCapabilitySetRev1() *CapabilitySet {
	return &CapabilitySet{
		CapabilitySetType:            CapabilitySetTypeBitmapCache,
		BitmapCacheCapabilitySetRev1: &BitmapCacheCapabilitySetRev1{},
	}
}

func (s *BitmapCacheCapabilitySetRev1) Serialize() []byte {
	buf := new(bytes.Buffer)

	buf.Write(make([]byte, 24)) // padding
	_ = binary.Write(buf, binary.LittleEndian, &s.Cache0Entries)
	_ = binary.Write(buf, binary.LittleEndian, &s.Cache0MaximumCellSize)
	_ = binary.Write(buf, binary.LittleEndian, &s.Cache1Entries)
	_ = binary.Write(buf, binary.LittleEndian, &s.Cache1MaximumCellSize)
	_ = binary.Write(buf, binary.LittleEndian, &s.Cache2Entries)
	_ = binary.Write(buf, binary.LittleEndian, &s.Cache2MaximumCellSize)

	return buf.Bytes()
}

func (s *BitmapCacheCapabilitySetRev1) Deserialize(wire io.Reader) error {
	var (
		padding [24]byte
		err     error
	)

	err = binary.Read(wire, binary.LittleEndian, &padding)
	if err != nil {
		return err
	}

	err = binary.Read(wire, binary.LittleEndian, &s.Cache0Entries)
	if err != nil {
		return err
	}

	err = binary.Read(wire, binary.LittleEndian, &s.Cache0MaximumCellSize)
	if err != nil {
		return err
	}

	err = binary.Read(wire, binary.LittleEndian, &s.Cache1Entries)
	if err != nil {
		return err
	}

	err = binary.Read(wire, binary.LittleEndian, &s.Cache1MaximumCellSize)
	if err != nil {
		return err
	}

	err = binary.Read(wire, binary.LittleEndian, &s.Cache2Entries)
	if err != nil {
		return err
	}

	err = binary.Read(wire, binary.LittleEndian, &s.Cache2MaximumCellSize)
	if err != nil {
		return err
	}

	return nil
}

type BitmapCacheCapabilitySetRev2 struct {
	CacheFlags           uint16
	NumCellCaches        uint8
	BitmapCache0CellInfo uint32
	BitmapCache1CellInfo uint32
	BitmapCache2CellInfo uint32
	BitmapCache3CellInfo uint32
	BitmapCache4CellInfo uint32
}

func NewBitmapCacheCapabilitySetRev2() *CapabilitySet {
	return &CapabilitySet{
		CapabilitySetType:            CapabilitySetTypeBitmapCacheRev2,
		BitmapCacheCapabilitySetRev2: &BitmapCacheCapabilitySetRev2{},
	}
}

func (s *BitmapCacheCapabilitySetRev2) Serialize() []byte {
	buf := new(bytes.Buffer)

	_ = binary.Write(buf, binary.LittleEndian, &s.CacheFlags)
	_ = binary.Write(buf, binary.LittleEndian, uint8(0)) // padding
	_ = binary.Write(buf, binary.LittleEndian, &s.NumCellCaches)
	_ = binary.Write(buf, binary.LittleEndian, &s.BitmapCache0CellInfo)
	_ = binary.Write(buf, binary.LittleEndian, &s.BitmapCache1CellInfo)
	_ = binary.Write(buf, binary.LittleEndian, &s.BitmapCache2CellInfo)
	_ = binary.Write(buf, binary.LittleEndian, &s.BitmapCache3CellInfo)
	_ = binary.Write(buf, binary.LittleEndian, &s.BitmapCache4CellInfo)
	buf.Write(make([]byte, 12)) // padding

	return buf.Bytes()
}

func (s *BitmapCacheCapabilitySetRev2) Deserialize(wire io.Reader) error {
	var err error

	err = binary.Read(wire, binary.LittleEndian, &s.CacheFlags)
	if err != nil {
		return err
	}

	var padding uint8
	err = binary.Read(wire, binary.LittleEndian, &padding)
	if err != nil {
		return err
	}

	err = binary.Read(wire, binary.LittleEndian, &s.NumCellCaches)
	if err != nil {
		return err
	}

	err = binary.Read(wire, binary.LittleEndian, &s.BitmapCache0CellInfo)
	if err != nil {
		return err
	}

	err = binary.Read(wire, binary.LittleEndian, &s.BitmapCache1CellInfo)
	if err != nil {
		return err
	}

	err = binary.Read(wire, binary.LittleEndian, &s.BitmapCache2CellInfo)
	if err != nil {
		return err
	}

	err = binary.Read(wire, binary.LittleEndian, &s.BitmapCache3CellInfo)
	if err != nil {
		return err
	}

	err = binary.Read(wire, binary.LittleEndian, &s.BitmapCache4CellInfo)
	if err != nil {
		return err
	}

	err = binary.Read(wire, binary.LittleEndian, &s.BitmapCache4CellInfo)
	if err != nil {
		return err
	}

	var padding2 [12]byte
	err = binary.Read(wire, binary.LittleEndian, &padding2)
	if err != nil {
		return err
	}

	return nil
}

type ColorCacheCapabilitySet struct {
	ColorTableCacheSize uint16
}

func (s *ColorCacheCapabilitySet) Serialize() []byte {
	buf := new(bytes.Buffer)

	binary.Write(buf, binary.LittleEndian, &s.ColorTableCacheSize)
	binary.Write(buf, binary.LittleEndian, uint16(0)) // padding

	return buf.Bytes()
}

func (s *ColorCacheCapabilitySet) Deserialize(wire io.Reader) error {
	var (
		padding uint16
		err     error
	)

	err = binary.Read(wire, binary.LittleEndian, &s.ColorTableCacheSize)
	if err != nil {
		return err
	}

	return binary.Read(wire, binary.LittleEndian, &padding)
}

type PointerCapabilitySet struct {
	ColorPointerFlag      uint16
	ColorPointerCacheSize uint16
	PointerCacheSize      uint16
	lengthCapability      uint16
}

func NewPointerCapabilitySet() *CapabilitySet {
	return &CapabilitySet{
		CapabilitySetType: CapabilitySetTypePointer,
		PointerCapabilitySet: &PointerCapabilitySet{
			ColorPointerFlag: 1, // color mouse cursors are supported
			PointerCacheSize: 25,
		},
	}
}

func (s *PointerCapabilitySet) Serialize() []byte {
	buf := new(bytes.Buffer)

	_ = binary.Write(buf, binary.LittleEndian, s.ColorPointerFlag)
	_ = binary.Write(buf, binary.LittleEndian, s.ColorPointerCacheSize)
	_ = binary.Write(buf, binary.LittleEndian, s.PointerCacheSize)

	return buf.Bytes()
}

func (s *PointerCapabilitySet) Deserialize(wire io.Reader) error {
	var err error

	err = binary.Read(wire, binary.LittleEndian, &s.ColorPointerFlag)
	if err != nil {
		return err
	}

	err = binary.Read(wire, binary.LittleEndian, &s.ColorPointerCacheSize)
	if err != nil {
		return err
	}

	if s.lengthCapability == 4 {
		return nil
	}

	err = binary.Read(wire, binary.LittleEndian, &s.PointerCacheSize)
	if err != nil {
		return err
	}

	return nil
}

type InputCapabilitySet struct {
	InputFlags          uint16
	KeyboardLayout      uint32
	KeyboardType        uint32
	KeyboardSubType     uint32
	KeyboardFunctionKey uint32
	ImeFileName         [64]byte
}

func NewInputCapabilitySet() *CapabilitySet {
	return &CapabilitySet{
		CapabilitySetType: CapabilitySetTypeInput,
		InputCapabilitySet: &InputCapabilitySet{
			InputFlags:          0x0001 | 0x0004 | 0x0010 | 0x0020, // INPUT_FLAG_SCANCODES, INPUT_FLAG_MOUSEX, INPUT_FLAG_UNICODE, INPUT_FLAG_FASTPATH_INPUT2
			KeyboardLayout:      0x00000409,                        // US
			KeyboardType:        0x00000004,                        // IBM enhanced (101- or 102-key) keyboard
			KeyboardFunctionKey: 12,
		},
	}
}

func (s *InputCapabilitySet) Serialize() []byte {
	buf := new(bytes.Buffer)

	_ = binary.Write(buf, binary.LittleEndian, s.InputFlags)
	_ = binary.Write(buf, binary.LittleEndian, uint16(0)) // padding
	_ = binary.Write(buf, binary.LittleEndian, s.KeyboardLayout)
	_ = binary.Write(buf, binary.LittleEndian, s.KeyboardType)
	_ = binary.Write(buf, binary.LittleEndian, s.KeyboardSubType)
	_ = binary.Write(buf, binary.LittleEndian, s.KeyboardFunctionKey)
	_ = binary.Write(buf, binary.LittleEndian, s.ImeFileName)

	return buf.Bytes()
}

func (s *InputCapabilitySet) Deserialize(wire io.Reader) error {
	var err error

	err = binary.Read(wire, binary.LittleEndian, &s.InputFlags)
	if err != nil {
		return err
	}

	var padding uint16
	err = binary.Read(wire, binary.LittleEndian, &padding)
	if err != nil {
		return err
	}

	err = binary.Read(wire, binary.LittleEndian, &s.KeyboardLayout)
	if err != nil {
		return err
	}

	err = binary.Read(wire, binary.LittleEndian, &s.KeyboardType)
	if err != nil {
		return err
	}

	err = binary.Read(wire, binary.LittleEndian, &s.KeyboardSubType)
	if err != nil {
		return err
	}

	err = binary.Read(wire, binary.LittleEndian, &s.KeyboardFunctionKey)
	if err != nil {
		return err
	}

	err = binary.Read(wire, binary.LittleEndian, &s.ImeFileName)
	if err != nil {
		return err
	}

	return nil
}

type BrushSupportLevel uint32

const (
	// BrushSupportLevelDefault BRUSH_DEFAULT
	BrushSupportLevelDefault BrushSupportLevel = 0

	// BrushSupportLevelColor8x8 BRUSH_COLOR_8x8
	BrushSupportLevelColor8x8 BrushSupportLevel = 1

	// BrushSupportLevelFull BRUSH_COLOR_FULL
	BrushSupportLevelFull BrushSupportLevel = 2
)

type BrushCapabilitySet struct {
	BrushSupportLevel BrushSupportLevel
}

func NewBrushCapabilitySet() *CapabilitySet {
	return &CapabilitySet{
		CapabilitySetType:  CapabilitySetTypeBrush,
		BrushCapabilitySet: &BrushCapabilitySet{},
	}
}

func (s *BrushCapabilitySet) Serialize() []byte {
	buf := new(bytes.Buffer)

	_ = binary.Write(buf, binary.LittleEndian, uint32(s.BrushSupportLevel))

	return buf.Bytes()
}

func (s *BrushCapabilitySet) Deserialize(wire io.Reader) error {
	return binary.Read(wire, binary.LittleEndian, &s.BrushSupportLevel)
}

type CacheDefinition struct {
	CacheEntries         uint16
	CacheMaximumCellSize uint16
}

func (d *CacheDefinition) Serialize() []byte {
	buf := new(bytes.Buffer)

	_ = binary.Write(buf, binary.LittleEndian, d.CacheEntries)
	_ = binary.Write(buf, binary.LittleEndian, d.CacheMaximumCellSize)

	return buf.Bytes()
}

func (d *CacheDefinition) Deserialize(wire io.Reader) error {
	var err error

	err = binary.Read(wire, binary.LittleEndian, &d.CacheEntries)
	if err != nil {
		return err
	}

	err = binary.Read(wire, binary.LittleEndian, &d.CacheMaximumCellSize)
	if err != nil {
		return err
	}

	return nil
}

type GlyphSupportLevel uint16

const (
	// GlyphSupportLevelNone GLYPH_SUPPORT_NONE
	GlyphSupportLevelNone GlyphSupportLevel = 0

	// GlyphSupportLevelPartial GLYPH_SUPPORT_PARTIAL
	GlyphSupportLevelPartial GlyphSupportLevel = 1

	// GlyphSupportLevelFull GLYPH_SUPPORT_FULL
	GlyphSupportLevelFull GlyphSupportLevel = 2

	// GlyphSupportLevelEncode GLYPH_SUPPORT_ENCODE
	GlyphSupportLevelEncode GlyphSupportLevel = 3
)

type GlyphCacheCapabilitySet struct {
	GlyphCache        [10]CacheDefinition
	FragCache         uint32
	GlyphSupportLevel GlyphSupportLevel
}

func NewGlyphCacheCapabilitySet() *CapabilitySet {
	return &CapabilitySet{
		CapabilitySetType:       CapabilitySetTypeGlyphCache,
		GlyphCacheCapabilitySet: &GlyphCacheCapabilitySet{},
	}
}

func (s *GlyphCacheCapabilitySet) Serialize() []byte {
	buf := new(bytes.Buffer)

	for i := range s.GlyphCache {
		buf.Write(s.GlyphCache[i].Serialize())
	}

	_ = binary.Write(buf, binary.LittleEndian, s.FragCache)
	_ = binary.Write(buf, binary.LittleEndian, s.GlyphSupportLevel)
	_ = binary.Write(buf, binary.LittleEndian, uint16(0)) // padding

	return buf.Bytes()
}

func (s *GlyphCacheCapabilitySet) Deserialize(wire io.Reader) error {
	var err error

	for i := range s.GlyphCache {
		err = s.GlyphCache[i].Deserialize(wire)
		if err != nil {
			return err
		}
	}

	err = binary.Read(wire, binary.LittleEndian, &s.FragCache)
	if err != nil {
		return err
	}

	err = binary.Read(wire, binary.LittleEndian, &s.GlyphSupportLevel)
	if err != nil {
		return err
	}

	var padding uint16
	err = binary.Read(wire, binary.LittleEndian, &padding)
	if err != nil {
		return err
	}

	return nil
}

type OffscreenBitmapCacheCapabilitySet struct {
	OffscreenSupportLevel uint32
	OffscreenCacheSize    uint16
	OffscreenCacheEntries uint16
}

func NewOffscreenBitmapCacheCapabilitySet() *CapabilitySet {
	return &CapabilitySet{
		CapabilitySetType:                 CapabilitySetTypeOffscreenBitmapCache,
		OffscreenBitmapCacheCapabilitySet: &OffscreenBitmapCacheCapabilitySet{},
	}
}

func (s *OffscreenBitmapCacheCapabilitySet) Serialize() []byte {
	buf := new(bytes.Buffer)

	_ = binary.Write(buf, binary.LittleEndian, s.OffscreenSupportLevel)
	_ = binary.Write(buf, binary.LittleEndian, s.OffscreenCacheSize)
	_ = binary.Write(buf, binary.LittleEndian, s.OffscreenCacheEntries)

	return buf.Bytes()
}

func (s *OffscreenBitmapCacheCapabilitySet) Deserialize(wire io.Reader) error {
	var err error

	err = binary.Read(wire, binary.LittleEndian, &s.OffscreenSupportLevel)
	if err != nil {
		return err
	}

	err = binary.Read(wire, binary.LittleEndian, &s.OffscreenCacheSize)
	if err != nil {
		return err
	}

	err = binary.Read(wire, binary.LittleEndian, &s.OffscreenCacheEntries)
	if err != nil {
		return err
	}

	return nil
}

type VirtualChannelCapabilitySet struct {
	Flags       uint32
	VCChunkSize uint32
}

func NewVirtualChannelCapabilitySet() *CapabilitySet {
	return &CapabilitySet{
		CapabilitySetType:           CapabilitySetTypeVirtualChannel,
		VirtualChannelCapabilitySet: &VirtualChannelCapabilitySet{},
	}
}

func (s *VirtualChannelCapabilitySet) Serialize() []byte {
	buf := new(bytes.Buffer)

	_ = binary.Write(buf, binary.LittleEndian, s.Flags)
	_ = binary.Write(buf, binary.LittleEndian, s.VCChunkSize)

	return buf.Bytes()
}

func (s *VirtualChannelCapabilitySet) Deserialize(wire io.Reader) error {
	var err error

	err = binary.Read(wire, binary.LittleEndian, &s.Flags)
	if err != nil {
		return err
	}

	err = binary.Read(wire, binary.LittleEndian, &s.VCChunkSize)
	if err != nil {
		return err
	}

	return nil
}

type DrawNineGridCacheCapabilitySet struct {
	drawNineGridSupportLevel uint32
	drawNineGridCacheSize    uint16
	drawNineGridCacheEntries uint16
}

func (s *DrawNineGridCacheCapabilitySet) Serialize() []byte {
	buf := new(bytes.Buffer)

	binary.Write(buf, binary.LittleEndian, s.drawNineGridSupportLevel)
	binary.Write(buf, binary.LittleEndian, s.drawNineGridCacheSize)
	binary.Write(buf, binary.LittleEndian, s.drawNineGridCacheEntries)

	return buf.Bytes()
}

func (s *DrawNineGridCacheCapabilitySet) Deserialize(wire io.Reader) error {
	var err error

	err = binary.Read(wire, binary.LittleEndian, &s.drawNineGridSupportLevel)
	if err != nil {
		return err
	}

	err = binary.Read(wire, binary.LittleEndian, &s.drawNineGridCacheSize)
	if err != nil {
		return err
	}

	err = binary.Read(wire, binary.LittleEndian, &s.drawNineGridCacheEntries)
	if err != nil {
		return err
	}

	return nil
}

type GDICacheEntries struct {
	GdipGraphicsCacheEntries        uint16
	GdipBrushCacheEntries           uint16
	GdipPenCacheEntries             uint16
	GdipImageCacheEntries           uint16
	GdipImageAttributesCacheEntries uint16
}

func (e *GDICacheEntries) Serialize() []byte {
	buf := new(bytes.Buffer)

	binary.Write(buf, binary.LittleEndian, e.GdipGraphicsCacheEntries)
	binary.Write(buf, binary.LittleEndian, e.GdipBrushCacheEntries)
	binary.Write(buf, binary.LittleEndian, e.GdipPenCacheEntries)
	binary.Write(buf, binary.LittleEndian, e.GdipImageCacheEntries)
	binary.Write(buf, binary.LittleEndian, e.GdipImageAttributesCacheEntries)

	return buf.Bytes()
}

func (e *GDICacheEntries) Deserialize(wire io.Reader) error {
	var err error

	err = binary.Read(wire, binary.LittleEndian, &e.GdipGraphicsCacheEntries)
	if err != nil {
		return err
	}

	err = binary.Read(wire, binary.LittleEndian, &e.GdipBrushCacheEntries)
	if err != nil {
		return err
	}

	err = binary.Read(wire, binary.LittleEndian, &e.GdipPenCacheEntries)
	if err != nil {
		return err
	}

	err = binary.Read(wire, binary.LittleEndian, &e.GdipImageCacheEntries)
	if err != nil {
		return err
	}

	err = binary.Read(wire, binary.LittleEndian, &e.GdipImageAttributesCacheEntries)
	if err != nil {
		return err
	}

	return nil
}

type GDICacheChunkSize struct {
	GdipGraphicsCacheChunkSize              uint16
	GdipObjectBrushCacheChunkSize           uint16
	GdipObjectPenCacheChunkSize             uint16
	GdipObjectImageAttributesCacheChunkSize uint16
}

func (s *GDICacheChunkSize) Serialize() []byte {
	buf := new(bytes.Buffer)

	binary.Write(buf, binary.LittleEndian, s.GdipGraphicsCacheChunkSize)
	binary.Write(buf, binary.LittleEndian, s.GdipObjectBrushCacheChunkSize)
	binary.Write(buf, binary.LittleEndian, s.GdipObjectPenCacheChunkSize)
	binary.Write(buf, binary.LittleEndian, s.GdipObjectImageAttributesCacheChunkSize)

	return buf.Bytes()
}

func (s *GDICacheChunkSize) Deserialize(wire io.Reader) error {
	var err error

	err = binary.Read(wire, binary.LittleEndian, &s.GdipGraphicsCacheChunkSize)
	if err != nil {
		return err
	}

	err = binary.Read(wire, binary.LittleEndian, &s.GdipObjectBrushCacheChunkSize)
	if err != nil {
		return err
	}

	err = binary.Read(wire, binary.LittleEndian, &s.GdipObjectPenCacheChunkSize)
	if err != nil {
		return err
	}

	err = binary.Read(wire, binary.LittleEndian, &s.GdipObjectImageAttributesCacheChunkSize)
	if err != nil {
		return err
	}

	return nil
}

type GDIImageCacheProperties struct {
	GdipObjectImageCacheChunkSize uint16
	GdipObjectImageCacheTotalSize uint16
	GdipObjectImageCacheMaxSize   uint16
}

func (p *GDIImageCacheProperties) Serialize() []byte {
	buf := new(bytes.Buffer)

	binary.Write(buf, binary.LittleEndian, p.GdipObjectImageCacheChunkSize)
	binary.Write(buf, binary.LittleEndian, p.GdipObjectImageCacheTotalSize)
	binary.Write(buf, binary.LittleEndian, p.GdipObjectImageCacheMaxSize)

	return buf.Bytes()
}

func (p *GDIImageCacheProperties) Deserialize(wire io.Reader) error {
	var err error

	err = binary.Read(wire, binary.LittleEndian, &p.GdipObjectImageCacheChunkSize)
	if err != nil {
		return err
	}

	err = binary.Read(wire, binary.LittleEndian, &p.GdipObjectImageCacheTotalSize)
	if err != nil {
		return err
	}

	err = binary.Read(wire, binary.LittleEndian, &p.GdipObjectImageCacheMaxSize)
	if err != nil {
		return err
	}

	return nil
}

type DrawGDIPlusCapabilitySet struct {
	drawGDIPlusSupportLevel  uint32
	GdipVersion              uint32
	drawGdiplusCacheLevel    uint32
	GdipCacheEntries         GDICacheEntries
	GdipCacheChunkSize       GDICacheChunkSize
	GdipImageCacheProperties GDIImageCacheProperties
}

func (s *DrawGDIPlusCapabilitySet) Serialize() []byte {
	buf := new(bytes.Buffer)

	binary.Write(buf, binary.LittleEndian, s.drawGDIPlusSupportLevel)
	binary.Write(buf, binary.LittleEndian, s.GdipVersion)
	binary.Write(buf, binary.LittleEndian, s.drawGdiplusCacheLevel)

	buf.Write(s.GdipCacheEntries.Serialize())
	buf.Write(s.GdipCacheChunkSize.Serialize())
	buf.Write(s.GdipImageCacheProperties.Serialize())

	return buf.Bytes()
}

func (s *DrawGDIPlusCapabilitySet) Deserialize(wire io.Reader) error {
	var err error

	err = binary.Read(wire, binary.LittleEndian, &s.drawGDIPlusSupportLevel)
	if err != nil {
		return err
	}

	err = binary.Read(wire, binary.LittleEndian, &s.GdipVersion)
	if err != nil {
		return err
	}

	err = binary.Read(wire, binary.LittleEndian, &s.drawGdiplusCacheLevel)
	if err != nil {
		return err
	}

	err = s.GdipCacheEntries.Deserialize(wire)
	if err != nil {
		return err
	}

	err = s.GdipCacheChunkSize.Deserialize(wire)
	if err != nil {
		return err
	}

	err = s.GdipImageCacheProperties.Deserialize(wire)
	if err != nil {
		return err
	}

	return nil
}

type SoundCapabilitySet struct {
	SoundFlags uint16
}

func NewSoundCapabilitySet() *CapabilitySet {
	return &CapabilitySet{
		CapabilitySetType:  CapabilitySetTypeSound,
		SoundCapabilitySet: &SoundCapabilitySet{},
	}
}

func (s *SoundCapabilitySet) Serialize() []byte {
	buf := new(bytes.Buffer)

	_ = binary.Write(buf, binary.LittleEndian, s.SoundFlags)
	_ = binary.Write(buf, binary.LittleEndian, uint16(0))

	return buf.Bytes()
}

func (s *SoundCapabilitySet) Deserialize(wire io.Reader) error {
	var err error

	err = binary.Read(wire, binary.LittleEndian, &s.SoundFlags)
	if err != nil {
		return err
	}

	var padding uint16
	err = binary.Read(wire, binary.LittleEndian, &padding)
	if err != nil {
		return err
	}

	return nil
}

type BitmapCacheHostSupportCapabilitySet struct{}

func NewBitmapCacheHostSupportCapabilitySet() *CapabilitySet {
	return &CapabilitySet{
		CapabilitySetType:                   CapabilitySetTypeBitmapCacheHostSupport,
		BitmapCacheHostSupportCapabilitySet: &BitmapCacheHostSupportCapabilitySet{},
	}
}

func (s *BitmapCacheHostSupportCapabilitySet) Deserialize(wire io.Reader) error {
	var (
		cacheVersion uint8
		padding1     uint8
		padding2     uint16
		err          error
	)

	err = binary.Read(wire, binary.LittleEndian, &cacheVersion)
	if err != nil {
		return err
	}

	err = binary.Read(wire, binary.LittleEndian, &padding1)
	if err != nil {
		return err
	}

	err = binary.Read(wire, binary.LittleEndian, &padding2)
	if err != nil {
		return err
	}

	return err
}

type ControlCapabilitySet struct{}

func (s *ControlCapabilitySet) Serialize() []byte {
	buf := new(bytes.Buffer)

	binary.Write(buf, binary.LittleEndian, uint16(0)) // controlFlags
	binary.Write(buf, binary.LittleEndian, uint16(0)) // remoteDetachFlag
	binary.Write(buf, binary.LittleEndian, uint16(2)) // controlInterest
	binary.Write(buf, binary.LittleEndian, uint16(2)) // detachInterest

	return buf.Bytes()
}

func (s *ControlCapabilitySet) Deserialize(wire io.Reader) error {
	padding := make([]byte, 8)

	return binary.Read(wire, binary.LittleEndian, &padding)
}

type WindowActivationCapabilitySet struct{}

func (s *WindowActivationCapabilitySet) Serialize() []byte {
	buf := new(bytes.Buffer)

	binary.Write(buf, binary.LittleEndian, uint16(0)) // helpKeyFlag
	binary.Write(buf, binary.LittleEndian, uint16(0)) // helpKeyIndexFlag
	binary.Write(buf, binary.LittleEndian, uint16(0)) // helpExtendedKeyFlag
	binary.Write(buf, binary.LittleEndian, uint16(0)) // windowManagerKeyFlag

	return buf.Bytes()
}

func (s *WindowActivationCapabilitySet) Deserialize(wire io.Reader) error {
	padding := make([]byte, 8)

	return binary.Read(wire, binary.LittleEndian, &padding)
}

type ShareCapabilitySet struct{}

func (s *ShareCapabilitySet) Serialize() []byte {
	buf := new(bytes.Buffer)

	binary.Write(buf, binary.LittleEndian, uint16(0)) // nodeID
	binary.Write(buf, binary.LittleEndian, uint16(0)) // pad2octets

	return buf.Bytes()
}

func (s *ShareCapabilitySet) Deserialize(wire io.Reader) error {
	padding := make([]byte, 4)

	return binary.Read(wire, binary.LittleEndian, &padding)
}

type FontCapabilitySet struct {
	fontSupportFlags uint16
}

func (s *FontCapabilitySet) Serialize() []byte {
	buf := new(bytes.Buffer)

	binary.Write(buf, binary.LittleEndian, s.fontSupportFlags)
	binary.Write(buf, binary.LittleEndian, uint16(0)) // padding

	return buf.Bytes()
}

func (s *FontCapabilitySet) Deserialize(wire io.Reader) error {
	padding := make([]byte, 2)

	err := binary.Read(wire, binary.LittleEndian, &s.fontSupportFlags)
	if err != nil {
		return err
	}

	return binary.Read(wire, binary.LittleEndian, &padding)
}

type MultifragmentUpdateCapabilitySet struct {
	MaxRequestSize uint32
}

func NewMultifragmentUpdateCapabilitySet() *CapabilitySet {
	return &CapabilitySet{
		CapabilitySetType:                CapabilitySetTypeMultifragmentUpdate,
		MultifragmentUpdateCapabilitySet: &MultifragmentUpdateCapabilitySet{},
	}
}

func (s *MultifragmentUpdateCapabilitySet) Serialize() []byte {
	buf := new(bytes.Buffer)

	binary.Write(buf, binary.LittleEndian, &s.MaxRequestSize)

	return buf.Bytes()
}

func (s *MultifragmentUpdateCapabilitySet) Deserialize(wire io.Reader) error {
	return binary.Read(wire, binary.LittleEndian, &s.MaxRequestSize)
}

type LargePointerCapabilitySet struct {
	LargePointerSupportFlags uint16
}

func (s *LargePointerCapabilitySet) Deserialize(wire io.Reader) error {
	return binary.Read(wire, binary.LittleEndian, &s.LargePointerSupportFlags)
}

type DesktopCompositionCapabilitySet struct {
	CompDeskSupportLevel uint16
}

func (s *DesktopCompositionCapabilitySet) Deserialize(wire io.Reader) error {
	return binary.Read(wire, binary.LittleEndian, &s.CompDeskSupportLevel)
}

type SurfaceCommandsCapabilitySet struct {
	CmdFlags uint32
}

func (s *SurfaceCommandsCapabilitySet) Deserialize(wire io.Reader) error {
	var (
		reserved uint32
		err      error
	)

	err = binary.Read(wire, binary.LittleEndian, &s.CmdFlags)
	if err != nil {
		return err
	}

	err = binary.Read(wire, binary.LittleEndian, &reserved)
	if err != nil {
		return err
	}

	return nil
}

type BitmapCodec struct {
	CodecGUID       [16]byte
	CodecID         uint8
	CodecProperties []byte
}

func (c *BitmapCodec) Deserialize(wire io.Reader) error {
	var err error

	err = binary.Read(wire, binary.LittleEndian, &c.CodecGUID)
	if err != nil {
		return err
	}

	err = binary.Read(wire, binary.LittleEndian, &c.CodecID)
	if err != nil {
		return err
	}

	var codecPropertiesLength uint16

	err = binary.Read(wire, binary.LittleEndian, &codecPropertiesLength)
	if err != nil {
		return err
	}

	c.CodecProperties = make([]byte, codecPropertiesLength)

	_, err = wire.Read(c.CodecProperties)
	if err != nil {
		return err
	}

	return nil
}

type BitmapCodecsCapabilitySet struct {
	BitmapCodecArray []BitmapCodec
}

func (s *BitmapCodecsCapabilitySet) Deserialize(wire io.Reader) error {
	var (
		bitmapCodecCount uint8
		err              error
	)

	err = binary.Read(wire, binary.LittleEndian, &bitmapCodecCount)
	if err != nil {
		return err
	}

	s.BitmapCodecArray = make([]BitmapCodec, bitmapCodecCount)

	for i := range s.BitmapCodecArray {
		err = s.BitmapCodecArray[i].Deserialize(wire)
		if err != nil {
			return err
		}
	}

	return nil
}

type RailCapabilitySet struct {
	RailSupportLevel uint32
}

func NewRailCapabilitySet() *CapabilitySet {
	return &CapabilitySet{
		CapabilitySetType: CapabilitySetTypeRail,
		RailCapabilitySet: &RailCapabilitySet{
			RailSupportLevel: 1, // TS_RAIL_LEVEL_SUPPORTED
		},
	}
}

func (s *RailCapabilitySet) Serialize() []byte {
	buf := new(bytes.Buffer)

	binary.Write(buf, binary.LittleEndian, s.RailSupportLevel)

	return buf.Bytes()
}

type WindowListCapabilitySet struct {
	WndSupportLevel     uint32
	NumIconCaches       uint8
	NumIconCacheEntries uint16
}

func NewWindowListCapabilitySet() *CapabilitySet {
	return &CapabilitySet{
		CapabilitySetType: CapabilitySetTypeWindow,
		WindowListCapabilitySet: &WindowListCapabilitySet{
			WndSupportLevel: 0, // TS_WINDOW_LEVEL_NOT_SUPPORTED
		},
	}
}

func (s *WindowListCapabilitySet) Serialize() []byte {
	buf := new(bytes.Buffer)

	binary.Write(buf, binary.LittleEndian, s.WndSupportLevel)
	binary.Write(buf, binary.LittleEndian, s.NumIconCaches)
	binary.Write(buf, binary.LittleEndian, s.NumIconCacheEntries)

	return buf.Bytes()
}
