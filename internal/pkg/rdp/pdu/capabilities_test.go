package pdu

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_GeneralCapabilitySet_Serialize(t *testing.T) {
	set := CapabilitySet{
		CapabilitySetType: CapabilitySetTypeGeneral,
		GeneralCapabilitySet: &GeneralCapabilitySet{
			OSMajorType: 1,
			OSMinorType: 3,
			ExtraFlags:  0x041d,
		},
	}

	expected := []byte{
		0x01, 0x00, 0x18, 0x00, 0x01, 0x00, 0x03, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0x00, 0x1d, 0x04,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	}
	actual := set.Serialize()

	require.Equal(t, expected, actual)
}

func Test_GeneralCapabilitySet_Serialize2(t *testing.T) {
	set := CapabilitySet{
		CapabilitySetType: CapabilitySetTypeGeneral,
		GeneralCapabilitySet: &GeneralCapabilitySet{
			OSMajorType: 1,
			OSMinorType: 3,
			ExtraFlags:  0x0415,
		},
	}

	expected, err := hex.DecodeString("010018000100030000020000000015040000000000000000")
	require.NoError(t, err)

	actual := set.Serialize()

	require.Equal(t, expected, actual)
}

func Test_BitmapCapabilitySet_Serialize(t *testing.T) {
	set := CapabilitySet{
		CapabilitySetType: CapabilitySetTypeBitmap,
		BitmapCapabilitySet: &BitmapCapabilitySet{
			PreferredBitsPerPixel: 0x18,
			Receive1BitPerPixel:   1,
			Receive4BitsPerPixel:  1,
			Receive8BitsPerPixel:  1,
			DesktopWidth:          1280,
			DesktopHeight:         1024,
			DesktopResizeFlag:     1,
		},
	}

	expected := []byte{
		0x02, 0x00, 0x1c, 0x00, 0x18, 0x00, 0x01, 0x00, 0x01, 0x00, 0x01, 0x00, 0x00, 0x05, 0x00, 0x04,
		0x00, 0x00, 0x01, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00,
	}
	actual := set.Serialize()

	require.Equal(t, expected, actual)
}

func Test_BitmapCapabilitySet_Serialize2(t *testing.T) {
	set := CapabilitySet{
		CapabilitySetType: CapabilitySetTypeBitmap,
		BitmapCapabilitySet: &BitmapCapabilitySet{
			PreferredBitsPerPixel: 0x18,
			Receive1BitPerPixel:   1,
			Receive4BitsPerPixel:  1,
			Receive8BitsPerPixel:  1,
			DesktopWidth:          1280,
			DesktopHeight:         800,
			DesktopResizeFlag:     0,
		},
	}

	expected, err := hex.DecodeString("02001c00180001000100010000052003000000000100000001000000")
	require.NoError(t, err)

	actual := set.Serialize()

	require.Equal(t, expected, actual)
}

func Test_OrderCapabilitySet_Serialize(t *testing.T) {
	set := CapabilitySet{
		CapabilitySetType: CapabilitySetTypeOrder,
		OrderCapabilitySet: &OrderCapabilitySet{
			OrderFlags: 0x002a,
			OrderSupport: [32]byte{
				0x01, 0x01, 0x01, 0x01, 0x01, 0x00, 0x00, 0x01, 0x01, 0x01, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01,
				0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x00, 0x01, 0x01, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00,
			},
			textFlags:        0x06a1,
			DesktopSaveSize:  0x38400,
			textANSICodePage: 0x04e4,
		},
	}

	expected := []byte{
		0x03, 0x00, 0x58, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x00, 0x14, 0x00, 0x00, 0x00, 0x01, 0x00,
		0x00, 0x00, 0x2a, 0x00, 0x01, 0x01, 0x01, 0x01, 0x01, 0x00, 0x00, 0x01, 0x01, 0x01, 0x00, 0x01,
		0x00, 0x00, 0x00, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x00, 0x01, 0x01, 0x01, 0x00,
		0x00, 0x00, 0x00, 0x00, 0xa1, 0x06, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x84, 0x03, 0x00,
		0x00, 0x00, 0x00, 0x00, 0xe4, 0x04, 0x00, 0x00,
	}
	actual := set.Serialize()

	require.Equal(t, expected, actual)
}

func Test_OrderCapabilitySet_Serialize2(t *testing.T) {
	set := CapabilitySet{
		CapabilitySetType: CapabilitySetTypeOrder,
		OrderCapabilitySet: &OrderCapabilitySet{
			OrderFlags:       0xa,
			OrderSupport:     [32]byte{},
			textFlags:        0,
			DesktopSaveSize:  0x38400,
			textANSICodePage: 0,
		},
	}

	expected, err := hex.DecodeString("030058000000000000000000000000000000000000000000010014000000010000000a0000000000000000000000000000000000000000000000000000000000000000000000000000000000008403000000000000000000")
	require.NoError(t, err)

	actual := set.Serialize()

	require.Equal(t, expected, actual)
}

func Test_BitmapCacheCapabilitySetRev1(t *testing.T) {
	set := CapabilitySet{
		CapabilitySetType:            CapabilitySetTypeBitmapCache,
		BitmapCacheCapabilitySetRev1: &BitmapCacheCapabilitySetRev1{},
	}

	expected, err := hex.DecodeString("04002800000000000000000000000000000000000000000000000000000000000000000000000000")
	require.NoError(t, err)

	actual := set.Serialize()

	require.Equal(t, expected, actual)
}

func Test_BitmapCacheCapabilitySetRev2(t *testing.T) {
	set := CapabilitySet{
		CapabilitySetType: CapabilitySetTypeBitmapCacheRev2,
		BitmapCacheCapabilitySetRev2: &BitmapCacheCapabilitySetRev2{
			CacheFlags:           0x0003,
			NumCellCaches:        3,
			BitmapCache0CellInfo: 0x00000078,
			BitmapCache1CellInfo: 0x00000078,
			BitmapCache2CellInfo: 0x800009fb,
		},
	}

	expected := []byte{
		0x13, 0x00, 0x28, 0x00, 0x03, 0x00, 0x00, 0x03, 0x78, 0x00, 0x00, 0x00, 0x78, 0x00, 0x00, 0x00,
		0xfb, 0x09, 0x00, 0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	}
	actual := set.Serialize()

	require.Equal(t, expected, actual)
}

func Test_ColorCacheCapabilitySet_Serialize(t *testing.T) {
	set := CapabilitySet{
		CapabilitySetType: CapabilitySetTypeColorCache,
		ColorCacheCapabilitySet: &ColorCacheCapabilitySet{
			ColorTableCacheSize: 6,
		},
	}

	expected := []byte{
		0x0a, 0x00, 0x08, 0x00, 0x06, 0x00, 0x00, 0x00,
	}
	actual := set.Serialize()

	require.Equal(t, expected, actual)
}

func Test_WindowActivationCapabilitySet_Serialize(t *testing.T) {
	set := CapabilitySet{
		CapabilitySetType:             CapabilitySetTypeActivation,
		WindowActivationCapabilitySet: &WindowActivationCapabilitySet{},
	}

	expected := []byte{
		0x07, 0x00, 0x0c, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	}
	actual := set.Serialize()

	require.Equal(t, expected, actual)
}

func Test_ControlCapabilitySet_Serialize(t *testing.T) {
	set := CapabilitySet{
		CapabilitySetType:    CapabilitySetTypeControl,
		ControlCapabilitySet: &ControlCapabilitySet{},
	}

	expected := []byte{
		0x05, 0x00, 0x0c, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0x00, 0x02, 0x00,
	}
	actual := set.Serialize()

	require.Equal(t, expected, actual)
}

func Test_PointerCapabilitySet_Serialize(t *testing.T) {
	set := CapabilitySet{
		CapabilitySetType: CapabilitySetTypePointer,
		PointerCapabilitySet: &PointerCapabilitySet{
			ColorPointerFlag:      1,
			ColorPointerCacheSize: 20,
			PointerCacheSize:      21,
		},
	}

	expected := []byte{
		0x08, 0x00, 0x0a, 0x00, 0x01, 0x00, 0x14, 0x00, 0x15, 0x00,
	}
	actual := set.Serialize()

	require.Equal(t, expected, actual)
}

func Test_ShareCapabilitySet_Serialize(t *testing.T) {
	set := CapabilitySet{
		CapabilitySetType:  CapabilitySetTypeShare,
		ShareCapabilitySet: &ShareCapabilitySet{},
	}

	expected := []byte{
		0x09, 0x00, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00,
	}
	actual := set.Serialize()

	require.Equal(t, expected, actual)
}

func Test_InputCapabilitySet_Serialize(t *testing.T) {
	set := CapabilitySet{
		CapabilitySetType: CapabilitySetTypeInput,
		InputCapabilitySet: &InputCapabilitySet{
			InputFlags:          0x0015,
			KeyboardLayout:      0x00000409,
			KeyboardType:        4,
			KeyboardFunctionKey: 12,
		},
	}

	expected := []byte{
		0x0d, 0x00, 0x58, 0x00, 0x15, 0x00, 0x00, 0x00, 0x09, 0x04, 0x00, 0x00, 0x04, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x0c, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	}
	actual := set.Serialize()

	require.Equal(t, expected, actual)
}

func Test_InputCapabilitySet_Serialize2(t *testing.T) {
	set := CapabilitySet{
		CapabilitySetType: CapabilitySetTypeInput,
		InputCapabilitySet: &InputCapabilitySet{
			InputFlags:          0x0015,
			KeyboardLayout:      0x00000409,
			KeyboardType:        4,
			KeyboardFunctionKey: 0,
		},
	}

	expected, err := hex.DecodeString("0d005800150000000904000004000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000")
	require.NoError(t, err)

	actual := set.Serialize()

	require.Equal(t, expected, actual)
}

func Test_SoundCapabilitySet_Serialize(t *testing.T) {
	set := CapabilitySet{
		CapabilitySetType: CapabilitySetTypeSound,
		SoundCapabilitySet: &SoundCapabilitySet{
			SoundFlags: 0x0001,
		},
	}

	expected := []byte{0x0c, 0x00, 0x08, 0x00, 0x01, 0x00, 0x00, 0x00}
	actual := set.Serialize()

	require.Equal(t, expected, actual)
}

func Test_SoundCapabilitySet_Serialize2(t *testing.T) {
	set := CapabilitySet{
		CapabilitySetType: CapabilitySetTypeSound,
		SoundCapabilitySet: &SoundCapabilitySet{
			SoundFlags: 0,
		},
	}

	expected, err := hex.DecodeString("0c00080000000000")
	require.NoError(t, err)

	actual := set.Serialize()

	require.Equal(t, expected, actual)
}

func Test_FontCapabilitySet(t *testing.T) {
	set := CapabilitySet{
		CapabilitySetType: CapabilitySetTypeFont,
		FontCapabilitySet: &FontCapabilitySet{
			fontSupportFlags: 0x0001,
		},
	}

	expected := []byte{0x0e, 0x00, 0x08, 0x00, 0x01, 0x00, 0x00, 0x00}
	actual := set.Serialize()

	require.Equal(t, expected, actual)
}

func Test_GlyphCacheCapabilitySet_Serialize(t *testing.T) {
	set := CapabilitySet{
		CapabilitySetType: CapabilitySetTypeGlyphCache,
		GlyphCacheCapabilitySet: &GlyphCacheCapabilitySet{
			GlyphCache: [10]CacheDefinition{
				{
					CacheEntries:         254,
					CacheMaximumCellSize: 4,
				},
				{
					CacheEntries:         254,
					CacheMaximumCellSize: 4,
				},
				{
					CacheEntries:         254,
					CacheMaximumCellSize: 8,
				},
				{
					CacheEntries:         254,
					CacheMaximumCellSize: 8,
				},
				{
					CacheEntries:         254,
					CacheMaximumCellSize: 16,
				},
				{
					CacheEntries:         254,
					CacheMaximumCellSize: 32,
				},
				{
					CacheEntries:         254,
					CacheMaximumCellSize: 64,
				},
				{
					CacheEntries:         254,
					CacheMaximumCellSize: 128,
				},
				{
					CacheEntries:         254,
					CacheMaximumCellSize: 256,
				},
				{
					CacheEntries:         64,
					CacheMaximumCellSize: 256,
				},
			},
			FragCache:         0x1000100,
			GlyphSupportLevel: 3,
		},
	}

	expected := []byte{
		0x10, 0x00, 0x34, 0x00, 0xfe, 0x00, 0x04, 0x00, 0xfe, 0x00, 0x04, 0x00, 0xfe, 0x00, 0x08, 0x00,
		0xfe, 0x00, 0x08, 0x00, 0xfe, 0x00, 0x10, 0x00, 0xfe, 0x00, 0x20, 0x00, 0xfe, 0x00, 0x40, 0x00,
		0xfe, 0x00, 0x80, 0x00, 0xfe, 0x00, 0x00, 0x01, 0x40, 0x00, 0x00, 0x01, 0x00, 0x01, 0x00, 0x01,
		0x03, 0x00, 0x00, 0x00,
	}
	actual := set.Serialize()

	require.Equal(t, expected, actual)
}

func Test_GlyphCacheCapabilitySet_Serialize2(t *testing.T) {
	set := CapabilitySet{
		CapabilitySetType: CapabilitySetTypeGlyphCache,
		GlyphCacheCapabilitySet: &GlyphCacheCapabilitySet{
			FragCache:         0,
			GlyphSupportLevel: 0,
		},
	}

	expected, err := hex.DecodeString("10003400000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000")
	require.NoError(t, err)

	actual := set.Serialize()

	require.Equal(t, expected, actual)
}

func Test_BrushCapabilitySet(t *testing.T) {
	set := CapabilitySet{
		CapabilitySetType: CapabilitySetTypeBrush,
		BrushCapabilitySet: &BrushCapabilitySet{
			BrushSupportLevel: 1,
		},
	}

	expected := []byte{0x0f, 0x00, 0x08, 0x00, 0x01, 0x00, 0x00, 0x00}
	actual := set.Serialize()

	require.Equal(t, expected, actual)
}

func Test_BrushCapabilitySet2(t *testing.T) {
	set := CapabilitySet{
		CapabilitySetType: CapabilitySetTypeBrush,
		BrushCapabilitySet: &BrushCapabilitySet{
			BrushSupportLevel: 0,
		},
	}

	expected, err := hex.DecodeString("0f00080000000000")
	require.NoError(t, err)

	actual := set.Serialize()

	require.Equal(t, expected, actual)
}

func Test_OffscreenBitmapCacheCapabilitySet_Serialize(t *testing.T) {
	set := CapabilitySet{
		CapabilitySetType: CapabilitySetTypeOffscreenBitmapCache,
		OffscreenBitmapCacheCapabilitySet: &OffscreenBitmapCacheCapabilitySet{
			OffscreenSupportLevel: 1,
			OffscreenCacheSize:    7680,
			OffscreenCacheEntries: 100,
		},
	}

	expected := []byte{0x11, 0x00, 0x0c, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x1e, 0x64, 0x00}
	actual := set.Serialize()

	require.Equal(t, expected, actual)
}

func Test_OffscreenBitmapCacheCapabilitySet_Serialize2(t *testing.T) {
	set := CapabilitySet{
		CapabilitySetType: CapabilitySetTypeOffscreenBitmapCache,
		OffscreenBitmapCacheCapabilitySet: &OffscreenBitmapCacheCapabilitySet{
			OffscreenSupportLevel: 0,
			OffscreenCacheSize:    0,
			OffscreenCacheEntries: 0,
		},
	}

	expected, err := hex.DecodeString("11000c000000000000000000")
	require.NoError(t, err)

	actual := set.Serialize()

	require.Equal(t, expected, actual)
}

func Test_VirtualChannelCapabilitySet_Serialize(t *testing.T) {
	set := CapabilitySet{
		CapabilitySetType: CapabilitySetTypeVirtualChannel,
		VirtualChannelCapabilitySet: &VirtualChannelCapabilitySet{
			Flags: 0x00000001,
		},
	}

	expected := []byte{0x14, 0x00, 0x0c, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	actual := set.Serialize()

	require.Equal(t, expected, actual)
}

func Test_VirtualChannelCapabilitySet_Serialize2(t *testing.T) {
	set := CapabilitySet{
		CapabilitySetType: CapabilitySetTypeVirtualChannel,
		VirtualChannelCapabilitySet: &VirtualChannelCapabilitySet{
			Flags: 0,
		},
	}

	expected, err := hex.DecodeString("14000c000000000000000000")
	require.NoError(t, err)

	actual := set.Serialize()

	require.Equal(t, expected, actual)
}

func Test_DrawNineGridCacheCapabilitySet_Serialize(t *testing.T) {
	set := CapabilitySet{
		CapabilitySetType: CapabilitySetTypeDrawNineGridCache,
		DrawNineGridCacheCapabilitySet: &DrawNineGridCacheCapabilitySet{
			drawNineGridSupportLevel: 2,
			drawNineGridCacheSize:    2560,
			drawNineGridCacheEntries: 256,
		},
	}

	expected := []byte{0x15, 0x00, 0x0c, 0x00, 0x02, 0x00, 0x00, 0x00, 0x00, 0x0a, 0x00, 0x01}
	actual := set.Serialize()

	require.Equal(t, expected, actual)
}

func Test_DrawGDIPlusCapabilitySet_Serialize(t *testing.T) {
	set := CapabilitySet{
		CapabilitySetType:        CapabilitySetTypeDrawGDIPlus,
		DrawGDIPlusCapabilitySet: &DrawGDIPlusCapabilitySet{},
	}

	expected := []byte{
		0x16, 0x00, 0x28, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	}
	actual := set.Serialize()

	require.Equal(t, expected, actual)
}
