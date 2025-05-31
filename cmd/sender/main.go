package main

import (
	"fmt"
	"log"
	"os"
	"sync/atomic"
	"time"

	"github.com/spf13/cobra"
	"github.com/kapishupadhyay22/SwiftTransfer/internal/chunker"
	"github.com/kapishupadhyay22/SwiftTransfer/internal/transfer"
)

var (
	nodes     []string
	chunkSize int
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "sender",
		Short: "Distributed File Transfer Sender",
	}

	var sendCmd = &cobra.Command{
		Use:   "send [file]",
		Short: "Send a file to distributed nodes",
		Args:  cobra.ExactArgs(1),
		Run:   sendFile,
	}

	sendCmd.Flags().StringSliceVarP(&nodes, "nodes", "n", []string{}, "Receiver nodes (comma separated)")
	sendCmd.Flags().IntVarP(&chunkSize, "chunk-size", "s", 4, "Chunk size in MB")

	rootCmd.AddCommand(sendCmd)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func sendFile(cmd *cobra.Command, args []string) {
	filePath := args[0]
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		log.Fatalf("File not found: %s", filePath)
	}

	if len(nodes) == 0 {
		log.Fatal("No receiver nodes specified")
	}

	log.Printf("Splitting %s into chunks...", filePath)
	chunks, err := chunker.SplitFile(filePath, chunkSize*1024*1024 - 500)  // considering the protocol overhead to be 500 bytes
	if err != nil {
		log.Fatalf("Error splitting file: %v", err)
	}

	totalChunks := len(chunks)
	log.Printf("Split into %d chunks, starting transfer...", totalChunks)

	manager := transfer.NewTransferManager(nodes, totalChunks)
	manager.StartWorkers(10) // Start 10 concurrent workers

	// Start progress monitoring
	go showProgress(manager.Progress, totalChunks)

	// Enqueue chunks
	for _, chunk := range chunks {
		manager.ChunkChan <- chunk
	}

	manager.Wait()
	log.Printf("\nTransfer completed successfully!")
}

func showProgress(progress *int32, total int) {
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	lastValue := int32(0)
	for range ticker.C {
		current := atomic.LoadInt32(progress)
		if current == int32(total) {
			fmt.Printf("\rProgress: 100%% (%d/%d chunks)", total, total)
			return
		}

		if current != lastValue {
			percent := float32(current) / float32(total) * 100
			fmt.Printf("\rProgress: %.1f%% (%d/%d chunks)", percent, current, total)
			lastValue = current
		}
	}
}
