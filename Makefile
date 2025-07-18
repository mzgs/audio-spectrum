# Makefile for Audio Spectrum Visualizer

BINARY_NAME=audio-spectrum
GO_FILES=$(shell find . -name '*.go' -type f)

# Build the application
build:
	go build -o $(BINARY_NAME) .

# Run the application with a test file
run: build
	./$(BINARY_NAME) input.mp3

# Install dependencies
deps:
	go mod download
	go mod tidy

# Clean build artifacts
clean:
	rm -f $(BINARY_NAME)
	rm -f spectrum_video.mp4
	rm -rf temp_frames*

# Format Go code
fmt:
	go fmt ./...

# Run tests
test:
	go test -v ./...

# Build for multiple platforms
build-all:
	GOOS=darwin GOARCH=amd64 go build -o $(BINARY_NAME)-darwin-amd64 .
	GOOS=darwin GOARCH=arm64 go build -o $(BINARY_NAME)-darwin-arm64 .
	GOOS=linux GOARCH=amd64 go build -o $(BINARY_NAME)-linux-amd64 .
	GOOS=windows GOARCH=amd64 go build -o $(BINARY_NAME)-windows-amd64.exe .

# Install the binary to GOPATH/bin
install: build
	go install

.PHONY: build run deps clean fmt test build-all install