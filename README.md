# SwiftTransfer - Distributed File Transfer System

**High-performance, resilient file transfers using Go and gRPC**

[![Go Report Card](https://goreportcard.com/badge/github.com/kapishupadhyay22/SwiftTransfer)](https://goreportcard.com/report/github.com/kapishupadhyay22/SwiftTransfer)

SwiftTransfer is a distributed file transfer system that splits files into chunks and transfers them concurrently to multiple nodes. Built with Go and gRPC, it features automatic retries, checksum verification, and progress tracking for reliable large file transfers.

---

## Key Features ✨

* **Parallel Transfers** - Concurrent chunk delivery to multiple nodes
* **Resilient Delivery** - Automatic retries with exponential backoff
* **Data Integrity** - SHA-256 checksum verification
* **Progress Tracking** - Real-time transfer monitoring
* **Distributed Architecture** - Scale across multiple machines
* **Simple CLI** - Easy-to-use command interface

---

Installation 🚀
Prerequisites
Go 1.24+
Protocol Buffer compiler (protoc)
Build from Source
Bash

# Clone repository
```
git clone [https://github.com/kapishupadhyay22/SwiftTransfer.git](https://github.com/kapishupadhyay22/SwiftTransfer.git)
cd SwiftTransfer
```
# Install dependencies
```
go get -u google.golang.org/grpc
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```
# Generate protobuf code
```
protoc --go_out=. --go_opt=paths=source_relative \
       --go-grpc_out=. --go-grpc_opt=paths=source_relative \
       proto/transfer.proto
```
# Build binaries
```
go build -o receiver ./cmd/receiver
go build -o sender ./cmd/sender
```
# Usage 📖
## 1. Start Receiver Service

### Create storage directory
```
mkdir -p chunks
```
### Start receiver (default port 50051)
```
./receiver
```
### Custom port and storage location
```
CHUNK_DIR=my_chunks PORT=6000 ./receiver
```
## 2. Send Files

### Send to single receiver
```
./sender send -n localhost:12345 -s 4 large_file.iso
(replace with the name of file you want to transfer)
```
### Send to multiple receivers
```
./sender send -n 192.168.1.100:50051,192.168.1.101:50051 -s 10 video.mp4
(replace with the name of file you want to transfer)
```
Options:

-n: Comma-separated receiver addresses

-s: Chunk size in MB (default: 4)

-z: Enable gzip compression (optional)

-c: Concurrent workers (default: 10)

### 3. Reassemble Files
```
go run scripts/assemble.go chunks/ reconstructed.iso
```
### 4. Verify Integrity
```
sha256sum original_file.iso reconstructed.iso
(replace with the name of file you want to transfer)
```
# Advanced Usage 🧠
Multiple Receiver Nodes


### Terminal 1 (Node 1)
```
CHUNK_DIR=chunks1 PORT=50051 ./receiver
```
### Terminal 2 (Node 2)
```
CHUNK_DIR=chunks2 PORT=50052 ./receiver
```
### Send file to both nodes
```
./sender send -n localhost:50051,localhost:50052 -s 8 data.bin
(replace with the name of file you want to transfer)
```
## Combine chunks from both nodes
```
mkdir all_chunks
cp chunks1/* all_chunks/
cp chunks2/* all_chunks/
```

```
go run scripts/assemble.go all_chunks/ restored.bin
```

Contributing 🤝
We welcome contributions! Please see our Contribution Guidelines for details. (LMAO, there's no guide, do whatever you want with the project and raise your pull requests)

License 📄
NO License!! I don't need it HAHA !! LOL , use it however you want

SwiftTransfer - Transfer large files at the speed of light! 