package security

import "strings"

var sensitiveHeaderKeys = []string{
	"authorization",
	"automizely-api-key",
	"aftership-api-key",
	"am-api-key",
}

func SensitiveHeaderKeys() []string {
	return sensitiveHeaderKeys
}

func IsSensitiveHeaderKey(key string) bool {
	for _, skey := range sensitiveHeaderKeys {
		if skey == strings.ToLower(key) {
			return true
		}
	}
	return false
}
