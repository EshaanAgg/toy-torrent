package cmd

import (
	"github.com/EshaanAgg/toy-bittorrent/app/types"
)

func HandleMagnetHandshake(args []string) {
	if len(args) != 1 {
		println("incorrect arguments passed. usage: go-torrent magnet_handshake <magnet-link>")
		return
	}

	magnetLink := args[0]
	m, err := types.NewMagnetURI(magnetLink)
	if err != nil {
		println("error creating MagnetURI:", err)
		return
	}

	peers, err := getPeers(m.TrackerURL, m.InfoHash, 999, true)
	if err != nil {
		println("error getting peers:", err)
		return
	}
	if len(peers) == 0 {
		println("no peers found")
		return
	}

	// Perform handshake with the first peer
	peer := peers[0]
	err = peer.PrepareToGetPieceData(m.InfoHash, true)
	if err != nil {
		println("error performing handshake:", err)
		return
	}
}
