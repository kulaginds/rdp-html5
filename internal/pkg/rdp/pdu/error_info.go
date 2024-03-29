package pdu

import (
	"encoding/binary"
	"fmt"
	"io"
)

type ErrorInfoPDUData struct {
	ErrorInfo uint32
}

func (pdu *ErrorInfoPDUData) Deserialize(wire io.Reader) error {
	return binary.Read(wire, binary.LittleEndian, &pdu.ErrorInfo)
}

var errorInfoMap = map[uint32]string{
	0x00000000: "ERRINFO_NONE",
	0x00000001: "ERRINFO_RPC_INITIATED_DISCONNECT",
	0x00000002: "ERRINFO_RPC_INITIATED_LOGOFF",
	0x00000003: "ERRINFO_IDLE_TIMEOUT",
	0x00000004: "ERRINFO_LOGON_TIMEOUT",
	0x00000005: "ERRINFO_DISCONNECTED_BY_OTHERCONNECTION",
	0x00000006: "ERRINFO_OUT_OF_MEMORY",
	0x00000007: "ERRINFO_SERVER_DENIED_CONNECTION",
	0x00000009: "ERRINFO_SERVER_INSUFFICIENT_PRIVILEGES",
	0x0000000A: "ERRINFO_SERVER_FRESH_CREDENTIALS_REQUIRED",
	0x0000000B: "ERRINFO_RPC_INITIATED_DISCONNECT_BYUSER",
	0x0000000C: "ERRINFO_LOGOFF_BY_USER",
	0x0000000F: "ERRINFO_CLOSE_STACK_ON_DRIVER_NOT_READY",
	0x00000010: "ERRINFO_SERVER_DWM_CRASH",
	0x00000011: "ERRINFO_CLOSE_STACK_ON_DRIVER_FAILURE",
	0x00000012: "ERRINFO_CLOSE_STACK_ON_DRIVER_IFACE_FAILURE",
	0x00000017: "ERRINFO_SERVER_WINLOGON_CRASH",
	0x00000018: "ERRINFO_SERVER_CSRSS_CRASH",
	0x00000019: "ERRINFO_SERVER_SHUTDOWN",
	0x0000001A: "ERRINFO_SERVER_REBOOT",
	0x00000100: "ERRINFO_LICENSE_INTERNAL",
	0x00000101: "ERRINFO_LICENSE_NO_LICENSE_SERVER",
	0x00000102: "ERRINFO_LICENSE_NO_LICENSE",
	0x00000103: "ERRINFO_LICENSE_BAD_CLIENT_MSG",
	0x00000104: "ERRINFO_LICENSE_HWID_DOESNT_MATCH_LICENSE",
	0x00000105: "ERRINFO_LICENSE_BAD_CLIENT_LICENSE",
	0x00000106: "ERRINFO_LICENSE_CANT_FINISH_PROTOCOL",
	0x00000107: "ERRINFO_LICENSE_CLIENT_ENDED_PROTOCOL",
	0x00000108: "ERRINFO_LICENSE_BAD_CLIENT_ENCRYPTION",
	0x00000109: "ERRINFO_LICENSE_CANT_UPGRADE_LICENSE",
	0x0000010A: "ERRINFO_LICENSE_NO_REMOTE_CONNECTIONS",
	0x00000400: "ERRINFO_CB_DESTINATION_NOT_FOUND",
	0x00000402: "ERRINFO_CB_LOADING_DESTINATION",
	0x00000404: "ERRINFO_CB_REDIRECTING_TO_DESTINATION",
	0x00000405: "ERRINFO_CB_SESSION_ONLINE_VM_WAKE",
	0x00000406: "ERRINFO_CB_SESSION_ONLINE_VM_BOOT",
	0x00000407: "ERRINFO_CB_SESSION_ONLINE_VM_NO_DNS",
	0x00000408: "ERRINFO_CB_DESTINATION_POOL_NOT_FREE",
	0x00000409: "ERRINFO_CB_CONNECTION_CANCELLED",
	0x00000410: "ERRINFO_CB_CONNECTION_ERROR_INVALID_SETTINGS",
	0x00000411: "ERRINFO_CB_SESSION_ONLINE_VM_BOOT_TIMEOUT",
	0x00000412: "ERRINFO_CB_SESSION_ONLINE_VM_SESSMON_FAILED",
	0x000010C9: "ERRINFO_UNKNOWNPDUTYPE2",
	0x000010CA: "ERRINFO_UNKNOWNPDUTYPE",
	0x000010CB: "ERRINFO_DATAPDUSEQUENCE",
	0x000010CD: "ERRINFO_CONTROLPDUSEQUENCE",
	0x000010CE: "ERRINFO_INVALIDCONTROLPDUACTION",
	0x000010CF: "ERRINFO_INVALIDINPUTPDUTYPE",
	0x000010D0: "ERRINFO_INVALIDINPUTPDUMOUSE",
	0x000010D1: "ERRINFO_INVALIDREFRESHRECTPDU",
	0x000010D2: "ERRINFO_CREATEUSERDATAFAILED",
	0x000010D3: "ERRINFO_CONNECTFAILED",
	0x000010D4: "ERRINFO_CONFIRMACTIVEWRONGSHAREID",
	0x000010D5: "ERRINFO_CONFIRMACTIVEWRONGORIGINATOR",
	0x000010DA: "ERRINFO_PERSISTENTKEYPDUBADLENGTH",
	0x000010DB: "ERRINFO_PERSISTENTKEYPDUILLEGALFIRST",
	0x000010DC: "ERRINFO_PERSISTENTKEYPDUTOOMANYTOTALKEYS",
	0x000010DD: "ERRINFO_PERSISTENTKEYPDUTOOMANYCACHEKEYS",
	0x000010DE: "ERRINFO_INPUTPDUBADLENGTH",
	0x000010DF: "ERRINFO_BITMAPCACHEERRORPDUBADLENGTH",
	0x000010E0: "ERRINFO_SECURITYDATATOOSHORT",
	0x000010E1: "ERRINFO_VCHANNELDATATOOSHORT",
	0x000010E2: "ERRINFO_SHAREDATATOOSHORT",
	0x000010E3: "ERRINFO_BADSUPRESSOUTPUTPDU",
	0x000010E5: "ERRINFO_CONFIRMACTIVEPDUTOOSHORT",
	0x000010E7: "ERRINFO_CAPABILITYSETTOOSMALL",
	0x000010E8: "ERRINFO_CAPABILITYSETTOOLARGE",
	0x000010E9: "ERRINFO_NOCURSORCACHE",
	0x000010EA: "ERRINFO_BADCAPABILITIES",
	0x000010EC: "ERRINFO_VIRTUALCHANNELDECOMPRESSIONERR",
	0x000010ED: "ERRINFO_INVALIDVCCOMPRESSIONTYPE",
	0x000010EF: "ERRINFO_INVALIDCHANNELID",
	0x000010F0: "ERRINFO_VCHANNELSTOOMANY",
	0x000010F3: "ERRINFO_REMOTEAPPSNOTENABLED",
	0x000010F4: "ERRINFO_CACHECAPNOTSET",
	0x000010F5: "ERRINFO_BITMAPCACHEERRORPDUBADLENGTH2",
	0x000010F6: "ERRINFO_OFFSCRCACHEERRORPDUBADLENGTH",
	0x000010F7: "ERRINFO_DNGCACHEERRORPDUBADLENGTH",
	0x000010F8: "ERRINFO_GDIPLUSPDUBADLENGTH",
	0x00001111: "ERRINFO_SECURITYDATATOOSHORT2",
	0x00001112: "ERRINFO_SECURITYDATATOOSHORT3",
	0x00001113: "ERRINFO_SECURITYDATATOOSHORT4",
	0x00001114: "ERRINFO_SECURITYDATATOOSHORT5",
	0x00001115: "ERRINFO_SECURITYDATATOOSHORT6",
	0x00001116: "ERRINFO_SECURITYDATATOOSHORT7",
	0x00001117: "ERRINFO_SECURITYDATATOOSHORT8",
	0x00001118: "ERRINFO_SECURITYDATATOOSHORT9",
	0x00001119: "ERRINFO_SECURITYDATATOOSHORT10",
	0x0000111A: "ERRINFO_SECURITYDATATOOSHORT11",
	0x0000111B: "ERRINFO_SECURITYDATATOOSHORT12",
	0x0000111C: "ERRINFO_SECURITYDATATOOSHORT13",
	0x0000111D: "ERRINFO_SECURITYDATATOOSHORT14",
	0x0000111E: "ERRINFO_SECURITYDATATOOSHORT15",
	0x0000111F: "ERRINFO_SECURITYDATATOOSHORT16",
	0x00001120: "ERRINFO_SECURITYDATATOOSHORT17",
	0x00001121: "ERRINFO_SECURITYDATATOOSHORT18",
	0x00001122: "ERRINFO_SECURITYDATATOOSHORT19",
	0x00001123: "ERRINFO_SECURITYDATATOOSHORT20",
	0x00001124: "ERRINFO_SECURITYDATATOOSHORT21",
	0x00001125: "ERRINFO_SECURITYDATATOOSHORT22",
	0x00001126: "ERRINFO_SECURITYDATATOOSHORT23",
	0x00001129: "ERRINFO_BADMONITORDATA",
	0x0000112A: "ERRINFO_VCDECOMPRESSEDREASSEMBLEFAILED",
	0x0000112B: "ERRINFO_VCDATATOOLONG",
	0x0000112C: "ERRINFO_BAD_FRAME_ACK_DATA",
	0x0000112D: "ERRINFO_GRAPHICSMODENOTSUPPORTED",
	0x0000112E: "ERRINFO_GRAPHICSSUBSYSTEMRESETFAILED",
	0x0000112F: "ERRINFO_GRAPHICSSUBSYSTEMFAILED",
	0x00001130: "ERRINFO_TIMEZONEKEYNAMELENGTHTOOSHORT",
	0x00001131: "ERRINFO_TIMEZONEKEYNAMELENGTHTOOLONG",
	0x00001132: "ERRINFO_DYNAMICDSTDISABLEDFIELDMISSING",
	0x00001133: "ERRINFO_VCDECODINGERROR",
	0x00001134: "ERRINFO_VIRTUALDESKTOPTOOLARGE",
	0x00001135: "ERRINFO_MONITORGEOMETRYVALIDATIONFAILED",
	0x00001136: "ERRINFO_INVALIDMONITORCOUNT",
	0x00001191: "ERRINFO_UPDATESESSIONKEYFAILED",
	0x00001192: "ERRINFO_DECRYPTFAILED",
	0x00001193: "ERRINFO_ENCRYPTFAILED",
	0x00001194: "ERRINFO_ENCPKGMISMATCH",
	0x00001195: "ERRINFO_DECRYPTFAILED2",
}

func (pdu *ErrorInfoPDUData) String() string {
	code, ok := errorInfoMap[pdu.ErrorInfo]
	if ok {
		return code
	}

	return fmt.Sprintf("unknown code: %d", pdu.ErrorInfo)
}
