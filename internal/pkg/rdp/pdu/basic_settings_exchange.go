package pdu

import (
	"bytes"
	"encoding/binary"
	"errors"
	"github.com/kulaginds/web-rdp-solution/internal/pkg/rdp/utf16"
	"io"
)

const (
	rdpVersion5Plus             = 0x00080004
	keyboardTypeIBM101or102Keys = 0x00000004
)

// earlyCapabilityFlags
const (
	ECFSupportErrInfoPDU        uint16 = 0x0001
	ECFWant32BPPSession         uint16 = 0x0002
	ECFSupportStatusInfoPDU     uint16 = 0x0004
	ECFStrongAsymmetricKeys     uint16 = 0x0008
	ECFUnused                   uint16 = 0x0010
	ECFValidConnectionType      uint16 = 0x0020
	ECFSupportMonitorLayoutPDU  uint16 = 0x0040
	ECFSupportNetCharAutodetect uint16 = 0x0080
	ECFSupportDynvcGFXProtocol  uint16 = 0x0100
	ECFSupportDynamicTimeZone   uint16 = 0x0200
	ECFSupportHeartbeatPDU      uint16 = 0x0400
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
		Version:                rdpVersion5Plus,
		DesktopWidth:           desktopWidth,
		DesktopHeight:          desktopHeight,
		ColorDepth:             0xCA01,     // RNS_UD_COLOR_8BPP
		SASSequence:            0xAA03,     // RNS_UD_SAS_DEL
		KeyboardLayout:         0x00000409, // US
		ClientBuild:            0xece,
		ClientName:             [32]byte{},
		KeyboardType:           keyboardTypeIBM101or102Keys,
		KeyboardSubType:        0x00000000,
		KeyboardFunctionKey:    12,
		ImeFileName:            [64]byte{},
		PostBeta2ColorDepth:    0xCA03, // RNS_UD_COLOR_16BPP_565
		ClientProductId:        0x0001,
		SerialNumber:           0x00000000,
		HighColorDepth:         0x0010, // HIGH_COLOR_16BPP
		SupportedColorDepths:   0x0002, // RNS_UD_16BPP_SUPPORT
		EarlyCapabilityFlags:   ECFSupportErrInfoPDU,
		ClientDigProductId:     [64]byte{},
		ConnectionType:         0x00,
		Pad1octet:              0x00,
		ServerSelectedProtocol: selectedProtocol,
	}

	copy(data.ClientName[:], utf16.Encode(projectName))

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

type ClientUserDataSet struct {
	ClientCoreData     *ClientCoreData
	ClientSecurityData *ClientSecurityData
	ClientNetworkData  *ClientNetworkData
	ClientClusterData  *ClientClusterData
}

func NewClientUserDataSet(selectedProtocol uint32,
	desktopWidth, desktopHeight uint16,
	channelNames []string) *ClientUserDataSet {
	return &ClientUserDataSet{
		ClientCoreData:     newClientCoreData(selectedProtocol, desktopWidth, desktopHeight),
		ClientSecurityData: newClientSecurityData(),
		ClientNetworkData:  newClientNetworkData(channelNames),
	}
}

func (data ClientCoreData) Serialize() []byte {
	const dataLen uint16 = 216

	buf := new(bytes.Buffer)

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
	buf := new(bytes.Buffer)

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

	chBuf := new(bytes.Buffer)

	for _, channelDef := range data.ChannelDefArray {
		chBuf.Write(channelDef.Serialize())
	}

	buf := new(bytes.Buffer)

	binary.Write(buf, binary.LittleEndian, uint16(0xC003))                // header type CS_NET
	binary.Write(buf, binary.LittleEndian, uint16(headerLen+chBuf.Len())) // packet size

	binary.Write(buf, binary.LittleEndian, data.ChannelCount)

	buf.Write(chBuf.Bytes())

	return buf.Bytes()
}

func (d ClientClusterData) Serialize() []byte {
	const dataLen uint16 = 12

	buf := new(bytes.Buffer)

	binary.Write(buf, binary.LittleEndian, uint16(0xC004)) // header type CS_CLUSTER
	binary.Write(buf, binary.LittleEndian, dataLen)        // packet size

	binary.Write(buf, binary.LittleEndian, d.Flags)
	binary.Write(buf, binary.LittleEndian, d.RedirectedSessionID)

	return buf.Bytes()
}

func (ud ClientUserDataSet) Serialize() []byte {
	buf := new(bytes.Buffer)

	buf.Write(ud.ClientCoreData.Serialize())

	if ud.ClientClusterData != nil {
		buf.Write(ud.ClientClusterData.Serialize())
	}

	buf.Write(ud.ClientSecurityData.Serialize())
	buf.Write(ud.ClientNetworkData.Serialize())

	return buf.Bytes()
}

type ServerCoreData struct {
	Version                  uint32
	ClientRequestedProtocols uint32
	EarlyCapabilityFlags     uint32

	DataLen uint16
}

func (d *ServerCoreData) Deserialize(wire io.Reader) error {
	var err error

	if err = binary.Read(wire, binary.LittleEndian, &d.Version); err != nil {
		return err
	}

	if d.DataLen == 4 {
		return nil
	}

	if err = binary.Read(wire, binary.LittleEndian, &d.ClientRequestedProtocols); err != nil {
		return err
	}

	if d.DataLen == 8 {
		return nil
	}

	if err = binary.Read(wire, binary.LittleEndian, &d.EarlyCapabilityFlags); err != nil {
		return err
	}

	return nil
}

type RSAPublicKey struct {
	Magic   uint32
	KeyLen  uint32
	BitLen  uint32
	DataLen uint32
	PubExp  uint32
	Modulus []byte
}

func (k *RSAPublicKey) Deserialize(wire io.Reader) error {
	var err error

	if err = binary.Read(wire, binary.LittleEndian, &k.Magic); err != nil {
		return err
	}

	if err = binary.Read(wire, binary.LittleEndian, &k.KeyLen); err != nil {
		return err
	}

	if err = binary.Read(wire, binary.LittleEndian, &k.BitLen); err != nil {
		return err
	}

	if err = binary.Read(wire, binary.LittleEndian, &k.DataLen); err != nil {
		return err
	}

	if err = binary.Read(wire, binary.LittleEndian, &k.PubExp); err != nil {
		return err
	}

	k.Modulus = make([]byte, k.KeyLen)

	if _, err = wire.Read(k.Modulus); err != nil {
		return err
	}

	return nil
}

type ServerProprietaryCertificate struct {
	DwSigAlgId        uint32
	DwKeyAlgId        uint32
	PublicKeyBlobType uint16
	PublicKeyBlobLen  uint16
	PublicKeyBlob     RSAPublicKey
	SignatureBlobType uint16
	SignatureBlobLen  uint16
	SignatureBlob     []byte
}

func (c *ServerProprietaryCertificate) Deserialize(wire io.Reader) error {
	var err error

	if err = binary.Read(wire, binary.LittleEndian, &c.DwSigAlgId); err != nil {
		return err
	}

	if err = binary.Read(wire, binary.LittleEndian, &c.DwKeyAlgId); err != nil {
		return err
	}

	if err = binary.Read(wire, binary.LittleEndian, &c.PublicKeyBlobType); err != nil {
		return err
	}

	if err = binary.Read(wire, binary.LittleEndian, &c.PublicKeyBlobLen); err != nil {
		return err
	}

	if err = c.PublicKeyBlob.Deserialize(wire); err != nil {
		return err
	}

	if err = binary.Read(wire, binary.LittleEndian, &c.SignatureBlobType); err != nil {
		return err
	}

	if err = binary.Read(wire, binary.LittleEndian, &c.SignatureBlobLen); err != nil {
		return err
	}

	c.SignatureBlob = make([]byte, c.SignatureBlobLen)

	if _, err = wire.Read(c.SignatureBlob); err != nil {
		return err
	}

	return nil
}

type ServerCertificate struct {
	DwVersion       uint32
	ProprietaryCert *ServerProprietaryCertificate
	X509Cert        []byte

	ServerCertLen uint32
}

func (c *ServerCertificate) Deserialize(wire io.Reader) error {
	var err error

	if err = binary.Read(wire, binary.LittleEndian, &c.DwVersion); err != nil {
		return err
	}

	if c.DwVersion&0x00000001 == 0x00000001 {
		c.ProprietaryCert = &ServerProprietaryCertificate{}

		return c.ProprietaryCert.Deserialize(wire)
	}

	c.X509Cert = make([]byte, c.ServerCertLen-4)

	if _, err = wire.Read(c.X509Cert); err != nil {
		return err
	}

	return nil
}

type ServerSecurityData struct {
	EncryptionMethod  uint32
	EncryptionLevel   uint32
	ServerRandomLen   uint32
	ServerCertLen     uint32
	ServerRandom      []byte
	ServerCertificate *ServerCertificate
}

func (d *ServerSecurityData) Deserialize(wire io.Reader) error {
	var err error

	if err = binary.Read(wire, binary.LittleEndian, &d.EncryptionMethod); err != nil {
		return err
	}

	if err = binary.Read(wire, binary.LittleEndian, &d.EncryptionLevel); err != nil {
		return err
	}

	if d.EncryptionMethod == 0 && d.EncryptionLevel == 0 { // ENCRYPTION_METHOD_NONE and ENCRYPTION_LEVEL_NONE
		return nil
	}

	if err = binary.Read(wire, binary.LittleEndian, &d.ServerRandomLen); err != nil {
		return err
	}

	if err = binary.Read(wire, binary.LittleEndian, &d.ServerCertLen); err != nil {
		return err
	}

	d.ServerRandom = make([]byte, d.ServerRandomLen)

	if _, err = wire.Read(d.ServerRandom); err != nil {
		return err
	}

	if d.ServerCertLen > 0 {
		d.ServerCertificate = &ServerCertificate{
			ServerCertLen: d.ServerCertLen,
		}

		return d.ServerCertificate.Deserialize(wire)
	}

	return nil
}

type ServerNetworkData struct {
	MCSChannelId   uint16
	ChannelCount   uint16
	ChannelIdArray []uint16
}

func (d *ServerNetworkData) Deserialize(wire io.Reader) error {
	var err error

	if err = binary.Read(wire, binary.LittleEndian, &d.MCSChannelId); err != nil {
		return err
	}

	if err = binary.Read(wire, binary.LittleEndian, &d.ChannelCount); err != nil {
		return err
	}

	if d.ChannelCount == 0 {
		return nil
	}

	d.ChannelIdArray = make([]uint16, d.ChannelCount)

	if err = binary.Read(wire, binary.LittleEndian, &d.ChannelIdArray); err != nil {
		return err
	}

	if d.ChannelCount%2 == 0 {
		return nil
	}

	padding := make([]byte, 2)

	if _, err = wire.Read(padding); err != nil {
		return err
	}

	return nil
}

type ServerMessageChannelData struct {
	MCSChannelID uint16
}

type ServerMultitransportChannelData struct {
	Flags uint32
}

type ServerUserData struct {
	ServerCoreData                  *ServerCoreData
	ServerNetworkData               *ServerNetworkData
	ServerSecurityData              *ServerSecurityData
	ServerMessageChannelData        *ServerMessageChannelData
	ServerMultitransportChannelData *ServerMultitransportChannelData
}

func (ud *ServerUserData) Deserialize(wire io.Reader) error {
	var (
		dataType uint16
		dataLen  uint16
		err      error
	)

	for {
		err = binary.Read(wire, binary.LittleEndian, &dataType)
		switch err {
		case nil: // pass
		case io.EOF:
			return nil
		default:
			return err
		}

		err = binary.Read(wire, binary.LittleEndian, &dataLen)
		if err != nil {
			return err
		}

		dataLen -= 4 // exclude User Data Header

		switch dataType {
		case 0x0C01:
			ud.ServerCoreData = &ServerCoreData{DataLen: dataLen}

			if err = ud.ServerCoreData.Deserialize(wire); err != nil {
				return err
			}
		case 0x0C02:
			ud.ServerSecurityData = &ServerSecurityData{}

			if err = ud.ServerSecurityData.Deserialize(wire); err != nil {
				return err
			}
		case 0x0C03:
			ud.ServerNetworkData = &ServerNetworkData{}

			if err = ud.ServerNetworkData.Deserialize(wire); err != nil {
				return err
			}
		case 0x0C04:
			ud.ServerMessageChannelData = &ServerMessageChannelData{}

			if err = binary.Read(wire, binary.LittleEndian, &ud.ServerMessageChannelData.MCSChannelID); err != nil {
				return err
			}
		case 0x0C08:
			ud.ServerMultitransportChannelData = &ServerMultitransportChannelData{}

			if err = binary.Read(wire, binary.LittleEndian, &ud.ServerMultitransportChannelData.Flags); err != nil {
				return err
			}
		default:
			return errors.New("unknown header type")
		}
	}
}
