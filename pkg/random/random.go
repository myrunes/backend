package random

import (
	crand "crypto/rand"
	"encoding/base64"
	mrand "math/rand"
)

// GetRandBase64Str creates a cryptographically randomly
// generated set of bytes with the length of lngth which
// is returned as base64 encoded string.
func GetRandBase64Str(lngth int) (string, error) {
	str := make([]byte, lngth)

	if _, err := crand.Read(str); err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(str), nil
}

// GetRandString returns a pseudo-random string with
// the given lngth from the set of characters passed.
func GetRandString(lngth int, set string) string {
	res := make([]byte, lngth)
	setlen := len(set)

	for i := 0; i < lngth; i++ {
		randn := mrand.Intn(setlen)
		res[i] = set[randn]
	}

	return string(res)
}
