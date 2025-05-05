package utils

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
)

func SHA1Hash(d []byte) (string, error) {
	h := sha1.New()
	_, err := h.Write(d)
	if err != nil {
		return "", fmt.Errorf("error writing to hash: %w", err)
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}
