package cmd

import (
	"fmt"
	"strconv"

	"github.com/EshaanAgg/toy-bittorrent/app/types"
	"github.com/EshaanAgg/toy-bittorrent/app/utils"
)

func HandleDownloadPiece(args []string) {
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

	pieceIdx, err := strconv.Atoi(args[3])
	if err != nil {
		fmt.Printf("error converting piece index to int: %v\n", err)
		return
	}
	if pieceIdx < 0 || pieceIdx >= len(fileInfo.InfoDict.Pieces) {
		fmt.Printf("piece index out of range: %d\n", pieceIdx)
		return
	}

	// Get the peers from the tracker
	peers, err := getPeersFromFile(fileInfo, true)
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
	peer := peers[0]
	err = peer.PrepareToGetPieceData(fileInfo.InfoHash, false)
	if err != nil {
		fmt.Printf("error preparing to get piece data: %v\n", err)
		return
	}

	// Register the piece with the peer
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
