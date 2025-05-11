package types

import (
	"github.com/EshaanAgg/toy-bittorrent/app/utils"
)

type Server struct {
	PeerID []byte
}

func NewServer() *Server {
	return &Server{
		PeerID: utils.GetRandomPeerID(),
	}
}
