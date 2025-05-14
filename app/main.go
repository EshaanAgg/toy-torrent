package main

import (
	"fmt"
	"os"

	"github.com/EshaanAgg/toy-bittorrent/app/cmd"
)

type HandlerFn func(args []string)

var handlers = map[string]HandlerFn{
	"decode":           cmd.HandleDecode,
	"info":             cmd.HandleInfo,
	"peers":            cmd.HandlePeers,
	"handshake":        cmd.HandleHandshake,
	"download_piece":   cmd.HandleDownloadPiece,
	"download":         cmd.HandleDownload,
	"magnet_parse":     cmd.HandleMagnetParse,
	"magnet_handshake": cmd.HandleMagnetHandshake,
}

func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		fmt.Printf("no arguments provided")
		return
	}

	handler, ok := handlers[args[0]]
	if !ok {
		fmt.Printf("unknown command: %s\n", args[0])
		return
	}
	handler(args[1:])
}
