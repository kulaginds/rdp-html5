package mcs

type protocol struct {
	x224Conn x224Conn

	connected bool
	channels  map[string]uint16
	userId    uint16

	skipChannelJoin bool
}

const globalChannel uint16 = 1003

func New(x224Conn x224Conn) *protocol {
	return &protocol{
		x224Conn: x224Conn,

		channels: map[string]uint16{
			"global": globalChannel,
		},

		skipChannelJoin: false,
	}
}
