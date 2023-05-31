package mcs

import "github.com/kulaginds/web-rdp-solution/internal/pkg/rdp/x224"

type Protocol struct {
	x224Conn *x224.Protocol
}

func New(x224Conn *x224.Protocol) *Protocol {
	return &Protocol{
		x224Conn: x224Conn,
	}
}
