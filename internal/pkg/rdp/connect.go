package rdp

import (
	"fmt"
	"github.com/kulaginds/web-rdp-solution/internal/pkg/rdp/pdu"
	"io"
	"log"
)

func (c *client) Connect() error {
	var err error

	if err = c.connectionInitiation(); err != nil {
		return fmt.Errorf("connection initiation: %w", err)
	}

	if err = c.basicSettingsExchange(); err != nil {
		return fmt.Errorf("basic settings exchange: %w", err)
	}

	if err = c.channelConnection(); err != nil {
		return fmt.Errorf("channel connection: %w", err)
	}

	if err = c.secureSettingsExchange(); err != nil {
		return fmt.Errorf("secure settings exchange: %w", err)
	}

	if err = c.licensing(); err != nil {
		return fmt.Errorf("licensing: %w", err)
	}

	if err = c.capabilitiesExchange(); err != nil {
		return fmt.Errorf("capabilities exchange: %w", err)
	}

	if err = c.connectionFinalization(); err != nil {
		return fmt.Errorf("connection finalizatioin: %w", err)
	}

	return nil
}

func (c *client) connectionInitiation() error {
	var err error

	req := pdu.ClientConnectionRequest{
		NegotiationRequest: pdu.NegotiationRequest{
			RequestedProtocols: c.selectedProtocol,
		},
	}

	var (
		resp pdu.ServerConnectionConfirm
		wire io.Reader
	)

	if wire, err = c.x224Layer.Connect(req.Serialize()); err != nil {
		return err
	}

	if err = resp.Deserialize(wire); err != nil {
		return err
	}

	if resp.Type.IsFailure() {
		return fmt.Errorf("negotiation faliure: failureCode = %d", resp.FailureCode())
	}

	c.serverNegotiationFlags = resp.Flags

	if !resp.SelectedProtocol().IsSSL() {
		return ErrUnsupportedRequestedProtocol
	}

	log.Println("Server negotiation flags: " + c.serverNegotiationFlags.String())

	if c.selectedProtocol.IsSSL() {
		return c.StartTLS()
	}

	return nil
}

func (c *client) basicSettingsExchange() error {
	clientUserDataSet := pdu.NewClientUserDataSet(uint32(c.selectedProtocol), c.desktopWidth, c.desktopHeight, c.channels)

	wire, err := c.mcsLayer.Connect(clientUserDataSet.Serialize())
	if err != nil {
		return err
	}

	var serverUserData pdu.ServerUserData
	err = serverUserData.Deserialize(wire)
	if err != nil {
		return err
	}

	c.initChannels(serverUserData.ServerNetworkData)

	log.Println("MCS: Server Connect Response: earlyCapabilityFlags: ", serverUserData.ServerCoreData.EarlyCapabilityFlags)

	// RNS_UD_SC_SKIP_CHANNELJOIN_SUPPORTED = 0x00000008
	c.skipChannelJoin = serverUserData.ServerCoreData.EarlyCapabilityFlags&0x8 == 0x8

	return nil
}

func (c *client) initChannels(serverNetworkData *pdu.ServerNetworkData) {
	if c.channels == nil {
		c.channelIDMap = make(map[string]uint16, len(c.channels))
	}

	for i, channelName := range c.channels {
		c.channelIDMap[channelName] = serverNetworkData.ChannelIdArray[i]
	}

	c.channelIDMap["global"] = serverNetworkData.MCSChannelId
}

func (c *client) channelConnection() error {
	err := c.mcsLayer.ErectDomain()
	if err != nil {
		return err
	}

	c.userID, err = c.mcsLayer.AttachUser()
	if err != nil {
		return err
	}

	c.channelIDMap["user"] = c.userID

	if c.skipChannelJoin {
		return nil
	}

	err = c.mcsLayer.JoinChannels(c.userID, c.channelIDMap)
	if err != nil {
		return err
	}

	return nil
}

func (c *client) secureSettingsExchange() error {
	clientInfoPDU := pdu.NewClientInfo(c.domain, c.username, c.password)

	if c.remoteApp != nil {
		clientInfoPDU.InfoPacket.Flags |= pdu.InfoFlagRail
	}

	log.Println("RDP: Client Info")

	if err := c.mcsLayer.Send(c.userID, c.channelIDMap["global"], clientInfoPDU.Serialize()); err != nil {
		return fmt.Errorf("client info: %w", err)
	}

	return nil
}

func (c *client) licensing() error {
	log.Println("RDP: Server License")

	_, wire, err := c.mcsLayer.Receive()
	if err != nil {
		return err
	}

	var resp pdu.ServerLicenseError
	if err = resp.Deserialize(wire); err != nil {
		return fmt.Errorf("server license error: %w", err)
	}

	if resp.Preamble.MsgType == 0x03 { // NEW_LICENSE
		return nil
	}

	if resp.Preamble.MsgType != 0xFF { // ERROR_ALERT
		return fmt.Errorf("unknown license msg type: %d", resp.Preamble.MsgType)
	}

	if resp.ValidClientMessage.ErrorCode != 0x00000007 { // STATUS_VALID_CLIENT
		return fmt.Errorf("unknown license error code: %d", resp.ValidClientMessage.ErrorCode)
	}

	if resp.ValidClientMessage.StateTransition != 0x00000002 { // ST_NO_TRANSITION
		return fmt.Errorf("unknown license state transition: %d", resp.ValidClientMessage.StateTransition)
	}

	return nil
}
