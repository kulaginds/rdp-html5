package mcs

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestClientDisconnectUltimatumRequest_Serialize(t *testing.T) {
	req := &ClientDisconnectUltimatumRequest{}

	expected := []byte{0x21, 0x80}
	actual := req.Serialize()

	require.Equal(t, expected, actual)
}
