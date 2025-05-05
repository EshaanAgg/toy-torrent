package bencode

func isDigit(b *byte) bool {
	return *b >= '0' && *b <= '9'
}
