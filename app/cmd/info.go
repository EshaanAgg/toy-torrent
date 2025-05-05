package cmd

import (
	"fmt"

	"github.com/EshaanAgg/toy-bittorrent/app/types"
)

// printPieceHashes prints the piece hashes from the pieces byte array.
// Each piece hash is 20 bytes long, so the length of the pieces byte array
// should be a multiple of 20.
func printPieceHashes(pieces []byte) {
	if len(pieces)%20 != 0 {
		fmt.Printf("pieces length is not a multiple of 20")
		return
	}

	fmt.Println("Piece Hashes:")
	for i := 0; i < len(pieces); i += 20 {
		pieceHash := pieces[i : i+20]
		fmt.Printf("%x\n", pieceHash)
	}
}

func HandleInfo(args []string) {
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

	fmt.Printf("Tracker URL: %s\n", fileInfo.TrackerURL)
	fmt.Printf("Length: %d\n", fileInfo.FileSize)
	fmt.Printf("Info Hash: %s\n", fileInfo.InfoHash)

	infoDict := fileInfo.InfoDict
	fmt.Printf("Piece Length: %d\n", infoDict.Map["piece length"].GetInteger().Value)
	pieces := infoDict.Map["pieces"].GetString().Value
	printPieceHashes(pieces)
}
