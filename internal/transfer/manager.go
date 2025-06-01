package transfer

import (
	"context"
	"log"
	"path/filepath"
	"sync"
	"sync/atomic"
	"time"

	"github.com/kapishupadhyay22/SwiftTransfer/internal/chunker"
	pb "github.com/kapishupadhyay22/SwiftTransfer/proto"
	"google.golang.org/grpc"
)

type TransferManager struct {
	Nodes     []string	// nodes contain the different receiver nodes
	ChunkChan chan chunker.FileChunk
	Wg        sync.WaitGroup
	Progress  *int32
}

func NewTransferManager(nodes []string, bufferSize int) *TransferManager {
	var progress int32
	return &TransferManager{
		Nodes:     nodes,
		ChunkChan: make(chan chunker.FileChunk, bufferSize),
		Progress:  &progress,
	}
}

func (tm *TransferManager) StartWorkers(workerCount int) {
	for i := 0; i < workerCount; i++ {
		tm.Wg.Add(1)
		go tm.worker()
	}
}

func (tm *TransferManager) worker() {
	defer tm.Wg.Done()

	for chunk := range tm.ChunkChan {
		success := false
		attempts := 0

		// Retry with exponential backoff
		for !success && attempts < 5 {
			node := tm.Nodes[attempts%len(tm.Nodes)]
			conn, err := grpc.Dial(node, grpc.WithInsecure())
			if err != nil {
				log.Printf("Connection failed to %s: %v", node, err)
				time.Sleep(time.Duration(attempts*attempts) * time.Second)
				attempts++
				continue
			}

			client := pb.NewFileTransferClient(conn)
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

			_, err = client.SendChunk(ctx, &pb.Chunk{
				Content:     chunk.Data,
				Checksum:    chunk.Checksum,
				Index:       int32(chunk.Index),
				Filename:    filepath.Base(chunk.FilePath),
				TotalChunks: int32(chunk.Total),
				FileId:      chunk.FileID,
			})
			cancel()

			if err != nil {
				log.Printf("Failed to send chunk %d: %v", chunk.Index, err)
				attempts++
				time.Sleep(time.Duration(attempts) * time.Second)
			} else {
				success = true
				atomic.AddInt32(tm.Progress, 1)
			}
			conn.Close()
		}

		if !success {
			log.Printf("Permanently failed to send chunk %d", chunk.Index)
		}
	}
}

func (tm *TransferManager) Wait() {
	close(tm.ChunkChan)
	tm.Wg.Wait()
}
