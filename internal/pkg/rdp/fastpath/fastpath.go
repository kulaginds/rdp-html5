package fastpath

import (
	"io"
)

type impl struct {
	conn io.ReadWriter

	updatePDUData []byte
}

func New(conn io.ReadWriter) *impl {
	return &impl{
		conn: conn,

		updatePDUData: make([]byte, 16*1024),
	}
}
