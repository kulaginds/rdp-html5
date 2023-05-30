package x224

import "github.com/kulaginds/web-rdp-solution/internal/pkg/rdp/tpkt"

type Protocol struct {
	tpktConn *tpkt.Protocol
}

func New(tpktConn *tpkt.Protocol) *Protocol {
	return &Protocol{
		tpktConn: tpktConn,
	}
}
