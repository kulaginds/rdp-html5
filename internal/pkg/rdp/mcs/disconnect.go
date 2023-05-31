package mcs

import (
	"fmt"
)

type ClientDisconnectUltimatumRequest struct{}

func (pdu *ClientDisconnectUltimatumRequest) Serialize() []byte {
	// per aligned RNUserRequested
	return []byte{
		0x21,
		0x80,
	}
}

func (p *Protocol) Disconnect() error {
	req := ClientDisconnectUltimatumRequest{}

	if err := p.x224Conn.Send(req.Serialize()); err != nil {
		return fmt.Errorf("client MCS disconnect ultimatum: %w", err)
	}

	return nil
}
