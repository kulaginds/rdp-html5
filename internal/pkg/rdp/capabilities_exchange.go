package rdp

import "github.com/kulaginds/rdp-html5/internal/pkg/rdp/pdu"

func (c *client) capabilitiesExchange() error {
	_, wire, err := c.mcsLayer.Receive()
	if err != nil {
		return err
	}

	var resp pdu.ServerDemandActive
	if err = resp.Deserialize(wire); err != nil {
		return err
	}

	c.shareID = resp.ShareID
	c.serverCapabilitySets = resp.CapabilitySets

	req := pdu.NewClientConfirmActive(resp.ShareID, c.userID, c.desktopWidth, c.desktopHeight, c.remoteApp != nil)

	return c.mcsLayer.Send(c.userID, c.channelIDMap["global"], req.Serialize())
}
