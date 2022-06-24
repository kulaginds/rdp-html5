#include "freerdp-rle.h"

#define BLACK_PIXEL 0x0000
#define WHITE_PIXEL 0xFFFF

#define UNROLL_BODY(_exp, _count)      \
	do                                 \
	{                                  \
		size_t x;                      \
		for (x = 0; x < (_count); x++) \
		{                              \
			do                         \
				_exp while (FALSE);    \
		}                              \
	} while (FALSE)

#define UNROLL_MULTIPLE(_condition, _exp, _count) \
	do                                            \
	{                                             \
		while ((_condition) >= _count)            \
		{                                         \
			UNROLL_BODY(_exp, _count);            \
			(_condition) -= _count;               \
		}                                         \
	} while (FALSE)

#define UNROLL(_condition, _exp)               \
	do                                         \
	{                                          \
		UNROLL_MULTIPLE(_condition, _exp, 16); \
		UNROLL_MULTIPLE(_condition, _exp, 4);  \
		UNROLL_MULTIPLE(_condition, _exp, 1);  \
	} while (FALSE)

#undef DESTWRITEPIXEL
#undef DESTREADPIXEL
#undef SRCREADPIXEL
#undef DESTNEXTPIXEL
#undef SRCNEXTPIXEL
#undef WRITEFGBGIMAGE
#undef WRITEFIRSTLINEFGBGIMAGE
#undef RLEDECOMPRESS
#undef RLEEXTRA
#undef WHITE_PIXEL
#define WHITE_PIXEL 0xFFFF
#define DESTWRITEPIXEL(_buf, _pix) write_pixel_16(_buf, _pix)
#define DESTREADPIXEL(_pix, _buf) _pix = ((UINT16*)(_buf))[0]
#define SRCREADPIXEL(_pix, _buf) _pix = (_buf)[0] | ((_buf)[1] << 8)
#define DESTNEXTPIXEL(_buf) _buf += 2
#define SRCNEXTPIXEL(_buf) _buf += 2
#define WRITEFGBGIMAGE WriteFgBgImage16to16
#define WRITEFIRSTLINEFGBGIMAGE WriteFirstLineFgBgImage16to16
#define RLEDECOMPRESS RleDecompress16to16
#define RLEEXTRA
#undef ENSURE_CAPACITY
#define ENSURE_CAPACITY(_start, _end, _size) ensure_capacity(_start, _end, _size, 2)

/**
 * Reads the supplied order header and extracts the compression
 * order code ID.
 */
UINT32 ExtractCodeId(BYTE bOrderHdr)
{
	if ((bOrderHdr & 0xC0) != 0xC0)
	{
		/* REGULAR orders
		 * (000x xxxx, 001x xxxx, 010x xxxx, 011x xxxx, 100x xxxx)
		 */
		return bOrderHdr >> 5;
	}
	else if ((bOrderHdr & 0xF0) == 0xF0)
	{
		/* MEGA and SPECIAL orders (0xF*) */
		return bOrderHdr;
	}
	else
	{
		/* LITE orders
		 * 1100 xxxx, 1101 xxxx, 1110 xxxx)
		 */
		return bOrderHdr >> 4;
	}
}

/**
 * Extract the run length of a compression order.
 */
UINT32 ExtractRunLength(UINT32 code, const BYTE* pbOrderHdr, const BYTE* pbEnd,
                                      UINT32* advance)
{
	UINT32 runLength = 0;
	UINT32 ladvance = 1;

	if (pbOrderHdr >= pbEnd)
		return 0;

	switch (code)
	{
		case REGULAR_FGBG_IMAGE:
			runLength = (*pbOrderHdr) & g_MaskRegularRunLength;

			if (runLength == 0)
			{
				if (pbOrderHdr + 1 >= pbEnd)
					return 0;
				runLength = (*(pbOrderHdr + 1)) + 1;
				ladvance += 1;
			}
			else
			{
				runLength = runLength * 8;
			}

			break;

		case LITE_SET_FG_FGBG_IMAGE:
			runLength = (*pbOrderHdr) & g_MaskLiteRunLength;

			if (runLength == 0)
			{
				if (pbOrderHdr + 1 >= pbEnd)
					return 0;

				runLength = (*(pbOrderHdr + 1)) + 1;
				ladvance += 1;
			}
			else
			{
				runLength = runLength * 8;
			}

			break;

		case REGULAR_BG_RUN:
		case REGULAR_FG_RUN:
		case REGULAR_COLOR_RUN:
		case REGULAR_COLOR_IMAGE:
			runLength = (*pbOrderHdr) & g_MaskRegularRunLength;

			if (runLength == 0)
			{
				/* An extended (MEGA) run. */
				if (pbOrderHdr + 1 >= pbEnd)
					return 0;

				runLength = (*(pbOrderHdr + 1)) + 32;
				ladvance += 1;
			}

			break;

		case LITE_SET_FG_FG_RUN:
		case LITE_DITHERED_RUN:
			runLength = (*pbOrderHdr) & g_MaskLiteRunLength;

			if (runLength == 0)
			{
				/* An extended (MEGA) run. */
				if (pbOrderHdr + 1 >= pbEnd)
					return 0;

				runLength = (*(pbOrderHdr + 1)) + 16;
				ladvance += 1;
			}

			break;

		case MEGA_MEGA_BG_RUN:
		case MEGA_MEGA_FG_RUN:
		case MEGA_MEGA_SET_FG_RUN:
		case MEGA_MEGA_DITHERED_RUN:
		case MEGA_MEGA_COLOR_RUN:
		case MEGA_MEGA_FGBG_IMAGE:
		case MEGA_MEGA_SET_FGBG_IMAGE:
		case MEGA_MEGA_COLOR_IMAGE:
			if (pbOrderHdr + 2 >= pbEnd)
				return 0;

			runLength = ((UINT16)pbOrderHdr[1]) | ((UINT16)(pbOrderHdr[2] << 8));
			ladvance += 2;
			break;
	}

	*advance = ladvance;
	return runLength;
}

void write_pixel_16(BYTE* _buf, UINT16 _pix)
{
	_buf[0] = _pix & 0xFF;
	_buf[1] = (_pix >> 8) & 0xFF;
}

BOOL ensure_capacity(const BYTE* start, const BYTE* end, size_t size, size_t base)
{
	const size_t available = (uintptr_t)end - (uintptr_t)start;
	const BOOL rc = available >= size * base;
	return rc && (start <= end);
}

/**
 * Write a foreground/background image to a destination buffer
 * for the first line of compressed data.
 */
BYTE* WRITEFIRSTLINEFGBGIMAGE(BYTE* pbDest, const BYTE* pbDestEnd, BYTE bitmask,
                                            PIXEL fgPel, UINT32 cBits)
{
	BYTE mask = 0x01;

	if (cBits > 8)
		return NULL;

	if (!ENSURE_CAPACITY(pbDest, pbDestEnd, cBits))
		return NULL;

	UNROLL(cBits, {
		UINT32 data;

		if (bitmask & mask)
			data = fgPel;
		else
			data = BLACK_PIXEL;

		DESTWRITEPIXEL(pbDest, data);
		DESTNEXTPIXEL(pbDest);
		mask = mask << 1;
	});
	return pbDest;
}

/**
 * Write a foreground/background image to a destination buffer.
 */
BYTE* WRITEFGBGIMAGE(BYTE* pbDest, const BYTE* pbDestEnd, UINT32 rowDelta,
                                   BYTE bitmask, PIXEL fgPel, INT32 cBits)
{
	PIXEL xorPixel;
	BYTE mask = 0x01;

	if (cBits > 8)
		return NULL;

	if (!ENSURE_CAPACITY(pbDest, pbDestEnd, cBits))
		return NULL;

	UNROLL(cBits, {
		UINT32 data;
		DESTREADPIXEL(xorPixel, pbDest - rowDelta);

		if (bitmask & mask)
			data = xorPixel ^ fgPel;
		else
			data = xorPixel;

		DESTWRITEPIXEL(pbDest, data);
		DESTNEXTPIXEL(pbDest);
		mask = mask << 1;
	});
	return pbDest;
}

/**
 * Decompress an RLE compressed bitmap.
 */
BOOL rle_decompress(const BYTE* pbSrcBuffer, UINT32 cbSrcBuffer, BYTE* pbDestBuffer,
                                 UINT32 rowDelta, UINT32 width, UINT32 height)
{
	const BYTE* pbSrc = pbSrcBuffer;
	const BYTE* pbEnd;
	const BYTE* pbDestEnd;
	BYTE* pbDest = pbDestBuffer;
	PIXEL temp;
	PIXEL fgPel = WHITE_PIXEL;
	BOOL fInsertFgPel = FALSE;
	BOOL fFirstLine = TRUE;
	BYTE bitmask;
	PIXEL pixelA, pixelB;
	UINT32 runLength;
	UINT32 code;
	UINT32 advance = 0;
	RLEEXTRA

	if ((rowDelta == 0) || (rowDelta < width))
		return FALSE;

	if (!pbSrcBuffer || !pbDestBuffer)
		return FALSE;

	pbEnd = pbSrcBuffer + cbSrcBuffer;
	pbDestEnd = pbDestBuffer + rowDelta * height;

	while (pbSrc < pbEnd)
	{
		/* Watch out for the end of the first scanline. */
		if (fFirstLine)
		{
			if ((UINT32)(pbDest - pbDestBuffer) >= rowDelta)
			{
				fFirstLine = FALSE;
				fInsertFgPel = FALSE;
			}
		}

		/*
		   Extract the compression order code ID from the compression
		   order header.
		*/
		code = ExtractCodeId(*pbSrc);

		/* Handle Background Run Orders. */
		if (code == REGULAR_BG_RUN || code == MEGA_MEGA_BG_RUN)
		{
			runLength = ExtractRunLength(code, pbSrc, pbEnd, &advance);
			pbSrc = pbSrc + advance;

			if (fFirstLine)
			{
				if (fInsertFgPel)
				{
					if (!ENSURE_CAPACITY(pbDest, pbDestEnd, 1))
						return FALSE;

					DESTWRITEPIXEL(pbDest, fgPel);
					DESTNEXTPIXEL(pbDest);
					runLength = runLength - 1;
				}

				if (!ENSURE_CAPACITY(pbDest, pbDestEnd, runLength))
					return FALSE;

				UNROLL(runLength, {
					DESTWRITEPIXEL(pbDest, BLACK_PIXEL);
					DESTNEXTPIXEL(pbDest);
				});
			}
			else
			{
				if (fInsertFgPel)
				{
					DESTREADPIXEL(temp, pbDest - rowDelta);

					if (!ENSURE_CAPACITY(pbDest, pbDestEnd, 1))
						return FALSE;

					DESTWRITEPIXEL(pbDest, temp ^ fgPel);
					DESTNEXTPIXEL(pbDest);
					runLength--;
				}

				if (!ENSURE_CAPACITY(pbDest, pbDestEnd, runLength))
					return FALSE;

				UNROLL(runLength, {
					DESTREADPIXEL(temp, pbDest - rowDelta);
					DESTWRITEPIXEL(pbDest, temp);
					DESTNEXTPIXEL(pbDest);
				});
			}

			/* A follow-on background run order will need a foreground pel inserted. */
			fInsertFgPel = TRUE;
			continue;
		}

		/* For any of the other run-types a follow-on background run
		    order does not need a foreground pel inserted. */
		fInsertFgPel = FALSE;

		switch (code)
		{
			/* Handle Foreground Run Orders. */
			case REGULAR_FG_RUN:
			case MEGA_MEGA_FG_RUN:
			case LITE_SET_FG_FG_RUN:
			case MEGA_MEGA_SET_FG_RUN:
				runLength = ExtractRunLength(code, pbSrc, pbEnd, &advance);
				pbSrc = pbSrc + advance;

				if (code == LITE_SET_FG_FG_RUN || code == MEGA_MEGA_SET_FG_RUN)
				{
					if (pbSrc >= pbEnd)
						return FALSE;
					SRCREADPIXEL(fgPel, pbSrc);
					SRCNEXTPIXEL(pbSrc);
				}

				if (!ENSURE_CAPACITY(pbDest, pbDestEnd, runLength))
					return FALSE;

				if (fFirstLine)
				{
					UNROLL(runLength, {
						DESTWRITEPIXEL(pbDest, fgPel);
						DESTNEXTPIXEL(pbDest);
					});
				}
				else
				{
					UNROLL(runLength, {
						DESTREADPIXEL(temp, pbDest - rowDelta);
						DESTWRITEPIXEL(pbDest, temp ^ fgPel);
						DESTNEXTPIXEL(pbDest);
					});
				}

				break;

			/* Handle Dithered Run Orders. */
			case LITE_DITHERED_RUN:
			case MEGA_MEGA_DITHERED_RUN:
				runLength = ExtractRunLength(code, pbSrc, pbEnd, &advance);
				pbSrc = pbSrc + advance;
				if (pbSrc >= pbEnd)
					return FALSE;
				SRCREADPIXEL(pixelA, pbSrc);
				SRCNEXTPIXEL(pbSrc);
				if (pbSrc >= pbEnd)
					return FALSE;
				SRCREADPIXEL(pixelB, pbSrc);
				SRCNEXTPIXEL(pbSrc);

				if (!ENSURE_CAPACITY(pbDest, pbDestEnd, runLength * 2))
					return FALSE;

				UNROLL(runLength, {
					DESTWRITEPIXEL(pbDest, pixelA);
					DESTNEXTPIXEL(pbDest);
					DESTWRITEPIXEL(pbDest, pixelB);
					DESTNEXTPIXEL(pbDest);
				});
				break;

			/* Handle Color Run Orders. */
			case REGULAR_COLOR_RUN:
			case MEGA_MEGA_COLOR_RUN:
				runLength = ExtractRunLength(code, pbSrc, pbEnd, &advance);
				pbSrc = pbSrc + advance;
				if (pbSrc >= pbEnd)
					return FALSE;
				SRCREADPIXEL(pixelA, pbSrc);
				SRCNEXTPIXEL(pbSrc);

				if (!ENSURE_CAPACITY(pbDest, pbDestEnd, runLength))
					return FALSE;

				UNROLL(runLength, {
					DESTWRITEPIXEL(pbDest, pixelA);
					DESTNEXTPIXEL(pbDest);
				});
				break;

			/* Handle Foreground/Background Image Orders. */
			case REGULAR_FGBG_IMAGE:
			case MEGA_MEGA_FGBG_IMAGE:
			case LITE_SET_FG_FGBG_IMAGE:
			case MEGA_MEGA_SET_FGBG_IMAGE:
				runLength = ExtractRunLength(code, pbSrc, pbEnd, &advance);
				pbSrc = pbSrc + advance;

				if (pbSrc >= pbEnd)
					return FALSE;
				if (code == LITE_SET_FG_FGBG_IMAGE || code == MEGA_MEGA_SET_FGBG_IMAGE)
				{
					SRCREADPIXEL(fgPel, pbSrc);
					SRCNEXTPIXEL(pbSrc);
				}

				if (fFirstLine)
				{
					while (runLength > 8)
					{
						bitmask = *pbSrc;
						pbSrc = pbSrc + 1;
						pbDest = WRITEFIRSTLINEFGBGIMAGE(pbDest, pbDestEnd, bitmask, fgPel, 8);

						if (!pbDest)
							return FALSE;

						runLength = runLength - 8;
					}
				}
				else
				{
					while (runLength > 8)
					{
						bitmask = *pbSrc;
						pbSrc = pbSrc + 1;
						pbDest = WRITEFGBGIMAGE(pbDest, pbDestEnd, rowDelta, bitmask, fgPel, 8);

						if (!pbDest)
							return FALSE;

						runLength = runLength - 8;
					}
				}

				if (runLength > 0)
				{
					bitmask = *pbSrc;
					pbSrc = pbSrc + 1;

					if (fFirstLine)
					{
						pbDest =
						    WRITEFIRSTLINEFGBGIMAGE(pbDest, pbDestEnd, bitmask, fgPel, runLength);
					}
					else
					{
						pbDest =
						    WRITEFGBGIMAGE(pbDest, pbDestEnd, rowDelta, bitmask, fgPel, runLength);
					}

					if (!pbDest)
						return FALSE;
				}

				break;

			/* Handle Color Image Orders. */
			case REGULAR_COLOR_IMAGE:
			case MEGA_MEGA_COLOR_IMAGE:
				runLength = ExtractRunLength(code, pbSrc, pbEnd, &advance);
				pbSrc = pbSrc + advance;
				if (!ENSURE_CAPACITY(pbDest, pbDestEnd, runLength))
					return FALSE;

				UNROLL(runLength, {
					if (pbSrc >= pbEnd)
						return FALSE;
					SRCREADPIXEL(temp, pbSrc);
					SRCNEXTPIXEL(pbSrc);
					DESTWRITEPIXEL(pbDest, temp);
					DESTNEXTPIXEL(pbDest);
				});
				break;

			/* Handle Special Order 1. */
			case SPECIAL_FGBG_1:
				pbSrc = pbSrc + 1;

				if (fFirstLine)
				{
					pbDest =
					    WRITEFIRSTLINEFGBGIMAGE(pbDest, pbDestEnd, g_MaskSpecialFgBg1, fgPel, 8);
				}
				else
				{
					pbDest =
					    WRITEFGBGIMAGE(pbDest, pbDestEnd, rowDelta, g_MaskSpecialFgBg1, fgPel, 8);
				}

				if (!pbDest)
					return FALSE;

				break;

			/* Handle Special Order 2. */
			case SPECIAL_FGBG_2:
				pbSrc = pbSrc + 1;

				if (fFirstLine)
				{
					pbDest =
					    WRITEFIRSTLINEFGBGIMAGE(pbDest, pbDestEnd, g_MaskSpecialFgBg2, fgPel, 8);
				}
				else
				{
					pbDest =
					    WRITEFGBGIMAGE(pbDest, pbDestEnd, rowDelta, g_MaskSpecialFgBg2, fgPel, 8);
				}

				if (!pbDest)
					return FALSE;

				break;

			/* Handle White Order. */
			case SPECIAL_WHITE:
				pbSrc = pbSrc + 1;

				if (!ENSURE_CAPACITY(pbDest, pbDestEnd, 1))
					return FALSE;

				DESTWRITEPIXEL(pbDest, WHITE_PIXEL);
				DESTNEXTPIXEL(pbDest);
				break;

			/* Handle Black Order. */
			case SPECIAL_BLACK:
				pbSrc = pbSrc + 1;

				if (!ENSURE_CAPACITY(pbDest, pbDestEnd, 1))
					return FALSE;

				DESTWRITEPIXEL(pbDest, BLACK_PIXEL);
				DESTNEXTPIXEL(pbDest);
				break;

			default:
				return FALSE;
		}
	}

	return TRUE;
}
