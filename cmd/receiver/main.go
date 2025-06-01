package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/kapishupadhyay22/SwiftTransfer/internal/transfer"
	pb "github.com/kapishupadhyay22/SwiftTransfer/proto"
	"google.golang.org/grpc"
)

func main() {
	port := "12345"
	if p := os.Getenv("PORT"); p != "" {
		port = p
	}

	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	server := grpc.NewServer()
	transferServer := transfer.NewServer("./chunks")
	pb.RegisterFileTransferServer(server, transferServer)

	log.Printf("Receiver listening on port %s", port)

	// Graceful shutdown
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, syscall.SIGINT, syscall.SIGTERM)
		<-sigint
		log.Println("Shutting down server...")
		server.GracefulStop()
	}()

	if err := server.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
