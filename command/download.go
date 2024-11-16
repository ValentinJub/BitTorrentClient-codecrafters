package command

import (
	"fmt"
	"math/rand/v2"
	"sync"
	"time"

	d "github.com/codecrafters-io/bittorrent-starter-go/decoder"
	"github.com/codecrafters-io/bittorrent-starter-go/utils"
)

// Download downloads a torrent file from a list of peers concurrently
func Download(t *d.TorrentFile, peers []string, outputFile string) error {
	// Time the execution of the function
	start := time.Now()
	fmt.Printf("List of peers: %v\n", peers)
	data := make(map[int][]byte) // Map to store the piece data, key is the piece index
	var mu sync.Mutex            // Mutex to protect the data map
	var wg sync.WaitGroup        // WaitGroup to wait for all goroutines to complete
	results := make(chan error, len(t.PieceHashes))

	// Download each piece concurrently
	for pIndex, pieceHash := range t.PieceHashes {
		wg.Add(1)
		go func(pIndex int, pieceHash string) {
			defer wg.Done()
			// Shuffle the peers slice to select a random peer
			shuffledPeers := make([]string, len(peers))
			copy(shuffledPeers, peers)
			rand.Shuffle(len(shuffledPeers), func(i, j int) {
				shuffledPeers[i], shuffledPeers[j] = shuffledPeers[j], shuffledPeers[i]
			})
			for _, peer := range shuffledPeers {
				pieceData, err := DownloadPiece(peer, t.Length, t.PieceLength, t.InfoHash, pieceHash, pIndex, pIndex == len(t.PieceHashes)-1)
				if err != nil { // If there is an error while downloading the piece, try the next peer
					fmt.Printf("error while downloading piece %d with peer %s, moving on\n", pIndex, peer)
					continue
				}
				mu.Lock()
				data[pIndex] = pieceData // Store the piece data that was downloaded successfully
				mu.Unlock()
				fmt.Printf("successfully downloaded piece %d with peer %s\n", pIndex, peer)
				results <- nil
				return
			}
			results <- fmt.Errorf("error while downloading piece %d, no peers succeeded", pIndex)
		}(pIndex, pieceHash)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	for err := range results {
		if err != nil {
			return err
		}
	}

	dataReconstructed := make([]byte, 0)
	for i := 0; i < len(t.PieceHashes); i++ {
		if d, ok := data[i]; !ok {
			return fmt.Errorf("error while downloading piece %d, piece data not found", i)
		} else {
			dataReconstructed = append(dataReconstructed, d...)
		}
	}

	err := utils.WriteFile(outputFile, dataReconstructed)
	if err != nil {
		return fmt.Errorf("error while writing to file: %v", err)
	}
	// Print the time taken to download the torrent
	// Convert the length from byte to megabyte
	torrentsize := float64(t.Length) / float64(1_000_000)
	fmt.Printf("Downloaded torrent of size %5fmb in %v\n", torrentsize, time.Since(start))
	return nil
}
