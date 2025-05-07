package types

import "fmt"

type Handshake struct {
	PeerID   []byte
	InfoHash []byte
}

func (h *Handshake) Bytes() []byte {
	var b []byte

	b = append(b, 19) // Protocol length
	b = append(b, "BitTorrent protocol"...)

	// 8 bytes reserved for future use
	for range 8 {
		b = append(b, 0)
	}

	b = append(b, h.InfoHash...)
	b = append(b, h.PeerID...)
	return b
}

func NewHandshakeFromBytes(data []byte) *Handshake {
	if len(data) < 68 {
		return nil
	}
	fmt.Println("Received Handshake: ", data)
	h := &Handshake{
		PeerID:   data[28:48],
		InfoHash: data[48:68],
	}
	fmt.Println("Peer ID: ", h.PeerID)
	fmt.Println("Info Hash: ", h.InfoHash)
	return h
}
