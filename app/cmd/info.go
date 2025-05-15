package cmd

import (
	"fmt"

	"github.com/EshaanAgg/toy-bittorrent/app/types"
)

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

	logInfo(fileInfo)
}
