package mcs

import (
	"io"

	"github.com/kulaginds/web-rdp-solution/internal/pkg/rdp/per"
)

type ServerSendDataIndication struct {
	Initiator uint16
	ChannelId uint16
}

func (d *ServerSendDataIndication) Deserialize(wire io.Reader) error {
	var err error

	d.Initiator, err = per.ReadInteger16(1001, wire)
	if err != nil {
		return err
	}

	d.ChannelId, err = per.ReadInteger16(0, wire)
	if err != nil {
		return err
	}

	_, err = per.ReadEnumerates(wire)
	if err != nil {
		return err
	}

	_, err = per.ReadLength(wire)
	if err != nil {
		return err
	}

	return nil
}

// Receive returns channelName, reader or error
func (p *Protocol) Receive() (string, io.Reader, error) {
	if !p.connected {
		return "", nil, ErrNotConnected
	}

	wire, err := p.x224Conn.Receive()
	if err != nil {
		return "", nil, err
	}

	var resp DomainPDU
	if err = resp.Deserialize(wire); err != nil {
		return "", nil, err
	}

	if resp.Application != SendDataIndication {
		return "", nil, ErrUnknownDomainApplication
	}

	for channelName, channelID := range p.channels {
		if channelID == resp.ServerSendDataIndication.ChannelId {
			return channelName, wire, nil
		}
	}

	return "", nil, ErrUnknownChannel
}
