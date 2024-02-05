package messages

import (
	"strconv"
	"strings"

	"github.com/MisterToaCarney/gosip/headers"
)

type ResponseMessage struct {
	StatusCode int
	StatusText string
	Headers    []headers.Header
	Body       string
}

func (response *ResponseMessage) ToString() string {
	sb := strings.Builder{}
	sb.WriteString("SIP/2.0 ")
	sb.WriteString(strconv.Itoa(response.StatusCode))
	sb.WriteString(" ")
	sb.WriteString(response.StatusText)
	sb.WriteString("\r\n")

	for _, header := range response.Headers {
		sb.WriteString(header.ToString())
	}
	sb.WriteString("\r\n")
	sb.WriteString(response.Body)

	return sb.String()
}

func NewResponseMessage(statusCode int, statusText string, messageHeaders []headers.Header, body string) *ResponseMessage {
	message := ResponseMessage{StatusCode: statusCode, StatusText: statusText, Headers: messageHeaders, Body: body}

	if headers.GetHeader(message.Headers, "Content-Length") == nil {
		message.Headers = append(message.Headers, headers.ContentLength(len(body)))
	}

	return &message
}
