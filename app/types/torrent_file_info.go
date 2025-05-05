package types

import (
	"fmt"
	"os"

	"github.com/EshaanAgg/toy-bittorrent/app/bencode"
	"github.com/EshaanAgg/toy-bittorrent/app/utils"
)

type TorrentFileInfo struct {
	TrackerURL string
	FileSize   int
	InfoHash   []byte
	InfoDict   *bencode.BencodeDictionary
}

// NewTorrentFileInfo creates a new TorrentFileInfo struct from the given torrent file path.
func NewTorrentFileInfo(torrentFilePath string) (*TorrentFileInfo, error) {
	fileContent, err := os.ReadFile(torrentFilePath)
	if err != nil {
		return nil, fmt.Errorf("error reading file: %v", err)
	}

	// Decode the file content as a dictionary
	bd, err := bencode.NewBencodeData(fileContent)
	if err != nil {
		return nil, fmt.Errorf("error decoding the file: %v", err)
	}

	// Access the nested elements directly as we can be assured
	// that the file passed is a valid torrent file
	d := bd.GetDictionary()
	trackerUrl := d.Map["announce"].GetString().Value
	infoDict := d.Map["info"].GetDictionary()
	fileSize := infoDict.Map["length"].GetInteger().Value

	infoHash, err := utils.SHA1Hash(infoDict.Encode())
	if err != nil {
		return nil, fmt.Errorf("error hashing the info dictionary: %v", err)
	}

	return &TorrentFileInfo{
		TrackerURL: string(trackerUrl),
		FileSize:   fileSize,
		InfoHash:   infoHash,
		InfoDict:   infoDict,
	}, nil
}

func (t *TorrentFileInfo) GetHexInfoHash() string {
	return fmt.Sprintf("%x", t.InfoHash)
}
