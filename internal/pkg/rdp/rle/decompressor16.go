package rle

type decompressor16 struct {
	// source bytes of compressed picture
	pbSrcBuffer []byte

	// destination bytes of decompressed rgb picture
	pbDestBuffer []byte

	// rowDelta one pixels row (width * 2 byte per pixel)
	rowDelta int
}

func newDecompressor16(pbSrcBuffer []byte, pbDestBuffer []byte, rowDelta int) *decompressor16 {
	return &decompressor16{
		pbSrcBuffer:  pbSrcBuffer,
		pbDestBuffer: pbDestBuffer,
		rowDelta:     rowDelta,
	}
}

func (d *decompressor16) Decompress() bool {
	// TODO: implement

	return false
}
