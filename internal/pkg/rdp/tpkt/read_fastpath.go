package tpkt

func (p *Protocol) StartHandleFastpath() {
	p.fastpathEnabled = true
}
