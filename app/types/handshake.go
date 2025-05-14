package types

type Handshake struct {
	PeerID   []byte
	InfoHash []byte
}

func (h *Handshake) Bytes() []byte {
	var b []byte

	b = append(b, 19) // Protocol length
	b = append(b, "BitTorrent protocol"...)

	// To signal support for extensions, set the 20th bit
	// from right to 1 (out of the 64 bits)
	// 20 => i = 5 | mask = 1 << 4
	for i := range 8 {
		if i == 5 {
			b = append(b, 1<<4)
		} else {
			b = append(b, 0)
		}
	}

	b = append(b, h.InfoHash...)
	b = append(b, h.PeerID...)
	return b
}

func NewHandshakeFromBytes(data []byte) *Handshake {
	if len(data) < 68 {
		return nil
	}

	h := &Handshake{
		InfoHash: data[28:48],
		PeerID:   data[48:68],
	}

	return h
}
