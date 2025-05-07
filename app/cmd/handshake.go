package cmd

import (
	"fmt"
	"net"

	"github.com/EshaanAgg/toy-bittorrent/app/types"
	"github.com/EshaanAgg/toy-bittorrent/app/utils"
)

func performHandshake(torrentFile string, peerAddress string, server *types.Server) (*types.Handshake, error) {
	fileInfo, err := types.NewTorrentFileInfo(torrentFile)
	if err != nil {
		return nil, fmt.Errorf("error creating TorrentFileInfo: %v", err)
	}

	handshake := types.Handshake{
		PeerID:   server.PeerID,
		InfoHash: fileInfo.InfoHash,
	}

	conn, err := net.Dial("tcp", peerAddress)
	if err != nil {
		return nil, fmt.Errorf("error connecting to peer: %v", err)
	}
	defer conn.Close()

	// Send the handshake message to the peer
	_, err = conn.Write(handshake.Bytes())
	if err != nil {
		return nil, fmt.Errorf("error sending handshake: %v", err)
	}

	// Read the response from the peer
	response := make([]byte, 68) // Handshake response size
	_, err = conn.Read(response)
	if err != nil {
		return nil, fmt.Errorf("error reading handshake response: %v", err)
	}
	recievedHandshake := types.NewHandshakeFromBytes(response)
	return recievedHandshake, nil
}

func HandleHandshake(args []string, server *types.Server) {
	if len(args) != 2 {
		fmt.Println("incorrect arguments passed. usage: go-torrent handshake <torrent-file> <peer-address>")
		return
	}

	torrentFile := args[0]
	peerAddress := args[1]

	handshake, err := performHandshake(torrentFile, peerAddress, server)
	if err != nil {
		fmt.Printf("error performing handshake: %v\n", err)
		return
	}

	fmt.Println("Peer ID: ", utils.BytesToHex(handshake.PeerID))
}
