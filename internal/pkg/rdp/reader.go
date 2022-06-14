package rdp

import (
	"context"
	"io"
	"log"
)

func (c *client) ReaderLoop(ctx context.Context) error {
	var (
		dataPDU *DataPDU
		wire    io.Reader
		err     error
	)

	for {
		select {
		case <-ctx.Done():
			return nil
		default: // pass
		}

		_, wire, err = c.mcsLayer.Receive()
		if err != nil {
			return err
		}

		dataPDU = &DataPDU{}
		if err = dataPDU.Deserialize(wire); err != nil {
			return err
		}

		pduType2 := dataPDU.ShareDataHeader.PDUType2

		switch {
		case pduType2.IsSynchronize():
			log.Println("server synchronize")
		case pduType2.IsControl():
			log.Printf("server control %d\n", dataPDU.ControlPDUData.Action)
		case pduType2.IsErrorInfo():
			log.Printf("server error info: %d\n", dataPDU.ErrorInfoPDUData.ErrorInfo)
		default:
			log.Printf("unknown server message with pduType2 = %d\n", pduType2)
		}
	}
}
