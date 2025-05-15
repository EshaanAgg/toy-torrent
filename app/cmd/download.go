package cmd

import (
	"fmt"

	"github.com/EshaanAgg/toy-bittorrent/app/types"
)

func HandleDownload(args []string) {
	if len(args) != 3 || args[0] != "-o" {
		println("incorrect arguments passed. usage: go-torrent download_file -o <output-file> <token-file>")
		return
	}

	// Create the torrent file info
	torrentFile := args[2]
	fileInfo, err := types.NewTorrentFileInfo(torrentFile)
	if err != nil {
		fmt.Printf("error creating TorrentFileInfo: %v\n", err)
		return
	}

	// Get peers and prepare them to send piece data
	peers, err := getPeersFromFile(fileInfo, true)
	if err != nil {
		fmt.Printf("error getting peers: %v\n", err)
		return
	}
	if len(peers) == 0 {
		fmt.Println("no peers found")
		return
	}
	for _, peer := range peers {
		err = peer.PrepareToGetPieceData(fileInfo.InfoHash)
		if err != nil {
			fmt.Printf("error preparing peer: %v\n", err)
			return
		}
	}

	downloadPieces(peers, fileInfo, args[1])
}
