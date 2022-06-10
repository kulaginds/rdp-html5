package gcc

import (
	"encoding/binary"
	"errors"
	"io"

	"github.com/kulaginds/web-rdp-solution/internal/pkg/rdp/per"
)

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

type ConferenceCreateResponseUserData struct {
	ServerCoreData                  *ServerCoreData
	ServerNetworkData               *ServerNetworkData
	ServerSecurityData              *ServerSecurityData
	ServerMessageChannelData        *ServerMessageChannelData
	ServerMultitransportChannelData *ServerMultitransportChannelData
}

func (ud *ConferenceCreateResponseUserData) Deserialize(wire io.Reader) error {
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

type ConferenceCreateResponse struct {
	UserData ConferenceCreateResponseUserData
}

func (r *ConferenceCreateResponse) Deserialize(wire io.Reader) error {
	_, err := per.ReadChoice(wire)
	if err != nil {
		return err
	}

	var objectIdentifier bool

	objectIdentifier, err = per.ReadObjectIdentifier(t124_02_98_oid, wire)
	if err != nil {
		return err
	}

	if !objectIdentifier {
		return errors.New("bad object identifier t124")
	}

	_, err = per.ReadLength(wire)
	if err != nil {
		return err
	}

	_, err = per.ReadChoice(wire)
	if err != nil {
		return err
	}

	_, err = per.ReadInteger16(1001, wire)
	if err != nil {
		return err
	}

	_, err = per.ReadInteger(wire)
	if err != nil {
		return err
	}

	_, err = per.ReadEnumerates(wire)
	if err != nil {
		return err
	}

	_, err = per.ReadNumberOfSet(wire)
	if err != nil {
		return err
	}

	_, err = per.ReadChoice(wire)
	if err != nil {
		return err
	}

	var octetStream bool

	octetStream, err = per.ReadOctetStream([]byte(h221SCKey), 4, wire)
	if err != nil {
		return err
	}

	if !octetStream {
		return errors.New("bad H221 SC_KEY")
	}

	_, err = per.ReadLength(wire)
	if err != nil {
		return err
	}

	return r.UserData.Deserialize(wire)
}
