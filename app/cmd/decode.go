package cmd

import (
	"fmt"

	"github.com/EshaanAgg/toy-bittorrent/app/bencode"
)

func HandleDecode(args []string) {
	if len(args) == 0 {
		fmt.Println("no data passed to decode. usage: kafka decode <bytes-to-decode>")
		return
	}

	d := args[0]
	bd, err := bencode.NewBencodeData([]byte(d))
	if err != nil {
		fmt.Printf("error decoding the passed data: %v", err)
		return
	}

	// Print the string representation of the decoded data to STDIN
	fmt.Println(bd.Value.String())
}
