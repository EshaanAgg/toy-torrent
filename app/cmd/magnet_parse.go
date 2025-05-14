package cmd

import (
	"fmt"

	"github.com/EshaanAgg/toy-bittorrent/app/types"
)

func HandleMagnetParse(args []string, s *types.Server) {
	if len(args) != 1 {
		fmt.Println("incorrect arguments passed. usage: go-torrent magnet_parse <magnet-link>")
		return
	}

	magnetLink := args[0]
	m, err := types.NewMagnetURI(magnetLink)
	if err != nil {
		fmt.Printf("error creating MagnetURI: %v\n", err)
		return
	}

	fmt.Printf("Tracker URL: %s\n", m.TrackerURL)
	fmt.Printf("Info Hash: %s\n", m.InfoHashHex)
}
