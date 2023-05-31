package mcs

import (
	"bytes"
	"fmt"
	"log"

	"github.com/kulaginds/web-rdp-solution/internal/pkg/rdp/per"
)

type ClientErectDomainRequest struct{}

func (pdu *ClientErectDomainRequest) Serialize() []byte {
	buf := new(bytes.Buffer)

	per.WriteInteger(0, buf)
	per.WriteInteger(0, buf)

	return buf.Bytes()
}

func (p *Protocol) ErectDomain() error {
	req := DomainPDU{
		Application:              erectDomainRequest,
		ClientErectDomainRequest: &ClientErectDomainRequest{},
	}

	log.Println("MCS: Erect Domain")

	if err := p.x224Conn.Send(req.Serialize()); err != nil {
		return fmt.Errorf("client MCS erect domain request: %w", err)
	}

	return nil
}
