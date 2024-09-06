package common

import (
	"net"
)

// IsIPv6 checks if the given string is a valid IPv6 address
func IsIPv6(ip string) bool {
	parsedIP := net.ParseIP(ip)
	return parsedIP != nil && parsedIP.To4() == nil
}
