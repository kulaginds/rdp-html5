package tpkt

import "net"

type protocol struct {
	conn net.Conn
}

func New(conn net.Conn) *protocol {
	return &protocol{
		conn: conn,
	}
}
