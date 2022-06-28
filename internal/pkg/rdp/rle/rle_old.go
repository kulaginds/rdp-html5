package rle

// BitmapDecompress16 decompress rgb 565 color to flipped vertically rgba 8888.
func BitmapDecompress16(
	pbSrcBuffer []byte,
	cbSrcBuffer int,
	pbDestBuffer []byte,
	rowDelta int,
) bool {
	d := newDecompressor16(pbSrcBuffer, cbSrcBuffer, pbDestBuffer, rowDelta)

	// todo: flipV
	// todo: to rgba

	return d.Decompress()
}
