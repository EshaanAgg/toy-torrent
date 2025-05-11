package cmd

import (
	"bytes"
	"fmt"
	"os"
	"strconv"

	"github.com/EshaanAgg/toy-bittorrent/app/types"
	"github.com/EshaanAgg/toy-bittorrent/app/utils"
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

	pieceIdx, err := strconv.Atoi(args[3])
	if err != nil {
		fmt.Printf("error converting piece index to int: %v\n", err)
		return
	}
	if pieceIdx < 0 || pieceIdx >= len(fileInfo.InfoDict.Pieces) {
		fmt.Printf("piece index out of range: %d\n", pieceIdx)
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
	peer := peers[0]
	err = peer.PrepareToGetPieceData(s, fileInfo.InfoHash)
	if err != nil {
		fmt.Printf("error preparing to get piece data: %v\n", err)
		return
	}
	pieceHash := fileInfo.InfoDict.Pieces[pieceIdx]

	go peer.RegisterPieceMessageHandler(func(sp *types.StoredPiece) {
		data := sp.GetData()

		// Create the output file and write the piece data to it
		err := utils.MakeFileWithData(args[1], data)
		if err != nil {
			fmt.Printf("error writing piece data to file: %v\n", err)
		}

		// Verify the piece hash
		hash, err := utils.SHA1Hash(data)
		if err != nil {
			fmt.Printf("error hashing piece data: %v\n", err)
		}

		if bytes.Equal(hash, pieceHash) {
			fmt.Printf("piece %d hash verified\n", pieceIdx)
		} else {
			fmt.Printf("piece %d hash verification failed\n", pieceIdx)
		}
		os.Exit(0) // Exit the program after writing the piece
	})

	// The length of the last piece may be less than the piece length
	pieceLength := fileInfo.InfoDict.PieceLength
	fileLength := fileInfo.InfoDict.Length
	if (pieceIdx+1)*pieceLength > fileLength {
		pieceLength = fileLength - (pieceIdx * pieceLength)
	}

	sp := peer.NewStoredPiece(uint32(pieceIdx), uint32(pieceLength))
	sp.Download(peer) // Start downloading the piece

	for {
	}
}
