package x224

import (
	"bytes"
	"encoding/binary"
)

const x224FixedPartLen = 6 // without length indicator (LI)

type RDPNegotiationType uint8

const (
	// RDPNegotiationTypeRequest TYPE_RDP_NEG_REQ
	RDPNegotiationTypeRequest RDPNegotiationType = 0x01

	// RDPNegotiationTypeResponse TYPE_RDP_NEG_RSP
	RDPNegotiationTypeResponse RDPNegotiationType = 0x02

	// RDPNegotiationTypeFailure TYPE_RDP_NEG_FAILURE
	RDPNegotiationTypeFailure RDPNegotiationType = 0x03
)

func (t RDPNegotiationType) IsRequest() bool {
	return t == RDPNegotiationTypeRequest
}

func (t RDPNegotiationType) IsResponse() bool {
	return t == RDPNegotiationTypeResponse
}

func (t RDPNegotiationType) IsFailure() bool {
	return t == RDPNegotiationTypeFailure
}

// RDPNegotiationRequestFlag Protocol flags.
type RDPNegotiationRequestFlag uint8

const (
	// RDPNegReqFlagRestrictedAdminModeRequired RESTRICTED_ADMIN_MODE_REQUIRED
	RDPNegReqFlagRestrictedAdminModeRequired RDPNegotiationRequestFlag = 0x01

	// RDPNegReqFlagRedirectedAuthenticationModeRequired REDIRECTED_AUTHENTICATION_MODE_REQUIRED
	RDPNegReqFlagRedirectedAuthenticationModeRequired RDPNegotiationRequestFlag = 0x02

	// RDPNegReqFlagCorrelationInfoPresent CORRELATION_INFO_PRESENT
	RDPNegReqFlagCorrelationInfoPresent RDPNegotiationRequestFlag = 0x08
)

func (f RDPNegotiationRequestFlag) IsRestrictedAdminModeRequired() bool {
	return f&RDPNegReqFlagRestrictedAdminModeRequired == RDPNegReqFlagRestrictedAdminModeRequired
}

func (f RDPNegotiationRequestFlag) IsRedirectedAuthenticationModeRequired() bool {
	return f&RDPNegReqFlagRedirectedAuthenticationModeRequired == RDPNegReqFlagRedirectedAuthenticationModeRequired
}

func (f RDPNegotiationRequestFlag) IsCorrelationInfoPresent() bool {
	return f&RDPNegReqFlagCorrelationInfoPresent == RDPNegReqFlagCorrelationInfoPresent
}

// RDPNegotiationProtocol Supported security protocol.
type RDPNegotiationProtocol uint32

const (
	// RDPNegotiationProtocolRDP PROTOCOL_RDP
	RDPNegotiationProtocolRDP RDPNegotiationProtocol = 0x00000000

	// RDPNegotiationProtocolSSL PROTOCOL_SSL
	RDPNegotiationProtocolSSL RDPNegotiationProtocol = 0x00000001

	// RDPNegotiationProtocolHybrid PROTOCOL_HYBRID
	RDPNegotiationProtocolHybrid RDPNegotiationProtocol = 0x00000002

	// RDPNegotiationProtocolRDSTLS PROTOCOL_RDSTLS
	RDPNegotiationProtocolRDSTLS RDPNegotiationProtocol = 0x00000004

	// RDPNegotiationProtocolHybridEx PROTOCOL_HYBRID_EX
	RDPNegotiationProtocolHybridEx RDPNegotiationProtocol = 0x00000008
)

func (p RDPNegotiationProtocol) IsRDP() bool {
	return p == RDPNegotiationProtocolRDP
}

func (p RDPNegotiationProtocol) IsSSL() bool {
	return p == RDPNegotiationProtocolSSL
}

func (p RDPNegotiationProtocol) IsHybrid() bool {
	return p == RDPNegotiationProtocolHybrid
}

func (p RDPNegotiationProtocol) IsRDSTLS() bool {
	return p == RDPNegotiationProtocolRDSTLS
}

func (p RDPNegotiationProtocol) IsHybridEx() bool {
	return p == RDPNegotiationProtocolHybridEx
}

// RDPNegotiationRequest RDP Negotiation Request (RDP_NEG_REQ).
type RDPNegotiationRequest struct {
	Flags              RDPNegotiationRequestFlag // protocol flags
	RequestedProtocols RDPNegotiationProtocol    // supported security protocols
}

func (r RDPNegotiationRequest) Serialize() []byte {
	const negReqLen = uint16(8)

	buf := bytes.NewBuffer(make([]byte, 0, negReqLen))

	buf.Write([]byte{
		byte(RDPNegotiationTypeRequest), // type TYPE_RDP_NEG_REQ
		byte(r.Flags),                   // flags
	})

	// length (always 8 bytes)
	_ = binary.Write(buf, binary.LittleEndian, negReqLen)

	// requestedProtocols
	_ = binary.Write(buf, binary.LittleEndian, r.RequestedProtocols)

	return buf.Bytes()
}

// RDPCorrelationInfo RDP Correlation Info (RDP_NEG_CORRELATION_INFO).
type RDPCorrelationInfo struct {
	correlationID []byte
}

func (i RDPCorrelationInfo) SetCorrelationID(correlationID []byte) error {
	if len(correlationID) != 16 {
		return ErrInvalidCorrelationID
	}

	// The first byte in the array SHOULD NOT have a value of 0x00 or 0xF4
	if correlationID[0] == 0x00 || correlationID[0] == 0xF4 {
		return ErrInvalidCorrelationID
	}

	// value 0x0D SHOULD NOT be contained in any of the bytes
	for _, b := range correlationID {
		if b == 0x0D {
			return ErrInvalidCorrelationID
		}
	}

	return nil
}

func (i RDPCorrelationInfo) Serialize() []byte {
	const corrInfoLen = uint16(36)

	buf := bytes.NewBuffer(make([]byte, 0, corrInfoLen))

	buf.Write([]byte{
		0x06, // type TYPE_RDP_CORRELATION_INFO
		0x00, // flags
	})

	// length (always 36 bytes)
	_ = binary.Write(buf, binary.LittleEndian, corrInfoLen)

	// correlationId
	if i.correlationID == nil {
		buf.Write(make([]byte, 16))
	} else {
		buf.Write(i.correlationID)
	}

	// reserved
	buf.Write(make([]byte, 16))

	return buf.Bytes()
}

// RDPNegotiationResponseFlag RDP Negotiation Response flags
type RDPNegotiationResponseFlag uint8

const (
	// RDPNegotiationResponseFlagECDBSupported EXTENDED_CLIENT_DATA_SUPPORTED
	RDPNegotiationResponseFlagECDBSupported RDPNegotiationResponseFlag = 0x01

	// RDPNegotiationResponseFlagGFXSupported DYNVC_GFX_PROTOCOL_SUPPORTED
	RDPNegotiationResponseFlagGFXSupported RDPNegotiationResponseFlag = 0x02

	// RDPNegotiationResponseFlagReserved NEGRSP_FLAG_RESERVED
	RDPNegotiationResponseFlagReserved RDPNegotiationResponseFlag = 0x04

	// RDPNegotiationResponseFlagAdminModeSupported RESTRICTED_ADMIN_MODE_SUPPORTED
	RDPNegotiationResponseFlagAdminModeSupported RDPNegotiationResponseFlag = 0x08

	// RDPNegotiationResponseFlagAuthModeSupported REDIRECTED_AUTHENTICATION_MODE_SUPPORTED
	RDPNegotiationResponseFlagAuthModeSupported RDPNegotiationResponseFlag = 0x10
)

func (f RDPNegotiationResponseFlag) IsExtendedClientDataSupported() bool {
	return f&RDPNegotiationResponseFlagECDBSupported == RDPNegotiationResponseFlagECDBSupported
}

func (f RDPNegotiationResponseFlag) IsGFXProtocolSupported() bool {
	return f&RDPNegotiationResponseFlagGFXSupported == RDPNegotiationResponseFlagGFXSupported
}

func (f RDPNegotiationResponseFlag) IsRestrictedAdminModeSupported() bool {
	return f&RDPNegotiationResponseFlagAdminModeSupported == RDPNegotiationResponseFlagAdminModeSupported
}

func (f RDPNegotiationResponseFlag) IsRedirectedAuthModeSupported() bool {
	return f&RDPNegotiationResponseFlagAuthModeSupported == RDPNegotiationResponseFlagAuthModeSupported
}

// RDPNegotiationFailureCode RDP Negotiation Failure failureCode
type RDPNegotiationFailureCode uint32

const (
	// RDPNegotiationFailureCodeSSLRequired SSL_REQUIRED_BY_SERVER
	RDPNegotiationFailureCodeSSLRequired RDPNegotiationFailureCode = 0x00000001

	// RDPNegotiationFailureCodeSSLNotAllowed SSL_NOT_ALLOWED_BY_SERVER
	RDPNegotiationFailureCodeSSLNotAllowed RDPNegotiationFailureCode = 0x00000002

	// RDPNegotiationFailureCodeSSLCertNotOnServer SSL_CERT_NOT_ON_SERVER
	RDPNegotiationFailureCodeSSLCertNotOnServer RDPNegotiationFailureCode = 0x00000003

	// RDPNegotiationFailureCodeInconsistentFlags INCONSISTENT_FLAGS
	RDPNegotiationFailureCodeInconsistentFlags RDPNegotiationFailureCode = 0x00000004

	// RDPNegotiationFailureCodeHybridRequired HYBRID_REQUIRED_BY_SERVER
	RDPNegotiationFailureCodeHybridRequired RDPNegotiationFailureCode = 0x00000005

	// RDPNegotiationFailureCodeSSLWithUserAuthRequired SSL_WITH_USER_AUTH_REQUIRED_BY_SERVER
	RDPNegotiationFailureCodeSSLWithUserAuthRequired RDPNegotiationFailureCode = 0x00000006
)
