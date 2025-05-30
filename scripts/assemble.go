package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: assemble <chunk_dir> <output_file>")
		os.Exit(1)
	}

	chunkDir := os.Args[1]
	outputFile := os.Args[2]

	files, err := ioutil.ReadDir(chunkDir)
	if err != nil {
		log.Fatal(err)
	}

	// Group chunks by file ID
	chunkGroups := make(map[string][]string)
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		if !strings.HasSuffix(file.Name(), ".chunk") {
			continue
		}

		parts := strings.Split(file.Name(), "_")
		if len(parts) < 2 {
			continue
		}

		fileID := parts[0]
		chunkGroups[fileID] = append(chunkGroups[fileID], file.Name())
	}

	// Process each file
	for fileID, chunkFiles := range chunkGroups {
		sort.Slice(chunkFiles, func(i, j int) bool {
			idx1, _ := strconv.Atoi(strings.Split(strings.TrimSuffix(chunkFiles[i], ".chunk"), "_")[1])
			idx2, _ := strconv.Atoi(strings.Split(strings.TrimSuffix(chunkFiles[j], ".chunk"), "_")[1])
			return idx1 < idx2
		})

		out, err := os.Create(outputFile + "_" + fileID)
		if err != nil {
			log.Printf("Error creating output file: %v", err)
			continue
		}
		defer out.Close()

		hash := sha256.New()
		multiWriter := io.MultiWriter(out, hash)

		for _, chunkFile := range chunkFiles {
			chunkPath := filepath.Join(chunkDir, chunkFile)
			data, err := ioutil.ReadFile(chunkPath)
			if err != nil {
				log.Printf("Error reading chunk %s: %v", chunkFile, err)
				continue
			}

			if _, err := multiWriter.Write(data); err != nil {
				log.Printf("Error writing chunk %s: %v", chunkFile, err)
			}
		}

		finalHash := hex.EncodeToString(hash.Sum(nil))
		log.Printf("Assembled %s with checksum: %s", outputFile+"_"+fileID, finalHash)
	}
}
