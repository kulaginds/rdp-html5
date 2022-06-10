package x224

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/kulaginds/web-rdp-solution/internal/pkg/rdp/headers"
)

// ClientConnectionRequestPDU Client X.224 Connection Request PDU
type ClientConnectionRequestPDU struct {
	RoutingToken       string // one of RoutingToken or Cookie ending CR+LF
	Cookie             string
	RDPNegReq          RDPNegotiationRequest // RDP Negotiation Request
	RDPCorrelationInfo RDPCorrelationInfo    // Correlation Info
}

func (pdu *ClientConnectionRequestPDU) Serialize() []byte {
	const (
		CRLF         = "\r\n"
		cookieHeader = "Cookie: mstshash="
	)

	buf := &bytes.Buffer{}

	// routingToken or cookie
	if pdu.RoutingToken != "" {
		buf.WriteString(strings.Trim(pdu.RoutingToken, CRLF) + CRLF)
	} else if pdu.Cookie != "" {
		buf.WriteString(cookieHeader + strings.Trim(pdu.Cookie, CRLF) + CRLF)
	}

	// rdpNegReq
	buf.Write(pdu.RDPNegReq.Serialize())

	// rdpCorrelationInfo
	if pdu.RDPNegReq.Flags.IsCorrelationInfoPresent() {
		buf.Write(pdu.RDPCorrelationInfo.Serialize())
	}

	return headers.WrapX224ConnectionRequestPDU(buf.Bytes())
}

type ServerConnectionConfirmPDU struct {
	Type  RDPNegotiationType
	Flags RDPNegotiationResponseFlag // RDP Negotiation Response flags
	data  uint32                     // RDP Negotiation Response selectedProtocol or RDP Negotiation Failure failureCode
}

func (pdu ServerConnectionConfirmPDU) SelectedProtocol() RDPNegotiationProtocol {
	return RDPNegotiationProtocol(pdu.data)
}

func (pdu ServerConnectionConfirmPDU) FailureCode() RDPNegotiationFailureCode {
	return RDPNegotiationFailureCode(pdu.data)
}

func (pdu *ServerConnectionConfirmPDU) Deserialize(wire io.Reader) error {
	const (
		fixedPartLen    uint8 = 0x06
		variablePartLen uint8 = 0x08
		packetLen             = fixedPartLen + variablePartLen
	)

	var li uint8

	binary.Read(wire, binary.LittleEndian, &li)

	if li != packetLen {
		return ErrSmallConnectionConfirmLength
	}

	packetData := make([]byte, li)

	wire.Read(packetData)

	ccCdt := packetData[0]

	if ccCdt&0xf0 != 0xd0 { // connection confirm code
		return ErrWrongConnectionConfirmCode
	}

	// skip unused fields
	packetData = packetData[fixedPartLen:]

	pdu.Type = RDPNegotiationType(packetData[0])

	if pdu.Type.IsResponse() {
		pdu.Flags = RDPNegotiationResponseFlag(packetData[1])
	}

	pdu.data = binary.LittleEndian.Uint32(packetData[4:8])

	return nil
}

func (p *protocol) Connect() error {
	if p.connected {
		return nil
	}

	if !p.requestedProtocols.IsSSL() {
		return ErrUnsupportedRequestedProtocol
	}

	var (
		wire io.Reader
		err  error
	)

	req := ClientConnectionRequestPDU{
		RDPNegReq: RDPNegotiationRequest{
			RequestedProtocols: p.requestedProtocols,
		},
	}

	log.Println("X224: Client Connection Request")

	if err = p.tpktConn.Send(req.Serialize()); err != nil {
		return fmt.Errorf("client connection request: %w", err)
	}

	log.Println("X224: Server Connection Confirm")

	wire, err = p.tpktConn.Receive()
	if err != nil {
		return fmt.Errorf("recieve connection response: %w", err)
	}

	var resp ServerConnectionConfirmPDU
	if err = resp.Deserialize(wire); err != nil {
		return fmt.Errorf("server connection confirm: %w", err)
	}

	if resp.Type.IsFailure() {
		return fmt.Errorf("negotiation faliure: failureCode = %d", resp.FailureCode())
	}

	p.ServerNegotiationFlags = resp.Flags

	if !resp.SelectedProtocol().IsSSL() {
		return ErrUnsupportedRequestedProtocol
	}

	if err = p.tpktConn.StartTLS(); err != nil {
		return err
	}

	p.connected = true

	return nil
}
