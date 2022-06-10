package per

import (
	"encoding/binary"
	"errors"
	"io"
)

func ReadChoice(r io.Reader) (uint8, error) {
	var choice uint8

	return choice >> 2, binary.Read(r, binary.BigEndian, &choice)
}

func ReadLength(r io.Reader) (int, error) {
	var (
		octet uint8
		size  int
		err   error
	)

	if err = binary.Read(r, binary.BigEndian, &octet); err != nil {
		return 0, err
	}

	if octet&0x80 != 0x80 {
		return int(octet), nil
	}

	octet &^= 0x80
	size = int(octet) << 8

	if err = binary.Read(r, binary.BigEndian, &octet); err != nil {
		return 0, err
	}

	size += int(octet)

	return size, nil
}

func ReadObjectIdentifier(oid [6]byte, r io.Reader) (bool, error) {
	size, err := ReadLength(r)
	if err != nil {
		return false, err
	}

	if size != 5 {
		return false, nil
	}

	var t12 uint8
	err = binary.Read(r, binary.BigEndian, &t12)
	if err != nil {
		return false, err
	}

	aOid := make([]byte, 6)
	aOid[0] = t12 >> 4
	aOid[1] = t12 & 0x0f

	for i := 2; i <= 5; i++ {
		err = binary.Read(r, binary.BigEndian, &aOid[i])
		if err != nil {
			return false, err
		}
	}

	for i := 0; i < len(aOid); i++ {
		if oid[i] != aOid[i] {
			return false, nil
		}
	}

	return true, nil
}

func ReadInteger16(minimum uint16, r io.Reader) (uint16, error) {
	var num uint16

	if err := binary.Read(r, binary.BigEndian, &num); err != nil {
		return 0, err
	}

	num += minimum

	return num, nil
}

func ReadInteger(r io.Reader) (int, error) {
	size, err := ReadLength(r)
	if err != nil {
		return 0, err
	}

	switch size {
	case 1:
		var num uint8

		return int(num), binary.Read(r, binary.BigEndian, &num)
	case 2:
		var num uint16

		return int(num), binary.Read(r, binary.BigEndian, &num)
	case 4:
		var num uint32

		return int(num), binary.Read(r, binary.BigEndian, &num)
	default:
		return 0, errors.New("bad integer length")
	}
}

func ReadEnumerates(r io.Reader) (uint8, error) {
	var num uint8

	return num, binary.Read(r, binary.BigEndian, &num)
}

func ReadNumberOfSet(r io.Reader) (uint8, error) {
	var num uint8

	return num, binary.Read(r, binary.BigEndian, &num)
}

func ReadOctetStream(octetStream []byte, minValue int, r io.Reader) (bool, error) {
	length, err := ReadLength(r)
	if err != nil {
		return false, err
	}

	size := length + minValue
	if size != len(octetStream) {
		return false, nil
	}

	var c uint8

	for i := 0; i < size; i++ {
		if err = binary.Read(r, binary.BigEndian, &c); err != nil {
			return false, err
		}

		if octetStream[i] != c {
			return false, nil
		}
	}

	return true, nil
}
