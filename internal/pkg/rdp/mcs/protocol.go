package mcs

import "github.com/kulaginds/web-rdp-solution/internal/pkg/rdp/x224"

type Protocol struct {
	x224Conn *x224.Protocol

	connected bool
	channels  map[string]uint16
	userId    uint16

	skipChannelJoin bool
}

const (
	ServerChannelID uint16 = 1002
)

func New(x224Conn *x224.Protocol) *Protocol {
	return &Protocol{
		x224Conn: x224Conn,

		channels: map[string]uint16{
			"global": 0,
		},

		skipChannelJoin: false,
	}
}
