package utf16

import (
	"bytes"
	"encoding/binary"
	"unicode/utf16"
)

func Encode(s string) []byte {
	buf := new(bytes.Buffer)

	for _, ch := range utf16.Encode([]rune(s)) {
		binary.Write(buf, binary.LittleEndian, ch)
	}

	return buf.Bytes()
}
