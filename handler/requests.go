package handler

import (
	"github.com/MisterToaCarney/gosip/headers"
	"github.com/MisterToaCarney/gosip/messages"
	"github.com/MisterToaCarney/gosip/parser"
	"github.com/MisterToaCarney/gosip/sdp"
	"github.com/MisterToaCarney/gosip/utils"
)

func RespondToRequest(request *messages.RequestMessage, myDetails utils.MyDetails) *messages.ResponseMessage {
	var response *messages.ResponseMessage
	switch request.Method {
	case "OPTIONS":
		response = RespondOPTIONS(request, myDetails)
	}

	return response
}

func RespondOPTIONS(request *messages.RequestMessage, myDetails utils.MyDetails) *messages.ResponseMessage {
	branchHeader := headers.GetHeader(request.Headers, "Via")
	if branchHeader == nil {
		return nil
	}

	fromHeader := headers.GetHeader(request.Headers, "From")
	if fromHeader == nil {
		return nil
	}

	toHeader := headers.GetHeader(request.Headers, "To")
	if toHeader == nil {
		return nil
	}

	if toHeader.Parameters["tag"] == "" {
		toHeader.Parameters["tag"] = headers.GenTag()
	}

	callIdHeader := headers.GetHeader(request.Headers, "Call-ID")
	if callIdHeader == nil {
		return nil
	}

	cseqHeader := headers.GetHeader(request.Headers, "CSeq")
	if cseqHeader == nil {
		return nil
	}
	parsedCseqHeader := parser.ParseCSeq(cseqHeader.Value)
	if parsedCseqHeader == nil {
		return nil
	}

	sdp := sdp.CreateSDP(myDetails.LocalHost, 9)

	headers := []headers.Header{
		headers.Via(headers.TransportToSipTransport[myDetails.Transport], myDetails.LocalHost, myDetails.LocalPort, branchHeader.Value),
		*fromHeader,
		*toHeader,
		headers.Contact(myDetails.Username, myDetails.LocalHost, myDetails.LocalPort),
		headers.CallID(callIdHeader.Value),
		headers.CSeq(parsedCseqHeader.CSeq, parsedCseqHeader.Method),
		headers.UserAgent(),
		headers.Allow(),
		headers.ContentType("application/sdp"),
		headers.ContentLength(len(sdp)),
	}

	return messages.NewResponseMessage(200, "OK", headers, sdp)
}
