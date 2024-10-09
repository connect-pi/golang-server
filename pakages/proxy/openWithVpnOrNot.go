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

	// Remove IPv6 port information if present
	url = strings.Split(url, "]:")[0]
	url = strings.Replace(url, "[", "", -1)
	isIpv6 := common.IsIPv6(url)

	// Convert IPv6 to IPv4 if necessary
	if isIpv6 {
		fmt.Println("is Ipv6: ", url)

		// Is iran ipv6
		// IsIranIpv6 := common.IsIranIpv6(url)
		// if IsIranIpv6 {
		// 	SetCache(url, true)
		// 	return true
		// }
		// fmt.Print("\niran Ipv6: ", common.IsIranIpv6(url))

		// ipObj := net.ParseIP(url)
		// if ipObj == nil {
		// 	fmt.Printf("\nInvalid IPv6 address\n")
		// 	SetCache(url, true)
		// 	return true
		// }

		// if ipv4 := ipObj.To4(); ipv4 != nil {
		// 	fmt.Printf("\nConverted IPv6 %s to IPv4: %s\n", ipObj.String(), ipv4.String())
		// 	url = ipv4.String() // Use the converted IPv4 address
		// } else {
		// 	// If it's an IPv6 and cannot be converted, return true
		// 	fmt.Printf("\nCould not convert IPv6 to IPv4")
		// 	SetCache(url, true)
		// 	return true
		// }

		SetCache(url, true)
		return true
	}

	// Fix url
	if common.IsIPv4(url) {
		// Ipv4
		url = strings.Split(url, ":")[0]
		fmt.Println("Cleaned ipv4: ", url)
	} else {
		// Domain
		url = strings.Replace(url, "http://", "", -1)
		url = strings.Replace(url, "https://", "", -1)
		url = strings.Split(url, "/")[0]
		url = strings.Split(url, ":")[0]
		fmt.Println("Cleaned Domain: ", url)
	}

	// Check the cache first
	cache.mu.RLock()
	if cachedResult, found := cache.domainCache[url]; found {
		cache.mu.RUnlock()
		fmt.Println("Cache hit: ", cachedResult)
		return cachedResult
	}
	cache.mu.RUnlock()

	// Check if the domain is valid
	isDomain := common.IsValidDomain(url)
	if isDomain {
		// Apply domain-specific rules
		rulesForDomain := rules.CheckRulesForDomain(url)
		if rulesForDomain != nil {
			fmt.Println("Rules for domain found: ", *rulesForDomain)
			SetCache(url, *rulesForDomain)
			return *rulesForDomain
		}

		// Check if it's an Iranian or local domain
		if strings.Contains(url, ".ir") || strings.Contains(url, ".local") {
			SetCache(url, false)
			return false
		}

		// Fetch the domain's IP address
		thisIp, ipErr := common.GetDomainIp(url)
		if ipErr != nil {
			fmt.Println("Error fetching IP: ", ipErr)
			SetCache(url, false)
			return false
		}
		url = thisIp
	}

	// Ipv4
	if common.IsIPv4(url) {
		// Apply IP-specific rules
		rulesForIp := rules.CheckRulesForIp(url)
		if rulesForIp != nil {
			fmt.Println("Rules for IP found: ", *rulesForIp)
			SetCache(url, *rulesForIp)
			return *rulesForIp
		}

		// Check if the IP belongs to Iran
		IsIranIp := common.IsIranIp(url)
		SetCache(url, !IsIranIp)
		return !IsIranIp
	}

	return true

}
