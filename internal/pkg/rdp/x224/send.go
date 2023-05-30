package x224

import "bytes"

func (p *Protocol) Send(pduData []byte) error {
	buf := bytes.NewBuffer(make([]byte, 0, 3+len(pduData)))

	buf.Write([]byte{
		0x02, // packet size
		0xF0, // message type TPDU_DATA
		0x80, // EOT flag is up, which indicates that current TPDU is the last data unit of a complete TPDU sequence
	})
	buf.Write(pduData)

	return p.tpktConn.Send(buf.Bytes())
}
