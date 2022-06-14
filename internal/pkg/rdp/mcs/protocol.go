package mcs

type protocol struct {
	x224Conn x224Conn

	connected bool
	channels  map[string]uint16
	userId    uint16

	skipChannelJoin bool
}

const (
	ServerChannelID uint16 = 1002
)

func New(x224Conn x224Conn) *protocol {
	return &protocol{
		x224Conn: x224Conn,

		channels: map[string]uint16{
			"global": 0,
		},

		skipChannelJoin: false,
	}
}
