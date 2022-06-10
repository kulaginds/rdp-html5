package gcc

import (
	"bytes"
	"encoding/binary"

	"github.com/kulaginds/web-rdp-solution/internal/pkg/rdp/per"
)

type ClientCoreData struct {
	Version                uint32
	DesktopWidth           uint16
	DesktopHeight          uint16
	ColorDepth             uint16
	SASSequence            uint16
	KeyboardLayout         uint32
	ClientBuild            uint32
	ClientName             [32]byte
	KeyboardType           uint32
	KeyboardSubType        uint32
	KeyboardFunctionKey    uint32
	ImeFileName            [64]byte
	PostBeta2ColorDepth    uint16
	ClientProductId        uint16
	SerialNumber           uint32
	HighColorDepth         uint16
	SupportedColorDepths   uint16
	EarlyCapabilityFlags   uint16
	ClientDigProductId     [64]byte
	ConnectionType         uint8
	Pad1octet              uint8
	ServerSelectedProtocol uint32
}

func newClientCoreData(selectedProtocol uint32, desktopWidth, desktopHeight uint16) *ClientCoreData {
	data := ClientCoreData{
		Version:        rdpVersion5Plus,
		DesktopWidth:   desktopWidth,
		DesktopHeight:  desktopHeight,
		ColorDepth:     0xCA01,     // RNS_UD_COLOR_8BPP
		SASSequence:    0xAA03,     // RNS_UD_SAS_DEL
		KeyboardLayout: 0x00000409, // US
		ClientBuild:    0xece,
		ClientName: [32]byte{
			'w', 'e', 'b', '-', 'r', 'd', 'p', '-', 's', 'o', 'l', 'u', 't', 'i', 'o', 'n',
		},
		KeyboardType:           keyboardTypeIBM101or102Keys,
		KeyboardSubType:        0x00000000,
		KeyboardFunctionKey:    12,
		ImeFileName:            [64]byte{},
		PostBeta2ColorDepth:    0xCA01, // RNS_UD_COLOR_8BPP
		ClientProductId:        0x0001,
		SerialNumber:           0x00000000,
		HighColorDepth:         0x0018,                            // HIGH_COLOR_24BPP
		SupportedColorDepths:   0x0004 | 0x0002 | 0x0001 | 0x0008, // RNS_UD_15BPP_SUPPORT, RNS_UD_16BPP_SUPPORT, RNS_UD_24BPP_SUPPORT, RNS_UD_32BPP_SUPPORT
		EarlyCapabilityFlags:   0x0001,                            // RNS_UD_CS_SUPPORT_ERRINFO_PDU
		ClientDigProductId:     [64]byte{},
		ConnectionType:         0x00,
		Pad1octet:              0x00,
		ServerSelectedProtocol: selectedProtocol,
	}

	return &data
}

const (
	// EncryptionFlag40Bit ENCRYPTION_FLAG_40BIT
	EncryptionFlag40Bit uint32 = 0x00000001

	// EncryptionFlag128Bit ENCRYPTION_FLAG_128BIT
	EncryptionFlag128Bit uint32 = 0x00000002

	// EncryptionFlag56Bit ENCRYPTION_FLAG_56BIT
	EncryptionFlag56Bit uint32 = 0x00000008

	// FIPSEncryptionFlag FIPS_ENCRYPTION_FLAG
	FIPSEncryptionFlag uint32 = 0x00000010
)

type ClientSecurityData struct {
	EncryptionMethods    uint32
	ExtEncryptionMethods uint32
}

func newClientSecurityData() *ClientSecurityData {
	data := ClientSecurityData{
		EncryptionMethods:    0,
		ExtEncryptionMethods: 0,
	}

	return &data
}

type ChannelDefinitionStructure struct {
	Name    [8]byte // seven ANSI chars with null-termination char in the end
	Options uint32
}

type ClientNetworkData struct {
	ChannelCount    uint32
	ChannelDefArray []ChannelDefinitionStructure
}

type ClientClusterData struct {
	Flags               uint32
	RedirectedSessionID uint32
}

type ConferenceCreateRequestUserData struct {
	ClientCoreData     *ClientCoreData
	ClientSecurityData *ClientSecurityData
	ClientNetworkData  *ClientNetworkData
	ClientClusterData  *ClientClusterData
}

type ConferenceCreateRequest struct {
	UserData ConferenceCreateRequestUserData
}

func NewConferenceCreateRequest(
	selectedProtocol uint32,
	desktopWidth, desktopHeight uint16,
	channelNames []string,
) *ConferenceCreateRequest {
	return &ConferenceCreateRequest{
		UserData: ConferenceCreateRequestUserData{
			ClientCoreData:     newClientCoreData(selectedProtocol, desktopWidth, desktopHeight),
			ClientSecurityData: newClientSecurityData(),
			ClientNetworkData:  newClientNetworkData(channelNames),
		},
	}
}

func (data ClientCoreData) Serialize() []byte {
	const dataLen uint16 = 216

	buf := bytes.NewBuffer(make([]byte, 0, dataLen))

	binary.Write(buf, binary.LittleEndian, uint16(0xC001)) // header type CS_CORE
	binary.Write(buf, binary.LittleEndian, dataLen)        // packet size

	binary.Write(buf, binary.LittleEndian, data.Version)
	binary.Write(buf, binary.LittleEndian, data.DesktopWidth)
	binary.Write(buf, binary.LittleEndian, data.DesktopHeight)
	binary.Write(buf, binary.LittleEndian, data.ColorDepth)
	binary.Write(buf, binary.LittleEndian, data.SASSequence)
	binary.Write(buf, binary.LittleEndian, data.KeyboardLayout)
	binary.Write(buf, binary.LittleEndian, data.ClientBuild)
	binary.Write(buf, binary.LittleEndian, data.ClientName)
	binary.Write(buf, binary.LittleEndian, data.KeyboardType)
	binary.Write(buf, binary.LittleEndian, data.KeyboardSubType)
	binary.Write(buf, binary.LittleEndian, data.KeyboardFunctionKey)
	binary.Write(buf, binary.LittleEndian, data.ImeFileName)
	binary.Write(buf, binary.LittleEndian, data.PostBeta2ColorDepth)
	binary.Write(buf, binary.LittleEndian, data.ClientProductId)
	binary.Write(buf, binary.LittleEndian, data.SerialNumber)
	binary.Write(buf, binary.LittleEndian, data.HighColorDepth)
	binary.Write(buf, binary.LittleEndian, data.SupportedColorDepths)
	binary.Write(buf, binary.LittleEndian, data.EarlyCapabilityFlags)
	binary.Write(buf, binary.LittleEndian, data.ClientDigProductId)
	binary.Write(buf, binary.LittleEndian, data.ConnectionType)
	binary.Write(buf, binary.LittleEndian, data.Pad1octet)
	binary.Write(buf, binary.LittleEndian, data.ServerSelectedProtocol)

	return buf.Bytes()
}

const (
	// EncryptionMethodFlag40Bit 40BIT_ENCRYPTION_FLAG
	EncryptionMethodFlag40Bit uint32 = 0x00000001

	// EncryptionMethodFlag56Bit 56BIT_ENCRYPTION_FLAG
	EncryptionMethodFlag56Bit uint32 = 0x00000008

	// EncryptionMethodFlag128Bit 128BIT_ENCRYPTION_FLAG
	EncryptionMethodFlag128Bit uint32 = 0x00000002

	// EncryptionMethodFlagFIPS FIPS_ENCRYPTION_FLAG
	EncryptionMethodFlagFIPS uint32 = 0x00000010
)

func (data ClientSecurityData) Serialize() []byte {
	const dataLen uint16 = 12

	buf := bytes.NewBuffer(make([]byte, 0, 6))

	binary.Write(buf, binary.LittleEndian, uint16(0xC002)) // header type CS_SECURITY
	binary.Write(buf, binary.LittleEndian, dataLen)        // packet size

	binary.Write(buf, binary.LittleEndian, data.EncryptionMethods)
	binary.Write(buf, binary.LittleEndian, data.ExtEncryptionMethods)

	return buf.Bytes()
}

func (s ChannelDefinitionStructure) Serialize() []byte {
	buf := &bytes.Buffer{}

	binary.Write(buf, binary.LittleEndian, s.Name)
	binary.Write(buf, binary.LittleEndian, s.Options)

	return buf.Bytes()
}

func newClientNetworkData(channelNames []string) *ClientNetworkData {
	data := ClientNetworkData{
		ChannelCount: uint32(len(channelNames)),
	}

	if data.ChannelCount == 0 {
		return &data
	}

	for _, channelName := range channelNames {
		channelDefinition := ChannelDefinitionStructure{}
		copy(channelDefinition.Name[:], channelName)

		data.ChannelDefArray = append(data.ChannelDefArray, channelDefinition)
	}

	return &data
}

func (data ClientNetworkData) Serialize() []byte {
	const headerLen = 8

	chBuf := &bytes.Buffer{}

	for _, channelDef := range data.ChannelDefArray {
		chBuf.Write(channelDef.Serialize())
	}

	buf := &bytes.Buffer{}

	binary.Write(buf, binary.LittleEndian, uint16(0xC003))                // header type CS_NET
	binary.Write(buf, binary.LittleEndian, uint16(headerLen+chBuf.Len())) // packet size

	binary.Write(buf, binary.LittleEndian, data.ChannelCount)

	buf.Write(chBuf.Bytes())

	return buf.Bytes()
}

func (d ClientClusterData) Serialize() []byte {
	const dataLen uint16 = 12

	buf := &bytes.Buffer{}

	binary.Write(buf, binary.LittleEndian, uint16(0xC004)) // header type CS_CLUSTER
	binary.Write(buf, binary.LittleEndian, dataLen)        // packet size

	binary.Write(buf, binary.LittleEndian, d.Flags)
	binary.Write(buf, binary.LittleEndian, d.RedirectedSessionID)

	return buf.Bytes()
}

func (ud ConferenceCreateRequestUserData) Serialize() []byte {
	buf := bytes.Buffer{}

	buf.Write(ud.ClientCoreData.Serialize())

	if ud.ClientClusterData != nil {
		buf.Write(ud.ClientClusterData.Serialize())
	}

	buf.Write(ud.ClientSecurityData.Serialize())
	buf.Write(ud.ClientNetworkData.Serialize())

	return buf.Bytes()
}

func (r *ConferenceCreateRequest) Serialize() []byte {
	buf := &bytes.Buffer{}

	userData := r.UserData.Serialize()

	per.WriteChoice(0, buf)
	per.WriteObjectIdentifier(t124_02_98_oid, buf)
	per.WriteLength(uint16(14+len(userData)), buf)

	per.WriteChoice(0, buf)
	per.WriteSelection(0x08, buf)

	per.WriteNumericString("1", 1, buf)
	per.WritePadding(1, buf)
	per.WriteNumberOfSet(1, buf)
	per.WriteChoice(0xc0, buf)
	per.WriteOctetStream(h221CSKey, 4, buf)
	per.WriteOctetStream(string(userData), 0, buf)

	return buf.Bytes()
}
