package cmd

import (
	"fmt"

	"github.com/EshaanAgg/toy-bittorrent/app/types"
)

func HandleDownloadPiece(args []string, s *types.Server) {
	if len(args) != 4 || args[0] != "-o" {
		fmt.Println("incorrect arguments passed. usage: go-torrent download_piece -o <output-file> <token-file> <piece-index>")
		return
	}

	torrentFile := args[2]
	fileInfo, err := types.NewTorrentFileInfo(torrentFile)
	if err != nil {
		fmt.Printf("error creating TorrentFileInfo: %v\n", err)
		return
	}

	peers, err := getPeers(fileInfo, s)
	if err != nil {
		fmt.Printf("error getting peers: %v\n", err)
		return
	}

	if len(peers) == 0 {
		fmt.Println("no peers found")
		return
	}

	// Take the first peer from the list to download the piece from
	// We can do this as all the peers have all the pieces.
	// TODO: Fix this assumption
}
