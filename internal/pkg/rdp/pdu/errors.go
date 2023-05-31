package pdu

import "errors"

var (
	ErrInvalidCorrelationID = errors.New("invalid correlationId")
	ErrDeactiateAll         = errors.New("deactivate all")
)
