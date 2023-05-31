package mcs

import "errors"

var (
	ErrChannelNotFound           = errors.New("channel not found")
	ErrUnknownConnectApplication = errors.New("unknown connect application")
	ErrUnknownDomainApplication  = errors.New("unknown domain application")
	ErrUnknownChannel            = errors.New("unknown channel")
	ErrDisconnectUltimatum       = errors.New("disconnect ultimatum")
)
