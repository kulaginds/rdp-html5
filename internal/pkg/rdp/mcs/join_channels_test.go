package mcs

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestClientMCSChannelJoinRequestPDU_Serialize from MS-RDPBCGR protocol examples 4.1.8.
// without TPKT and X224 headers
func TestClientMCSChannelJoinRequestPDU_Serialize(t *testing.T) {
	testCases := []struct {
		name string

		req ClientChannelJoinRequest

		expected []byte
	}{
		{
			name: "join channel 1007",

			req: ClientChannelJoinRequest{
				Initiator: 1007,
				ChannelId: 1007,
			},

			expected: []byte{
				0x38, 0x00, 0x06, 0x03, 0xef,
			},
		},
		{
			name: "join channel 1003",

			req: ClientChannelJoinRequest{
				Initiator: 1007,
				ChannelId: 1003,
			},

			expected: []byte{0x38, 0x00, 0x06, 0x03, 0xeb},
		},
		{
			name: "join channel 1004",

			req: ClientChannelJoinRequest{
				Initiator: 1007,
				ChannelId: 1004,
			},

			expected: []byte{0x38, 0x00, 0x06, 0x03, 0xec},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := tc.req.Serialize()

			require.Equal(t, tc.expected, actual)
		})
	}
}

// TestServerMCSChannelJoinConfirmPDU_Deserialize from MS-RDPBCGR protocol examples 4.1.8.
// without TPKT and X224 headers
func TestServerMCSChannelJoinConfirmPDU_Deserialize(t *testing.T) {
	testCases := []struct {
		name string

		input []byte

		expected ServerChannelJoinConfirm
	}{
		{
			name: "confirm join 1007",

			input: []byte{
				0x3e, 0x00, 0x00, 0x06, 0x03, 0xef, 0x03, 0xef,
			},

			expected: ServerChannelJoinConfirm{
				Result:    0x00,
				Initiator: 1007,
				Requested: 1007,
				ChannelId: 1007,
			},
		},
		{
			name: "confirm join 1003",

			input: []byte{0x3e, 0x00, 0x00, 0x06, 0x03, 0xeb, 0x03, 0xeb},

			expected: ServerChannelJoinConfirm{
				Result:    0x00,
				Initiator: 1007,
				Requested: 1003,
				ChannelId: 1003,
			},
		},
		{
			name: "confirm join 1004",

			input: []byte{0x3e, 0x00, 0x00, 0x06, 0x03, 0xec, 0x03, 0xec},

			expected: ServerChannelJoinConfirm{
				Result:    0x00,
				Initiator: 1007,
				Requested: 1004,
				ChannelId: 1004,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var actual ServerChannelJoinConfirm

			require.NoError(t, actual.Deserialize(bytes.NewBuffer(tc.input)))
			require.Equal(t, tc.expected, actual)
		})
	}
}
