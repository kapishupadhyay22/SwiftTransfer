package chunker

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
	"path/filepath"
	"time"
)

type FileChunk struct {
	Index    int
	Data     []byte
	Checksum string
	FilePath string
	Total    int
	FileID   string
}

func SplitFile(filePath string, chunkSize int) ([]FileChunk, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	fileInfo, _ := file.Stat()
	totalChunks := int((fileInfo.Size() + int64(chunkSize) - 1) / int64(chunkSize))
	fileID := generateFileID(filePath)

	var chunks []FileChunk
	buffer := make([]byte, chunkSize)
	idx := 0

	for {
		bytesRead, err := file.Read(buffer)
		if err != nil {
			if err == io.EOF {
				break	// end of file, nothing more to read
			}
			return nil, err
		}

		chunkData := make([]byte, bytesRead)
		copy(chunkData, buffer[:bytesRead])

		hash := sha256.Sum256(chunkData)
		checksum := hex.EncodeToString(hash[:])

		chunks = append(chunks, FileChunk{
			Index:    idx,
			Data:     chunkData,
			Checksum: checksum,
			FilePath: filePath,
			Total:    totalChunks,
			FileID:   fileID,
		})

		idx++
	}

	return chunks, nil
}

func generateFileID(filePath string) string {
	absPath, _ := filepath.Abs(filePath)
	hash := sha256.Sum256([]byte(absPath + "-" + time.Now().Format(time.RFC3339Nano)))
	return hex.EncodeToString(hash[:])
}
