package pdu

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_NewClientUserDataSet(t *testing.T) {
	r := require.New(t)

	input := NewClientUserDataSet(0x00000000, 1280, 1024, []string{"rdpdr", "cliprdp", "rdpsnd"})
	input.ClientCoreData.ColorDepth = 0xca01
	input.ClientCoreData.SASSequence = 0xaa03
	input.ClientCoreData.KeyboardLayout = 0x409
	input.ClientCoreData.ClientBuild = 3790
	input.ClientCoreData.ClientName = [32]byte{'E', 0x0, 'L', 0x0, 'T', 0x0, 'O', 0x0, 'N', 0x0, 'S', 0x0, '-', 0x0, 'D', 0x0, 'E', 0x0, 'V', 0x0, '2', 0x0}
	input.ClientCoreData.PostBeta2ColorDepth = 0xca01
	input.ClientCoreData.SupportedColorDepths = 0x0007
	input.ClientCoreData.EarlyCapabilityFlags = 0x0001
	input.ClientCoreData.HighColorDepth = 0x0018
	input.ClientCoreData.ClientDigProductId = [64]byte{'6', 0x0, '9', 0x0, '7', 0x0, '1', 0x0, '2', 0x0, '-', 0x0, '7', 0x0, '8', 0x0, '3', 0x0, '-', 0x0, '0', 0x0, '3', 0x0, '5', 0x0, '7', 0x0, '9', 0x0, '7', 0x0, '4', 0x0, '-', 0x0, '4', 0x0, '2', 0x0, '7', 0x0, '1', 0x0, '4', 0x0}
	input.ClientCoreData.ServerSelectedProtocol = 0x00000000
	input.ClientClusterData = &ClientClusterData{Flags: 0x0000000d}
	input.ClientSecurityData.EncryptionMethods = 0x0000001b
	input.ClientNetworkData.ChannelCount = 3
	input.ClientNetworkData.ChannelDefArray = []ChannelDefinitionStructure{
		{
			Name:    [8]byte{'r', 'd', 'p', 'd', 'r'},
			Options: 0x80800000,
		},
		{
			Name:    [8]byte{'c', 'l', 'i', 'p', 'r', 'd', 'r'},
			Options: 0xc0a00000,
		},
		{
			Name:    [8]byte{'r', 'd', 'p', 's', 'n', 'd'},
			Options: 0xc0000000,
		},
	}

	expected := []byte{
		0x01, 0xc0, 0xd8, 0x00, 0x04, 0x00, 0x08, 0x00, 0x00, 0x05, 0x00, 0x04,
		0x01, 0xca, 0x03, 0xaa, 0x09, 0x04, 0x00, 0x00, 0xce, 0x0e, 0x00, 0x00, 0x45, 0x00, 0x4c, 0x00,
		0x54, 0x00, 0x4f, 0x00, 0x4e, 0x00, 0x53, 0x00, 0x2d, 0x00, 0x44, 0x00, 0x45, 0x00, 0x56, 0x00,
		0x32, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x04, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x0c, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0xca, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x18, 0x00, 0x07, 0x00, 0x01, 0x00, 0x36, 0x00, 0x39, 0x00, 0x37, 0x00, 0x31, 0x00, 0x32, 0x00,
		0x2d, 0x00, 0x37, 0x00, 0x38, 0x00, 0x33, 0x00, 0x2d, 0x00, 0x30, 0x00, 0x33, 0x00, 0x35, 0x00,
		0x37, 0x00, 0x39, 0x00, 0x37, 0x00, 0x34, 0x00, 0x2d, 0x00, 0x34, 0x00, 0x32, 0x00, 0x37, 0x00,
		0x31, 0x00, 0x34, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x04, 0xc0, 0x0c, 0x00,
		0x0d, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xc0, 0x0c, 0x00, 0x1b, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x03, 0xc0, 0x2c, 0x00, 0x03, 0x00, 0x00, 0x00, 0x72, 0x64, 0x70, 0x64,
		0x72, 0x00, 0x00, 0x00, 0x00, 0x00, 0x80, 0x80, 0x63, 0x6c, 0x69, 0x70, 0x72, 0x64, 0x72, 0x00,
		0x00, 0x00, 0xa0, 0xc0, 0x72, 0x64, 0x70, 0x73, 0x6e, 0x64, 0x00, 0x00, 0x00, 0x00, 0x00, 0xc0,
	}

	r.Equal(expected, input.Serialize())
}
