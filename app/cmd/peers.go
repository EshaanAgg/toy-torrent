package cmd

import (
	"fmt"

	"github.com/EshaanAgg/toy-bittorrent/app/types"
)

// Parses the torrent file and retrieves the peers from the tracker.
func getPeers(torrentFile string, server *types.Server) ([]*types.Peer, error) {
	fileInfo, err := types.NewTorrentFileInfo(torrentFile)
	if err != nil {
		return nil, fmt.Errorf("error creating TorrentFileInfo: %v", err)
	}

	req := types.TrackerGetRequest{
		TrackerURL: fileInfo.TrackerURL,
		InfoHash:   fileInfo.InfoHash,
		PeerID:     string(server.PeerID),
		Port:       6881,
		Uploaded:   0,
		Downloaded: 0,
		Compact:    1,
		Left:       fileInfo.FileSize,
	}
	resp, err := req.MakeRequest()
	if err != nil {
		return nil, fmt.Errorf("error making request to tracker: %v", err)
	}
	return resp.Peers, nil

}

func HandlePeers(args []string, server *types.Server) {
	// Validate the number of arguments passed to the info command
	if len(args) == 0 {
		fmt.Println("no data passed to info. usage: go-torrent peers <path-to-file>")
		return
	}

	if len(args) > 1 {
		fmt.Println("too many arguments passed to info. usage: go-torrent peers <path-to-file>")
		return
	}

	peers, err := getPeers(args[0], server)
	if err != nil {
		fmt.Printf("error getting peers: %v\n", err)
		return
	}
	for _, peer := range peers {
		fmt.Printf("%s:%d\n", peer.IP, peer.Port)
	}
}
