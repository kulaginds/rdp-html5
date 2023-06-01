package pdu

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/kulaginds/rdp-html5/internal/pkg/rdp/mcs"
)

func TestClientSynchronizePDU_Serialize(t *testing.T) {
	synchronize := NewSynchronize(66538, 1007)

	req := mcs.DomainPDU{
		Application: mcs.SendDataRequest,
		ClientSendDataRequest: &mcs.ClientSendDataRequest{
			Initiator: 1007,
			ChannelId: 1003, // global
			Data:      synchronize.Serialize(),
		},
	}

	expected := []byte{
		0x64, 0x00, 0x06, 0x03, 0xeb, 0x70, 0x16, 0x16, 0x00, 0x17, 0x00, 0xef, 0x03, 0xea, 0x03, 0x01,
		0x00, 0x00, 0x01, 0x08, 0x00, 0x1f, 0x00, 0x00, 0x00, 0x01, 0x00, 0xea, 0x03,
	}
	actual := req.Serialize()

	require.Equal(t, expected, actual)
}

func TestClientSynchronizePDU_Serialize2(t *testing.T) {
	synchronize := NewSynchronize(66538, 1004)

	req := mcs.DomainPDU{
		Application: mcs.SendDataRequest,
		ClientSendDataRequest: &mcs.ClientSendDataRequest{
			Initiator: 1004,
			ChannelId: 1003, // global
			Data:      synchronize.Serialize(),
		},
	}

	expected, err := hex.DecodeString("64000303eb701616001700ec03ea030100000108001f0000000100ea03")
	require.NoError(t, err)

	actual := req.Serialize()

	require.Equal(t, expected, actual)
}

func TestClientControlPDU_Serialize_Cooperate(t *testing.T) {
	control := NewControl(66538, 1007, ControlActionCooperate)

	req := mcs.DomainPDU{
		Application: mcs.SendDataRequest,
		ClientSendDataRequest: &mcs.ClientSendDataRequest{
			Initiator: 1007,
			ChannelId: 1003, // global
			Data:      control.Serialize(),
		},
	}

	expected := []byte{
		0x64, 0x00, 0x06, 0x03, 0xeb, 0x70, 0x1a, 0x1a, 0x00, 0x17, 0x00, 0xef, 0x03, 0xea, 0x03, 0x01,
		0x00, 0x00, 0x01, 0x0c, 0x00, 0x14, 0x00, 0x00, 0x00, 0x04, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00,
	}
	actual := req.Serialize()

	require.Equal(t, expected, actual)
}

func TestClientControlPDU_Serialize_Cooperate2(t *testing.T) {
	control := NewControl(66538, 1004, ControlActionCooperate)

	req := mcs.DomainPDU{
		Application: mcs.SendDataRequest,
		ClientSendDataRequest: &mcs.ClientSendDataRequest{
			Initiator: 1004,
			ChannelId: 1003, // global
			Data:      control.Serialize(),
		},
	}

	expected, err := hex.DecodeString("64000303eb701a1a001700ec03ea03010000010c00140000000400000000000000")
	require.NoError(t, err)

	actual := req.Serialize()

	require.Equal(t, expected, actual)
}

func TestClientControlPDU_Serialize_RequestControl(t *testing.T) {
	control := NewControl(66538, 1007, ControlActionRequestControl)

	req := mcs.DomainPDU{
		Application: mcs.SendDataRequest,
		ClientSendDataRequest: &mcs.ClientSendDataRequest{
			Initiator: 1007,
			ChannelId: 1003, // global
			Data:      control.Serialize(),
		},
	}

	expected := []byte{
		0x64, 0x00, 0x06, 0x03, 0xeb, 0x70, 0x1a, 0x1a, 0x00, 0x17, 0x00, 0xef, 0x03, 0xea, 0x03, 0x01,
		0x00, 0x00, 0x01, 0x0c, 0x00, 0x14, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00,
	}
	actual := req.Serialize()

	require.Equal(t, expected, actual)
}

func TestClientControlPDU_Serialize_RequestControl2(t *testing.T) {
	control := NewControl(66538, 1004, ControlActionRequestControl)

	req := mcs.DomainPDU{
		Application: mcs.SendDataRequest,
		ClientSendDataRequest: &mcs.ClientSendDataRequest{
			Initiator: 1004,
			ChannelId: 1003, // global
			Data:      control.Serialize(),
		},
	}

	expected, err := hex.DecodeString("64000303eb701a1a001700ec03ea03010000010c00140000000100000000000000")
	require.NoError(t, err)

	actual := req.Serialize()

	require.Equal(t, expected, actual)
}

func TestClientFontListPDU_Serialize(t *testing.T) {
	fontlist := NewFontList(66538, 1007)

	req := mcs.DomainPDU{
		Application: mcs.SendDataRequest,
		ClientSendDataRequest: &mcs.ClientSendDataRequest{
			Initiator: 1007,
			ChannelId: 1003, // global
			Data:      fontlist.Serialize(),
		},
	}

	expected := []byte{
		0x64, 0x00, 0x06, 0x03, 0xeb, 0x70, 0x1a, 0x1a, 0x00, 0x17, 0x00, 0xef, 0x03, 0xea, 0x03, 0x01,
		0x00, 0x00, 0x01, 0x0c, 0x00, 0x27, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x03, 0x00, 0x32,
		0x00,
	}
	actual := req.Serialize()

	require.Equal(t, expected, actual)
}

func TestClientFontListPDU_Serialize2(t *testing.T) {
	fontlist := NewFontList(66538, 1004)

	req := mcs.DomainPDU{
		Application: mcs.SendDataRequest,
		ClientSendDataRequest: &mcs.ClientSendDataRequest{
			Initiator: 1004,
			ChannelId: 1003, // global
			Data:      fontlist.Serialize(),
		},
	}

	expected, err := hex.DecodeString("64000303eb701a1a001700ec03ea03010000010c00270000000000000003003200")
	require.NoError(t, err)

	actual := req.Serialize()

	require.Equal(t, expected, actual)
}
