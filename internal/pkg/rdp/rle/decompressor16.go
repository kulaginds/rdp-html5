package rle

type decompressor16 struct {
	// source bytes of compressed picture
	pbSrcBuffer []byte

	// length of pbSrcBuffer
	cbSrcBuffer int

	// destination bytes of decompressed rgb picture
	pbDestBuffer []byte

	// rowDelta one pixels row (width * 2 byte per pixel)
	rowDelta int

	// input offset
	i int

	// output offset
	o int
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
	fgPel := WhitePixel
	fInsertFgPel := false
	fFirstLine := true

	var (
		bitmask        byte
		pixelA, pixelB Pixel

		runLength uint16
		code      Code
	)

	for d.i < d.cbSrcBuffer {
		//
		// Watch out for the end of the first scanline.
		//
		if fFirstLine {
			if d.i >= d.rowDelta {
				fFirstLine = false
				fInsertFgPel = false
			}
		}

		//
		// Extract the compression order code ID from the compression
		// order header.
		//
		code = d.extractCodeId()

		//
		// Handle Background Run Orders.
		//
		if code == RegularBackgroundRun || code == MegaMegaBackgroundRun {
			runLength = d.extractRunLength(code)
			d.advanceOverOrderHeader(code)

			if fFirstLine {
				if fInsertFgPel {
					d.writePixel(fgPel)
					d.nextDestPixel()
					runLength = runLength - 1
				}

				for runLength > 0 {
					d.writePixel(BlackPixel)
					d.nextDestPixel()
					runLength = runLength - 1
				}
			} else {
				if fInsertFgPel {
					WritePixel(
					pbDest,
					ReadPixel(pbDest - rowDelta) XOR fgPel
					);
					pbDest = NextPixel(pbDest);
					runLength = runLength - 1;
				}

				for runLength > 0 {
					WritePixel(pbDest, ReadPixel(pbDest - rowDelta));
					pbDest = NextPixel(pbDest);
					runLength = runLength - 1;
				}
			}

			//
			// A follow-on background run order will need a
			// foreground pel inserted.
			//
			fInsertFgPel = true
			continue;
		}
	}

	return false
}

func (d *decompressor16) extractCodeId() Code {
	header := d.pbSrcBuffer[d.i]

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

func (d *decompressor16) extractRunLength(code Code) uint16 {
	// todo: implement
	return 0
}

func (d *decompressor16) advanceOverOrderHeader(code Code) {
	// todo: implement
}

func (d *decompressor16) writePixel(pixel Pixel) {
	// todo: implement
}

func (d *decompressor16) nextDestPixel() {
	// todo: implement
}
