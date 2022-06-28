package rle

type Code uint8

const (
	// RegularBackgroundRun REGULAR_BG_RUN
	// If run length is zero, then run length is encoded in the byte following and MUST be incremented by 32 to final value.
	RegularBackgroundRun Code = 0x0

	// RegularForegroundRun REGULAR_FG_RUN
	// If run length is zero, then run length is encoded in the byte following and MUST be incremented by 32 to final value.
	RegularForegroundRun Code = 0x1

	// RegularColorRun REGULAR_COLOR_RUN
	// If run length is zero, then run length is encoded in the byte following and MUST be incremented by 32 to final value.
	RegularColorRun Code = 0x3

	// RegularForegroundBackgroundImageRun REGULAR_FGBG_IMAGE
	// The run length is encoded in the five low-order bits of the order header byte and MUST be multiplied by 8 to give the final value.
	// If this value is zero, then the run length is encoded in the byte following the order header and MUST be incremented by 1 to give the final value.
	RegularForegroundBackgroundImageRun Code = 0x2

	// RegularColorImageRun REGULAR_COLOR_IMAGE
	// If this value is zero, then the run length is encoded in the byte following the order header and MUST be incremented by 32 to give the final value.
	RegularColorImageRun Code = 0x4
)

const (
	// LiteSetForegroundRun LITE_SET_FG_FG_RUN
	// If run length is zero, then run length is encoded in the byte following and MUST be incremented by 16 to final value.
	LiteSetForegroundRun Code = 0xC

	// LiteDitheredRun LITE_DITHERED_RUN
	// If run length is zero, then run length is encoded in the byte following and MUST be incremented by 16 to final value.
	LiteDitheredRun Code = 0xE

	// LiteSetForegroundForegroundBackgroundImageRun LITE_SET_FG_FGBG_IMAGE
	// The run length is encoded in the four low- order bits of the order header byte and MUST be multiplied by 8 to give the final value.
	// If this value is zero, then the run length is encoded in the byte following the order header and MUST be incremented by 1 to give the final value.
	LiteSetForegroundForegroundBackgroundImageRun Code = 0xD
)

const (
	// MegaMegaBackgroundRun MEGA_MEGA_BG_RUN
	// Mega form.
	MegaMegaBackgroundRun Code = 0xF0

	// MegaMegaForegroundRun MEGA_MEGA_FG_RUN
	MegaMegaForegroundRun Code = 0xF1

	// MegaMegaSetForegroundRun MEGA_MEGA_SET_FG_RUN
	MegaMegaSetForegroundRun Code = 0xF6

	// MegaMegaDitheredRun MEGA_MEGA_DITHERED_RUN
	MegaMegaDitheredRun Code = 0xF8

	// MegaMegaColorRun MEGA_MEGA_COLOR_RUN
	MegaMegaColorRun Code = 0xF3

	// MegaMegaForegroundBackgroundImageRun MEGA_MEGA_FGBG_IMAGE
	MegaMegaForegroundBackgroundImageRun Code = 0xF2

	// MegaMegaSetForegroundBackgroundImage MEGA_MEGA_SET_FGBG_IMAGE
	MegaMegaSetForegroundBackgroundImage Code = 0xF7

	// MegaMegaColorImage MEGA_MEGA_COLOR_IMAGE
	MegaMegaColorImage Code = 0xF4
)

const (
	// SpecialForegroundBackground1Run SPECIAL_FGBG_1
	// 8-bit bitmask of 0x03.
	SpecialForegroundBackground1Run Code = 0xF9

	// SpecialForegroundBackground2Run SPECIAL_FGBG_2
	// 8-bit bitmask of 0x05.
	SpecialForegroundBackground2Run Code = 0xFA

	// WhiteRun WHITE
	// Encodes a single white pixel.
	WhiteRun Code = 0xFD

	// BlackRun BLACK
	// Encodes a single black pixel.
	BlackRun Code = 0xFE
)

type Pixel uint16

const (
	WhitePixel Pixel = 0xFFFF
	BlackPixel Pixel = 0x0000
)
