package parser

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/MisterToaCarney/gosip/headers"
	"github.com/MisterToaCarney/gosip/messages"
	"github.com/MisterToaCarney/gosip/utils"
)

type ParsedStatusLine struct {
	Version    string
	StatusCode int
	StatusText string
}

type ParsedRequestLine struct {
	Method  string
	URI     string
	Version string
}

var headerRe = regexp.MustCompile(`^ *(.+?) *: *(.+)`)
var extraLineRe = regexp.MustCompile(`^\s+(.+)`)
var paramRe = regexp.MustCompile(`^ *(\w+)=(.+)`)

func ParseMessage(message string) (*messages.ResponseMessage, *messages.RequestMessage) {
	lines := strings.Split(message, "\r\n")
	statusLine, requestLine := ParseFirstLine(lines[0])
	headers := ParseHeaders(lines[1:])

	if statusLine != nil {
		return messages.NewResponseMessage(statusLine.StatusCode, statusLine.StatusText, headers, ""), nil
	}
	if requestLine != nil {
		return nil, messages.NewRequestMessage(requestLine.Method, requestLine.URI, headers, "")
	}
	return nil, nil
}

func ParseFirstLine(line string) (*ParsedStatusLine, *ParsedRequestLine) {
	splits := strings.Split(line, " ")
	if splits[0] == "SIP/2.0" { // Line is a status line
		statusCode, err := strconv.Atoi(splits[1])
		if err != nil {
			return nil, nil
		}
		return &ParsedStatusLine{Version: splits[0], StatusCode: statusCode, StatusText: splits[2]}, nil
	} else if splits[2] == "SIP/2.0" { // Line is a request line
		return nil, &ParsedRequestLine{Method: splits[0], URI: splits[1], Version: splits[2]}
	} else {
		return nil, nil
	}
}

func ParseHeaders(lines []string) []headers.Header {
	mergedLines := make([][]string, 0, len(lines))
	currentKey, currentValue := "", ""
	for _, line := range lines {
		match := headerRe.FindStringSubmatch(line)
		if len(match) == 3 {
			if currentKey != "" && currentValue != "" {
				mergedLines = append(mergedLines, []string{currentKey, currentValue})
				currentKey, currentValue = match[1], match[2]
			} else {
				currentKey, currentValue = match[1], match[2]
			}
		}

		match = extraLineRe.FindStringSubmatch(line)
		if len(match) == 2 {
			currentValue += match[1]
		}
	}
	if currentKey != "" && currentValue != "" {
		mergedLines = append(mergedLines, []string{currentKey, currentValue})
	}

	out := make([]headers.Header, 0, 32)
	for _, line := range mergedLines {
		name, rawValue := line[0], line[1]
		header := headers.NewHeader(name)
		if utils.SliceContains([]any{"WWW-Authenticate", "Authorization", "Proxy-Authenticate", "Proxy-Authorization"}, header.Name) {
			header.Value = rawValue
			out = append(out, header)
			continue
		} else {
			parameterList := strings.Split(rawValue, ";")
			header.Value = parameterList[0]
			parameterList = parameterList[1:]
			for _, rawParam := range parameterList {
				match := paramRe.FindStringSubmatch(rawParam)
				if len(match) != 3 {
					continue
				}
				header.Parameters[match[1]] = match[2]
			}
			out = append(out, header)
		}
	}
	return out
}
