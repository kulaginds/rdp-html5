package mcs

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// TestClientMCSErectDomainRequestPDU_Serialize from MS-RDPBCGR protocol examples 4.1.5.
// without TPKT and X224 headers
func TestClientMCSErectDomainRequestPDU_Serialize(t *testing.T) {
	req := DomainPDU{
		Application:              erectDomainRequest,
		ClientErectDomainRequest: &ClientErectDomainRequest{},
	}

	expected := []byte{0x04, 0x01, 0x00, 0x01, 0x00}
	actual := req.Serialize()

	require.Equal(t, expected, actual)
}
