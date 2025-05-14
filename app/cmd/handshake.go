package cmd

import (
	"fmt"

	"github.com/EshaanAgg/toy-bittorrent/app/types"
	"github.com/EshaanAgg/toy-bittorrent/app/utils"
)

func HandleHandshake(args []string) {
	if len(args) != 2 {
		fmt.Println("incorrect arguments passed. usage: go-torrent handshake <torrent-file> <peer-host:port>")
		return
	}

	torrentFile := args[0]
	fileInfo, err := types.NewTorrentFileInfo(torrentFile)
	if err != nil {
		fmt.Printf("error creating TorrentFileInfo: %v\n", err)
		return
	}

	peer, err := types.NewPeerFromAddr(args[1])
	if err != nil {
		fmt.Printf("error creating peer: %v\n", err)
		return
	}

	handshake, err := peer.PerformHandshake(fileInfo.InfoHash)
	if err != nil {
		fmt.Printf("error performing handshake: %v\n", err)
		return
	}
	fmt.Println("Peer ID:", utils.BytesToHex(handshake.PeerID))
}
