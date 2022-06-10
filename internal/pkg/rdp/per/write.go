package per

import (
	"encoding/binary"
	"io"
)

func WriteChoice(choice uint8, w io.Writer) {
	_, _ = w.Write([]byte{
		choice << 2,
	})
}

func WriteObjectIdentifier(oid [6]byte, w io.Writer) {
	WriteLength(5, w)

	_, _ = w.Write([]byte{
		(oid[0] << 4) & (oid[1] & 0x0f),
		oid[2],
		oid[3],
		oid[4],
		oid[5],
	})
}

func WriteLength(value uint16, w io.Writer) {
	if value > 0x7f {
		_ = binary.Write(w, binary.BigEndian, value|0x8000)
		return
	}

	_, _ = w.Write([]byte{uint8(value)})
}

func WriteSelection(selection uint8, w io.Writer) {
	_, _ = w.Write([]byte{
		selection,
	})
}

func WriteNumericString(nStr string, minValue int, w io.Writer) {
	length := len(nStr)
	mLength := minValue

	if length-minValue >= 0 {
		mLength = length - minValue
	}

	result := make([]byte, 0, mLength)

	for i := 0; i < length; i += 2 {
		c1 := nStr[i]
		c2 := byte(0x30)

		if i+1 < length {
			c2 = nStr[i+1]
		}

		c1 = (c1 - 0x30) % 10
		c2 = (c2 - 0x30) % 10

		result = append(result, (c1<<4)|c2)
	}

	WriteLength(uint16(mLength), w)
	_, _ = w.Write(result)
}

func WritePadding(length int, w io.Writer) {
	_, _ = w.Write(make([]byte, length))
}

func WriteNumberOfSet(numberOfSet uint8, w io.Writer) {
	_, _ = w.Write([]byte{numberOfSet})
}

func WriteOctetStream(oStr string, minValue int, w io.Writer) {
	length := len(oStr)
	mLength := minValue

	if length-minValue >= 0 {
		mLength = length - minValue
	}

	result := make([]byte, 0, mLength)
	for i := 0; i < length; i++ {
		result = append(result, oStr[i])
	}

	WriteLength(uint16(mLength), w)
	_, _ = w.Write(result)
}

func WriteInteger(value int, w io.Writer) {
	if value <= 0xff {
		WriteLength(1, w)
		w.Write([]byte{uint8(value)})

		return
	}

	if value < 0xffff {
		WriteLength(2, w)
		binary.Write(w, binary.BigEndian, uint16(value))

		return
	}

	WriteLength(4, w)
	binary.Write(w, binary.BigEndian, uint32(value))
}

func WriteInteger16(value, minimum uint16, w io.Writer) {
	value -= minimum

	binary.Write(w, binary.BigEndian, value)
}
