package common

import (
	"regexp"
)

var domainPattern = regexp.MustCompile(`^(?:[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?\.)+[a-zA-Z]{2,}$`)

// IsValidDomain checks if a domain name is valid
func IsValidDomain(domain string) bool {
	return domainPattern.MatchString(domain)
}
