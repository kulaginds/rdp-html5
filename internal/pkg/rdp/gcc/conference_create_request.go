package gcc

import (
	"bytes"

	"github.com/kulaginds/web-rdp-solution/internal/pkg/rdp/per"
)

type ConferenceCreateRequest struct {
	UserData []byte
}

func NewConferenceCreateRequest(userData []byte) *ConferenceCreateRequest {
	return &ConferenceCreateRequest{
		UserData: userData,
	}
}

func (r *ConferenceCreateRequest) Serialize() []byte {
	buf := new(bytes.Buffer)

	per.WriteChoice(0, buf)
	per.WriteObjectIdentifier(t124_02_98_oid, buf)
	per.WriteLength(uint16(14+len(r.UserData)), buf)

	per.WriteChoice(0, buf)
	per.WriteSelection(0x08, buf)

	per.WriteNumericString("1", 1, buf)
	per.WritePadding(1, buf)
	per.WriteNumberOfSet(1, buf)
	per.WriteChoice(0xc0, buf)
	per.WriteOctetStream(h221CSKey, 4, buf)
	per.WriteOctetStream(string(r.UserData), 0, buf)

	return buf.Bytes()
}
