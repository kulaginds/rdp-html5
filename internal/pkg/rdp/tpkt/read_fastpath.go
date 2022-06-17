package tpkt

func (p *protocol) StartHandleFastpath() {
	p.fastpathEnabled = true
}
