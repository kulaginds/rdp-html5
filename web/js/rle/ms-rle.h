typedef unsigned char BYTE;
typedef unsigned int UINT;
typedef void VOID;
typedef unsigned short PIXEL;
typedef unsigned char BOOL;
typedef unsigned short UINT16;
typedef unsigned int UINT32;

#define TRUE 0xFF
#define FALSE 0x00
#define AND &
#define OR |
#define XOR ^

//
// Bitmasks
//
BYTE g_MaskBit0 = 0x01; // Least significant bit
BYTE g_MaskBit1 = 0x02;
BYTE g_MaskBit2 = 0x04;
BYTE g_MaskBit3 = 0x08;
BYTE g_MaskBit4 = 0x10;
BYTE g_MaskBit5 = 0x20;
BYTE g_MaskBit6 = 0x40;
BYTE g_MaskBit7 = 0x80; // Most significant bit

BYTE g_MaskRegularRunLength = 0x1F;
BYTE g_MaskLiteRunLength = 0x0F;

BYTE g_MaskSpecialFgBg1 = 0x03;
BYTE g_MaskSpecialFgBg2 = 0x05;

#define REGULAR_FGBG_IMAGE 0x2
#define LITE_SET_FG_FGBG_IMAGE 0xD

#define REGULAR_BG_RUN 0x0
#define MEGA_MEGA_BG_RUN 0xF0

#define REGULAR_FG_RUN 0x1
#define MEGA_MEGA_FG_RUN 0xF1

#define LITE_SET_FG_FG_RUN 0xC
#define MEGA_MEGA_SET_FG_RUN 0xF6

#define LITE_DITHERED_RUN 0xE
#define MEGA_MEGA_DITHERED_RUN 0xF8

#define REGULAR_COLOR_RUN 0x3
#define MEGA_MEGA_COLOR_RUN 0xF3

#define MEGA_MEGA_FGBG_IMAGE 0xF2
#define MEGA_MEGA_SET_FGBG_IMAGE 0xF7

#define REGULAR_COLOR_IMAGE 0x4
#define MEGA_MEGA_COLOR_IMAGE 0xF4

#define SPECIAL_FGBG_1 0xF9
#define SPECIAL_FGBG_2 0xFA

#define WHITE 0xFD
#define BLACK 0xFE
