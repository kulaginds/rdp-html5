package fastpath

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestInputEventPDU_Serialize(t *testing.T) {
	event := NewInputEventPDU([]byte{0x30, 0x35, 0x6b, 0x5b, 0xb5, 0x34, 0xc8, 0x47, 0x26, 0x18, 0x5e, 0x76, 0x0e, 0xde, 0x28})
	event.flags = 0x1 | 0x2 // FASTPATH_INPUT_SECURE_CHECKSUM | FASTPATH_INPUT_ENCRYPTED

	expected := []byte{
		0xc4, 0x11, 0x30, 0x35, 0x6b, 0x5b, 0xb5, 0x34, 0xc8, 0x47, 0x26, 0x18, 0x5e, 0x76, 0x0e, 0xde,
		0x28,
	}
	actual := event.Serialize()

	require.Equal(t, expected, actual)
}
