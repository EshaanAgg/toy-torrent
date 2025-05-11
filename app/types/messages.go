package types

import (
	"encoding/binary"
	"fmt"
)

type PieceMessage struct {
	PieceIndex uint32 // The index of the piece
	Begin      uint32 // The offset within the piece
	Block      []byte // The block of data
}

func NewPieceMessage(data []byte) (*PieceMessage, error) {
	if len(data) < 9 {
		return nil, fmt.Errorf("data too short for piece message, expected at least 9 bytes, got %d", len(data))
	}
	if data[0] != PIECE_MESSAGE_ID {
		return nil, fmt.Errorf("invalid message ID for piece message, expected %d, got %d", PIECE_MESSAGE_ID, data[0])
	}

	pieceIdx := binary.BigEndian.Uint32(data[1:5])
	begin := binary.BigEndian.Uint32(data[5:9])
	block := data[9:]

	return &PieceMessage{
		PieceIndex: pieceIdx,
		Begin:      begin,
		Block:      block,
	}, nil
}
