package tpkt

import (
	"io"
)

type Protocol struct {
	conn            io.ReadWriteCloser
	fastpathEnabled bool
}

func New(conn io.ReadWriteCloser) *Protocol {
	return &Protocol{
		conn: conn,
	}
}
