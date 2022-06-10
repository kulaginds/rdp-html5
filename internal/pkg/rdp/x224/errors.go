package x224

import "errors"

var (
	ErrNotConnected                 = errors.New("not connected")
	ErrInvalidCorrelationID         = errors.New("invalid correlationId")
	ErrUnsupportedRequestedProtocol = errors.New("unsupported requested protocol")
	ErrSmallConnectionConfirmLength = errors.New("small connection confirm length")
	ErrWrongDataLength              = errors.New("wrong data length")
	ErrWrongConnectionConfirmCode   = errors.New("wrong connection confirm code")
)
