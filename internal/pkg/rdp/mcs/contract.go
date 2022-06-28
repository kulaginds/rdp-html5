package mcs

import "io"

type x224Conn interface {
	Receive() (io.Reader, error)
	Send(pduData []byte) error
}
