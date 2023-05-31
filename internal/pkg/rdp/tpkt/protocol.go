package tpkt

import (
	"io"
)

type Protocol struct {
	conn io.ReadWriteCloser
}

func New(conn io.ReadWriteCloser) *Protocol {
	return &Protocol{
		conn: conn,
	}
}
