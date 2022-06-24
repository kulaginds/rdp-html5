#include <stdlib.h>
#include <stdint.h>

typedef unsigned char BOOL;
typedef unsigned char BYTE;
typedef unsigned short UINT16;
typedef unsigned int UINT32;
typedef int INT32;
typedef UINT32 PIXEL;

BYTE g_MaskSpecialFgBg1 = 0x03;
BYTE g_MaskSpecialFgBg2 = 0x05;

BYTE g_MaskRegularRunLength = 0x1F;
BYTE g_MaskLiteRunLength = 0x0F;

#define TRUE 0xFF
#define FALSE 0x00

#define REGULAR_BG_RUN 0x00
#define MEGA_MEGA_BG_RUN 0xF0
#define REGULAR_FG_RUN 0x01
#define MEGA_MEGA_FG_RUN 0xF1
#define LITE_SET_FG_FG_RUN 0x0C
#define MEGA_MEGA_SET_FG_RUN 0xF6
#define LITE_DITHERED_RUN 0x0E
#define MEGA_MEGA_DITHERED_RUN 0xF8
#define REGULAR_COLOR_RUN 0x03
#define MEGA_MEGA_COLOR_RUN 0xF3
#define REGULAR_FGBG_IMAGE 0x02
#define MEGA_MEGA_FGBG_IMAGE 0xF2
#define LITE_SET_FG_FGBG_IMAGE 0x0D
#define MEGA_MEGA_SET_FGBG_IMAGE 0xF7
#define REGULAR_COLOR_IMAGE 0x04
#define MEGA_MEGA_COLOR_IMAGE 0xF4
#define SPECIAL_FGBG_1 0xF9
#define SPECIAL_FGBG_2 0xFA
#define SPECIAL_WHITE 0xFD
#define SPECIAL_BLACK 0xFE
