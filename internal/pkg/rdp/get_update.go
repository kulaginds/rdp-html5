package rdp

import (
	"errors"
	"fmt"
	"log"

	"github.com/kulaginds/rdp-html5/internal/pkg/rdp/fastpath"
	"github.com/kulaginds/rdp-html5/internal/pkg/rdp/pdu"
)

func (c *client) GetUpdate() (*fastpath.UpdatePDU, error) {
	protocol, err := receiveProtocol(c.buffReader)
	if err != nil {
		return nil, err
	}

	if protocol.IsX224() {
		err = c.getX224Update()
		switch {
		case err == nil: // pass
		case errors.Is(err, pdu.ErrDeactiateAll):
			return nil, err

		default:
			return nil, fmt.Errorf("get X.224 update: %w", err)
		}

		return c.GetUpdate()
	}

	return c.fastPath.Receive()
}

func (c *client) getX224Update() error {
	channelID, wire, err := c.mcsLayer.Receive()
	if err != nil {
		return err
	}

	if channelID == c.channelIDMap["rail"] {
		err = c.handleRail(wire)
		if err != nil {
			return err
		}

		return nil
	}

	var data pdu.Data
	if err = data.Deserialize(wire); err != nil {
		return err
	}

	if data.ShareDataHeader.PDUType2.IsErrorInfo() {
		log.Printf("received error info: %s\n", data.ErrorInfoPDUData.String())
	}

	return nil
}
