package comparison

import "strings"

func diffRunes(a, b byte) int {
	return int(b) - int(a)
}

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
