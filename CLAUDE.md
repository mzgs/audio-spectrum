# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build and Development Commands

### Library Development
```bash
go mod download        # Download dependencies
go mod tidy           # Clean up go.mod and go.sum
go test ./...         # Run tests
```

### Build CLI Example
```bash
cd examples/cli
go build -o audio-spectrum
./audio-spectrum input.mp3  # Run with specific audio file
```

### Dependencies
```bash
make deps              # Download and tidy Go dependencies
go mod download        # Alternative: download dependencies
go mod tidy           # Clean up go.mod and go.sum
```

### Code Quality
```bash
make fmt              # Format all Go code files
go fmt ./...          # Alternative format command
make test             # Run tests (if any exist)
```

### Clean
```bash
make clean            # Remove binary, output videos, and temp frames
```

### Cross-platform Build
```bash
make build-all        # Build for Darwin (amd64/arm64), Linux, and Windows
```

## Architecture Overview

This is a Go-based audio spectrum visualizer that generates videos from audio files. The codebase follows a modular design:

### Core Components

1. **spectrum.go**: Public API for the library. Provides `Generate()` and `GenerateWithDefaults()` functions along with configuration validation.

2. **visualizer.go**: Contains the core `Visualizer` struct and `VisualizerConfig` that manage:
   - Audio extraction via FFmpeg
   - FFT computation using go-dsp library
   - Frame generation orchestration
   - Video encoding pipeline

3. **draw.go**: Implements all visualization rendering logic using the fogleman/gg graphics library. Contains different visualization types (bars, circular, wave, radial, etc.) and color schemes.

### Key Processing Flow

1. Extract audio data from input file using FFmpeg
2. Compute FFT spectrum data for all time windows
3. Generate frames either sequentially (fast) or in parallel
4. Encode frames to video using FFmpeg with audio track

### External Dependencies

- **FFmpeg**: Required system dependency for audio/video processing
- **github.com/fogleman/gg**: 2D graphics rendering
- **github.com/mjibson/go-dsp**: Digital signal processing (FFT calculations)

### Important Implementation Details

- The visualizer pre-computes all FFT data before rendering for optimal performance
- Supports two processing methods: "fast" (sequential) and "parallel" (concurrent frame generation)
- Uses direct FFmpeg pipes for efficient frame streaming without intermediate files
- All visualization types share common rendering infrastructure in draw.go