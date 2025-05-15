package cmd

import (
	"fmt"
	"os"

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

	// Create the pieces
	var pieces []*pieceToDownload
	for i := range len(fileInfo.InfoDict.Pieces) {
		length := getPieceLength(fileInfo, i)
		pieces = append(pieces, &pieceToDownload{
			index:  uint32(i),
			length: length,
			hash:   fileInfo.InfoDict.Pieces[i],
		})
	}

	fileData := downloadPieces(peers, pieces)
	outputFile := args[1]
	err = os.WriteFile(outputFile, fileData, 0644)
	if err != nil {
		fmt.Printf("error writing final file: %v\n", err)
		return
	}

	fmt.Printf("Downloaded file saved to '%s'\n", outputFile)
}
