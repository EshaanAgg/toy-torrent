package types

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"sync/atomic"

	"github.com/EshaanAgg/toy-bittorrent/app/utils"
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

func (pb *pieceBlock) makeRequest(p *Peer, pieceIdx uint32) error {
	message := make([]byte, 13)

	message[0] = REQUEST_MESSAGE_ID
	binary.BigEndian.PutUint32(message[1:], pieceIdx)
	binary.BigEndian.PutUint32(message[5:], pb.byteOffset)
	binary.BigEndian.PutUint32(message[9:], pb.length)

	err := p.SendMessage(message)
	if err != nil {
		return fmt.Errorf("error sending request message: %w", err)
	}

	return nil
}

// StoredPiece represents a piece of data that needs to be downloaded.
type StoredPiece struct {
	Index  uint32
	Length uint32 // The total length of the piece, in bytes
	Hash   []byte

	NumberOfBlocks     uint32 // The number of blocks in the piece
	Blocks             []*pieceBlock
	RecievedBlockCount atomic.Uint32 // The number of blocks received

	peerConn net.Conn
}

func (p *Peer) DownloadPiece(index, length uint32, hash []byte) (*StoredPiece, error) {
	if p.assignedPiece != nil {
		return nil, fmt.Errorf("peer already has an assigned piece")
	}

	sp := &StoredPiece{
		Index:  index,
		Length: length,
		Hash:   hash,

		NumberOfBlocks: (length + BLOCK_SIZE - 1) / BLOCK_SIZE,
		Blocks:         make([]*pieceBlock, 0),
		peerConn:       p.conn,
	}

	currentOffset := uint32(0)

	for range sp.NumberOfBlocks {
		blockLength := BLOCK_SIZE
		if currentOffset+BLOCK_SIZE > length {
			blockLength = length - currentOffset
		}

		// Create a block and start a goroutine to make the request
		block := newPieceBlock(uint32(currentOffset), uint32(blockLength))
		sp.Blocks = append(sp.Blocks, block)
		currentOffset += blockLength
	}

	// Register the piece with the peer
	p.assignedPiece = sp
	p.Log("assigned piece %d", sp.Index)

	sp.makeInitialDownloadRequests(p)
	err := p.getCompletePiece()
	if err != nil {
		return nil, fmt.Errorf("error getting complete piece: %w", err)
	}

	return sp, nil
}

func (sp *StoredPiece) makeInitialDownloadRequests(p *Peer) {
	for _, block := range sp.Blocks {
		err := block.makeRequest(p, sp.Index)
		if err != nil {
			fmt.Printf("error making request for block: %v\n", err)
		}
	}
	p.Log("%d blocks requested for piece %d", len(sp.Blocks), sp.Index)
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

func (sp *StoredPiece) VerifyHash() error {
	hash, err := utils.SHA1Hash(sp.GetData())
	if err != nil {
		return fmt.Errorf("error hashing piece data: %w", err)
	}

	if !bytes.Equal(hash, sp.Hash) {
		return fmt.Errorf("piece hash verification failed, expected %x, got %x", sp.Hash, hash)
	}

	return nil
}
