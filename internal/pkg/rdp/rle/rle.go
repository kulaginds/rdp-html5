package rle

func Decompress(pbSrcBuffer []byte, pbDestBuffer []byte, rowDelta int) bool {
	d := newDecompressor16(pbSrcBuffer, pbDestBuffer, rowDelta)

	return d.Decompress()
}
