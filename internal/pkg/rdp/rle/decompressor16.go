package rle

type decompressor16 struct {
	pbSrcBuffer  []byte
	cbSrcBuffer  int
	pbDestBuffer []byte
	rowDelta     int
}

func newDecompressor16(
	pbSrcBuffer []byte,
	cbSrcBuffer int,
	pbDestBuffer []byte,
	rowDelta int,
) *decompressor16 {
	return &decompressor16{
		pbSrcBuffer:  pbSrcBuffer,
		cbSrcBuffer:  cbSrcBuffer,
		pbDestBuffer: pbDestBuffer,
		rowDelta:     rowDelta,
	}
}

func (d *decompressor16) Decompress() bool {
	codeID := extractCodeId(d.)

	return false
}

func extractCodeId(header uint8) Code {
	if (header & 0xC0) != 0xC0 { // don't have high 2 bytes, which have lite order
		/* REGULAR orders
		 * (000x xxxx, 001x xxxx, 010x xxxx, 011x xxxx, 100x xxxx)
		 */
		return Code(header >> 5)
	}

	if (header & 0xF0) == 0xF0 {
		/* MEGA and SPECIAL orders (0xF*) */
		return Code(header)
	}

	/* LITE orders
	 * 1100 xxxx, 1101 xxxx, 1110 xxxx)
	 */
	return Code(header >> 4)
}
