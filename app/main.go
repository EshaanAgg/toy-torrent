package main

import (
	"fmt"
	"os"

	"github.com/EshaanAgg/toy-bittorrent/app/cmd"
	"github.com/EshaanAgg/toy-bittorrent/app/types"
)

func main() {
	args := os.Args[1:]
	s := types.NewServer()

	if len(args) == 0 {
		fmt.Printf("no arguments provided")
		return
	}

	switch args[0] {
	case "decode":
		cmd.HandleDecode(args[1:])

	case "info":
		cmd.HandleInfo(args[1:])

	case "peers":
		cmd.HandlePeers(args[1:], s)

	case "handshake":
		cmd.HandleHandshake(args[1:], s)

	case "download_piece":
		cmd.HandleDownloadPiece(args[1:], s)

	default:
		fmt.Printf("unrecognized command '%s'", args[0])
	}
}
