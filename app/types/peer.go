package types

import (
	"fmt"
	"net"
	"strconv"
)

type PeerState int

const (
	PEER_STATE_IDLE PeerState = iota
	PEER_STATE_HANDSHAKED
)

type Peer struct {
	IP    string
	Port  int
	State PeerState
}

func NewPeer(ip string, port int) *Peer {
	return &Peer{
		IP:    ip,
		Port:  port,
		State: PEER_STATE_IDLE,
	}
}

func NewPeerFromAddr(addr string) (*Peer, error) {
	host, port, err := net.SplitHostPort(addr)
	if err != nil {
		return nil, fmt.Errorf("invalid address format: %v", err)
	}
	portInt, err := strconv.Atoi(port)
	if err != nil {
		return nil, fmt.Errorf("invalid port number: %v", err)
	}
	return &Peer{
		IP:    host,
		Port:  portInt,
		State: PEER_STATE_IDLE,
	}, nil
}

func (p *Peer) GetConnection(s *Server) (net.Conn, error) {
	// Check if the peer is already connected
	if conn, exists := s.PeerConnections[p]; exists {
		return conn, nil
	}

	// Create a new connection to the peer
	conn, err := net.Dial("tcp", net.JoinHostPort(p.IP, fmt.Sprintf("%d", p.Port)))
	if err != nil {
		return nil, err
	}
	s.PeerConnections[p] = conn
	return conn, nil
}

func (p *Peer) CloseConnection(s *Server) error {
	if conn, exists := s.PeerConnections[p]; exists {
		err := conn.Close()
		if err != nil {
			return err
		}
		delete(s.PeerConnections, p)
	}
	return nil
}

func (p *Peer) PerformHandshake(s *Server, infoHash []byte) (*Handshake, error) {
	conn, err := p.GetConnection(s)
	if err != nil {
		return nil, fmt.Errorf("error getting connection: %v", err)
	}

	handshake := Handshake{
		PeerID:   s.PeerID,
		InfoHash: infoHash,
	}
	_, err = conn.Write(handshake.Bytes())
	if err != nil {
		return nil, fmt.Errorf("error sending handshake: %v", err)
	}

	response := make([]byte, 68) // Handshake response size
	_, err = conn.Read(response)
	if err != nil {
		return nil, fmt.Errorf("error reading handshake response: %v", err)
	}
	recievedHandshake := NewHandshakeFromBytes(response)
	p.State = PEER_STATE_HANDSHAKED
	return recievedHandshake, nil
}
