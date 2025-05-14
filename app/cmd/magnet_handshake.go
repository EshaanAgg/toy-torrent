package cmd

import (
	"fmt"

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
	handshake, err := peer.PerformHandshake(m.InfoHash)
	if err != nil {
		println("error performing handshake:", err)
		return
	}

	fmt.Printf("Peer ID: %x\n", handshake.PeerID)

	// TODO: Send a bitfield message to the peer
	err = peer.BlockTillBitFieldMessage()
	if err != nil {
		fmt.Printf("error waiting for bitfield message: %v\n", err)
		return
	}

	if handshake.SupportsExtensions {
		extHandshake, err := peer.PerformExtensionHandshake()
		if err != nil {
			fmt.Printf("error performing extension handshake: %v\n", err)
		}
		if mtExtensionId, ok := extHandshake.ExtensionMap["ut_metadata"]; ok {
			fmt.Printf("Peer Metadata Extension ID: %d\n", mtExtensionId)
		}
	}
}
