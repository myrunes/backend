// Package random generates cryptographically
// random values.
package random

import (
	"crypto/rand"
	"encoding/base64"
	"math/big"
)

// ByteArray returns a cryptographically random byte
// array of the given length.
func ByteArray(lngth int) ([]byte, error) {
	arr := make([]byte, lngth)
	_, err := rand.Read(arr)
	return arr, err
}

// Base64 creates a cryptographically randomly
// generated set of bytes with the length of lngth which
// is returned as base64 encoded string.
func Base64(lngth int) (string, error) {
	str, err := ByteArray(lngth)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(str), nil
}

// String returns a cryptographically random string with
// the given lngth from the set of characters passed.
func String(lngth int, set string) (string, error) {
	res := make([]byte, lngth)
	setlen := big.NewInt(int64(len(set)))

	for i := 0; i < lngth; i++ {
		randn, err := rand.Int(rand.Reader, setlen)
		if err != nil {
			return "", err
		}
		res[i] = set[randn.Int64()]
	}

	return string(res), nil
}
