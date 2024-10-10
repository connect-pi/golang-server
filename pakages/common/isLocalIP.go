package common

import "net"

func IsLocalIP(ipString string) bool {
	if ipString == "0.0.0.0" {
		return true
	}

	ip := net.ParseIP(ipString)

	// Check if the IP is a loopback address (localhost)
	if ip.IsLoopback() {
		return true
	}

	// Check if the IP is a private address (for IPv4 and IPv6)
	if ip.IsPrivate() {
		return true
	}

	// Redundant check for loopback addresses (already checked above)
	if ip.IsLoopback() {
		return true
	}

	// If it's not loopback or private, return false
	return false
}
