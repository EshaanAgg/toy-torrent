package cmd

import (
	"fmt"

	"github.com/EshaanAgg/toy-bittorrent/app/types"
	"github.com/EshaanAgg/toy-bittorrent/app/utils"
)

func HandlePeers(args []string) {
	// Validate the number of arguments passed to the info command
	if len(args) == 0 {
		fmt.Println("no data passed to info. usage: go-torrent peers <path-to-file>")
		return
	}

	if len(args) > 1 {
		fmt.Println("too many arguments passed to info. usage: go-torrent peers <path-to-file>")
		return
	}

	fileInfo, err := types.NewTorrentFileInfo(args[0])
	if err != nil {
		fmt.Printf("error creating TorrentFileInfo: %v\n", err)
		return
	}

	req := types.TrackerGetRequest{
		TrackerURL: fileInfo.TrackerURL,
		InfoHash:   fileInfo.InfoHash,
		PeerID:     utils.GetRandomPeerID(),
		Port:       6881,
		Uploaded:   0,
		Downloaded: 0,
		Left:       fileInfo.FileSize,
		Compact:    1,
	}
	resp, err := req.MakeRequest()
	if err != nil {
		fmt.Printf("error making request to tracker: %v\n", err)
		return
	}

	for _, peer := range resp.Peers {
		fmt.Printf("%s:%d\n", peer.IP, peer.Port)
	}
}
