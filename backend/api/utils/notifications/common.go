package notifications

import "strings"

func isInvalidToken(errStr string) bool {
	es := strings.ToLower(errStr)
	return strings.Contains(es, "registration-token-not-registered") ||
		strings.Contains(es, "mismatchsenderid") ||
		strings.Contains(es, "invalid-registration-token")
}
