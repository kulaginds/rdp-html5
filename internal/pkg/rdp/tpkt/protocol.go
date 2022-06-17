package tpkt

import (
	"io"
)

type protocol struct {
	conn            io.ReadWriteCloser
	fastpathEnabled bool
}

func New(conn io.ReadWriteCloser) *protocol {
	return &protocol{
		conn: conn,
	}
}
