package pdu

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestClientConnectionRequestPDU_Serialize from MS-RDPBCGR Protocol examples 4.1.1.
// without TPKT header
func TestClientConnectionRequestPDU_Serialize(t *testing.T) {
	var req ClientConnectionRequest

	req.Cookie = "eltons"
	req.NegotiationRequest.RequestedProtocols = NegotiationProtocolRDP

	actual := req.Serialize()
	expected := []byte{
		0x43, 0x6f, 0x6f, 0x6b, 0x69, 0x65, 0x3a, 0x20, 0x6d, 0x73, 0x74, 0x73, 0x68, 0x61, 0x73, 0x68,
		0x3d, 0x65, 0x6c, 0x74, 0x6f, 0x6e, 0x73, 0x0d, 0x0a, 0x01, 0x00, 0x08, 0x00, 0x00, 0x00, 0x00,
		0x00,
	}

	require.Equal(t, expected, actual)
}

// TestServerConnectionConfirmPDU_Deserialize from MS-RDPBCGR Protocol examples 4.1.2.
// without TPKT header
func TestServerConnectionConfirmPDU_Deserialize(t *testing.T) {
	var actual ServerConnectionConfirm

	expected := ServerConnectionConfirm{
		Type:   NegotiationTypeResponse,
		Flags:  0,
		length: 8,
		data:   uint32(NegotiationProtocolRDP),
	}

	input := bytes.NewBuffer([]byte{
		0x02, 0x00, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00,
	})

	require.NoError(t, actual.Deserialize(input))
	require.Equal(t, expected, actual)
}
