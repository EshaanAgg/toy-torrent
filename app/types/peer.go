package types

import (
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"strconv"
)

// Peer represents a remote peer in the network.
type Peer struct {
	IP   string
	Port int

	conn     net.Conn
	logger   *log.Logger
	pieceMap map[uint32]*StoredPiece
}

// NewPeerFromAddr initializes a Peer and establishes a TCP connection to it.
func NewPeerFromAddr(addr string) (*Peer, error) {
	host, port, err := net.SplitHostPort(addr)
	if err != nil {
		return nil, fmt.Errorf("invalid address format: %v", err)
	}
	portInt, err := strconv.Atoi(port)
	if err != nil {
		return nil, fmt.Errorf("invalid port number: %v", err)
	}

	// Create a new connection to the peer
	conn, err := net.Dial("tcp", net.JoinHostPort(host, fmt.Sprintf("%d", portInt)))
	if err != nil {
		return nil, fmt.Errorf("error connecting to address %s: %w", addr, err)
	}

	logger := log.New(conn, fmt.Sprintf("[%s:%d] ", host, portInt), 0)
	logger.SetOutput(log.Writer())

	return &Peer{
		IP:       host,
		Port:     portInt,
		conn:     conn,
		logger:   logger,
		pieceMap: make(map[uint32]*StoredPiece),
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
		return fmt.Errorf("error sending message: %v", err)
	}

	return nil
}

// RecieveMessage reads a message from the peer.
// It first reads the 4-byte length prefix, then reads the message of that length.
func (p *Peer) RecieveMessage() ([]byte, error) {
	lengthPrefix := make([]byte, 4)
	_, err := p.conn.Read(lengthPrefix)
	if err != nil {
		return nil, fmt.Errorf("error reading message length: %v", err)
	}
	length := binary.BigEndian.Uint32(lengthPrefix)

	// Keep-alive message, so recursively call itself
	if length == 0 {
		return p.RecieveMessage()
	}

	message := make([]byte, length)
	_, err = p.conn.Read(message)
	if err != nil {
		return nil, fmt.Errorf("error reading message body: %v", err)
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
		return fmt.Errorf("error while sending interested message: %v", err)
	}

	return nil
}

// CloseConnection closes the TCP connection to the peer.
func (p *Peer) CloseConnection() error {
	err := p.conn.Close()
	if err != nil {
		return err
	}
	p.conn = nil
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
	p.Log("Initializing handshake")

	_, err := p.conn.Write(handshake.Bytes())
	if err != nil {
		return nil, fmt.Errorf("error sending handshake: %v", err)
	}

	response := make([]byte, 68) // Handshake response size
	_, err = p.conn.Read(response)
	if err != nil {
		return nil, fmt.Errorf("error reading handshake response: %v", err)
	}
	recievedHandshake := NewHandshakeFromBytes(response)
	p.Log("Handshaking completed")

	return recievedHandshake, nil
}

// SendInterested sends an "interested" message to the peer indicating that
// we want to download data from them. It also waits for the unchoke message
// indicating that the peer is ready to send us data.
func (p *Peer) SendInterested() error {
	interested := []byte{INTERESTED_MESSAGE_ID}
	err := p.SendMessage(interested)
	if err != nil {
		return fmt.Errorf("error sending interested message: %v", err)
	}
	p.Log("Sent interested message")

	err = p.blockTillUnchokeMessage()
	if err != nil {
		return fmt.Errorf("error while wating for unchoke message: %v", err)
	}
	return nil
}

// blockTillBitFieldMessage blocks the peer until a bitfield message (ID = 5) is received.
// It skips over other messages and discards their payloads.
func (p *Peer) blockTillBitFieldMessage() error {
	for {
		msg, err := p.RecieveMessage()
		if err != nil {
			return fmt.Errorf("error receiving message: %v", err)
		}

		if msg[0] == BITFIELD_MESSAGE_ID {
			// TODO: Parse the response to get all the pieces present
			p.Log("Recieved bitfield message")
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
			return fmt.Errorf("error receiving message: %v", err)
		}

		if msg[0] == UNCHOKE_MESSAGE_ID {
			p.Log("Recieved unchoke message")
			return nil
		}

		p.Log("Recieved message bytes: %q while waiting for unchoke message", msg)
	}
}

type CompletePieceCallback func(sp *StoredPiece)

// RegisterPieceMessageHandler registers a callback to be called when a piece message is received.
// It continuously listens for incoming messages and processes them.
// When a piece is downloaded completely, it calls the provided callback function.
func (p *Peer) RegisterPieceMessageHandler(callback CompletePieceCallback) {
	for {
		message, err := p.RecieveMessage()
		if err != nil {
			p.Log("error receiving message: %v", err)
		}

		if message[0] != PIECE_MESSAGE_ID {
			p.Log("recieved message with ID = %d, expected PIECE_MESSAGE_ID = %d", message[0], PIECE_MESSAGE_ID)
			continue
		}

		m, err := NewPieceMessage(message)
		if err != nil {
			p.Log("error creating PieceMessage: %v", err)
		}
		p.Log("recieved piece message for piece_idx = %d, begin_idx = %d, block_len = %d", m.PieceIndex, m.Begin/BLOCK_SIZE, len(m.Block))

		piece, ok := p.pieceMap[m.PieceIndex]
		if !ok {
			p.Log("piece with index %d not found in piece map", m.PieceIndex)
			continue
		}

		err = piece.HandlePieceMessage(m)
		if err != nil {
			p.Log("error handling piece message: %v", err)
		}

		if piece.IsComplete() {
			p.Log("piece %d is complete", piece.Index)
			callback(piece)
			delete(p.pieceMap, piece.Index)
		}
	}
}
