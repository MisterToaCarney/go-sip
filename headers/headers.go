package headers

import (
	"strings"
)

type Header struct {
	Name       string
	Value      string
	Parameters map[string]string
}

func NewHeader(name string) Header {
	return Header{Name: name, Value: "", Parameters: make(map[string]string)}
}

func (header Header) ToString() string {
	var sb strings.Builder
	sb.WriteString(header.Name)
	sb.WriteString(": ")
	sb.WriteString(header.Value)
	for key, value := range header.Parameters {
		sb.WriteString(";")
		sb.WriteString(key)
		sb.WriteString("=")
		sb.WriteString(value)
	}
	sb.WriteString("\r\n")
	return sb.String()
}

func (header *Header) Update(newHeader Header) *Header {
	header.Name = newHeader.Name
	header.Value = newHeader.Value
	header.Parameters = newHeader.Parameters
	return header
}
