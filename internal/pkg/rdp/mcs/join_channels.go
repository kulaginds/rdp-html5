package mcs

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"

	"github.com/kulaginds/rdp-html5/internal/pkg/rdp/per"
)

type ClientChannelJoinRequest struct {
	Initiator uint16
	ChannelId uint16
}

func (pdu *ClientChannelJoinRequest) Serialize() []byte {
	buf := new(bytes.Buffer)

	per.WriteInteger16(pdu.Initiator, 1001, buf)
	per.WriteInteger16(pdu.ChannelId, 0, buf)

	return buf.Bytes()
}

type ServerChannelJoinConfirm struct {
	Result    uint8
	Initiator uint16
	Requested uint16
	ChannelId uint16
}

func (pdu *ServerChannelJoinConfirm) Deserialize(wire io.Reader) error {
	var err error

	pdu.Result, err = per.ReadEnumerates(wire)
	if err != nil {
		return err
	}

	pdu.Initiator, err = per.ReadInteger16(1001, wire)
	if err != nil {
		return err
	}

	pdu.Requested, err = per.ReadInteger16(0, wire)
	if err != nil {
		return err
	}

	// optional
	pdu.ChannelId, err = per.ReadInteger16(0, wire)
	switch {
	case errors.Is(err, nil), // pass
		errors.Is(err, io.EOF):
	default:
		return err
	}

	return nil
}

func (p *Protocol) JoinChannels(userID uint16, channelIDMap map[string]uint16) error {
	if len(channelIDMap) == 0 {
		return nil
	}

	for channelName, channelID := range channelIDMap {
		req := DomainPDU{
			Application: channelJoinRequest,
			ClientChannelJoinRequest: &ClientChannelJoinRequest{
				Initiator: userID,
				ChannelId: channelID,
			},
		}

		log.Printf("MCS: Channel Join Request: %s\n", channelName)

		if err := p.x224Conn.Send(req.Serialize()); err != nil {
			return fmt.Errorf("client MCS channel join request for channel=%s: %w", channelName, err)
		}

		log.Printf("MCS: Channel Join Confirm: %s\n", channelName)

		wire, err := p.x224Conn.Receive()
		if err != nil {
			return err
		}

		var resp DomainPDU
		if err = resp.Deserialize(wire); err != nil {
			return fmt.Errorf("server MCS channel join confirm reponse: %w", err)
		}

		if resp.ServerChannelJoinConfirm.Result != RTSuccessful {
			return fmt.Errorf(
				"unsuccessful MCS channel join confirm for channel=%s; result=%d",
				channelName,
				resp.ServerChannelJoinConfirm.Result,
			)
		}
	}

	return nil
}
