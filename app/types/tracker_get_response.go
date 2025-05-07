package types

import (
	"fmt"

	"github.com/EshaanAgg/toy-bittorrent/app/bencode"
)

// TrackerResponse is a struct that holds the response from
// the tracker.
type TrackerGetResponse struct {
	Interval int
	Peers    []*Peer
}

func NewTrackerGetResponse(data []byte) (*TrackerGetResponse, error) {
	d, err := bencode.NewBencodeData(data)
	if err != nil {
		return nil, fmt.Errorf("error decoding bencode data: %w", err)
	}

	dict := d.GetDictionary()
	interval := dict.Map["interval"].GetInteger().Value

	peerString := dict.Map["peers"].GetString().Value
	if len(peerString)%6 != 0 {
		return nil, fmt.Errorf("peers length is not a multiple of 6")
	}
	peers := make([]*Peer, 0)
	for i := 0; i < len(peerString); i += 6 {
		ip := fmt.Sprintf("%d.%d.%d.%d", peerString[i], peerString[i+1], peerString[i+2], peerString[i+3])
		port := int(peerString[i+4])<<8 + int(peerString[i+5])
		peers = append(peers, &Peer{
			IP:   ip,
			Port: port,
		})
	}

	return &TrackerGetResponse{
		Interval: interval,
		Peers:    peers,
	}, nil
}
