package cmd

import (
	"fmt"
	"os"

	"github.com/EshaanAgg/toy-bittorrent/app/bencode"
	"github.com/EshaanAgg/toy-bittorrent/app/utils"
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

func printInfoRelatedFields(infoDict *bencode.BencodeDictionary) {
	infoHash, err := utils.SHA1Hash(infoDict.Encode())
	if err != nil {
		fmt.Printf("error hashing the info dictionary: %v", err)
		return
	}

	fmt.Printf("Info Hash: %s\n", infoHash)
	fmt.Printf("Piece Length: %d\n", infoDict.Map["piece length"].GetInteger().Value)
	pieces := infoDict.Map["pieces"].GetString().Value
	printPieceHashes(pieces)
}

func HandleInfo(args []string) {
	// Validate the number of arguments passed to the info command
	if len(args) == 0 {
		fmt.Println("no data passed to info. usage: kafka info <path-to-file>")
		return
	}

	if len(args) > 1 {
		fmt.Println("too many arguments passed to info. usage: kafka info <path-to-file>")
		return
	}

	// Read the file content
	filePath := args[0]
	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Printf("error reading file: %v", err)
		return
	}

	// Decode the file content as a dictionary
	bd, err := bencode.NewBencodeData(fileContent)
	if err != nil {
		fmt.Printf("error decoding the passed data: %v", err)
		return
	}

	// Access the nested elements directly as we can be assured
	// that the file passed is a valid torrent file
	d := bd.GetDictionary()
	trackerUrl := d.Map["announce"].GetString().Value
	fmt.Printf("Tracker URL: %s\n", trackerUrl)

	infoDict := d.Map["info"].GetDictionary()
	fileSize := infoDict.Map["length"].GetInteger().Value
	fmt.Printf("Length: %d\n", fileSize)

	printInfoRelatedFields(infoDict)
}
