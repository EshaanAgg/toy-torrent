package types

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
)

// Peer represents a remote peer in the network.
type Peer struct {
	IP   string
	Port int

	conn                net.Conn
	logger              *log.Logger
	shutMessageListener chan bool

	assignedPiece *StoredPiece
}

// NewPeerFromAddr initializes a Peer and establishes a TCP connection to it.
func NewPeerFromAddr(addr string) (*Peer, error) {
	host, port, err := net.SplitHostPort(addr)
	if err != nil {
		return nil, fmt.Errorf("invalid address format: %w", err)
	}
	portInt, err := strconv.Atoi(port)
	if err != nil {
		return nil, fmt.Errorf("invalid port number: %w", err)
	}

	// Create a new connection to the peer
	conn, err := net.Dial("tcp", net.JoinHostPort(host, fmt.Sprintf("%d", portInt)))
	if err != nil {
		return nil, fmt.Errorf("error connecting to address %s: %w", addr, err)
	}

	logger := log.New(conn, fmt.Sprintf("[Peer %d] ", getPeerID(addr)), 0)
	logger.SetOutput(log.Writer())

	return &Peer{
		IP:                  host,
		Port:                portInt,
		conn:                conn,
		logger:              logger,
		shutMessageListener: make(chan bool),
	}, nil
}

// SendMessage sends a message to the peer with a 4-byte length prefix.
func (p *Peer) SendMessage(messageBytes []byte) error {
	// Prepend the length of the message as a 4-byte integer
	length := uint32(len(messageBytes))
	data := make([]byte, 4)
	binary.BigEndian.PutUint32(data, length)

	// Append the message bytes and send
	data = append(data, messageBytes...)
	_, err := p.conn.Write(data)
	if err != nil {
		return fmt.Errorf("error sending message: %w", err)
	}

	return nil
}

func (p *Peer) readExactBytes(n uint32) ([]byte, error) {
	data := make([]byte, n)
	_, err := io.ReadFull(p.conn, data)
	if err != nil {
		return nil, fmt.Errorf("error reading %d bytes: %w", n, err)
	}
	return data, nil
}

// RecieveMessage reads a message from the peer.
// It first reads the 4-byte length prefix, then reads the message of that length.
func (p *Peer) RecieveMessage() ([]byte, error) {
	lengthPrefix := make([]byte, 4)
	_, err := p.conn.Read(lengthPrefix)
	if err != nil {
		return nil, fmt.Errorf("error reading message length: %w", err)
	}
	length := binary.BigEndian.Uint32(lengthPrefix)

	// Keep-alive message, so recursively call itself
	if length == 0 {
		return p.RecieveMessage()
	}

	message, err := p.readExactBytes(length)
	if err != nil {
		return nil, fmt.Errorf("error reading message: %w", err)
	}
	return message, nil
}
