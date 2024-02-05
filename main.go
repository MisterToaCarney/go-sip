package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/MisterToaCarney/gosip/handler"
	"github.com/MisterToaCarney/gosip/headers"
	"github.com/MisterToaCarney/gosip/messages"
	"github.com/MisterToaCarney/gosip/parser"
	"github.com/MisterToaCarney/gosip/utils"
)

func readRawHeader(reader *bufio.Reader) (string, error) {
	recievedMessage := ""
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			return "", err
		}
		recievedMessage += line
		if line == "\r\n" {
			return recievedMessage, nil
		}
	}
}

func readBody(reader *bufio.Reader, unparsedContentLength string) string {
	contentLength, err := strconv.Atoi(unparsedContentLength)
	if err != nil {
		fmt.Println("Unable to parse Content-Length header. Skipping body.")
		return ""
	}

	body := make([]byte, contentLength)
	io.ReadFull(reader, body)
	bodyStr := string(body)
	return bodyStr
}

func readSocket(conn net.Conn, requests chan *messages.RequestMessage, responses chan *messages.ResponseMessage) {
	reader := bufio.NewReader(conn)
	for {
		conn.SetReadDeadline(time.Time{})
		_, err := reader.Peek(1)
		if err != nil {
			panic(err)
		}
		conn.SetReadDeadline(time.Now().Add(5 * time.Second))
		rawHeader, err := readRawHeader(reader)
		if err != nil && !errors.Is(err, os.ErrDeadlineExceeded) {
			panic(err)
		}

		incomingResponseMessage, incomingRequestMessage := parser.ParseMessage(rawHeader)

		var body string

		if incomingRequestMessage != nil {
			contentLengthHeader := headers.GetHeader(incomingRequestMessage.Headers, "Content-Length")
			if contentLengthHeader != nil {
				body = readBody(reader, contentLengthHeader.Value)
				incomingRequestMessage.Body = body
				requests <- incomingRequestMessage
			}
		} else if incomingResponseMessage != nil {
			contentLengthHeader := headers.GetHeader(incomingResponseMessage.Headers, "Content-Length")
			if contentLengthHeader != nil {
				body = readBody(reader, contentLengthHeader.Value)
				incomingResponseMessage.Body = body
				responses <- incomingResponseMessage
			}
		}

		if len(rawHeader) > 0 {
			fmt.Println("Recieved:")
			fmt.Print(rawHeader)
			fmt.Print(body)
		}

	}
}

func writeSocket(conn net.Conn, message string) {
	fmt.Println("Sending:")
	fmt.Print(message)
	_, err := fmt.Fprint(conn, message)
	if err != nil {
		panic(err)
	}
}

func readTerminal(outgoingRequests chan messages.RequestMessage, d utils.MyDetails) {
	reader := bufio.NewReader(os.Stdin)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error occured whilst reading from terminal", err)
			continue
		}

		line = strings.ReplaceAll(line, "\r", "")
		line = strings.ReplaceAll(line, "\n", "")

		switch line {
		case "/call":
			fmt.Println("Call initiated")
			outgoingRequests <- *messages.Invite(d, "0277786643")
		default:
			fmt.Println("Command not recognized")
		}

	}
}

func SendRequest(conn net.Conn, request *messages.RequestMessage, requestsByCseq map[string]messages.RequestMessage) {
	if request == nil {
		fmt.Println("Warning: Tried sending a nil request. Doing nothing.")
		return
	}

	outgoingCseqHeader := headers.GetHeader(request.Headers, "CSeq")
	if outgoingCseqHeader != nil {
		requestsByCseq[outgoingCseqHeader.Value] = *request
	}

	writeSocket(conn, request.ToString())
}

func SendResponse(conn net.Conn, response *messages.ResponseMessage, responsesByCseq map[string]messages.ResponseMessage) {
	if response == nil {
		fmt.Println("Warning: Tried sending a nil response. Doing nothing.")
		return
	}

	outgoingCseqHeader := headers.GetHeader(response.Headers, "CSeq")
	if outgoingCseqHeader != nil {
		responsesByCseq[outgoingCseqHeader.Value] = *response
	}

	writeSocket(conn, response.ToString())
}

func main() {
	outgoingUserRequests := make(chan messages.RequestMessage)
	incomingRequests := make(chan *messages.RequestMessage)
	incomingResponses := make(chan *messages.ResponseMessage)
	requestsByCseq := make(map[string]messages.RequestMessage)
	responsesByCseq := make(map[string]messages.ResponseMessage)

	myDetails := GetConfig()

	conn, err := net.Dial(myDetails.Transport, myDetails.RemoteHost+":"+myDetails.RemotePort)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	addr, port := utils.ParseNetString(conn.LocalAddr())
	myDetails.LocalHost = addr
	myDetails.LocalPort = port

	go readSocket(conn, incomingRequests, incomingResponses)
	go readTerminal(outgoingUserRequests, myDetails)

	request := messages.BeginRegister(myDetails)
	SendRequest(conn, request, requestsByCseq)

	for {
		select {
		case incomingResponse := <-incomingResponses:
			incomingCseqHeader := headers.GetHeader(incomingResponse.Headers, "CSeq")
			var previousRequest messages.RequestMessage
			if incomingCseqHeader != nil {
				previousRequest = requestsByCseq[incomingCseqHeader.Value]
			}
			outgoingRequest := handler.RespondToResponse(incomingResponse, previousRequest, myDetails)
			if outgoingRequest != nil {
				SendRequest(conn, outgoingRequest, requestsByCseq)
			}

		case incomingRequest := <-incomingRequests:
			outgoingResponse := handler.RespondToRequest(incomingRequest, myDetails)
			if outgoingResponse != nil {
				SendResponse(conn, outgoingResponse, responsesByCseq)
			}

		case outgoingUserRequest := <-outgoingUserRequests:
			SendRequest(conn, &outgoingUserRequest, requestsByCseq)
		}
	}
}
