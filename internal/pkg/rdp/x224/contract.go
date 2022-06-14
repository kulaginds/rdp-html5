package x224

import (
	"io"
)

type tpktConn interface {
	Receive() (io.Reader, error)
	Send(pduData []byte) error
	Close() error
}
