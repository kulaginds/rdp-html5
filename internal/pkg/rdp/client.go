package rdp

import (
	"fmt"
	"net"
	"time"

	"github.com/kulaginds/web-rdp-solution/internal/pkg/rdp/fastpath"
	"github.com/kulaginds/web-rdp-solution/internal/pkg/rdp/mcs"
	"github.com/kulaginds/web-rdp-solution/internal/pkg/rdp/tpkt"
	"github.com/kulaginds/web-rdp-solution/internal/pkg/rdp/x224"
)

type RemoteApp struct {
	App        string
	WorkingDir string
	Args       string
}

type client struct {
	conn      net.Conn
	tpktLayer *tpkt.Protocol
	x224Layer *x224.Protocol
	mcsLayer  *mcs.Protocol
	fastPath  *fastpath.Protocol

	domain   string
	username string
	password string

	desktopWidth, desktopHeight uint16

	remoteApp *RemoteApp
	railState RailState

	selectedProtocol       NegotiationProtocol
	serverNegotiationFlags NegotiationResponseFlag
	channels               []string
	shareID                uint32
}

const (
	tcpConnectionTimeout = 15 * time.Second
)

func NewClient(
	hostname, username, password string,
	desktopWidth, desktopHeight int,
) (*client, error) {
	c := client{
		domain:   "",
		username: username,
		password: password,

		desktopWidth:  uint16(desktopWidth),
		desktopHeight: uint16(desktopHeight),

		selectedProtocol: NegotiationProtocolSSL,
	}

	var err error

	c.conn, err = net.DialTimeout("tcp", hostname, tcpConnectionTimeout)
	if err != nil {
		return nil, fmt.Errorf("tcp connect: %w", err)
	}

	c.tpktLayer = tpkt.New(&c)
	c.x224Layer = x224.New(c.tpktLayer)
	c.mcsLayer = mcs.New(c.x224Layer)
	c.fastPath = fastpath.New(&c)

	return &c, nil
}
