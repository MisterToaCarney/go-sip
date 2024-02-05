package parser

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/MisterToaCarney/gosip/headers"
)

var wwwAuthSplitRe = regexp.MustCompile(`^ *(\w+) +(.+)`)
var parameterRe = regexp.MustCompile(`^ *(\w+)=(.+)`)

func ParseWWWAuthenticate(value string) headers.WWWAuthenticateHeader {
	header := headers.WWWAuthenticateHeader{Parameters: make(map[string]string)}
	splitSubmatches := wwwAuthSplitRe.FindStringSubmatch(value)
	if len(splitSubmatches) != 3 {
		return headers.WWWAuthenticateHeader{}
	}
	header.Method = splitSubmatches[1]
	rawParameters := strings.Split(splitSubmatches[2], ",")
	for _, rawParameter := range rawParameters {
		parameterMatch := parameterRe.FindStringSubmatch(rawParameter)
		if len(parameterMatch) != 3 {
			continue
		}
		key := parameterMatch[1]
		value := strings.ReplaceAll(parameterMatch[2], `"`, "")
		header.Parameters[key] = value
	}

	return header
}

func ParseCSeq(value string) *headers.CSeqHeader {
	splits := strings.Split(value, " ")
	seq, err := strconv.Atoi(splits[0])
	if err != nil {
		return nil
	}
	return &headers.CSeqHeader{Method: splits[1], CSeq: seq}
}
