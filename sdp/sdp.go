package sdp

import (
	"fmt"
	"strings"
	"time"
)

func Timestamp1900() int64 {
	t := time.Now().Unix()
	return t + 2208988800
}

func CreateSDP(localAddr string, offerPort int) string {
	timestamp := Timestamp1900()

	sb := strings.Builder{}
	sb.WriteString("v=0\r\n")
	sb.WriteString(fmt.Sprintf("o=- %v %v IN IP4 %v\r\n", timestamp, timestamp, localAddr))
	sb.WriteString("s=-\r\n")
	sb.WriteString(fmt.Sprintf("c=IN IP4 %v\r\n", localAddr))
	sb.WriteString("t=0 0\r\n")
	sb.WriteString(fmt.Sprintf("m=audio %v RTP/AVP 8\r\n", offerPort))
	sb.WriteString("a=rtpmap:8 PCMA/8000\r\n")
	sb.WriteString("a=sendrecv\r\n")

	return sb.String()
}
