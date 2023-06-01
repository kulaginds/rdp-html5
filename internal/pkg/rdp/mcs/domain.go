package mcs

import (
	"bytes"
	"fmt"
	"io"

	"github.com/kulaginds/rdp-html5/internal/pkg/rdp/per"
)

type DomainPDUApplication uint8

const (
	plumbDomainIndication DomainPDUApplication = iota
	erectDomainRequest
	mergeChannelsRequest
	mergeChannelsConfirm
	purgeChannelsIndication
	mergeTokensRequest
	mergeTokensConfirm
	purgeTokensIndication
	disconnectProviderUltimatum
	rejectMCSPDUUltimatum
	attachUserRequest
	attachUserConfirm
	detachUserRequest
	detachUserIndication
	channelJoinRequest
	channelJoinConfirm
	channelLeaveRequest
	channelConveneRequest
	channelConveneConfirm
	channelDisbandRequest
	channelDisbandIndication
	channelAdmitRequest
	channelAdmitIndication
	channelExpelRequest
	channelExpelIndication
	SendDataRequest
	SendDataIndication
	uniformSendDataRequest
	uniformSendDataIndication
	tokenGrabRequest
	tokenGrabConfirm
	tokenInhibitRequest
	tokenInhibitConfirm
	tokenGiveRequest
	tokenGiveIndication
	tokenGiveResponse
	tokenGiveConfirm
	tokenPleaseRequest
	tokenPleaseIndication
	tokenReleaseRequest
	tokenReleaseConfirm
	tokenTestRequest
	tokenTestConfirm
)

type DomainPDU struct {
	Application DomainPDUApplication

	ClientErectDomainRequest *ClientErectDomainRequest
	ClientAttachUserRequest  *ClientAttachUserRequest
	ClientChannelJoinRequest *ClientChannelJoinRequest
	ClientSendDataRequest    *ClientSendDataRequest

	ServerAttachUserConfirm  *ServerAttachUserConfirm
	ServerChannelJoinConfirm *ServerChannelJoinConfirm
	ServerSendDataIndication *ServerSendDataIndication
}

func (pdu *DomainPDU) Serialize() []byte {
	buf := new(bytes.Buffer)

	per.WriteChoice(uint8(pdu.Application<<2), buf)

	var data []byte

	switch pdu.Application {
	case attachUserRequest:
		data = pdu.ClientAttachUserRequest.Serialize()
	case erectDomainRequest:
		data = pdu.ClientErectDomainRequest.Serialize()
	case channelJoinRequest:
		data = pdu.ClientChannelJoinRequest.Serialize()
	case SendDataRequest:
		data = pdu.ClientSendDataRequest.Serialize()
	}

	buf.Write(data)

	return buf.Bytes()
}

func (pdu *DomainPDU) Deserialize(wire io.Reader) error {
	var (
		application uint8
		err         error
	)

	application, err = per.ReadChoice(wire)
	if err != nil {
		return err
	}
	pdu.Application = DomainPDUApplication(application >> 2)

	switch pdu.Application {
	case attachUserConfirm:
		pdu.ServerAttachUserConfirm = &ServerAttachUserConfirm{}

		return pdu.ServerAttachUserConfirm.Deserialize(wire)
	case channelJoinConfirm:
		pdu.ServerChannelJoinConfirm = &ServerChannelJoinConfirm{}

		return pdu.ServerChannelJoinConfirm.Deserialize(wire)
	case SendDataIndication:
		pdu.ServerSendDataIndication = &ServerSendDataIndication{}

		return pdu.ServerSendDataIndication.Deserialize(wire)
	case SendDataRequest:
		pdu.ClientSendDataRequest = &ClientSendDataRequest{}

		return pdu.ClientSendDataRequest.Deserialize(wire)
	case disconnectProviderUltimatum:
		return ErrDisconnectUltimatum
	}

	return fmt.Errorf("%w: application=%v", ErrUnknownDomainApplication, pdu.Application)
}
