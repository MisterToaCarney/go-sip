package handler

import (
	"github.com/MisterToaCarney/gosip/headers"
	"github.com/MisterToaCarney/gosip/messages"
	"github.com/MisterToaCarney/gosip/parser"
	"github.com/MisterToaCarney/gosip/utils"
)

func RespondToResponse(response *messages.ResponseMessage, lastRequest messages.RequestMessage, myDetails utils.MyDetails) *messages.RequestMessage {
	var message *messages.RequestMessage
	switch response.StatusCode {
	case 401:
		message = Respond401(response, lastRequest, myDetails)
	}

	return message
}

func Respond401(response *messages.ResponseMessage, lastRequest messages.RequestMessage, myDetails utils.MyDetails) *messages.RequestMessage {
	cseqHeader := headers.GetHeader(response.Headers, "CSeq")
	if cseqHeader == nil {
		return nil
	}
	parsedCseqHeader := parser.ParseCSeq(cseqHeader.Value)
	if parsedCseqHeader == nil {
		return nil
	}

	var message *messages.RequestMessage

	switch parsedCseqHeader.Method {
	case "REGISTER":
		message = Respond401Register(response, lastRequest, myDetails)
	case "INVITE":
		message = Respond401Invite(response, lastRequest, myDetails)
	}

	return message
}
