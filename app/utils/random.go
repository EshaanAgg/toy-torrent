package utils

func GetRandomClientID() string {
	clientID := make([]byte, 20)
	for i := range 20 {
		clientID[i] = byte('a' + i%26)
	}
	return string(clientID)
}
