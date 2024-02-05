package utils

import (
	"net"
	"strings"
)

type MyDetails struct {
	Scheme     string
	Transport  string
	Username   string
	Password   string
	RemoteHost string
	RemotePort string
	LocalHost  string
	LocalPort  string
}

func SliceContains(slice []any, target any) bool {
	for _, item := range slice {
		if item == target {
			return true
		}
	}
	return false
}

func ParseNetString(addr net.Addr) (string, string) {
	parts := strings.Split(addr.String(), ":")
	return parts[0], parts[1]
}
