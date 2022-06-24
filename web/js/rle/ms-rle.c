#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include "ms-rle.h"

//
// Returns the color depth (in bytes per pixel) that was selected
// for the RDP connection.
//
UINT
GetColorDepth()
{
 return 16;
}

//
// PIXEL is a dynamic type that is sized based on the current color
// depth being used for the RDP connection.
//
// if (GetColorDepth() == 8) then PIXEL is an 8-bit unsigned integer
// if (GetColorDepth() == 15) then PIXEL is a 16-bit unsigned integer
// if (GetColorDepth() == 16) then PIXEL is a 16-bit unsigned integer
// if (GetColorDepth() == 24) then PIXEL is a 24-bit unsigned integer
//

//
// Writes a pixel to the specified buffer.
//
VOID
WritePixel(
 BYTE* pbBuffer,
 PIXEL pixel
 )
{
    pbBuffer[0] = pixel & 0xFF;
    pbBuffer[1] = (pixel >> 8) & 0xFF;
}

//
// Reads a pixel from the specified buffer.
//
PIXEL
ReadPixel(
 BYTE* pbBuffer
 )
{
    return (PIXEL) (pbBuffer)[0] | ((pbBuffer)[1] << 8);
}

//
// Returns the size of a pixel in bytes.
//
UINT
GetPixelSize()
{
 UINT colorDepth = GetColorDepth();

 if (colorDepth == 8)
 {
     return 1;
 }

 if (colorDepth == 15 || colorDepth == 16)
 {
     return 2;
 }

 if (colorDepth == 24)
 {
     return 3;
 }

 return 0;
}

//
// Returns a pointer to the next pixel in the specified buffer.
//
BYTE*
NextPixel(
 BYTE* pbBuffer
 )
{
 return pbBuffer + GetPixelSize();
}

//
// Reads the supplied order header and extracts the compression
// order code ID.
//
UINT
ExtractCodeId(
 BYTE bOrderHdr
 )
{
    if ((bOrderHdr & 0xC0) != 0xC0)
    {
        /* REGULAR orders
         * (000x xxxx, 001x xxxx, 010x xxxx, 011x xxxx, 100x xxxx)
         */
        return bOrderHdr >> 5;
    }

    if ((bOrderHdr & 0xF0) == 0xF0)
    {
        /* MEGA and SPECIAL orders (0xF*) */
        return bOrderHdr;
    }

    /* LITE orders
     * 1100 xxxx, 1101 xxxx, 1110 xxxx)
     */
    return bOrderHdr >> 4;
}

//
// Returns TRUE if the supplied code identifier is for a regular-form
// standard compression order. For example IsRegularCode(0x01) returns
// TRUE as 0x01 is the code ID for a Regular Foreground Run Order.
//
BOOL
IsRegularCode(
 UINT codeId
 )
{
    switch (codeId) {
    case REGULAR_BG_RUN:
    case REGULAR_FG_RUN:
    case REGULAR_COLOR_RUN:
    case REGULAR_COLOR_IMAGE:
    case REGULAR_FGBG_IMAGE:
        return TRUE;
    }

    return FALSE;
}

//
// Returns TRUE if the supplied code identifier is for a lite-form
// standard compression order. For example IsLiteCode(0x0E) returns
// TRUE as 0x0E is the code ID for a Lite Dithered Run Order.
//
BOOL
IsLiteCode(
 UINT codeId
 )
{
    switch (codeId) {
    case LITE_SET_FG_FG_RUN:
    case LITE_DITHERED_RUN:
    case LITE_SET_FG_FGBG_IMAGE:
        return TRUE;
    }

    return FALSE;
}

//
// Returns TRUE if the supplied code identifier is for a MEGA_MEGA
// type extended compression order. For example IsMegaMegaCode(0xF0)
// returns TRUE as 0xF0 is the code ID for a MEGA_MEGA Background
// Run Order.
//
BOOL
IsMegaMegaCode(
 UINT codeId
 )
{
    switch (codeId) {
    case MEGA_MEGA_BG_RUN:
    case MEGA_MEGA_FG_RUN:
    case MEGA_MEGA_SET_FG_RUN:
    case MEGA_MEGA_DITHERED_RUN:
    case MEGA_MEGA_COLOR_RUN:
    case MEGA_MEGA_FGBG_IMAGE:
    case MEGA_MEGA_SET_FGBG_IMAGE:
    case MEGA_MEGA_COLOR_IMAGE:
        return TRUE;
    }

    return FALSE;
}

//
// Returns a black pixel.
//
PIXEL
GetColorBlack()
{
 UINT colorDepth = GetColorDepth();

 if (colorDepth == 8)
 {
     return (PIXEL) 0x00;
 }

 if (colorDepth == 15)
 {
     return (PIXEL) 0x0000;
 }

 if (colorDepth == 16)
 {
     return (PIXEL) 0x0000;
 }

 if (colorDepth == 24)
 {
     return (PIXEL) 0x000000;
 }

 return 0;
}

//
// Returns a white pixel.
//
PIXEL
GetColorWhite()
{
 UINT colorDepth = GetColorDepth();

 if (colorDepth == 8)
 {
     //
     // Palette entry #255 holds black.
     //
     return (PIXEL) 0xFF;
 }

 if (colorDepth == 15)
 {
     //
     // 5 bits per RGB component:
     // 0111 1111 1111 1111 (binary)
     //
     return (PIXEL) 0x7FFF;
 }

 if (colorDepth == 16)
 {
     //
     // 5 bits for red, 6 bits for green, 5 bits for green:
     // 1111 1111 1111 1111 (binary)
     //
     return (PIXEL) 0xFFFF;
 }

 if (colorDepth == 24)
 {
     //
     // 8 bits per RGB component:
     // 1111 1111 1111 1111 1111 1111 (binary)
     //
     return (PIXEL) 0xFFFFFF;
 }

 return 0;
}

//
// Returns a pointer to the data that follows the compression
// order header and optional run length.
//
BYTE*
AdvanceOverOrderHeader(
 UINT codeId,
 BYTE* pbOrderHdr
 )
{
    UINT32 advance = 1;
    UINT runLength;

    if (IsRegularCode(codeId)) {
        runLength = *pbOrderHdr AND g_MaskRegularRunLength;

        if (runLength == 0) {
            advance += 1; // mega length
        }
    }

    if (IsLiteCode(codeId)) {
        runLength = *pbOrderHdr AND g_MaskLiteRunLength;

        if (runLength == 0) {
            advance += 1; // mega length
        }
    }

    if (IsMegaMegaCode(codeId)) {
        advance += 2;
    }

    // if special - nothing to do

    pbOrderHdr += advance;

    return pbOrderHdr;
}

//
// Extract the run length of a Regular-Form Foreground/Background
// Image Order.
//
UINT
ExtractRunLengthRegularFgBg(
 BYTE* pbOrderHdr
 )
{
 UINT runLength;

 runLength = *pbOrderHdr AND g_MaskRegularRunLength;
 if (runLength == 0)
 {
     runLength = *(pbOrderHdr + 1) + 1;
 }
 else
 {
     runLength = runLength * 8;
 }

 return runLength;
}

//
// Extract the run length of a Lite-Form Foreground/Background
// Image Order.
//
UINT
ExtractRunLengthLiteFgBg(
 BYTE* pbOrderHdr
 )
{
 UINT runLength;

 runLength = *pbOrderHdr AND g_MaskLiteRunLength;
 if (runLength == 0)
 {
     runLength = *(pbOrderHdr + 1) + 1;
 }
 else
 {
     runLength = runLength * 8;
 }

 return runLength;
}

//
// Extract the run length of a regular-form compression order.
//
UINT
ExtractRunLengthRegular(
 BYTE* pbOrderHdr
 )
{
 UINT runLength;

 runLength = *pbOrderHdr AND g_MaskRegularRunLength;
 if (runLength == 0)
 {
     //
     // An extended (MEGA) run.
     //
     runLength = *(pbOrderHdr + 1) + 32;
 }

 return runLength;
}

//
// Extract the run length of a lite-form compression order.
//
UINT
ExtractRunLengthLite(
 BYTE* pbOrderHdr
 )
{
 UINT runLength;

 runLength = *pbOrderHdr AND g_MaskLiteRunLength;
 if (runLength == 0)
 {
     //
     // An extended (MEGA) run.
     //
     runLength = *(pbOrderHdr + 1) + 16;
 }

 return runLength;
}

//
// Extract the run length of a MEGA_MEGA-type compression order.
//
UINT
ExtractRunLengthMegaMega(
 BYTE* pbOrderHdr
 )
{
 UINT runLength;

 pbOrderHdr = pbOrderHdr + 1;
 runLength = ((UINT16) pbOrderHdr[0]) OR ((UINT16) pbOrderHdr[1] << 8);

 return runLength;
}

//
// Extract the run length of a compression order.
//
UINT
ExtractRunLength(
 UINT code,
 BYTE* pbOrderHdr
 )
{
 UINT runLength;

 if (code == REGULAR_FGBG_IMAGE)
 {
     runLength = ExtractRunLengthRegularFgBg(pbOrderHdr);
 }
 else if (code == LITE_SET_FG_FGBG_IMAGE)
 {
     runLength = ExtractRunLengthLiteFgBg(pbOrderHdr);
 }
 else if (IsRegularCode(code))
 {
     runLength = ExtractRunLengthRegular(pbOrderHdr);
 }
 else if (IsLiteCode(code))
 {
     runLength = ExtractRunLengthLite(pbOrderHdr);
 }
 else if (IsMegaMegaCode(code))
 {
     runLength = ExtractRunLengthMegaMega(pbOrderHdr);
 }
 else
 {
     runLength = 0;
 }

 return runLength;
}

//
// Write a foreground/background image to a destination buffer.
//
BYTE*
WriteFgBgImage(
 BYTE* pbDest,
 UINT rowDelta,
 BYTE bitmask,
 PIXEL fgPel,
 UINT cBits
 )
{
 PIXEL xorPixel;

 xorPixel = ReadPixel(pbDest - rowDelta);
 if (bitmask AND g_MaskBit0)
 {
     WritePixel(pbDest, xorPixel XOR fgPel);
 }
 else
 {
     WritePixel(pbDest, xorPixel);
 }
 pbDest = NextPixel(pbDest);
 cBits = cBits - 1;

 if (cBits > 0)
 {
     xorPixel = ReadPixel(pbDest - rowDelta);
     if (bitmask AND g_MaskBit1)
     {
         WritePixel(pbDest, xorPixel XOR fgPel);
     }
     else
     {
         WritePixel(pbDest, xorPixel);
     }
     pbDest = NextPixel(pbDest);
     cBits = cBits - 1;

     if (cBits > 0)
     {
         xorPixel = ReadPixel(pbDest - rowDelta);
         if (bitmask AND g_MaskBit2)
         {
             WritePixel(pbDest, xorPixel XOR fgPel);
         }
         else
         {
             WritePixel(pbDest, xorPixel);
         }
         pbDest = NextPixel(pbDest);
         cBits = cBits - 1;

         if (cBits > 0)
         {
             xorPixel = ReadPixel(pbDest - rowDelta);
             if (bitmask AND g_MaskBit3)
             {
                 WritePixel(pbDest, xorPixel XOR fgPel);
             }
             else
             {
                 WritePixel(pbDest, xorPixel);
             }
             pbDest = NextPixel(pbDest);
             cBits = cBits - 1;

             if (cBits > 0)
             {
                 xorPixel = ReadPixel(pbDest - rowDelta);
                 if (bitmask AND g_MaskBit4)
                 {
                     WritePixel(pbDest, xorPixel XOR fgPel);
                 }
                 else
                 {
                     WritePixel(pbDest, xorPixel);
                 }
                 pbDest = NextPixel(pbDest);
                 cBits = cBits - 1;

                 if (cBits > 0)
                 {
                     xorPixel = ReadPixel(pbDest - rowDelta);
                     if (bitmask AND g_MaskBit5)
                     {
                         WritePixel(pbDest, xorPixel XOR fgPel);
                     }
                     else
                     {
                         WritePixel(pbDest, xorPixel);
                     }
                     pbDest = NextPixel(pbDest);
                     cBits = cBits - 1;

                     if (cBits > 0)
                     {
                         xorPixel = ReadPixel(pbDest - rowDelta);
                         if (bitmask AND g_MaskBit6)
                         {
                             WritePixel(pbDest, xorPixel XOR fgPel);
                         }
                         else
                         {
                             WritePixel(pbDest, xorPixel);
                         }
                         pbDest = NextPixel(pbDest);
                         cBits = cBits - 1;

                         if (cBits > 0)
                         {
                             xorPixel = ReadPixel(pbDest - rowDelta);
                             if (bitmask AND g_MaskBit7)
                             {
                                 WritePixel(pbDest, xorPixel XOR fgPel);
                             }
                             else
                             {
                                 WritePixel(pbDest, xorPixel);
                             }
                             pbDest = NextPixel(pbDest);
                         }
                     }
                 }
             }
         }
     }
 }

 return pbDest;
}

//
// Write a foreground/background image to a destination buffer
// for the first line of compressed data.
//
BYTE*
WriteFirstLineFgBgImage(
 BYTE* pbDest,
 BYTE bitmask,
 PIXEL fgPel,
 UINT cBits
 )
{
 if (bitmask AND g_MaskBit0)
 {
     WritePixel(pbDest, fgPel);
 }
 else
 {
     WritePixel(pbDest, GetColorBlack());
 }
 pbDest = NextPixel(pbDest);
 cBits = cBits - 1;

 if (cBits > 0)
 {
     if (bitmask AND g_MaskBit1)
     {
         WritePixel(pbDest, fgPel);
     }
     else
     {
         WritePixel(pbDest, GetColorBlack());
     }
     pbDest = NextPixel(pbDest);
     cBits = cBits - 1;

     if (cBits > 0)
     {
         if (bitmask AND g_MaskBit2)
         {
             WritePixel(pbDest, fgPel);
         }
         else
         {
             WritePixel(pbDest, GetColorBlack());
         }
         pbDest = NextPixel(pbDest);
         cBits = cBits - 1;

         if (cBits > 0)
         {
             if (bitmask AND g_MaskBit3)
             {
                 WritePixel(pbDest, fgPel);
             }
             else
             {
                 WritePixel(pbDest, GetColorBlack());
             }
             pbDest = NextPixel(pbDest);
             cBits = cBits - 1;

             if (cBits > 0)
             {
                 if (bitmask AND g_MaskBit4)
                 {
                     WritePixel(pbDest, fgPel);
                 }
                 else
                 {
                     WritePixel(pbDest, GetColorBlack());
                 }
                 pbDest = NextPixel(pbDest);
                 cBits = cBits - 1;

                 if (cBits > 0)
                 {
                     if (bitmask AND g_MaskBit5)
                     {
                         WritePixel(pbDest, fgPel);
                     }
                     else
                     {
                         WritePixel(pbDest, GetColorBlack());
                     }
                     pbDest = NextPixel(pbDest);
                     cBits = cBits - 1;

                     if (cBits > 0)
                     {
                         if (bitmask AND g_MaskBit6)
                         {
                             WritePixel(pbDest, fgPel);
                         }
                         else
                         {
                             WritePixel(pbDest, GetColorBlack());
                         }
                         pbDest = NextPixel(pbDest);
                         cBits = cBits - 1;

                         if (cBits > 0)
                         {
                             if (bitmask AND g_MaskBit7)
                             {
                                 WritePixel(pbDest, fgPel);
                             }
                             else
                             {
                                 WritePixel(pbDest, GetColorBlack());
                             }
                             pbDest = NextPixel(pbDest);
                         }
                     }
                 }
             }
         }
     }
 }

 return pbDest;
}

//
// Decompress an RLE compressed bitmap.
//
BOOL
RleDecompress(
 BYTE* pbSrcBuffer, // Source buffer containing compressed bitmap
 UINT cbSrcBuffer, // Size of source buffer in bytes
 BYTE* pbDestBuffer, // Destination buffer
 UINT rowDelta // Scanline length in bytes
 )
{
    BYTE* pbSrc = pbSrcBuffer;
    BYTE* pbEnd = pbSrcBuffer + cbSrcBuffer;
    BYTE* pbDest = pbDestBuffer;

    PIXEL fgPel = GetColorWhite();
    BOOL fInsertFgPel = FALSE;
    BOOL fFirstLine = TRUE;

    BYTE bitmask;
    PIXEL pixelA, pixelB;

    UINT runLength;
    UINT code;

    while (pbSrc < pbEnd)
    {
        //
        // Watch out for the end of the first scanline.
        //
        if (fFirstLine)
        {
         if (pbDest - pbDestBuffer >= rowDelta)
         {
             fFirstLine = FALSE;
             fInsertFgPel = FALSE;
         }
        }

        //
        // Extract the compression order code ID from the compression
        // order header.
        //
        code = ExtractCodeId(*pbSrc);

        //
        // Handle Background Run Orders.
        //
        if (code == REGULAR_BG_RUN OR
         code == MEGA_MEGA_BG_RUN)
        {
         runLength = ExtractRunLength(code, pbSrc);
         pbSrc = AdvanceOverOrderHeader(code, pbSrc);

         if (fFirstLine)
         {
             if (fInsertFgPel)
             {
                 WritePixel(pbDest, fgPel);
                 pbDest = NextPixel(pbDest);
                 runLength = runLength - 1;
             }
             while (runLength > 0)
             {
                 WritePixel(pbDest, GetColorBlack());
                 pbDest = NextPixel(pbDest);
                 runLength = runLength - 1;
             }
         }
         else
         {
             if (fInsertFgPel)
             {
                 WritePixel(
                     pbDest,
                     ReadPixel(pbDest - rowDelta) XOR fgPel
                     );
                 pbDest = NextPixel(pbDest);
                 runLength = runLength - 1;
             }

             while (runLength > 0)
             {
                 WritePixel(pbDest, ReadPixel(pbDest - rowDelta));
                 pbDest = NextPixel(pbDest);
                 runLength = runLength - 1;
             }
         }

         //
         // A follow-on background run order will need a
         // foreground pel inserted.
         //
         fInsertFgPel = TRUE;
         continue;
        }

        //
        // For any of the other run-types a follow-on background run
        // order does not need a foreground pel inserted.
        //
        fInsertFgPel = FALSE;

        //
        // Handle Foreground Run Orders.
        //
        if (code == REGULAR_FG_RUN OR
         code == MEGA_MEGA_FG_RUN OR
         code == LITE_SET_FG_FG_RUN OR
         code == MEGA_MEGA_SET_FG_RUN)
        {
         runLength = ExtractRunLength(code, pbSrc);
         pbSrc = AdvanceOverOrderHeader(code, pbSrc);

         if (code == LITE_SET_FG_FG_RUN OR
             code == MEGA_MEGA_SET_FG_RUN)
         {
             fgPel = ReadPixel(pbSrc);
             pbSrc = NextPixel(pbSrc);
         }

         while (runLength > 0)
         {
             if (fFirstLine)
             {
                 WritePixel(pbDest, fgPel);
                 pbDest = NextPixel(pbDest);
             }
             else
             {
                 WritePixel(
                     pbDest,
                     ReadPixel(pbDest - rowDelta) XOR fgPel
                     );
                 pbDest = NextPixel(pbDest);
             }

             runLength = runLength - 1;
         }

         continue;
        }

        //
        // Handle Dithered Run Orders.
        //
        if (code == LITE_DITHERED_RUN OR
         code == MEGA_MEGA_DITHERED_RUN)
        {
         runLength = ExtractRunLength(code, pbSrc);
         pbSrc = AdvanceOverOrderHeader(code, pbSrc);

         pixelA = ReadPixel(pbSrc);
         pbSrc = NextPixel(pbSrc);
         pixelB = ReadPixel(pbSrc);
         pbSrc = NextPixel(pbSrc);

         while (runLength > 0)
         {
             WritePixel(pbDest, pixelA);
             pbDest = NextPixel(pbDest);
             WritePixel(pbDest, pixelB);
             pbDest = NextPixel(pbDest);

             runLength = runLength - 1;
         }

         continue;
        }

        //
        // Handle Color Run Orders.
        //
        if (code == REGULAR_COLOR_RUN OR
         code == MEGA_MEGA_COLOR_RUN)
        {
         runLength = ExtractRunLength(code, pbSrc);
         pbSrc = AdvanceOverOrderHeader(code, pbSrc);

         pixelA = ReadPixel(pbSrc);
         pbSrc = NextPixel(pbSrc);

         while (runLength > 0)
         {
             WritePixel(pbDest, pixelA);
             pbDest = NextPixel(pbDest);

             runLength = runLength - 1;
         }

         continue;
        }

        //
        // Handle Foreground/Background Image Orders.
        //
        if (code == REGULAR_FGBG_IMAGE OR
         code == MEGA_MEGA_FGBG_IMAGE OR
         code == LITE_SET_FG_FGBG_IMAGE OR
         code == MEGA_MEGA_SET_FGBG_IMAGE)
        {
         runLength = ExtractRunLength(code, pbSrc);
         pbSrc = AdvanceOverOrderHeader(code, pbSrc);

         if (code == LITE_SET_FG_FGBG_IMAGE OR
             code == MEGA_MEGA_SET_FGBG_IMAGE)
         {
             fgPel = ReadPixel(pbSrc);
             pbSrc = NextPixel(pbSrc);
         }

         while (runLength > 8)
         {
             bitmask = *pbSrc;
             pbSrc = pbSrc + 1;

             if (fFirstLine)
             {
                 pbDest = WriteFirstLineFgBgImage(
                     pbDest,
                     bitmask,
                     fgPel,
                     8
                     );
             }
             else
             {
                 pbDest = WriteFgBgImage(
                     pbDest,
                     rowDelta,
                     bitmask,
                     fgPel,
                     8
                     );
             }

             runLength = runLength - 8;
         }

         if (runLength > 0)
         {
             bitmask = *pbSrc;
             pbSrc = pbSrc + 1;

             if (fFirstLine)
             {
                 pbDest = WriteFirstLineFgBgImage(
                     pbDest,
                     bitmask,
                     fgPel,
                     runLength
                     );
             }
             else
             {
                 pbDest = WriteFgBgImage(
                     pbDest,
                     rowDelta,
                     bitmask,
                     fgPel,
                     runLength
                     );
             }
         }

         continue;
        }

        //
        // Handle Color Image Orders.
        //
        if (code == REGULAR_COLOR_IMAGE OR
         code == MEGA_MEGA_COLOR_IMAGE)
        {
         UINT byteCount;

         runLength = ExtractRunLength(code, pbSrc);
         pbSrc = AdvanceOverOrderHeader(code, pbSrc);

         byteCount = runLength * GetPixelSize();

         while (byteCount > 0)
         {
             *pbDest = *pbSrc;
             pbDest = pbDest + 1;
             pbSrc = pbSrc + 1;

             byteCount = byteCount - 1;
         }

         continue;
        }

        //
        // Handle Special Order 1.
        //
        if (code == SPECIAL_FGBG_1)
        {
        pbSrc = AdvanceOverOrderHeader(code, pbSrc);

         if (fFirstLine)
         {
             pbDest = WriteFirstLineFgBgImage(
                 pbDest,
                 g_MaskSpecialFgBg1,
                 fgPel,
                 8
                 );
         }
         else
         {
             pbDest = WriteFgBgImage(
                 pbDest,
                 rowDelta,
                 g_MaskSpecialFgBg1,
                 fgPel,
                 8
                 );
         }

         continue;
        }

        //
        // Handle Special Order 2.
        //
        if (code == SPECIAL_FGBG_2)
        {
         pbSrc = AdvanceOverOrderHeader(code, pbSrc);

         if (fFirstLine)
         {
             pbDest = WriteFirstLineFgBgImage(
                 pbDest,
                 g_MaskSpecialFgBg2,
                 fgPel,
                 8
                 );
         }
         else
         {
             pbDest = WriteFgBgImage(
                 pbDest,
                 rowDelta,
                 g_MaskSpecialFgBg2,
                 fgPel,
                 8
                 );
         }

         continue;
        }

        //
        // Handle White Order.
        //
        if (code == WHITE)
        {
         pbSrc = AdvanceOverOrderHeader(code, pbSrc);

         WritePixel(pbDest, GetColorWhite());
         pbDest = NextPixel(pbDest);

         continue;
        }

        //
        // Handle Black Order.
        //
        if (code == BLACK)
        {
         pbSrc = AdvanceOverOrderHeader(code, pbSrc);

         WritePixel(pbDest, GetColorBlack());
         pbDest = NextPixel(pbDest);

         continue;
        }

        return FALSE;
    }

    return TRUE;
}

VOID
flipV(BYTE* inA, UINT width, UINT height, BYTE* tmp)
{
    UINT rowDelta = width * 2;
    UINT half = height / 2;
    UINT bottomLine = rowDelta * (height - 1);
    UINT topLine = 0;
    UINT i;

    for (i = 0; i < half; ++i)
    {
        memcpy(tmp, inA + topLine, rowDelta);
        memcpy(inA + topLine, inA + bottomLine, rowDelta);
        memcpy(inA + bottomLine, tmp, rowDelta);

        topLine += rowDelta;
        bottomLine -= rowDelta;
    }
}

VOID buf2RGBA(BYTE* inA, BYTE* outA)
{
    UINT16 pel = *inA;
    pel |= *(inA + 1) << 8;

    PIXEL pelR = (pel & 0xF800) >> 11;
    PIXEL pelG = (pel & 0x7E0) >> 5;
    PIXEL pelB = pel & 0x1F;

    // 565 -> 888
    pelR = (pelR << 3 & ~0x7) | (pelR >> 2);
    pelG = (pelG << 2 & ~0x3) | (pelG >> 4);
    pelB = (pelB << 3 & ~0x7) | (pelB >> 2);

    *(outA) = pelR;
    *(outA + 1) = pelG;
    *(outA + 2) = pelB;
    *(outA + 3) = 255;
}

VOID rgb2rgba(BYTE* inA, UINT inLength, BYTE* outA)
{
    UINT inI = 0;
    UINT outI = 0;

    while (inI < inLength) {
        buf2RGBA(inA+inI, outA+outI);
        inI += 2;
        outI += 4;
    }
}

BOOL
RleDecompressAndFlipAndRGBA(
 BYTE* pbSrcBuffer, // Source buffer containing compressed bitmap
 UINT cbSrcBuffer, // Size of source buffer in bytes
 BYTE* pbDestBuffer, // Destination buffer
 UINT rowDelta, // Scanline length in bytes
 UINT ouputSize,
 UINT width,
 UINT height,
 BYTE* flipVTempPtr,
 BYTE* pbResultBuffer
 )
{
    BOOL result = RleDecompress(pbSrcBuffer, cbSrcBuffer, pbDestBuffer, rowDelta);

    if (result == FALSE)
    {
        return result;
    }

    flipV(pbDestBuffer, width, height, flipVTempPtr);
    rgb2rgba(pbDestBuffer, ouputSize, pbResultBuffer);
//    memcpy(pbDestBuffer, temp, ouputSize);

    return TRUE;
}
