package rdp

import (
	"bytes"
	"encoding/binary"
	"strings"

	"github.com/kulaginds/web-rdp-solution/internal/pkg/rdp/headers"
	"github.com/kulaginds/web-rdp-solution/internal/pkg/rdp/utf16"
)

type SystemTime struct {
	Year         uint16
	Month        uint16
	DayOfWeek    uint16
	Day          uint16
	Hour         uint16
	Minute       uint16
	Second       uint16
	Milliseconds uint16
}

func (t *SystemTime) Serialize() []byte {
	buf := &bytes.Buffer{}

	binary.Write(buf, binary.LittleEndian, t.Year)
	binary.Write(buf, binary.LittleEndian, t.Month)
	binary.Write(buf, binary.LittleEndian, t.DayOfWeek)
	binary.Write(buf, binary.LittleEndian, t.Day)
	binary.Write(buf, binary.LittleEndian, t.Hour)
	binary.Write(buf, binary.LittleEndian, t.Minute)
	binary.Write(buf, binary.LittleEndian, t.Second)
	binary.Write(buf, binary.LittleEndian, t.Milliseconds)

	return buf.Bytes()
}

type TimeZoneInformation struct {
	Bias         uint32
	StandardName [64]byte
	StandardDate SystemTime
	StandardBias uint32
	DaylightName [64]byte
	DaylightDate SystemTime
	DaylightBias uint32
}

func (i *TimeZoneInformation) Serialize() []byte {
	buf := &bytes.Buffer{}

	binary.Write(buf, binary.LittleEndian, i.Bias)
	binary.Write(buf, binary.LittleEndian, i.StandardName)

	buf.Write(i.StandardDate.Serialize())

	binary.Write(buf, binary.LittleEndian, i.StandardBias)
	binary.Write(buf, binary.LittleEndian, i.DaylightName)

	buf.Write(i.DaylightDate.Serialize())

	binary.Write(buf, binary.LittleEndian, i.DaylightBias)

	return buf.Bytes()
}

type AddressFamily uint16

const (
	// AddressFamilyINET AF_INET IPv4
	AddressFamilyINET AddressFamily = 0x00002

	// AddressFamilyINET6 AF_INET6 IPv6
	AddressFamilyINET6 AddressFamily = 0x0017
)

type ExtendedInfoPacket struct {
	PerformanceFlags uint32
}

func (p *ExtendedInfoPacket) Serialize() []byte {
	buf := &bytes.Buffer{}

	binary.Write(buf, binary.LittleEndian, uint16(0x0002)) // ClientAddressFamily = AF_INET
	binary.Write(buf, binary.LittleEndian, uint16(2))      // cbClientAddress
	buf.Write([]byte{0, 0})                                // ClientAddress
	binary.Write(buf, binary.LittleEndian, uint16(2))      // cbClientDir
	buf.Write([]byte{0, 0})                                // ClientDir
	buf.Write(make([]byte, 172))                           // ClientTimeZone
	binary.Write(buf, binary.LittleEndian, uint32(0))      // ClientSessionId
	binary.Write(buf, binary.LittleEndian, p.PerformanceFlags)

	return buf.Bytes()
}

type ClientInfoPacket struct {
	CodePage       uint32
	Flags          InfoFlag
	Domain         string
	Username       string
	Password       string
	AlternateShell string
	WorkingDir     string
	ExtraInfo      ExtendedInfoPacket
}

func (p *ClientInfoPacket) Serialize() []byte {
	cbDomain := uint16(0)
	cbUserName := uint16(0)
	cbPassword := uint16(0)
	cbAlternateShell := uint16(0)
	cbWorkingDir := uint16(0)

	domain := []byte{0x00, 0x00}
	username := []byte{0x00, 0x00}
	password := []byte{0x00, 0x00}
	alternateShell := []byte{0x00, 0x00}
	workingDir := []byte{0x00, 0x00}

	if len(p.Domain) > 0 {
		domain = utf16.Encode(strings.Trim(p.Domain, " ") + "\x00")
		cbDomain = uint16(len(domain) - 2)
	}

	if len(p.Username) > 0 {
		username = utf16.Encode(strings.Trim(p.Username, " ") + "\x00")
		cbUserName = uint16(len(username) - 2)
	}

	if len(p.Password) > 0 {
		password = utf16.Encode(strings.Trim(p.Password, " ") + "\x00")
		cbPassword = uint16(len(password) - 2)
	}

	if len(p.AlternateShell) > 0 {
		alternateShell = utf16.Encode(strings.Trim(p.AlternateShell, " ") + "\x00")
		cbAlternateShell = uint16(len(alternateShell) - 2)
	}

	if len(p.WorkingDir) > 0 {
		workingDir = utf16.Encode(strings.Trim(p.WorkingDir, " ") + "\x00")
		cbWorkingDir = uint16(len(workingDir) - 2)
	}

	buf := &bytes.Buffer{}

	binary.Write(buf, binary.LittleEndian, p.CodePage)
	binary.Write(buf, binary.LittleEndian, uint32(p.Flags))
	binary.Write(buf, binary.LittleEndian, cbDomain)
	binary.Write(buf, binary.LittleEndian, cbUserName)
	binary.Write(buf, binary.LittleEndian, cbPassword)
	binary.Write(buf, binary.LittleEndian, cbAlternateShell)
	binary.Write(buf, binary.LittleEndian, cbWorkingDir)

	buf.Write(domain)
	buf.Write(username)
	buf.Write(password)
	buf.Write(alternateShell)
	buf.Write(workingDir)

	buf.Write(p.ExtraInfo.Serialize())

	return buf.Bytes()
}

type ClientInfoPDU struct {
	InfoPacket ClientInfoPacket
}

type InfoFlag uint32

const (
	// InfoFlagMouse INFO_MOUSE
	InfoFlagMouse InfoFlag = 0x00000001

	// InfoFlagDisableCtrlAltDel INFO_DISABLECTRLALTDEL
	InfoFlagDisableCtrlAltDel InfoFlag = 0x00000002

	// InfoFlagAutoLogon INFO_AUTOLOGON
	InfoFlagAutoLogon InfoFlag = 0x00000008

	// InfoFlagUnicode INFO_UNICODE
	InfoFlagUnicode InfoFlag = 0x00000010

	// InfoFlagMaximizeShell INFO_MAXIMIZESHELL
	InfoFlagMaximizeShell InfoFlag = 0x00000020

	// InfoFlagLogonNotify INFO_LOGONNOTIFY
	InfoFlagLogonNotify InfoFlag = 0x00000040

	// InfoFlagCompression INFO_COMPRESSION
	InfoFlagCompression InfoFlag = 0x00000080

	// InfoFlagEnableWindowsKey INFO_ENABLEWINDOWSKEY
	InfoFlagEnableWindowsKey InfoFlag = 0x00000100

	// InfoFlagRemoteConsoleAudio INFO_REMOTECONSOLEAUDIO
	InfoFlagRemoteConsoleAudio InfoFlag = 0x00002000

	// InfoFlagForceEncryptedCSPDU INFO_FORCE_ENCRYPTED_CS_PDU
	InfoFlagForceEncryptedCSPDU InfoFlag = 0x00004000

	// InfoFlagRail INFO_RAIL
	InfoFlagRail InfoFlag = 0x00008000

	// InfoFlagLogonErrors INFO_LOGONERRORS
	InfoFlagLogonErrors InfoFlag = 0x00010000

	// InfoFlagMouseHasWheel INFO_MOUSE_HAS_WHEEL
	InfoFlagMouseHasWheel InfoFlag = 0x00020000

	// InfoFlagPasswordIsSCPIN INFO_PASSWORD_IS_SC_PIN
	InfoFlagPasswordIsSCPIN InfoFlag = 0x00040000

	// InfoFlagNoAudioPlayback INFO_NOAUDIOPLAYBACK
	InfoFlagNoAudioPlayback InfoFlag = 0x00080000

	// InfoFlagUsingSavedCreds INFO_USING_SAVED_CREDS
	InfoFlagUsingSavedCreds InfoFlag = 0x00100000

	// InfoFlagAudioCapture INFO_AUDIOCAPTURE
	InfoFlagAudioCapture InfoFlag = 0x00200000

	// InfoFlagVideoDisable INFO_VIDEO_DISABLE
	InfoFlagVideoDisable InfoFlag = 0x00400000

	// InfoFlagHiDefRailSupported INFO_HIDEF_RAIL_SUPPORTED
	InfoFlagHiDefRailSupported InfoFlag = 0x02000000
)

const (
	CompressionTypeMask  uint32 = 0x00001E00
	CompressionType8K    uint32 = 0x0
	CompressionType64K   uint32 = 0x1
	CompressionTypeRDP6  uint32 = 0x2
	CompressionTypeRDP61 uint32 = 0x3
)

func NewClientInfoPDU(domain, username, password string) *ClientInfoPDU {
	return &ClientInfoPDU{
		InfoPacket: ClientInfoPacket{
			Flags:    InfoFlagMouse | InfoFlagUnicode | InfoFlagAutoLogon | InfoFlagDisableCtrlAltDel | InfoFlagEnableWindowsKey,
			Domain:   domain,
			Username: username,
			Password: password,
			ExtraInfo: ExtendedInfoPacket{
				// PERF_DISABLE_WALLPAPER, PERF_DISABLE_FULLWINDOWDRAG, PERF_DISABLE_MENUANIMATIONS,
				// PERF_DISABLE_THEMING, PERF_DISABLE_CURSOR_SHADOW, PERF_DISABLE_CURSORSETTINGS
				PerformanceFlags: 0x00000001 | 0x00000002 | 0x00000004 | 0x00000008 | 0x00000020 | 0x00000040,
			},
		},
	}
}

func (pdu *ClientInfoPDU) Serialize() []byte {
	return headers.WrapSecurityFlag(
		0x0040, // SEC_INFO_PKT
		pdu.InfoPacket.Serialize(),
	)
}
