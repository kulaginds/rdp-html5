package rdp

import (
	"io"

	"github.com/kulaginds/web-rdp-solution/internal/pkg/rdp/fastpath"
	"github.com/kulaginds/web-rdp-solution/internal/pkg/rdp/tpkt"
)

type tpktLayer interface {
	Receive() (io.Reader, error)
	Send(pduData []byte) error
	StartHandleFastpath()
	ReceiveProtocol() (tpkt.ProtocolCode, error)
}

type x224Layer interface {
	Connect() error
	Close() error
	Send(pduData []byte) error
	Receive() (io.Reader, error)
}

type mcsLayer interface {
	Connect(selectedProtocol uint32, desktopWidth, desktopHeight uint16, channelNames []string) error
	Disconnect() error
	ErectDomain() error
	AttachUser() error
	JoinChannels() error
	Send(channelName string, pduData []byte) error
	Receive() (string, io.Reader, error)
	UserId() uint16
}

type fastPath interface {
	Send(pdu *fastpath.InputEventPDU) error
	Receive(fpOutputHeader uint8) (*fastpath.UpdatePDU, error)
}
