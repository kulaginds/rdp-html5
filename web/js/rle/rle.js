//
// Bitmasks
//
const g_MaskBit0 = 0x01; // Least significant bit
const g_MaskBit1 = 0x02;
const g_MaskBit2 = 0x04;
const g_MaskBit3 = 0x08;
const g_MaskBit4 = 0x10;
const g_MaskBit5 = 0x20;
const g_MaskBit6 = 0x40;
const g_MaskBit7 = 0x80; // Most significant bit

const g_MaskRegularRunLength = 0x1F;
const g_MaskLiteRunLength = 0x0F;

const g_MaskSpecialFgBg1 = 0x03;
const g_MaskSpecialFgBg2 = 0x05;

function RLEDecompressor(colorDepth) {
    this.colorDepth = colorDepth;
}

//
// Writes a pixel to the specified buffer.
//
RLEDecompressor.prototype.WritePixel = function (pbBuffer, pixel) {
    // TODO: implement
};

//
// Reads a pixel from the specified buffer.
//
RLEDecompressor.prototype.ReadPixel = function (pbBuffer) {
    // returns pixel
    // TODO: implement
};

//
// Returns the size of a pixel in bytes.
//
RLEDecompressor.prototype.GetPixelSize = function () {
    if (this.colorDepth === 8)
    {
        return 1;
    }
    else if (this.colorDepth === 15 || this.colorDepth === 16)
    {
        return 2;
    }
    else if (this.colorDepth === 24)
    {
        return 3;
    }
};

//
// Returns a pointer to the next pixel in the specified buffer.
//
RLEDecompressor.prototype.NextPixel = function (pbBuffer) {
    // return pbBuffer + GetPixelSize();
    // TODO: implement
};

//
// Reads the supplied order header and extracts the compression
// order code ID.
//
RLEDecompressor.prototype.ExtractCodeId = function (bOrderHdr) {
    // TODO: implement
};

//
// Returns a pointer to the data that follows the compression
// order header and optional run length.
//
RLEDecompressor.prototype.AdvanceOverOrderHeader = function (codeId, pbOrderHdr) {
    // TODO: implement
};

//
// Returns TRUE if the supplied code identifier is for a regular-form
// standard compression order. For example IsRegularCode(0x01) returns
// TRUE as 0x01 is the code ID for a Regular Foreground Run Order.
//
RLEDecompressor.prototype.IsRegularCode = function (codeId) {
    // TODO: implement
};

//
// Returns TRUE if the supplied code identifier is for a lite-form
// standard compression order. For example IsLiteCode(0x0E) returns
// TRUE as 0x0E is the code ID for a Lite Dithered Run Order.
//
RLEDecompressor.prototype.IsLiteCode = function (codeId) {
    // TODO: implement
};

//
// Returns TRUE if the supplied code identifier is for a MEGA_MEGA
// type extended compression order. For example IsMegaMegaCode(0xF0)
// returns TRUE as 0xF0 is the code ID for a MEGA_MEGA Background
// Run Order.
//
RLEDecompressor.prototype.IsMegaMegaCode = function (codeId) {
    // TODO: implement
};

//
// Returns a black pixel.
//
RLEDecompressor.prototype.GetColorBlack = function () {
    if (this.colorDepth === 8)
    {
        return 0x00;
    }
    else if (this.colorDepth === 15)
    {
        return0x0000;
    }
    else if (this.colorDepth === 16)
    {
        return 0x0000;
    }
    else if (this.colorDepth === 24)
    {
        return 0x000000;
    }
};

//
// Returns a white pixel.
//
RLEDecompressor.prototype.GetColorWhite = function () {
    if (this.colorDepth === 8)
    {
        //
        // Palette entry #255 holds black.
        //
        return 0xFF;
    }
    else if (this.colorDepth === 15)
    {
        //
        // 5 bits per RGB component:
        // 0111 1111 1111 1111 (binary)
        //
        return 0x7FFF;
    }
    else if (this.colorDepth === 16)
    {
        //
        // 5 bits for red, 6 bits for green, 5 bits for green:
        // 1111 1111 1111 1111 (binary)
        //
        return 0xFFFF;
    }
    else if (this.colorDepth === 24)
    {
        //
        // 8 bits per RGB component:
        // 1111 1111 1111 1111 1111 1111 (binary)
        //
        return 0xFFFFFF;
    }
};

//
// Extract the run length of a Regular-Form Foreground/Background
// Image Order.
//
RLEDecompressor.prototype.ExtractRunLengthRegularFgBg = function (pbOrderHdr) {
    // UINT runLength;
    //
    // runLength = *pbOrderHdr AND g_MaskRegularRunLength;
    // if (runLength == 0)
    // {
    //     runLength = *(pbOrderHdr + 1) + 1;
    // }
    // else
    // {
    //     runLength = runLength * 8;
    // }
    //
    // return runLength;
    // TODO: implement
};

//
// Extract the run length of a Lite-Form Foreground/Background
// Image Order.
//
RLEDecompressor.prototype.ExtractRunLengthLiteFgBg = function (pbOrderHdr) {
    // UINT runLength;
    //
    // runLength = *pbOrderHdr AND g_MaskLiteRunLength;
    // if (runLength == 0)
    // {
    //     runLength = *(pbOrderHdr + 1) + 1;
    // }
    // else
    // {
    //     runLength = runLength * 8;
    // }
    //
    // return runLength;
    // TODO: implement
};

//
// Extract the run length of a regular-form compression order.
//
RLEDecompressor.prototype.ExtractRunLengthRegular = function (pbOrderHdr) {
    // UINT runLength;
    //
    // runLength = *pbOrderHdr AND g_MaskRegularRunLength;
    // if (runLength == 0)
    // {
    //     //
    //     // An extended (MEGA) run.
    //     //
    //     runLength = *(pbOrderHdr + 1) + 32;
    // }
    //
    // return runLength;
    // TODO: implement
};

//
// Extract the run length of a lite-form compression order.
//
RLEDecompressor.prototype.ExtractRunLengthLite = function (pbOrderHdr) {
    // UINT runLength;
    //
    // runLength = *pbOrderHdr AND g_MaskLiteRunLength;
    // if (runLength == 0)
    // {
    //     //
    //     // An extended (MEGA) run.
    //     //
    //     runLength = *(pbOrderHdr + 1) + 16;
    // }
    //
    // return runLength;
    // TODO: implement
};

//
// Extract the run length of a MEGA_MEGA-type compression order.
//
RLEDecompressor.prototype.ExtractRunLengthMegaMega = function (pbOrderHdr) {
    // UINT runLength;
    //
    // pbOrderHdr = pbOrderHdr + 1;
    // runLength = ((UINT16) pbOrderHdr[0]) OR ((UINT16) pbOrderHdr[1] << 8);
    //
    // return runLength;
    // TODO: implement
};

const REGULAR_FGBG_IMAGE = 0x2;
const LITE_SET_FG_FGBG_IMAGE = 0xD;

//
// Extract the run length of a compression order.
//
RLEDecompressor.prototype.ExtractRunLength = function (code, pbOrderHdr) {
    if (code === REGULAR_FGBG_IMAGE)
    {
        return this.ExtractRunLengthRegularFgBg(pbOrderHdr);
    }

    if (code === LITE_SET_FG_FGBG_IMAGE)
    {
        return this.ExtractRunLengthLiteFgBg(pbOrderHdr);
    }

    if (this.IsRegularCode(code))
    {
        return this.ExtractRunLengthRegular(pbOrderHdr);
    }

    if (this.IsLiteCode(code))
    {
        return this.ExtractRunLengthLite(pbOrderHdr);
    }

    if (this.IsMegaMegaCode(code))
    {
        return this.ExtractRunLengthMegaMega(pbOrderHdr);
    }

    return 0;
};

//
// Write a foreground/background image to a destination buffer.
//
RLEDecompressor.prototype.WriteFgBgImage = function (pbDest, rowDelta, bitmask, fgPel, cBits) {
//     PIXEL xorPixel;
//
//     xorPixel = ReadPixel(pbDest - rowDelta);
//     if (bitmask AND g_MaskBit0)
//     {
//         WritePixel(pbDest, xorPixel XOR fgPel);
//     }
//     else
//     {
//         WritePixel(pbDest, xorPixel);
//     }
//     pbDest = NextPixel(pbDest);
//     cBits = cBits - 1;
//
//     if (cBits > 0)
//     {
//         xorPixel = ReadPixel(pbDest - rowDelta);
//         if (bitmask AND g_MaskBit1)
//         {
//             WritePixel(pbDest, xorPixel XOR fgPel);
//         }
//     else
//         {
//             WritePixel(pbDest, xorPixel);
//         }
//         pbDest = NextPixel(pbDest);
//         cBits = cBits - 1;
//
//         if (cBits > 0)
//         {
//             xorPixel = ReadPixel(pbDest - rowDelta);
//             if (bitmask AND g_MaskBit2)
//             {
//                 WritePixel(pbDest, xorPixel XOR fgPel);
//             }
//         else
//             {
//                 WritePixel(pbDest, xorPixel);
//             }
//             pbDest = NextPixel(pbDest);
//             cBits = cBits - 1;
//
//             if (cBits > 0)
//             {
//                 xorPixel = ReadPixel(pbDest - rowDelta);
//                 if (bitmask AND g_MaskBit3)
//                 {
//                     WritePixel(pbDest, xorPixel XOR fgPel);
//                 }
//             else
//                 {
//                     WritePixel(pbDest, xorPixel);
//                 }
//                 pbDest = NextPixel(pbDest);
//                 cBits = cBits - 1;
//
//                 if (cBits > 0)
//                 {
//                     xorPixel = ReadPixel(pbDest - rowDelta);
//                     if (bitmask AND g_MaskBit4)
//                     {
//                         WritePixel(pbDest, xorPixel XOR fgPel);
//                     }
//                 else
//                     {
//                         WritePixel(pbDest, xorPixel);
//                     }
//                     pbDest = NextPixel(pbDest);
//                     cBits = cBits - 1;
//
//                     if (cBits > 0)
//                     {
//                         xorPixel = ReadPixel(pbDest - rowDelta);
//                         if (bitmask AND g_MaskBit5)
//                         {
//                             WritePixel(pbDest, xorPixel XOR fgPel);
//                         }
//                     else
//                         {
//                             WritePixel(pbDest, xorPixel);
//                         }
//                         pbDest = NextPixel(pbDest);
//                         cBits = cBits - 1;
//
//                         if (cBits > 0)
//                         {
//                             xorPixel = ReadPixel(pbDest - rowDelta);
//                             if (bitmask AND g_MaskBit6)
//                             {
//                                 WritePixel(pbDest, xorPixel XOR fgPel);
//                             }
//                         else
//                             {
//                                 WritePixel(pbDest, xorPixel);
//                             }
//                             pbDest = NextPixel(pbDest);
//                             cBits = cBits - 1;
//
//                             if (cBits > 0)
//                             {
//                                 xorPixel = ReadPixel(pbDest - rowDelta);
//                                 if (bitmask AND g_MaskBit7)
//                                 {
//                                     WritePixel(pbDest, xorPixel XOR fgPel);
//                                 }
//                             else
//                                 {
//                                     WritePixel(pbDest, xorPixel);
//                                 }
//                                 pbDest = NextPixel(pbDest);
//                             }
//                         }
//                     }
//                 }
//             }
//         }
//     }
//
//     return pbDest;
    // TODO: implement
};

//
// Write a foreground/background image to a destination buffer
// for the first line of compressed data.
//
RLEDecompressor.prototype.WriteFirstLineFgBgImage = function (pbDest, bitmask, fgPel, cBits) {
//     if (bitmask AND g_MaskBit0)
//     {
//         WritePixel(pbDest, fgPel);
//     }
//     else
//     {
//         WritePixel(pbDest, GetColorBlack());
//     }
//     pbDest = NextPixel(pbDest);
//     cBits = cBits - 1;
//
//     if (cBits > 0)
//     {
//         if (bitmask AND g_MaskBit1)
//         {
//             WritePixel(pbDest, fgPel);
//         }
//     else
//         {
//             WritePixel(pbDest, GetColorBlack());
//         }
//         pbDest = NextPixel(pbDest);
//         cBits = cBits - 1;
//
//         if (cBits > 0)
//         {
//             if (bitmask AND g_MaskBit2)
//             {
//                 WritePixel(pbDest, fgPel);
//             }
//         else
//             {
//                 WritePixel(pbDest, GetColorBlack());
//             }
//             pbDest = NextPixel(pbDest);
//             cBits = cBits - 1;
//
//             if (cBits > 0)
//             {
//                 if (bitmask AND g_MaskBit3)
//                 {
//                     WritePixel(pbDest, fgPel);
//                 }
//             else
//                 {
//                     WritePixel(pbDest, GetColorBlack());
//                 }
//                 pbDest = NextPixel(pbDest);
//                 cBits = cBits - 1;
//
//                 if (cBits > 0)
//                 {
//                     if (bitmask AND g_MaskBit4)
//                     {
//                         WritePixel(pbDest, fgPel);
//                     }
//                 else
//                     {
//                         WritePixel(pbDest, GetColorBlack());
//                     }
//                     pbDest = NextPixel(pbDest);
//                     cBits = cBits - 1;
//
//                     if (cBits > 0)
//                     {
//                         if (bitmask AND g_MaskBit5)
//                         {
//                             WritePixel(pbDest, fgPel);
//                         }
//                     else
//                         {
//                             WritePixel(pbDest, GetColorBlack());
//                         }
//                         pbDest = NextPixel(pbDest);
//                         cBits = cBits - 1;
//
//                         if (cBits > 0)
//                         {
//                             if (bitmask AND g_MaskBit6)
//                             {
//                                 WritePixel(pbDest, fgPel);
//                             }
//                         else
//                             {
//                                 WritePixel(pbDest, GetColorBlack());
//                             }
//                             pbDest = NextPixel(pbDest);
//                             cBits = cBits - 1;
//
//                             if (cBits > 0)
//                             {
//                                 if (bitmask AND g_MaskBit7)
//                                 {
//                                     WritePixel(pbDest, fgPel);
//                                 }
//                             else
//                                 {
//                                     WritePixel(pbDest, GetColorBlack());
//                                 }
//                                 pbDest = NextPixel(pbDest);
//                             }
//                         }
//                     }
//                 }
//             }
//         }
//     }
//
//     return pbDest;
    // TODO: implement
};

//
// Decompress an RLE compressed bitmap.
//
RLEDecompressor.prototype.Decompress = function (pbSrcBuffer, cbSrcBuffer, pbDestBuffer, rowDelta) {
    // BYTE* pbSrc = pbSrcBuffer;
    // BYTE* pbEnd = pbSrcBuffer + cbSrcBuffer;
    // BYTE* pbDest = pbDestBuffer;
    //
    // PIXEL fgPel = GetColorWhite();
    // BOOL fInsertFgPel = FALSE;
    // BOOL fFirstLine = TRUE;
    //
    // BYTE bitmask;
    // PIXEL pixelA, pixelB;
    //
    // UINT runLength;
    // UINT code;
    //
    // while (pbSrc < pbEnd)
    // {
    //     //
    //     // Watch out for the end of the first scanline.
    //     //
    //     if (fFirstLine)
    //     {
    //         if (pbDest - pbDestBuffer >= rowDelta)
    //         {
    //             fFirstLine = FALSE;
    //             fInsertFgPel = FALSE;
    //         }
    //     }
    //
    //     //
    //     // Extract the compression order code ID from the compression
    //     // order header.
    //     //
    //     code = ExtractCodeId(*pbSrc);
    //
    //     //
    //     // Handle Background Run Orders.
    //     //
    //     if (code == REGULAR_BG_RUN OR
    //     code == MEGA_MEGA_BG_RUN)
    //     {
    //         runLength = ExtractRunLength(code, pbSrc);
    //         pbSrc = AdvanceOverOrderHeader(code, pbSrc);
    //
    //         if (fFirstLine)
    //         {
    //             if (fInsertFgPel)
    //             {
    //                 WritePixel(pbDest, fgPel);
    //                 pbDest = NextPixel(pbDest);
    //                 runLength = runLength - 1;
    //             }
    //             while (runLength > 0)
    //             {
    //                 WritePixel(pbDest, GetColorBlack());
    //                 pbDest = NextPixel(pbDest);
    //                 runLength = runLength - 1;
    //             }
    //         }
    //         else
    //         {
    //             if (fInsertFgPel)
    //             {
    //                 WritePixel(
    //                     pbDest,
    //                     ReadPixel(pbDest - rowDelta) XOR fgPel
    //             );
    //                 pbDest = NextPixel(pbDest);
    //                 runLength = runLength - 1;
    //             }
    //
    //             while (runLength > 0)
    //             {
    //                 WritePixel(pbDest, ReadPixel(pbDest - rowDelta));
    //                 pbDest = NextPixel(pbDest);
    //                 runLength = runLength - 1;
    //             }
    //         }
    //
    //         //
    //         // A follow-on background run order will need a
    //         // foreground pel inserted.
    //         //
    //         fInsertFgPel = TRUE;
    //         continue;
    //     }
    //
    //     //
    //     // For any of the other run-types a follow-on background run
    //     // order does not need a foreground pel inserted.
    //     //
    //     fInsertFgPel = FALSE;
    //
    //     //
    //     // Handle Foreground Run Orders.
    //     //
    //     if (code == REGULAR_FG_RUN OR
    //     code == MEGA_MEGA_FG_RUN OR
    //     code == LITE_SET_FG_FG_RUN OR
    //     code == MEGA_MEGA_SET_FG_RUN)
    //     {
    //         runLength = ExtractRunLength(code, pbSrc);
    //         pbSrc = AdvanceOverOrderHeader(code, pbSrc);
    //
    //         if (code == LITE_SET_FG_FG_RUN OR
    //         code == MEGA_MEGA_SET_FG_RUN)
    //         {
    //             fgPel = ReadPixel(pbSrc);
    //             pbSrc = NextPixel(pbSrc);
    //         }
    //
    //         while (runLength > 0)
    //         {
    //             if (fFirstLine)
    //             {
    //                 WritePixel(pbDest, fgPel);
    //                 pbDest = NextPixel(pbDest);
    //             }
    //             else
    //             {
    //                 WritePixel(
    //                     pbDest,
    //                     ReadPixel(pbDest - rowDelta) XOR fgPel
    //             );
    //                 pbDest = NextPixel(pbDest);
    //             }
    //
    //             runLength = runLength - 1;
    //         }
    //
    //         continue;
    //     }
    //
    //     //
    //     // Handle Dithered Run Orders.
    //     //
    //     if (code == LITE_DITHERED_RUN OR
    //     code == MEGA_MEGA_DITHERED_RUN)
    //     {
    //         runLength = ExtractRunLength(code, pbSrc);
    //         pbSrc = AdvanceOverOrderHeader(code, pbSrc);
    //
    //         pixelA = ReadPixel(pbSrc);
    //         pbSrc = NextPixel(pbSrc);
    //         pixelB = ReadPixel(pbSrc);
    //         pbSrc = NextPixel(pbSrc);
    //
    //         while (runLength > 0)
    //         {
    //             WritePixel(pbDest, pixelA);
    //             pbDest = NextPixel(pbDest);
    //             WritePixel(pbDest, pixelB);
    //             pbDest = NextPixel(pbDest);
    //
    //             runLength = runLength - 1;
    //         }
    //
    //         continue;
    //     }
    //
    //     //
    //     // Handle Color Run Orders.
    //     //
    //     if (code == REGULAR_COLOR_RUN OR
    //     code == MEGA_MEGA_COLOR_RUN)
    //     {
    //         runLength = ExtractRunLength(code, pbSrc);
    //         pbSrc = AdvanceOverOrderHeader(code, pbSrc);
    //
    //         pixelA = ReadPixel(pbSrc);
    //         pbSrc = NextPixel(pbSrc);
    //
    //         while (runLength > 0)
    //         {
    //             WritePixel(pbDest, pixelA);
    //             pbDest = NextPixel(pbDest);
    //
    //             runLength = runLength - 1;
    //         }
    //
    //         continue;
    //     }
    //
    //     //
    //     // Handle Foreground/Background Image Orders.
    //     //
    //     if (code == REGULAR_FGBG_IMAGE OR
    //     code == MEGA_MEGA_FGBG_IMAGE OR
    //     code == LITE_SET_FG_FGBG_IMAGE OR
    //     code == MEGA_MEGA_SET_FGBG_IMAGE)
    //     {
    //         runLength = ExtractRunLength(code, pbSrc);
    //         pbSrc = AdvanceOverOrderHeader(code, pbSrc);
    //
    //         if (code == LITE_SET_FG_FGBG_IMAGE OR
    //         code == MEGA_MEGA_SET_FGBG_IMAGE)
    //         {
    //             fgPel = ReadPixel(pbSrc);
    //             pbSrc = NextPixel(pbSrc);
    //         }
    //
    //         while (runLength > 8)
    //         {
    //             bitmask = *pbSrc;
    //             pbSrc = pbSrc + 1;
    //
    //             if (fFirstLine)
    //             {
    //                 pbDest = WriteFirstLineFgBgImage(
    //                     pbDest,
    //                     bitmask,
    //                     fgPel,
    //                     8
    //                 );
    //             }
    //             else
    //             {
    //                 pbDest = WriteFgBgImage(
    //                     pbDest,
    //                     rowDelta,
    //                     bitmask,
    //                     fgPel,
    //                     8
    //                 );
    //             }
    //
    //             runLength = runLength - 8;
    //         }
    //
    //         if (runLength > 0)
    //         {
    //             bitmask = *pbSrc;
    //             pbSrc = pbSrc + 1;
    //
    //             if (fFirstLine)
    //             {
    //                 pbDest = WriteFirstLineFgBgImage(
    //                     pbDest,
    //                     bitmask,
    //                     fgPel,
    //                     runLength
    //                 );
    //             }
    //             else
    //             {
    //                 pbDest = WriteFgBgImage(
    //                     pbDest,
    //                     rowDelta,
    //                     bitmask,
    //                     fgPel,
    //                     runLength
    //                 );
    //             }
    //         }
    //
    //         continue;
    //     }
    //
    //     //
    //     // Handle Color Image Orders.
    //     //
    //     if (code == REGULAR_COLOR_IMAGE OR
    //     code == MEGA_MEGA_COLOR_IMAGE)
    //     {
    //         UINT byteCount;
    //
    //         runLength = ExtractRunLength(code, pbSrc);
    //         pbSrc = AdvanceOverOrderHeader(code, pbSrc);
    //
    //         byteCount = runLength * GetColorDepth();
    //
    //         while (byteCount > 0)
    //         {
    //         *pbDest = *pbSrc;
    //             pbDest = pbDest + 1;
    //             pbSrc = pbSrc + 1;
    //
    //             byteCount = byteCount - 1;
    //         }
    //
    //         continue;
    //     }
    //
    //     //
    //     // Handle Special Order 1.
    //     //
    //     if (code == SPECIAL_FGBG_1)
    //     {
    //         if (fFirstLine)
    //         {
    //             pbDest = WriteFirstLineFgBgImage(
    //                 pbDest,
    //                 g_MaskSpecialFgBg1,
    //                 fgPel,
    //                 8
    //             );
    //         }
    //         else
    //         {
    //             pbDest = WriteFgBgImage(
    //                 pbDest,
    //                 rowDelta,
    //                 g_MaskSpecialFgBg1,
    //                 fgPel,
    //                 8
    //             );
    //         }
    //
    //         continue;
    //     }
    //
    //     //
    //     // Handle Special Order 2.
    //     //
    //     if (code == SPECIAL_FGBG_2)
    //     {
    //         if (fFirstLine)
    //         {
    //             pbDest = WriteFirstLineFgBgImage(
    //                 pbDest,
    //                 g_MaskSpecialFgBg2,
    //                 fgPel,
    //                 8
    //             );
    //         }
    //         else
    //         {
    //             pbDest = WriteFgBgImage(
    //                 pbDest,
    //                 rowDelta,
    //                 g_MaskSpecialFgBg2,
    //                 fgPel,
    //                 8
    //             );
    //         }
    //
    //         continue;
    //     }
    //
    //     //
    //     // Handle White Order.
    //     //
    //     if (code == WHITE)
    //     {
    //         WritePixel(pbDest, GetColorWhite());
    //         pbDest = NextPixel(pbDest);
    //
    //         continue;
    //     }
    //
    //     //
    //     // Handle Black Order.
    //     //
    //     if (code == BLACK)
    //     {
    //         WritePixel(pbDest, GetColorBlack());
    //         pbDest = NextPixel(pbDest);
    //
    //         continue;
    //     }
    // }
};
