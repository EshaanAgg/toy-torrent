package types

import (
	"encoding/binary"
	"errors"
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

func (p *Peer) PrepareToGetPieceData(s *Server, infoHash []byte) error {
	_, err := p.PerformHandshake(s, infoHash)
	if err != nil {
		return fmt.Errorf("error performing handshake: %w", err)
	}

	err = p.blockTillBitFieldMessage()
	if err != nil {
		return fmt.Errorf("error while waiting for bitfield message: %w", err)
	}

	err = p.SendInterested()
	if err != nil {
		return fmt.Errorf("error while sending interested message: %w", err)
	}

	p.Log("completed initialization. ready to download pieces")
	return nil
}

func (p *Peer) Log(s string, vals ...any) {
	p.logger.Printf(s+"\n", vals...)
}

// PerformHandshake sends a handshake to the peer and waits for a response.
// Returns the received handshake or an error.
func (p *Peer) PerformHandshake(s *Server, infoHash []byte) (*Handshake, error) {
	handshake := Handshake{
		PeerID:   s.PeerID,
		InfoHash: infoHash,
	}

	_, err := p.conn.Write(handshake.Bytes())
	if err != nil {
		return nil, fmt.Errorf("error sending handshake: %w", err)
	}

	response, err := p.readExactBytes(68) // 68 bytes for the handshake response
	if err != nil {
		return nil, fmt.Errorf("error reading handshake response: %w", err)
	}
	recievedHandshake := NewHandshakeFromBytes(response)

	return recievedHandshake, nil
}

// SendInterested sends an "interested" message to the peer indicating that
// we want to download data from them. It also waits for the unchoke message
// indicating that the peer is ready to send us data.
func (p *Peer) SendInterested() error {
	interested := []byte{INTERESTED_MESSAGE_ID}
	err := p.SendMessage(interested)
	if err != nil {
		return fmt.Errorf("error sending interested message: %w", err)
	}

	err = p.blockTillUnchokeMessage()
	if err != nil {
		return fmt.Errorf("error while wating for unchoke message: %w", err)
	}
	return nil
}

// blockTillBitFieldMessage blocks the peer until a bitfield message (ID = 5) is received.
// It skips over other messages and discards their payloads.
func (p *Peer) blockTillBitFieldMessage() error {
	for {
		msg, err := p.RecieveMessage()
		if err != nil {
			return fmt.Errorf("error receiving message: %w", err)
		}

		if msg[0] == BITFIELD_MESSAGE_ID {
			// TODO: Parse the response to get all the pieces present
			return nil

		}

		p.Log("Recieved message bytes: %q while waiting for bitfield message", msg)
	}
}

// blockTillUnchokeMessage blocks the peer until an unchoke message (ID = 1) is received.
// It skips over other messages and discards their payloads.
func (p *Peer) blockTillUnchokeMessage() error {
	for {
		msg, err := p.RecieveMessage()
		if err != nil {
			return fmt.Errorf("error receiving message: %w", err)
		}

		if msg[0] == UNCHOKE_MESSAGE_ID {
			return nil
		}

		p.Log("Recieved message bytes: %q while waiting for unchoke message", msg)
	}
}

// getCompletePiece continuously listens for incoming messages and processes them.
// It returns when the piece download is complete or an error occurs.
func (p *Peer) getCompletePiece() error {
	if p.assignedPiece == nil {
		return fmt.Errorf("no piece assigned to peer")
	}
	piece := p.assignedPiece

	for {
		message, err := p.RecieveMessage()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return fmt.Errorf("peer closed connection: %w", err)
			}
			return fmt.Errorf("error receiving message: %w", err)
		}

		if message[0] != PIECE_MESSAGE_ID {
			return fmt.Errorf("recieved message with ID = %d, expected PIECE_MESSAGE_ID = %d", message[0], PIECE_MESSAGE_ID)
		}

		m, err := NewPieceMessage(message)
		if err != nil {
			return fmt.Errorf("error creating PieceMessage: %w", err)
		}

		err = piece.HandlePieceMessage(m)
		if err != nil {
			return fmt.Errorf("error handling piece message: %w", err)
		}

		if piece.IsComplete() {
			p.Log("piece %d completed", piece.Index)
			err := piece.VerifyHash()
			if err != nil {
				return fmt.Errorf("error verifying piece hash: %w", err)
			}
			p.Log("piece %d hash verified", piece.Index)
			return nil
		}
	}
}
