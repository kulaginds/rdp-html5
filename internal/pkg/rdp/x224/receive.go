package x224

import (
	"io"
)

func (p *Protocol) Receive() (io.Reader, error) {
	wire, err := p.tpktConn.Receive()
	if err != nil {
		return nil, err
	}

	var resp Data

	if err = resp.Deserialize(wire); err != nil {
		return nil, err
	}

	return wire, nil
}
