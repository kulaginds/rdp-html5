package fastpath

import (
	"io"
)

type Protocol struct {
	conn io.ReadWriter

	updatePDUData []byte
}

func New(conn io.ReadWriter) *Protocol {
	return &Protocol{
		conn: conn,

		updatePDUData: make([]byte, 16*1024),
	}
}
