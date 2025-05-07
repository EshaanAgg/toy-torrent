package utils

import "fmt"

func BytesToHex(b []byte) string {
	return fmt.Sprintf("%x", b)
}
