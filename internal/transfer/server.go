package transfer

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"

	pb "github.com/kapishupadhyay22/SwiftTransfer/proto"
)

type Server struct {
	pb.UnimplementedFileTransferServer
	mu             sync.Mutex
	receivedChunks map[string]map[int]bool // fileID -> chunk index
	chunkDirectory string
}

func NewServer(chunkDir string) *Server {
	os.MkdirAll(chunkDir, 0755)
	return &Server{
		receivedChunks: make(map[string]map[int]bool),
		chunkDirectory: chunkDir,
	}
}

func (s *Server) SendChunk(ctx context.Context, chunk *pb.Chunk) (*pb.Ack, error) {
	// Verify checksum
	hash := sha256.Sum256(chunk.Content)
	if hex.EncodeToString(hash[:]) != chunk.Checksum {
		return &pb.Ack{Success: false, Message: "Checksum mismatch"}, nil
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	// Track received chunks
	if s.receivedChunks[chunk.FileId] == nil {
		s.receivedChunks[chunk.FileId] = make(map[int]bool)
	}
	s.receivedChunks[chunk.FileId][int(chunk.Index)] = true

	// Save chunk to disk
	chunkPath := filepath.Join(s.chunkDirectory, fmt.Sprintf("%s_%d.chunk", chunk.FileId, chunk.Index))
	if err := ioutil.WriteFile(chunkPath, chunk.Content, 0644); err != nil {
		return &pb.Ack{Success: false, Message: err.Error()}, nil
	}

	return &pb.Ack{Success: true, Message: "Chunk received"}, nil
}

func (s *Server) HealthCheck(ctx context.Context, req *pb.HealthRequest) (*pb.HealthResponse, error) {
	return &pb.HealthResponse{Ready: true}, nil
}
