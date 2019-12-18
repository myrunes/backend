// Package comparison provides general functionalities
// for comparing different types and forms of data.
package comparison

import "strings"

// diffRunes returns the numerical
// difference of the byte valeus
// between b and a (b - a).
func diffRunes(a, b byte) int {
	return int(b) - int(a)
}

// Alphabetically returns true, if string
// a is, if sorted alphabetically, indexed
// before string b.
//
// This can be used, for example, to sort a
// slice of strings alphabetically:
//
//   sort.Slice(s, func (a, b int) bool {
//     return comparison.Alphabetically(s[a], s[b])
//   })
func Alphabetically(a, b string) bool {
	a = strings.ToLower(a)
	b = strings.ToLower(b)

	minLen := len(a)

	lenB := len(b)
	if lenB < minLen {
		minLen = lenB
	}

	for i := 0; i < minLen; i++ {
		diff := diffRunes(a[i], b[i])
		if diff != 0 {
			return diff > 0
		}
	}

	return false
}
