package fastpath

import (
	"io"
)

type impl struct {
	conn io.ReadWriter
}

func New(conn io.ReadWriter) *impl {
	return &impl{
		conn: conn,
	}
}
