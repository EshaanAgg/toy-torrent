package types

import (
	"net"

	"github.com/EshaanAgg/toy-bittorrent/app/utils"
)

type Server struct {
	PeerID          []byte
	PeerConnections map[*Peer]net.Conn
}

func NewServer() *Server {
	return &Server{
		PeerID: utils.GetRandomPeerID(),
	}
}
