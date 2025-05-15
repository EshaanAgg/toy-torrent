package cmd

import (
	"fmt"

	"github.com/EshaanAgg/toy-bittorrent/app/types"
)

// getPieceLength calculates the length of a piece in a torrent file.
// It takes into account the last piece which may be shorter than the others.
func getPieceLength(fileInfo *types.TorrentFileInfo, pieceIdx int) uint32 {
	pieceLength := fileInfo.InfoDict.PieceLength
	fileLength := fileInfo.InfoDict.Length
	if (pieceIdx+1)*pieceLength > fileLength {
		pieceLength = fileLength - (pieceIdx * pieceLength)
	}
	return uint32(pieceLength)
}

// getPeers makes a request to the tracker to get a list of peers.
// makeConnection is a boolean that indicates whether to make a connection to the
// generated peers or not.
func getPeers(trackerURL string, infoHash []byte, leftLength int, makeConnection bool) ([]*types.Peer, error) {
	req := types.TrackerGetRequest{
		TrackerURL: trackerURL,
		InfoHash:   infoHash,
		PeerID:     string(types.SERVER_PEER_ID),
		Port:       6881,
		Uploaded:   0,
		Downloaded: 0,
		Compact:    1,
		Left:       leftLength,
	}
	resp, err := req.MakeRequest(makeConnection)
	if err != nil {
		return nil, fmt.Errorf("error making request to tracker: %w", err)
	}

	return resp.Peers, nil
}

// getPeersFromFile is a wrapper function around getPeers.
func getPeersFromFile(fileInfo *types.TorrentFileInfo, makeConnection bool) ([]*types.Peer, error) {
	peers, err := getPeers(fileInfo.TrackerURL, fileInfo.InfoHash, fileInfo.InfoDict.Length, makeConnection)
	if err != nil {
		return nil, fmt.Errorf("error getting peers from tracker: %w", err)
	}
	return peers, nil
}

func logInfo(fileInfo *types.TorrentFileInfo) {
	infoDict := fileInfo.InfoDict

	fmt.Printf("Tracker URL: %s\n", fileInfo.TrackerURL)
	fmt.Printf("Length: %d\n", infoDict.Length)
	fmt.Printf("Info Hash: %s\n", fileInfo.GetHexInfoHash())
	fmt.Printf("Piece Length: %d\n", infoDict.PieceLength)
	fmt.Println("Piece Hashes:")
	for _, p := range infoDict.Pieces {
		fmt.Printf("%x\n", p)
	}
}
