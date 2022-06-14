package tpkt

import (
	"io"
)

type protocol struct {
	conn io.ReadWriteCloser
}

func New(conn io.ReadWriteCloser) *protocol {
	return &protocol{
		conn: conn,
	}
}
