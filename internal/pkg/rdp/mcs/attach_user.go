package mcs

import (
	"fmt"
	"io"
	"log"

	"github.com/kulaginds/rdp-html5/internal/pkg/rdp/per"
)

type ClientAttachUserRequest struct{}

func (pdu *ClientAttachUserRequest) Serialize() []byte {
	// empty structure in T.125, section 7, page 18
	return nil
}

type ServerAttachUserConfirm struct {
	Result    uint8
	Initiator uint16
}

func (pdu *ServerAttachUserConfirm) Deserialize(wire io.Reader) error {
	var err error

	pdu.Result, err = per.ReadEnumerates(wire)
	if err != nil {
		return err
	}

	pdu.Initiator, err = per.ReadInteger16(1001, wire)
	if err != nil {
		return err
	}

	return nil
}

func (p *Protocol) AttachUser() (uint16, error) {
	req := DomainPDU{
		Application:             attachUserRequest,
		ClientAttachUserRequest: &ClientAttachUserRequest{},
	}

	log.Println("MCS: Attach User Request")

	if err := p.x224Conn.Send(req.Serialize()); err != nil {
		return 0, fmt.Errorf("client MCS attach user request: %w", err)
	}

	log.Println("MCS: Attach User Confirm")

	wire, err := p.x224Conn.Receive()
	if err != nil {
		return 0, err
	}

	var resp DomainPDU
	if err = resp.Deserialize(wire); err != nil {
		return 0, fmt.Errorf("server MCS attach user confirm reponse: %w", err)
	}

	if resp.ServerAttachUserConfirm.Result != RTSuccessful {
		return 0, fmt.Errorf("unsuccessful MCS attach user; result=%d", resp.ServerAttachUserConfirm.Result)
	}

	return resp.ServerAttachUserConfirm.Initiator, nil
}
