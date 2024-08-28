package proxy

import (
	"fmt"
	"project/pakages/common"
	"project/pakages/proxy/rules"
	"strings"
	"sync"
)

// Cache structure
type Cache struct {
	domainCache map[string]bool
	mu          sync.RWMutex
}

// Global cache instance
var cache = Cache{
	domainCache: make(map[string]bool),
}

// SetCache function to store results in the cache
func SetCache(domain string, result bool) {
	cache.mu.Lock()
	defer cache.mu.Unlock()
	cache.domainCache[domain] = result
}

// Check for open with vpn
func OpenWithVpn(domain string) bool {
	// Fix domain
	domain = strings.Replace(domain, "http://", "", -1)
	domain = strings.Replace(domain, "https://", "", -1)
	domain = strings.Split(domain, "/")[0]
	domain = strings.Split(domain, ":")[0]

	fmt.Print("\ndomain: ", domain)

	// Check cache
	cache.mu.RLock()
	if cachedResult, found := cache.domainCache[domain]; found {
		cache.mu.RUnlock()
		fmt.Print("\nCache hit: ", cachedResult)
		return cachedResult
	}
	cache.mu.RUnlock()

	// Check domain rules
	rulesForDomain := rules.CheckRulesForDomain(domain)
	if rulesForDomain != nil {
		fmt.Print("\nrulesForDomain: ", *rulesForDomain)
		SetCache(domain, *rulesForDomain)
		return *rulesForDomain
	}

	// Validate domain
	if !common.IsValidDomain(domain) {
		SetCache(domain, false)
		return false
	}

	// Check .ir/local domains
	if strings.Contains(domain, ".ir") || strings.Contains(domain, ".local") {
		SetCache(domain, false)
		return false
	}

	// Check iran ip
	ip, ipErr := common.GetDomainIp(domain)
	if ipErr != nil {
		fmt.Print("\nGet ip error: ", ipErr, "\n")
		SetCache(domain, false)
		return false
	}

	// Check ip rules
	rulesForIp := rules.CheckRulesForIp(ip)
	if rulesForIp != nil {
		fmt.Print("\nrulesForIp: ", *rulesForIp)
		SetCache(domain, *rulesForIp)
		return *rulesForIp
	}

	fmt.Print("\nip: ", ip)

	isNotIranIp := !common.IsIranIp(ip)
	SetCache(domain, isNotIranIp)

	return isNotIranIp
}
