package cmd

import (
	"fmt"

	"github.com/EshaanAgg/toy-bittorrent/app/types"
)

// Parses the torrent file and retrieves the peers from the tracker.
func getPeers(fileInfo *types.TorrentFileInfo, makeConnection bool) ([]*types.Peer, error) {
	req := types.TrackerGetRequest{
		TrackerURL: fileInfo.TrackerURL,
		InfoHash:   fileInfo.InfoHash,
		PeerID:     string(types.SERVER_PEER_ID),
		Port:       6881,
		Uploaded:   0,
		Downloaded: 0,
		Compact:    1,
		Left:       fileInfo.InfoDict.Length,
	}
	resp, err := req.MakeRequest(makeConnection)
	if err != nil {
		return nil, fmt.Errorf("error making request to tracker: %w", err)
	}

	return resp.Peers, nil
}

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
		fmt.Printf("error creating TorrentFileInfo: %v", err)
		return
	}

	peers, err := getPeers(fileInfo, false)
	if err != nil {
		fmt.Printf("error getting peers: %v\n", err)
		return
	}
	for _, peer := range peers {
		fmt.Printf("%s:%d\n", peer.IP, peer.Port)
	}
}
