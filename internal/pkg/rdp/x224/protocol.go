package x224

type protocol struct {
	tpktConn tpktConn

	requestedProtocols     RDPNegotiationProtocol
	connected              bool
	ServerNegotiationFlags RDPNegotiationResponseFlag
}

func New(tpktConn tpktConn, requestedProtocols RDPNegotiationProtocol) *protocol {
	return &protocol{
		tpktConn: tpktConn,

		requestedProtocols: requestedProtocols,
		connected:          false,
	}
}
