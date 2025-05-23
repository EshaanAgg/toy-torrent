package types

import (
	"fmt"
	"os"

	"github.com/EshaanAgg/toy-bittorrent/app/bencode"
	"github.com/EshaanAgg/toy-bittorrent/app/utils"
)

type InfoDict struct {
	Length      int
	Name        string
	PieceLength int
	Pieces      [][]byte // Hashes of each piece
}

type TorrentFileInfo struct {
	TrackerURL string
	InfoDict   *InfoDict
	InfoHash   []byte
}

func piecesFromString(pieces []byte) [][]byte {
	// Each piece hash is 20 bytes long
	pieceLength := 20
	numPieces := len(pieces) / pieceLength

	pieceHashes := make([][]byte, numPieces)
	for i := range numPieces {
		start := i * pieceLength
		end := start + pieceLength
		pieceHashes[i] = []byte(pieces[start:end])
	}

	return pieceHashes
}

// NewTorrentFileInfo creates a new TorrentFileInfo struct from the given torrent file path.
func NewTorrentFileInfo(torrentFilePath string) (*TorrentFileInfo, error) {
	fileContent, err := os.ReadFile(torrentFilePath)
	if err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	bd, err := bencode.NewBencodeData(fileContent)
	if err != nil {
		return nil, fmt.Errorf("error decoding the file: %w", err)
	}

	// Access the nested elements directly as we can be assured
	// that the file passed is a valid torrent file
	d := bd.GetDictionary()
	infoDict := d.Map["info"].GetDictionary()

	// Parse the info dictionary to get the info hash
	infoHash, err := utils.SHA1Hash(infoDict.Encode())
	if err != nil {
		return nil, fmt.Errorf("error hashing the info dictionary: %w", err)
	}

	return &TorrentFileInfo{
		TrackerURL: string(d.Map["announce"].GetString().Value),
		InfoHash:   infoHash,
		InfoDict: &InfoDict{
			Length:      infoDict.Map["length"].GetInteger().Value,
			Name:        string(infoDict.Map["name"].GetString().Value),
			PieceLength: infoDict.Map["piece length"].GetInteger().Value,
			Pieces:      piecesFromString(infoDict.Map["pieces"].GetString().Value),
		},
	}, nil
}

func NewTorrentFileInfoFromMagnet(magnet *MagnetURI, infoDict *bencode.BencodeDictionary) (*TorrentFileInfo, error) {
	return &TorrentFileInfo{
		TrackerURL: magnet.TrackerURL,
		InfoHash:   magnet.InfoHash,
		InfoDict: &InfoDict{
			Length:      infoDict.Map["length"].GetInteger().Value,
			Name:        string(infoDict.Map["name"].GetString().Value),
			PieceLength: infoDict.Map["piece length"].GetInteger().Value,
			Pieces:      piecesFromString(infoDict.Map["pieces"].GetString().Value),
		},
	}, nil
}

func (t *TorrentFileInfo) GetHexInfoHash() string {
	return fmt.Sprintf("%x", t.InfoHash)
}
