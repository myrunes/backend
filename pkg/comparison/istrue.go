package comparison

import "strings"

func IsTrue(s string) bool {
	return s == "1" || strings.ToLower(s) == "true"
}
