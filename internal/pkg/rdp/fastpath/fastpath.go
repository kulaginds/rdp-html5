package fastpath

import (
	"io"
	"sync"
)

type impl struct {
	conn io.ReadWriter

	updatePDUDataPool sync.Pool
}

func New(conn io.ReadWriter) *impl {
	return &impl{
		conn: conn,

		updatePDUDataPool: sync.Pool{
			New: func() interface{} { return make([]byte, 16*1024) },
		},
	}
}
