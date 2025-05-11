package cmd

import (
	"fmt"

	"github.com/EshaanAgg/toy-bittorrent/app/types"
)

func HandleInfo(args []string, _s *types.Server) {
	// Validate the number of arguments passed to the info command
	if len(args) == 0 {
		fmt.Println("no data passed to info. usage: go-torrent info <path-to-file>")
		return
	}

	if len(args) > 1 {
		fmt.Println("too many arguments passed to info. usage: go-torrent info <path-to-file>")
		return
	}

	fileInfo, err := types.NewTorrentFileInfo(args[0])
	if err != nil {
		fmt.Printf("error creating TorrentFileInfo: %v\n", err)
		return
	}

	infoDict := fileInfo.InfoDict

	fmt.Printf("Tracker URL: %s\n", fileInfo.TrackerURL)
	fmt.Printf("Length: %d\n", infoDict.Length)
	fmt.Printf("Info Hash: %s\n", fileInfo.GetHexInfoHash())
	fmt.Printf("Piece Length: %d\n", infoDict.PieceLength)
	fmt.Println("Piece Hashes:")
	for _, p := range infoDict.Pieces {
		fmt.Printf("%x\n", p)
	}
}
