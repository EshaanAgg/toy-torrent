package main

import (
	"fmt"
	"os"

	"github.com/EshaanAgg/toy-bittorrent/app/cmd"
	"github.com/EshaanAgg/toy-bittorrent/app/types"
)

type HandlerFn func(args []string, server *types.Server)

var handlers = map[string]HandlerFn{
	"decode":         cmd.HandleDecode,
	"info":           cmd.HandleInfo,
	"peers":          cmd.HandlePeers,
	"handshake":      cmd.HandleHandshake,
	"download_piece": cmd.HandleDownloadPiece,
	"download":       cmd.HandleDownload,
	"magnet_parse":   cmd.HandleMagnetParse,
}

func main() {
	args := os.Args[1:]
	s := types.NewServer()

	if len(args) == 0 {
		fmt.Printf("no arguments provided")
		return
	}

	handler, ok := handlers[args[0]]
	if !ok {
		fmt.Printf("unknown command: %s\n", args[0])
		return
	}
	handler(args[1:], s)
}
