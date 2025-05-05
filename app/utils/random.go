package utils

import "math/rand"

func GetRandomPeerID() string {
	clientID := make([]byte, 20)
	for i := range 20 {
		rnd := rand.Intn(26)
		clientID[i] = byte('a' + rnd)
	}
	return string(clientID)
}
