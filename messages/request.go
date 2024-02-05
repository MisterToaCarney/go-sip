package messages

import (
	"fmt"
	"strings"

	"github.com/MisterToaCarney/gosip/headers"
	"github.com/MisterToaCarney/gosip/sdp"
	"github.com/MisterToaCarney/gosip/utils"
)

type RequestMessage struct {
	Method  string
	URI     string
	Headers []headers.Header
	Body    string
}

func (request *RequestMessage) ToString() string {
	var sb strings.Builder
	sb.WriteString(request.Method)
	sb.WriteString(" ")
	sb.WriteString(request.URI)
	sb.WriteString(" SIP/2.0\r\n")

	for _, header := range request.Headers {
		sb.WriteString(header.ToString())
	}
	sb.WriteString("\r\n")
	sb.WriteString(request.Body)
	return sb.String()
}

func NewRequestMessage(method, uri string, messageHeaders []headers.Header, body string) *RequestMessage {
	message := RequestMessage{Headers: messageHeaders, Method: method, URI: uri}
	if headers.GetHeader(message.Headers, "Content-Length") == nil {
		message.Headers = append(message.Headers, headers.ContentLength(len(body)))
	}
	message.Body = body
	return &message
}

func BeginRegister(d utils.MyDetails) *RequestMessage {
	method := "REGISTER"
	uri := d.Scheme + ":" + d.RemoteHost
	headers := []headers.Header{
		headers.Via(headers.TransportToSipTransport[d.Transport], d.LocalHost, d.LocalPort, ""),
		headers.Contact(d.Username, d.LocalHost, d.LocalPort),
		headers.MaxForwards(70),
		headers.To(d.Username, d.RemoteHost, ""),
		headers.From(d.Username, d.RemoteHost, "generate_tag"),
		headers.CallID(""),
		headers.CSeq(15001, method),
		headers.UserAgent(),
		headers.Allow(),
		headers.ContentLength(0),
	}
	return NewRequestMessage(method, uri, headers, "")
}

func RegisterWithAuth(d utils.MyDetails, challenge headers.WWWAuthenticateHeader, challengeResponse, callId, fromTag, toTag string, cseq int) *RequestMessage {
	method := "REGISTER"
	uri := d.Scheme + ":" + d.RemoteHost
	headers := []headers.Header{
		headers.Via(headers.TransportToSipTransport[d.Transport], d.LocalHost, d.LocalPort, ""),
		headers.Contact(d.Username, d.LocalHost, d.LocalPort),
		headers.MaxForwards(70),
		headers.DigestAuthorization(
			d.Username, challenge.Parameters["realm"], challenge.Parameters["nonce"], uri,
			challengeResponse, challenge.Parameters["algorithm"],
		),
		headers.To(d.Username, d.RemoteHost, ""),
		headers.From(d.Username, d.RemoteHost, fromTag),
		headers.CallID(callId),
		headers.CSeq(cseq, method),
		headers.UserAgent(),
		headers.Allow(),
		headers.ContentLength(0),
	}
	return NewRequestMessage(method, uri, headers, "")
}

func Invite(d utils.MyDetails, phoneNumber string) *RequestMessage {
	method := "INVITE"
	uri := fmt.Sprintf("sip:%v@%v", phoneNumber, d.RemoteHost)
	body := sdp.CreateSDP(d.LocalHost, 5678)
	headers := []headers.Header{
		headers.Via(headers.TransportToSipTransport[d.Transport], d.LocalHost, d.LocalPort, ""),
		headers.Contact(d.Username, d.LocalHost, d.LocalPort),
		headers.MaxForwards(70),
		headers.To(phoneNumber, d.RemoteHost, ""),
		headers.From(d.Username, d.RemoteHost, "generate_tag"),
		headers.CallID(""),
		headers.CSeq(120, method),
		headers.UserAgent(),
		headers.Allow(),
		headers.ContentType("application/sdp"),
		headers.ContentLength(len(body)),
	}

	return NewRequestMessage(method, uri, headers, body)
}
