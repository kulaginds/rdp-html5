package rdp

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"strings"
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

	c.tpktLayer.StartHandleFastpath()

	return nil
}

// ClientConnectionRequestPDU Client X.224 Connection Request PDU
type ClientConnectionRequestPDU struct {
	RoutingToken       string // one of RoutingToken or Cookie ending CR+LF
	Cookie             string
	NegotiationRequest NegotiationRequest // RDP Negotiation Request
	CorrelationInfo    CorrelationInfo    // Correlation Info
}

func (pdu *ClientConnectionRequestPDU) Serialize() []byte {
	const (
		CRLF         = "\r\n"
		cookieHeader = "Cookie: mstshash="
	)

	buf := new(bytes.Buffer)

	// routingToken or cookie
	if pdu.RoutingToken != "" {
		buf.WriteString(strings.Trim(pdu.RoutingToken, CRLF) + CRLF)
	} else if pdu.Cookie != "" {
		buf.WriteString(cookieHeader + strings.Trim(pdu.Cookie, CRLF) + CRLF)
	}

	// rdpNegReq
	buf.Write(pdu.NegotiationRequest.Serialize())

	// rdpCorrelationInfo
	if pdu.NegotiationRequest.Flags.IsCorrelationInfoPresent() {
		buf.Write(pdu.CorrelationInfo.Serialize())
	}

	return buf.Bytes()
}

type ServerConnectionConfirmPDU struct {
	Type   NegotiationType
	Flags  NegotiationResponseFlag // RDP Negotiation Response flags
	length uint16
	data   uint32 // RDP Negotiation Response selectedProtocol or RDP Negotiation Failure failureCode
}

func (pdu *ServerConnectionConfirmPDU) SelectedProtocol() NegotiationProtocol {
	return NegotiationProtocol(pdu.data)
}

func (pdu *ServerConnectionConfirmPDU) FailureCode() NegotiationFailureCode {
	return NegotiationFailureCode(pdu.data)
}

func (pdu *ServerConnectionConfirmPDU) Deserialize(wire io.Reader) error {
	err := binary.Read(wire, binary.LittleEndian, &pdu.Type)
	if err != nil {
		return err
	}

	err = binary.Read(wire, binary.LittleEndian, &pdu.Flags)
	if err != nil {
		return err
	}

	err = binary.Read(wire, binary.LittleEndian, &pdu.length)
	if err != nil {
		return err
	}

	err = binary.Read(wire, binary.LittleEndian, &pdu.data)
	if err != nil {
		return err
	}

	return nil
}

func (c *client) connectionInitiation() error {
	var err error

	req := ClientConnectionRequestPDU{
		NegotiationRequest: NegotiationRequest{
			RequestedProtocols: c.selectedProtocol,
		},
	}

	var (
		resp ServerConnectionConfirmPDU
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

	if c.remoteApp != nil {
		clientInfoPDU.InfoPacket.Flags |= InfoFlagRail
	}

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
