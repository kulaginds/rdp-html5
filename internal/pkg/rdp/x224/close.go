package x224

func (p *protocol) Close() error {
	if !p.connected {
		return nil
	}

	if err := p.tpktConn.Close(); err != nil {
		return err
	}

	p.connected = false

	return nil
}
