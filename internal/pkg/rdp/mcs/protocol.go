package mcs

import "github.com/kulaginds/rdp-html5/internal/pkg/rdp/x224"

type Protocol struct {
	x224Conn *x224.Protocol
}

func New(x224Conn *x224.Protocol) *Protocol {
	return &Protocol{
		x224Conn: x224Conn,
	}
}
