package tpkt

func (p *protocol) Close() error {
	return p.conn.Close()
}
