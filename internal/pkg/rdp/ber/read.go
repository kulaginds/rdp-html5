package ber

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"

	"github.com/kulaginds/rdp-html5/internal/pkg/rdp/asn1"
)

func ReadApplicationTag(r io.Reader) (uint8, error) {
	var (
		identifier uint8
		tag        uint8
		err        error
	)

	err = binary.Read(r, binary.BigEndian, &identifier)
	if err != nil {
		return 0, err
	}

	if identifier != (asn1.ClassApplication|asn1.PCConstruct)|asn1.TagMask {
		return 0, errors.New("ReadApplicationTag invalid data")
	}

	err = binary.Read(r, binary.BigEndian, &tag)
	if err != nil {
		return 0, err
	}

	return tag, nil
}

func ReadLength(r io.Reader) (uint16, error) {
	var (
		ret  uint16
		size uint8
		err  error
	)

	err = binary.Read(r, binary.BigEndian, &size)
	if err != nil {
		return 0, err
	}

	if size&0x80 > 0 {
		size = size &^ 0x80

		if size == 1 {
			err = binary.Read(r, binary.BigEndian, &size)
			if err != nil {
				return 0, err
			}

			ret = uint16(size)
		} else if size == 2 {
			err = binary.Read(r, binary.BigEndian, &ret)
			if err != nil {
				return 0, err
			}
		} else {
			return 0, errors.New("BER length may be 1 or 2")
		}
	} else {
		ret = uint16(size)
	}

	return ret, nil
}

func berPC(pc bool) uint8 {
	if pc {
		return asn1.PCConstruct
	}
	return asn1.PCPrimitive
}

func ReadUniversalTag(tag uint8, pc bool, r io.Reader) (bool, error) {
	var bb uint8

	err := binary.Read(r, binary.BigEndian, &bb)
	if err != nil {
		return false, err
	}

	return bb == (asn1.ClassUniversal|berPC(pc))|(asn1.TagMask&tag), nil
}

func ReadEnumerated(r io.Reader) (uint8, error) {
	universalTag, err := ReadUniversalTag(asn1.TagEnumerated, false, r)
	if err != nil {
		return 0, err
	}

	if !universalTag {
		return 0, errors.New("invalid ber tag")
	}

	length, err := ReadLength(r)
	if err != nil {
		return 0, err
	}

	if length != 1 {
		return 0, fmt.Errorf("enumerate size is wrong, get %v, expect 1", length)
	}

	var enumerated uint8

	return enumerated, binary.Read(r, binary.BigEndian, &enumerated)
}

func ReadInteger(r io.Reader) (int, error) {
	universalTag, err := ReadUniversalTag(asn1.TagInteger, false, r)
	if err != nil {
		return 0, err
	}

	if !universalTag {
		return 0, errors.New("Bad integer tag")
	}

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
	case 3:
		var (
			int1 uint8
			int2 uint16
		)

		err = binary.Read(r, binary.BigEndian, &int1)
		if err != nil {
			return 0, err
		}

		err = binary.Read(r, binary.BigEndian, &int2)
		if err != nil {
			return 0, err
		}

		return int(int1)<<0x10 + int(int2), nil
	case 4:
		var num uint32

		return int(num), binary.Read(r, binary.BigEndian, &num)
	default:
		return 0, errors.New("wrong size")
	}
}
