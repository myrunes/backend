// Package comparison provides general functionalities
// for comparing different types and forms of data.
package comparison

import "strings"

// IsTrue returns true, if the passed string
// lowercased either equals "1" or "true".
func IsTrue(s string) bool {
	return s == "1" || strings.ToLower(s) == "true"
}
