package main

import (
	"fmt"
	"os"

	"github.com/EshaanAgg/toy-bittorrent/app/cmd"
)

func main() {
	args := os.Args[1:]

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
		cmd.HandlePeers(args[1:])

	default:
		fmt.Printf("unrecognized command '%s'", args[0])
	}
}
