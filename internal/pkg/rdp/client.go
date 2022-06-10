package rdp

import (
	"fmt"
	"net"
	"time"

	"github.com/kulaginds/web-rdp-solution/internal/pkg/rdp/mcs"
	"github.com/kulaginds/web-rdp-solution/internal/pkg/rdp/tpkt"
	"github.com/kulaginds/web-rdp-solution/internal/pkg/rdp/x224"
)

type client struct {
	x224Layer x224Layer
	mcsLayer  mcsLayer

	domain   string
	username string
	password string

	desktopWidth, desktopHeight uint16

	selectedProtocol       x224.RDPNegotiationProtocol
	serverNegotiationFlags x224.RDPNegotiationResponseFlag
	channels               []string
}

const (
	tcpConnectionTimeout = 15 * time.Second
)

func NewClient(
	hostname, domain, username, password string,
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

	conn, err := net.DialTimeout("tcp", hostname, tcpConnectionTimeout)
	if err != nil {
		return nil, fmt.Errorf("tcp connect: %w", err)
	}

	c.x224Layer = x224.New(tpkt.New(conn), c.selectedProtocol)
	c.mcsLayer = mcs.New(c.x224Layer)

	return &c, nil
}
