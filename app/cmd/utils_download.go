package cmd

import (
	"bytes"
	"fmt"
	"sync"
	"time"

	"github.com/EshaanAgg/toy-bittorrent/app/types"
)

type pieceToDownload struct {
	index  uint32
	length uint32
	hash   []byte
}

type pieceResult struct {
	index uint32
	data  []byte
	err   error
}

// downloadPieces downloads all pieces from peers concurrently. It uses a
// worker pool to limit the number of concurrent downloads and handles
// failed downloads by requeuing them for retry.
func downloadPieces(peers []*types.Peer, pieces []*pieceToDownload) []byte {
	var wg sync.WaitGroup
	numWorkers := len(peers)
	pieceQueue := make(chan *pieceToDownload, len(pieces))
	results := make(chan pieceResult, len(pieces))

	// Feed the queue with all pieces
	for _, piece := range pieces {
		pieceQueue <- piece
	}
	close(pieceQueue)

	// Worker function
	worker := func(peer *types.Peer) {
		defer wg.Done()
		for piece := range pieceQueue {
			sp, err := peer.DownloadPiece(piece.index, piece.length, piece.hash)
			if err != nil {
				fmt.Printf("error with piece %d: %v, retrying\n", piece.index, err)

				// Requeue the piece
				go func(p *pieceToDownload) {
					time.Sleep(500 * time.Millisecond) // Short delay to avoid busy-loop
					pieceQueue <- p
				}(piece)

				continue
			}
			results <- pieceResult{index: piece.index, data: sp.GetData()}
		}
	}

	// Start workers
	wg.Add(numWorkers)
	for _, peer := range peers {
		go worker(peer)
	}

	wg.Wait()
	close(results)

	// Collect results
	filePieces := make([][]byte, len(pieces))
	for res := range results {
		filePieces[res.index] = res.data
	}

	finalData := bytes.Join(filePieces, []byte{})
	return finalData
}
