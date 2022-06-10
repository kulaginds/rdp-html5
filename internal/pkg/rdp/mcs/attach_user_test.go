package mcs

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestClientMCSAttachUserRequestPDU_Serialize(t *testing.T) {
	req := DomainPDU{
		Application:             attachUserRequest,
		ClientAttachUserRequest: &ClientAttachUserRequest{},
	}

	expected := []byte{0x28}
	actual := req.Serialize()

	require.Equal(t, expected, actual)
}

func TestServerMCSAttachUserConfirmPDU_Deserialize(t *testing.T) {
	var actual DomainPDU

	expected := DomainPDU{
		Application: attachUserConfirm,
		ServerAttachUserConfirm: &ServerAttachUserConfirm{
			Result:    0x00,
			Initiator: 1007,
		},
	}

	input := bytes.NewBuffer([]byte{0x2e, 0x00, 0x00, 0x06})

	require.NoError(t, actual.Deserialize(input))
	require.Equal(t, expected, actual)
}
