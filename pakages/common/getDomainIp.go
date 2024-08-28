package common

import (
	"fmt"
	"net"
)

// resolveDomain looks up the IP addresses of a domain with caching
func resolveDomain(domain string) ([]net.IP, error) {
	// Perform the DNS lookup
	ips, err := net.LookupIP(domain)
	if err != nil {
		return nil, err
	}

	return ips, nil
}

// getDomainIp returns just one IP address of the domain as a string
func GetDomainIp(domain string) (string, error) {
	ips, err := resolveDomain(domain)
	if err != nil {
		return "", err
	}

	if len(ips) == 0 {
		return "", fmt.Errorf("no IP addresses found for domain %s", domain)
	}

	return ips[0].String(), nil
}
