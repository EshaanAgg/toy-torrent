package cmd

import (
	"fmt"

	"github.com/EshaanAgg/toy-bittorrent/app/types"
)

func HandleDownload(args []string, s *types.Server) {
	if len(args) != 3 || args[0] != "-o" {
		println("incorrect arguments passed. usage: go-torrent download_file -o <output-file> <token-file>")
		return
	}

	// Parse the torrent file and the output file
	torrentFile := args[2]
	fileInfo, err := types.NewTorrentFileInfo(torrentFile)
	if err != nil {
		fmt.Printf("error creating TorrentFileInfo: %v\n", err)
		return
	}

	peers, err := getPeers(fileInfo, s, true)
	if err != nil {
		fmt.Printf("error getting peers: %v\n", err)
		return
	}
	if len(peers) == 0 {
		fmt.Println("no peers found")
		return
	}

	// pieceCnt := len(fileInfo.InfoDict.Pieces)
	// fmt.Printf("pieces count: %d, peers count: %d\n", pieceCnt, len(peers))
	// wg := &sync.WaitGroup{}

	// // Assign each of the piece to a peer in round robin fashion
	// for idx, pieceHash := range fileInfo.InfoDict.Pieces {
	// 	peerIdx := idx % len(peers)
	// 	peer := peers[peerIdx]

	// 	// If this is the first time accessing the peer, prepare it to get piece data
	// 	if peerIdx < len(peers) {
	// 		err = peer.PrepareToGetPieceData(s, fileInfo.InfoHash)
	// 		if err != nil {
	// 			println("error preparing to get piece data: ", err)
	// 			return
	// 		}
	// 		peer.SetWg(wg)
	// 		go peer.RegisterPieceMessageHandler()
	// 	}

	// 	// Create a new piece and assign it to the peer
	// 	pieceLen := getPieceLength(fileInfo, idx)
	// 	go peer.NewStoredPiece(uint32(idx), pieceLen, pieceHash)
	// }

	// wg.Wait()

	// // All the pieces have been downloaded, contact them to get the file
	// fileData := make([]byte, 0)
	// for idx := range pieceCnt {
	// 	peerIdx := idx % len(peers)
	// 	d, err := peers[peerIdx].GetPieceData(uint32(idx))
	// 	if err != nil {
	// 		fmt.Printf("error getting piece data: %v\n", err)
	// 		return
	// 	}
	// 	fileData = append(fileData, d...)
	// }

	// // Write the file data to the output file
	// err = utils.MakeFileWithData(args[1], fileData)
	// if err != nil {
	// 	fmt.Printf("error writing data to file: %v\n", err)
	// }
}
