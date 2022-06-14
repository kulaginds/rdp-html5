package rdp

import (
	"fmt"
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

	return nil
}

func (c *client) connectionInitiation() error {
	var err error

	if err = c.x224Layer.Connect(); err != nil {
		return err
	}

	if c.selectedProtocol.IsSSL() {
		return c.StartTLS()
	}

	return nil
}

func (c *client) basicSettingsExchange() error {
	return c.mcsLayer.Connect(uint32(c.selectedProtocol), c.desktopWidth, c.desktopHeight, c.channels)
}

func (c *client) channelConnection() error {
	err := c.mcsLayer.ErectDomain()
	if err != nil {
		return err
	}

	err = c.mcsLayer.AttachUser()
	if err != nil {
		return err
	}

	err = c.mcsLayer.JoinChannels()
	if err != nil {
		return err
	}

	return nil
}

func (c *client) secureSettingsExchange() error {
	clientInfoPDU := NewClientInfoPDU(c.domain, c.username, c.password)

	log.Println("RDP: Client Info")

	if err := c.mcsLayer.Send("global", clientInfoPDU.Serialize()); err != nil {
		return fmt.Errorf("client info: %w", err)
	}

	return nil
}

func (c *client) licensing() error {
	log.Println("RDP: Server License Error")

	_, wire, err := c.mcsLayer.Receive()
	if err != nil {
		return err
	}

	var resp ServerLicenseErrorPDU
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
