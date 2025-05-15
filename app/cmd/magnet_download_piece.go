package cmd

import (
	"fmt"
	"strconv"

	"github.com/EshaanAgg/toy-bittorrent/app/types"
	"github.com/EshaanAgg/toy-bittorrent/app/utils"
)

func HandleMagnetDownloadPiece(args []string) {
	if len(args) != 4 || args[0] != "-o" {
		fmt.Println("incorrect arguments passed. usage: go-torrent magnet_download_piece -o <output-file> <magnet-file> <piece-index>")
		return
	}

	// Get the magnet file
	magnetFile := args[2]
	m, err := types.NewMagnetURI(magnetFile)
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

	// Prepare to get piece data from the first peer, and also get the info file
	peer := peers[0]
	err = peer.PrepareToGetPieceData(m.InfoHash, false)
	if err != nil {
		fmt.Printf("error performing handshake: %v\n", err)
		return
	}
	fileInfo, err := peer.GetInfoFile(m)
	if err != nil {
		fmt.Printf("error getting info file: %v\n", err)
		return
	}

	// Parse the piece index
	pieceIdx, err := strconv.Atoi(args[3])
	if err != nil {
		fmt.Printf("error converting piece index to int: %v\n", err)
		return
	}
	if pieceIdx < 0 || pieceIdx >= len(fileInfo.InfoDict.Pieces) {
		fmt.Printf("piece index out of range: %d\n", pieceIdx)
		return
	}

	pieceLen := getPieceLength(fileInfo, pieceIdx)
	pieceHash := fileInfo.InfoDict.Pieces[pieceIdx]

	var sp *types.StoredPiece

	for {
		sp, err = peer.DownloadPiece(uint32(pieceIdx), pieceLen, pieceHash)
		if err != nil {
			fmt.Printf("error downloading piece: %v\n", err)
			continue
		}

		break
	}

	err = utils.MakeFileWithData(args[1], sp.GetData())
	if err != nil {
		fmt.Printf("error writing data to file: %v\n", err)
	}
}
