package cmd

import (
	"fmt"

	"github.com/EshaanAgg/toy-bittorrent/app/types"
)

func HandleMagnetDownload(args []string) {
	if len(args) != 3 || args[0] != "-o" {
		fmt.Println("usage: go-torrent magnet_download -o <output-file> <magnet-url>")
		return
	}

	m, err := types.NewMagnetURI(args[2])
	if err != nil {
		fmt.Printf("error creating MagnetURI: %v\n", err)
		return
	}

	// Get the peers from the tracker
	peers, err := getPeers(m.TrackerURL, m.InfoHash, 999, true)
	if err != nil {
		fmt.Printf("error getting peers: %v\n", err)
		return
	}
	if len(peers) == 0 {
		fmt.Println("no peers found")
		return
	}

	var fileInfo *types.TorrentFileInfo

	// Prepare the peers to get piece data
	for _, peer := range peers {
		fileInfo, err = peer.PrepareToGetPieceData_Magnet(m)
		if err != nil {
			fmt.Printf("error performing handshake: %v\n", err)
			return
		}
	}

	downloadPieces(peers, fileInfo, args[1])
}
