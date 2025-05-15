package cmd

import (
	"fmt"
	"os"

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
