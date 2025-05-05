package utils

import (
	"crypto/sha1"
	"fmt"
)

func SHA1Hash(d []byte) ([]byte, error) {
	h := sha1.New()
	_, err := h.Write(d)
	if err != nil {
		return nil, fmt.Errorf("error writing to hash: %w", err)
	}
	return h.Sum(nil), nil
}
