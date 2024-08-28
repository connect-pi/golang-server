package proxy

import (
	"fmt"
	"net"
	"sync"
	"time"
)

var (
	cache         = make(map[string][]net.IP)  // Cache for domain to IPs
	cacheMutex    = &sync.RWMutex{}            // Mutex to handle concurrent access
	cacheDuration = 60 * time.Minute           // Duration to cache the result
	cacheExpiry   = make(map[string]time.Time) // Expiry times for cache entries
)

// resolveDomain looks up the IP addresses of a domain with caching
func resolveDomain(domain string) ([]net.IP, error) {
	cacheMutex.RLock()
	if ips, found := cache[domain]; found && time.Now().Before(cacheExpiry[domain]) {
		cacheMutex.RUnlock()
		return ips, nil
	}
	cacheMutex.RUnlock()

	// Perform the DNS lookup
	ips, err := net.LookupIP(domain)
	if err != nil {
		return nil, err
	}

	// Update the cache
	cacheMutex.Lock()
	cache[domain] = ips
	cacheExpiry[domain] = time.Now().Add(cacheDuration)
	cacheMutex.Unlock()

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
