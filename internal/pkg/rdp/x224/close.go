package x224

func (p *protocol) Close() error {
	if !p.connected {
		return nil
	}

	p.connected = false

	return nil
}
