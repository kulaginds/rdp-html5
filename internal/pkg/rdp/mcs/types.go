package mcs

import (
	"bytes"
	"io"

	"github.com/kulaginds/web-rdp-solution/internal/pkg/rdp/ber"
)

const (
	RTSuccessful uint8 = iota
	RTDomainMerging
	RTDomainNotHierarchical
	RTNoSuchChannel
	RTNoSuchDomain
	RTNoSuchUser
	RTNotAdmitted
	RTOtherUserId
	RTParametersUnacceptable
	RTTokenNotAvailable
	RTTokenNotPossessed
	RTTooManyChannels
	RTTooManyTokens
	RTTooManyUsers
	RTUnspecifiedFailure
	RTUserRejected
)

const (
	RNDomainDisconnected uint8 = iota
	RNProviderInitiated
	RNTokenPurged
	RNUserRequested
	RNChannelPurged
)

type domainParameters struct {
	maxChannelIds   int
	maxUserIds      int
	maxTokenIds     int
	numPriorities   int
	minThroughput   int
	maxHeight       int
	maxMCSPDUsize   int
	protocolVersion int
}

func (params domainParameters) Serialize() []byte {
	buf := &bytes.Buffer{}

	ber.WriteInteger(params.maxChannelIds, buf)
	ber.WriteInteger(params.maxUserIds, buf)
	ber.WriteInteger(params.maxTokenIds, buf)
	ber.WriteInteger(params.numPriorities, buf)
	ber.WriteInteger(params.minThroughput, buf)
	ber.WriteInteger(params.maxHeight, buf)
	ber.WriteInteger(params.maxMCSPDUsize, buf)
	ber.WriteInteger(params.protocolVersion, buf)

	return buf.Bytes()
}

func (params *domainParameters) Deserialize(wire io.Reader) error {
	var err error

	params.maxChannelIds, err = ber.ReadInteger(wire)
	if err != nil {
		return err
	}

	params.maxUserIds, err = ber.ReadInteger(wire)
	if err != nil {
		return err
	}

	params.maxTokenIds, err = ber.ReadInteger(wire)
	if err != nil {
		return err
	}

	params.numPriorities, err = ber.ReadInteger(wire)
	if err != nil {
		return err
	}

	params.minThroughput, err = ber.ReadInteger(wire)
	if err != nil {
		return err
	}

	params.maxHeight, err = ber.ReadInteger(wire)
	if err != nil {
		return err
	}

	params.maxMCSPDUsize, err = ber.ReadInteger(wire)
	if err != nil {
		return err
	}

	params.protocolVersion, err = ber.ReadInteger(wire)
	if err != nil {
		return err
	}

	return nil
}
