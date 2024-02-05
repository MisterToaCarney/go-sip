package headers

import (
	crand "crypto/rand"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"math/rand"
	"time"

	"github.com/google/uuid"
)

const (
	TCP = "SIP/2.0/TCP"
	UDP = "SIP/2.0/UDP"
)

var TransportToSipTransport = map[string]string{
	"tcp": "SIP/2.0/TCP",
	"udp": "SIP/2.0/UDP",
}

func GenBranch() string {
	nanoBytes := make([]byte, 8)
	randBytes := make([]byte, 4)
	binary.LittleEndian.PutUint64(nanoBytes, uint64(time.Now().UnixNano()))
	binary.LittleEndian.PutUint32(randBytes, rand.Uint32())
	return base64.URLEncoding.EncodeToString(append(nanoBytes, randBytes...))
}

func GenTag() string {
	nanoBytes := make([]byte, 8)
	randBytes := make([]byte, 4)
	binary.LittleEndian.PutUint64(nanoBytes, uint64(time.Now().UnixMicro()))
	crand.Read(randBytes)
	return hex.EncodeToString(append(nanoBytes, randBytes...))
}

func GenCallID() string {
	return uuid.NewString()
}

func GetHeaderIndex(headers []Header, name string) int {
	for i, header := range headers {
		if header.Name == name {
			return i
		}
	}
	return -1
}

func GetHeaderIndices(headers []Header, name string) []int {
	indicies := make([]int, 0, 16)
	for i, header := range headers {
		if header.Name == name {
			indicies = append(indicies, i)
		}
	}
	return indicies
}

func GetHeader(headers []Header, name string) *Header {
	i := GetHeaderIndex(headers, name)
	if i == -1 {
		return nil
	}
	return &headers[i]
}

func GetHeaders(headers []Header, name string) []*Header {
	out := make([]*Header, 0, 16)
	indices := GetHeaderIndices(headers, name)
	for _, index := range indices {
		out = append(out, &headers[index])
	}
	return out
}
