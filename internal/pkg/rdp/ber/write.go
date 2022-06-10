package ber

import (
	"encoding/binary"
	"io"
)

func WriteBoolean(b bool, w io.Writer) {
	bb := uint8(0)
	if b {
		bb = uint8(0xff)
	}
	w.Write([]byte{0x01}) // tag boolean
	WriteLength(1, w)
	w.Write([]byte{bb})
}

func WriteInteger(n int, w io.Writer) {
	w.Write([]byte{0x02}) // tag integer
	if n <= 0xff {
		WriteLength(1, w)
		w.Write([]byte{uint8(n)})
	} else if n <= 0xffff {
		WriteLength(2, w)
		binary.Write(w, binary.BigEndian, uint16(n))
	} else {
		WriteLength(4, w)
		binary.Write(w, binary.BigEndian, uint32(n))
	}
}

func WriteOctetString(str []byte, w io.Writer) {
	w.Write([]byte{0x04}) // tag octet string
	WriteLength(len(str), w)
	w.Write(str)
}

func WriteSequence(data []byte, w io.Writer) {
	w.Write([]byte{0x30}) // tag sequence
	WriteLength(len(data), w)
	w.Write(data)
}

func WriteApplicationTag(tag uint8, size int, w io.Writer) {
	if tag > 30 {
		w.Write([]byte{
			0x7f, // leading octet for tags with number greater than or equal to 31
			tag,
		})
		WriteLength(size, w)
	} else {
		w.Write([]byte{tag})
		WriteLength(size, w)
	}
}

func WriteLength(size int, w io.Writer) {
	if size > 0x7f {
		w.Write([]byte{0x82})
		binary.Write(w, binary.BigEndian, uint16(size))
	} else {
		w.Write([]byte{uint8(size)})
	}
}
