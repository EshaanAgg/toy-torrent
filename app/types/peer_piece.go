package types

import (
	"errors"
	"fmt"
	"io"
)

// PrepareToGetPieceData is to be used for stages without magnet links.
func (p *Peer) PrepareToGetPieceData(infoHash []byte) error {
	_, err := p.PerformHandshake(infoHash)
	if err != nil {
		return fmt.Errorf("error performing handshake: %w", err)
	}

	// TODO: Send a bitfield message to the peer
	err = p.BlockTillBitFieldMessage()
	if err != nil {
		return fmt.Errorf("error waiting for bitfield message: %w", err)
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
