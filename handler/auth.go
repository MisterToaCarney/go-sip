package handler

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"

	"github.com/MisterToaCarney/gosip/headers"
	"github.com/MisterToaCarney/gosip/messages"
	"github.com/MisterToaCarney/gosip/parser"
	"github.com/MisterToaCarney/gosip/sdp"
	"github.com/MisterToaCarney/gosip/utils"
)

func DoDigestChallenge(username, realm, password, method, digestURI, nonce string) string {
	a1 := []byte(fmt.Sprintf("%s:%s:%s", username, realm, password))
	a2 := []byte(fmt.Sprintf("%s:%s", method, digestURI))
	ha1 := md5.Sum(a1)
	ha2 := md5.Sum(a2)
	str_ha1 := hex.EncodeToString(ha1[:])
	str_ha2 := hex.EncodeToString(ha2[:])
	a3 := []byte(fmt.Sprintf("%s:%s:%s", str_ha1, nonce, str_ha2))
	response := md5.Sum(a3)
	strResponse := hex.EncodeToString(response[:])
	return strResponse
}

func Respond401Register(response *messages.ResponseMessage, lastRequest messages.RequestMessage, myDetails utils.MyDetails) *messages.RequestMessage {
	authHeader := headers.GetHeader(response.Headers, "WWW-Authenticate")
	if authHeader == nil {
		return nil
	}

	challenge := parser.ParseWWWAuthenticate(authHeader.Value)
	if challenge.Method != "Digest" || challenge.Parameters["algorithm"] != "MD5" {
		return nil
	}

	digestAnswer := DoDigestChallenge(
		myDetails.Username, challenge.Parameters["realm"], myDetails.Password, "REGISTER",
		lastRequest.URI, challenge.Parameters["nonce"],
	)

	cseqHeader := headers.GetHeader(response.Headers, "CSeq")
	if cseqHeader == nil {
		return nil
	}

	parsedCseqHeader := parser.ParseCSeq(cseqHeader.Value)
	if parsedCseqHeader == nil {
		return nil
	}

	cseq := parsedCseqHeader.CSeq + 1

	var header *headers.Header
	header = headers.GetHeader(response.Headers, "Call-ID")
	if header == nil {
		return nil
	}
	callId := header.Value

	header = headers.GetHeader(response.Headers, "From")
	if header == nil {
		return nil
	}
	fromTag := header.Parameters["tag"]

	header = headers.GetHeader(response.Headers, "To")
	if header == nil {
		return nil
	}
	toTag := header.Parameters["tag"]

	message := messages.RegisterWithAuth(
		myDetails,
		challenge,
		digestAnswer,
		callId,
		fromTag,
		toTag,
		cseq,
	)

	return message
}

func Respond401Invite(response *messages.ResponseMessage, lastRequest messages.RequestMessage, myDetails utils.MyDetails) *messages.RequestMessage {
	message := lastRequest

	authHeader := headers.GetHeader(response.Headers, "WWW-Authenticate")
	if authHeader == nil {
		return nil
	}
	challenge := parser.ParseWWWAuthenticate(authHeader.Value)
	if challenge.Method != "Digest" || challenge.Parameters["algorithm"] != "MD5" {
		return nil
	}
	digestAnswer := DoDigestChallenge(
		myDetails.Username, challenge.Parameters["realm"], myDetails.Password, message.Method,
		lastRequest.URI, challenge.Parameters["nonce"],
	)

	cseqHeader := headers.GetHeader(response.Headers, "CSeq")
	if cseqHeader == nil {
		return nil
	}
	parsedCseqHeader := parser.ParseCSeq(cseqHeader.Value)
	if parsedCseqHeader == nil {
		return nil
	}
	headers.GetHeader(message.Headers, "CSeq").Update(headers.CSeq(parsedCseqHeader.CSeq+1, parsedCseqHeader.Method))

	toHeader := headers.GetHeader(response.Headers, "To")
	if toHeader != nil {
		headers.GetHeader(message.Headers, "To").Parameters["tag"] = toHeader.Parameters["tag"]
	}

	digestHeader := headers.DigestAuthorization(
		myDetails.Username, challenge.Parameters["realm"], challenge.Parameters["nonce"],
		message.URI, digestAnswer, challenge.Parameters["algorithm"],
	)

	headers := make([]headers.Header, 0, 20)
	headers = append(headers, message.Headers[:2]...)
	headers = append(headers, digestHeader)
	headers = append(headers, message.Headers[2:]...)

	message.Headers = headers

	message.Body = sdp.CreateSDP(myDetails.LocalHost, 5678)

	return &message
}
