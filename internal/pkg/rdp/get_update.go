package rdp

import (
	"github.com/kulaginds/web-rdp-solution/internal/pkg/rdp/pdu"
	"io"
	"log"

	"github.com/kulaginds/web-rdp-solution/internal/pkg/rdp/fastpath"
)

func (c *client) GetUpdate() (*fastpath.UpdatePDU, error) {
	protocol, err := c.tpktLayer.ReceiveProtocol()
	if err != nil {
		return nil, err
	}

	if protocol.IsX224() {
		var (
			channelID uint16
			wire      io.Reader
		)

		channelID, wire, err = c.mcsLayer.Receive()
		if err != nil {
			return nil, err
		}

		if channelID == c.channelIDMap["rail"] {
			err = c.handleRail(wire)
			if err != nil {
				return nil, err
			}

			return c.GetUpdate()
		}

		var data pdu.Data
		if err = data.Deserialize(wire); err != nil {
			return nil, err
		}

		if data.ShareDataHeader.PDUType2.IsErrorInfo() {
			log.Printf("received error info: %d\n", data.ErrorInfoPDUData.String())
		}

		return c.GetUpdate()
	}

	return c.fastPath.Receive(uint8(protocol))
}
