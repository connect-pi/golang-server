package rules

func CheckRulesForDomain(domain string) *bool {
	// Perform the check
	var result *bool
	if CombinedRules.Domain.Has("on", domain) {
		res := true
		result = &res
	} else if CombinedRules.Domain.Has("off", domain) {
		res := false
		result = &res
	}

	return result
}

func CheckRulesForIp(ip string) *bool {

	// Perform the check
	var result *bool
	if CombinedRules.IP.Has("on", ip) {
		res := true
		result = &res
	} else if CombinedRules.IP.Has("off", ip) {
		res := false
		result = &res
	}

	return result
}
