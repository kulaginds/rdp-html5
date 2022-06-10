package gcc

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

var (
	t124_02_98_oid = [6]byte{0, 0, 20, 124, 0, 1}
	h221CSKey      = "Duca"
	h221SCKey      = "McDn"
)
