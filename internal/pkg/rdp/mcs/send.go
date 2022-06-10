package mcs

import (
	"bytes"
	"fmt"

	"github.com/kulaginds/web-rdp-solution/internal/pkg/rdp/per"
)

type ClientSendDataRequest struct {
	Initiator uint16
	ChannelId uint16
	Data      []byte
}

func (d *ClientSendDataRequest) Serialize() []byte {
	buf := &bytes.Buffer{}

	per.WriteInteger16(d.Initiator, 1001, buf)
	per.WriteInteger16(d.ChannelId, 0, buf)
	buf.WriteByte(0x70) // magic word
	per.WriteLength(uint16(len(d.Data)), buf)

	buf.Write(d.Data)

	return buf.Bytes()
}

func (p *protocol) Send(channelName string, pduData []byte) error {
	if !p.connected {
		return ErrNotConnected
	}

	channelID, ok := p.channels[channelName]
	if !ok {
		return ErrChannelNotFound
	}

	req := DomainPDU{
		Application: sendDataRequest,
		ClientSendDataRequest: &ClientSendDataRequest{
			Initiator: p.userId,
			ChannelId: channelID,
			Data:      pduData,
		},
	}

	if err := p.x224Conn.Send(req.Serialize()); err != nil {
		return fmt.Errorf("client MCS send data request for channel=%s: %w", channelName, err)
	}

	return nil
}
