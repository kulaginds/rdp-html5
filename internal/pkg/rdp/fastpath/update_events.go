package fastpath

import (
	"encoding/binary"
	"io"
)

type PaletteEntry struct {
	Red   uint8
	Green uint8
	Blue  uint8
}

func (e *PaletteEntry) Deserialize(wire io.Reader) error {
	var err error

	err = binary.Read(wire, binary.LittleEndian, &e.Red)
	if err != nil {
		return err
	}

	err = binary.Read(wire, binary.LittleEndian, &e.Green)
	if err != nil {
		return err
	}

	err = binary.Read(wire, binary.LittleEndian, &e.Green)
	if err != nil {
		return err
	}

	return nil
}

type paletteUpdateData struct {
	PaletteEntries []PaletteEntry
}

func (d *paletteUpdateData) Deserialize(wire io.Reader) error {
	var err error

	var updateType uint16
	err = binary.Read(wire, binary.LittleEndian, &updateType)
	if err != nil {
		return err
	}

	var padding uint16
	err = binary.Read(wire, binary.LittleEndian, &padding)
	if err != nil {
		return err
	}

	var numberColors uint16
	err = binary.Read(wire, binary.LittleEndian, &numberColors)
	if err != nil {
		return err
	}

	d.PaletteEntries = make([]PaletteEntry, numberColors)

	for i := 0; i < len(d.PaletteEntries); i++ {
		err = d.PaletteEntries[i].Deserialize(wire)
		if err != nil {
			return err
		}
	}

	return nil
}

type CompressedDataHeader struct {
	CbCompMainBodySize uint16
	CbScanWidth        uint16
	CbUncompressedSize uint16
}

func (h *CompressedDataHeader) Deserialize(wire io.Reader) error {
	var err error

	var cbCompFirstRowSize uint16
	err = binary.Read(wire, binary.LittleEndian, &cbCompFirstRowSize)
	if err != nil {
		return err
	}

	err = binary.Read(wire, binary.LittleEndian, &h.CbCompMainBodySize)
	if err != nil {
		return err
	}

	err = binary.Read(wire, binary.LittleEndian, &h.CbScanWidth)
	if err != nil {
		return err
	}

	err = binary.Read(wire, binary.LittleEndian, &h.CbUncompressedSize)
	if err != nil {
		return err
	}

	return nil
}

type BitmapDataFlag uint16

const (
	// BitmapDataFlagCompression BITMAP_COMPRESSION
	BitmapDataFlagCompression BitmapDataFlag = 0x0001

	// BitmapDataFlagNoHDR NO_BITMAP_COMPRESSION_HDR
	BitmapDataFlagNoHDR BitmapDataFlag = 0x0400
)

type BitmapData struct {
	DestLeft         uint16
	DestTop          uint16
	DestRight        uint16
	DestBottom       uint16
	Width            uint16
	Height           uint16
	BitsPerPixel     uint16
	Flags            BitmapDataFlag
	BitmapLength     uint16
	BitmapComprHdr   *CompressedDataHeader
	BitmapDataStream []byte
}

func (d *BitmapData) Deserialize(wire io.Reader) error {
	var err error

	err = binary.Read(wire, binary.LittleEndian, &d.DestLeft)
	if err != nil {
		return err
	}

	err = binary.Read(wire, binary.LittleEndian, &d.DestTop)
	if err != nil {
		return err
	}

	err = binary.Read(wire, binary.LittleEndian, &d.DestRight)
	if err != nil {
		return err
	}

	err = binary.Read(wire, binary.LittleEndian, &d.DestBottom)
	if err != nil {
		return err
	}

	err = binary.Read(wire, binary.LittleEndian, &d.Width)
	if err != nil {
		return err
	}

	err = binary.Read(wire, binary.LittleEndian, &d.Height)
	if err != nil {
		return err
	}

	err = binary.Read(wire, binary.LittleEndian, &d.BitsPerPixel)
	if err != nil {
		return err
	}

	err = binary.Read(wire, binary.LittleEndian, &d.Flags)
	if err != nil {
		return err
	}

	err = binary.Read(wire, binary.LittleEndian, &d.BitmapLength)
	if err != nil {
		return err
	}

	if d.Flags&BitmapDataFlagNoHDR != BitmapDataFlagNoHDR && d.Flags&BitmapDataFlagCompression == BitmapDataFlagCompression {
		err = d.BitmapComprHdr.Deserialize(wire)
		if err != nil {
			return err
		}

		d.BitmapLength -= 8
	}

	d.BitmapDataStream = make([]byte, d.BitmapLength)

	_, err = wire.Read(d.BitmapDataStream)
	if err != nil {
		return err
	}

	return nil
}

type bitmapUpdateData struct {
	Rectangles []BitmapData
}

func (d *bitmapUpdateData) Deserialize(wire io.Reader) error {
	var err error

	var updateType uint16
	err = binary.Read(wire, binary.LittleEndian, &updateType)
	if err != nil {
		return err
	}

	var numberRectangles uint16
	err = binary.Read(wire, binary.LittleEndian, &numberRectangles)
	if err != nil {
		return err
	}

	d.Rectangles = make([]BitmapData, numberRectangles)

	for i := range d.Rectangles {
		err = d.Rectangles[i].Deserialize(wire)
		if err != nil {
			return err
		}
	}

	return nil
}

type pointerPositionUpdateData struct {
	xPos uint16
	yPos uint16
}

func (d *pointerPositionUpdateData) Deserialize(wire io.Reader) error {
	var err error

	err = binary.Read(wire, binary.LittleEndian, &d.xPos)
	if err != nil {
		return err
	}

	err = binary.Read(wire, binary.LittleEndian, &d.yPos)
	if err != nil {
		return err
	}

	return nil
}

type colorPointerUpdateData struct {
	cacheIndex    uint16
	xPos          uint16
	yPos          uint16
	width         uint16
	height        uint16
	lengthAndMask uint16
	lengthXorMask uint16
	xorMaskData   []byte
	andMaskData   []byte
}

func (d *colorPointerUpdateData) Deserialize(wire io.Reader) error {
	var err error

	err = binary.Read(wire, binary.LittleEndian, &d.cacheIndex)
	if err != nil {
		return err
	}

	err = binary.Read(wire, binary.LittleEndian, &d.xPos)
	if err != nil {
		return err
	}

	err = binary.Read(wire, binary.LittleEndian, &d.yPos)
	if err != nil {
		return err
	}

	err = binary.Read(wire, binary.LittleEndian, &d.width)
	if err != nil {
		return err
	}

	err = binary.Read(wire, binary.LittleEndian, &d.height)
	if err != nil {
		return err
	}

	err = binary.Read(wire, binary.LittleEndian, &d.lengthAndMask)
	if err != nil {
		return err
	}

	err = binary.Read(wire, binary.LittleEndian, &d.lengthXorMask)
	if err != nil {
		return err
	}

	if d.lengthXorMask > 0 {
		d.xorMaskData = make([]byte, d.lengthXorMask)
		_, err = wire.Read(d.xorMaskData)
		if err != nil {
			return err
		}
	}

	if d.lengthAndMask > 0 {
		d.andMaskData = make([]byte, d.lengthAndMask)
		_, err = wire.Read(d.andMaskData)
		if err != nil {
			return err
		}
	}

	var padding uint8

	return binary.Read(wire, binary.LittleEndian, &padding)
}
