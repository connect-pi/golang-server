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

// Check for open with VPN
func OpenWithVpnOrNot(url string) bool {
	fmt.Println("domain: ", url)

	// Check the cache first
	cache.mu.RLock()
	if cachedResult, found := cache.domainCache[url]; found {
		cache.mu.RUnlock()
		fmt.Println("Cache hit: ", cachedResult)
		return cachedResult
	}
	cache.mu.RUnlock()

	// Remove IPv6 port information if present
	originalUrl := url
	url = strings.Split(url, "]:")[0]
	url = strings.Replace(url, "[", "", -1)
	isIpv6 := common.IsIPv6(url)

	// For localhost
	if url == "::1" {
		SetCache(originalUrl, false)
		return false
	}

	// Convert IPv6 to IPv4 if necessary
	if isIpv6 {
		fmt.Println("is Ipv6: ", isIpv6)

		// Local ip
		isLocalIP := common.IsLocalIP(url)
		fmt.Println("is local IPv6: ", isLocalIP)
		if isLocalIP {
			SetCache(originalUrl, false)
			return false
		}

		SetCache(originalUrl, true)
		return true
	}

	// Fix url
	if common.IsIPv4(url) {
		// Ipv4
		url = strings.Split(url, ":")[0]
		fmt.Println("Cleaned ipv4: ", url)

		// Local ip
		isLocalIP := common.IsLocalIP(url)
		fmt.Println("is local: ", isLocalIP)
		if isLocalIP {
			SetCache(originalUrl, false)
			return false
		}
	} else {
		// Domain
		url = strings.Replace(url, "http://", "", -1)
		url = strings.Replace(url, "https://", "", -1)
		url = strings.Split(url, "/")[0]
		url = strings.Split(url, ":")[0]
		fmt.Println("Cleaned Domain: ", url)
	}

	// Check if the domain is valid
	isDomain := common.IsValidDomain(url)
	if isDomain {
		// Apply domain-specific rules
		rulesForDomain := rules.CheckRulesForDomain(url)
		if rulesForDomain != nil {
			fmt.Println("Rules for domain found: ", *rulesForDomain)
			SetCache(originalUrl, *rulesForDomain)
			return *rulesForDomain
		}

		// Check if it's an Iranian or local domain
		if strings.Contains(url, ".ir") || strings.Contains(url, ".local") {
			SetCache(originalUrl, false)
			return false
		}

		// Fetch the domain's IP address
		thisIp, ipErr := common.GetDomainIp(url)
		if ipErr != nil {
			fmt.Println("Error fetching IP: ", ipErr)
			SetCache(originalUrl, false)
			return false
		}
		url = thisIp
	}

	// Ipv4
	if common.IsIPv4(url) {
		// Apply IP-specific ruleso
		rulesForIp := rules.CheckRulesForIp(url)
		if rulesForIp != nil {
			fmt.Println("Rules for IP found: ", *rulesForIp)
			SetCache(originalUrl, *rulesForIp)
			return *rulesForIp
		}

		// Check if the IP belongs to Iran
		IsIranIp := common.IsIranIp(url)
		SetCache(originalUrl, !IsIranIp)
		return !IsIranIp
	}

	SetCache(originalUrl, true)
	return true
}
