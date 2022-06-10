package rdp

import "io"

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
