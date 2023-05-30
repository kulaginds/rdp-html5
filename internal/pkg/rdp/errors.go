package rdp

import "errors"

var (
	ErrInvalidCorrelationID         = errors.New("invalid correlationId")
	ErrUnsupportedRequestedProtocol = errors.New("unsupported requested protocol")
)
