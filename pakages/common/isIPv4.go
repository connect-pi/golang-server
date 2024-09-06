package common

import "strings"

func IsIPv4(address string) bool {
	return strings.Count(address, ".") == 3
}
