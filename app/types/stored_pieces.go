package types

import (
	"encoding/binary"
	"fmt"
	"net"
	"sync/atomic"
)

// pieceBlock represents a block of data received from a peer.
// Each pieceBlock should be managed only one co-routine at a time.
type pieceBlock struct {
	byteOffset uint32 // The offset of the block in the piece
	length     uint32 // The length of the block in bytes

	data     []byte
	recieved bool
}

func newPieceBlock(byteOffset, length uint32) *pieceBlock {
	return &pieceBlock{
		byteOffset: byteOffset,
		length:     length,
		data:       make([]byte, length),
		recieved:   false,
	}
}

func (pb *pieceBlock) setData(data []byte) {
	pb.data = data
	pb.recieved = true
}

func (pb *pieceBlock) makeRequest(conn net.Conn, pieceIdx uint32) error {
	message := make([]byte, 13)

	message[0] = REQUEST_MESSAGE_ID
	binary.BigEndian.PutUint32(message[1:], pieceIdx)
	binary.BigEndian.PutUint32(message[5:], pb.byteOffset)
	binary.BigEndian.PutUint32(message[9:], pb.length)

	_, err := conn.Write(message)
	if err != nil {
		return fmt.Errorf("error sending request message: %v", err)
	}

	return nil
}

// StoredPiece represents a piece of data that needs to be downloaded.
type StoredPiece struct {
	Index          uint32
	Length         uint32 // The total length of the piece, in bytes
	NumberOfBlocks uint32 // The number of blocks in the piece

	Blocks             []*pieceBlock
	RecievedBlockCount atomic.Uint32 // The number of blocks received
}

func NewStoredPiece(index, length uint32) *StoredPiece {
	sp := &StoredPiece{
		Index:          index,
		Length:         length,
		NumberOfBlocks: (length + BLOCK_SIZE - 1) / BLOCK_SIZE,
		Blocks:         make([]*pieceBlock, 0),
	}

	currentOffset := uint32(0)
	for range sp.NumberOfBlocks {
		blockLength := BLOCK_SIZE
		if currentOffset+BLOCK_SIZE > length {
			blockLength = length - currentOffset
		}
		sp.Blocks = append(sp.Blocks, newPieceBlock(uint32(currentOffset), uint32(blockLength)))
		currentOffset += blockLength
	}

	return sp
}

func (sp *StoredPiece) IsComplete() bool {
	return sp.RecievedBlockCount.Load() == sp.NumberOfBlocks
}

func (sp *StoredPiece) HandlePieceMessage(m *PieceMessage) error {
	if m.PieceIndex != uint32(sp.Index) {
		return fmt.Errorf("called handlePieceMessage function for StoredPiece with index = %d, while the message was for pieceIdx = %d", sp.Index, m.PieceIndex)
	}

	blockIdx := m.Begin / BLOCK_SIZE
	if blockIdx >= sp.NumberOfBlocks {
		return fmt.Errorf("The message has offset = %d (block index = %d), whereas this piece only has %d blocks", m.Begin, blockIdx, sp.NumberOfBlocks)
	}

	if sp.Blocks[blockIdx].length != uint32(len(m.Block)) {
		return fmt.Errorf("The block at index %d is expected to have length %d, but the message data has length %d", blockIdx, sp.Blocks[blockIdx].length, len(m.Block))
	}

	sp.Blocks[blockIdx].setData(m.Block)
	sp.RecievedBlockCount.Add(1)

	return nil
}

func (sp *StoredPiece) GetData() []byte {
	if !sp.IsComplete() {
		fmt.Printf("warning: called GetData on an incompletely downloaded StoredPiece, this may cause incorrect data being returned")
	}

	d := make([]byte, 0)
	for _, b := range sp.Blocks {
		d = append(d, b.data...)
	}
	return d
}
