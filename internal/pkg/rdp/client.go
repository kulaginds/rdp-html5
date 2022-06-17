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

type client struct {
	conn      net.Conn
	tpktLayer tpktLayer
	x224Layer x224Layer
	mcsLayer  mcsLayer
	fastPath  fastPath

	domain   string
	username string
	password string

	desktopWidth, desktopHeight uint16

	selectedProtocol       x224.RDPNegotiationProtocol
	serverNegotiationFlags x224.RDPNegotiationResponseFlag
	channels               []string
	shareID                uint32
}

const (
	tcpConnectionTimeout = 15 * time.Second
)

func NewClient(
	hostname, username, password string,
	desktopWidth, desktopHeight uint16,
) (*client, error) {
	c := client{
		domain:   "",
		username: username,
		password: password,

		desktopWidth:  desktopWidth,
		desktopHeight: desktopHeight,

		selectedProtocol: x224.RDPNegotiationProtocolSSL,
	}

	var err error

	c.conn, err = net.DialTimeout("tcp", hostname, tcpConnectionTimeout)
	if err != nil {
		return nil, fmt.Errorf("tcp connect: %w", err)
	}

	c.tpktLayer = tpkt.New(&c)
	c.x224Layer = x224.New(c.tpktLayer, c.selectedProtocol)
	c.mcsLayer = mcs.New(c.x224Layer)
	c.fastPath = fastpath.New(&c)

	return &c, nil
}
