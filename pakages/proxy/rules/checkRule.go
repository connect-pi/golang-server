package rules

func CheckRulesForDomain(domain string) *bool {
	// Perform the check
	var result *bool
	if combinedRules.Domain.Has("on", domain) {
		res := true
		result = &res
	} else if combinedRules.Domain.Has("off", domain) {
		res := false
		result = &res
	}

	return result
}

func CheckRulesForIp(ip string) *bool {

	// Perform the check
	var result *bool
	if combinedRules.IP.Has("on", ip) {
		res := true
		result = &res
	} else if combinedRules.IP.Has("off", ip) {
		res := false
		result = &res
	}

	return result
}
