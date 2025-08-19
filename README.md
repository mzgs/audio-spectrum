# Audio Spectrum Visualizer

A high-performance Go library for generating audio spectrum visualization videos. Create stunning visual representations of audio files with various visualization styles and color schemes.

## Features

- 🚀 **High Performance** - Optimized FFT processing and parallel frame generation
- 🎨 **Multiple Visualizations** - 8 different visualization types (bars, circular, wave, radial, etc.)
- 🌈 **Rich Color Schemes** - 15 built-in color schemes
- 🎬 **Flexible Output** - Customizable resolution, frame rate, and duration
- 📦 **Easy Integration** - Simple API for use in your Go projects
- 🛠️ **FFmpeg Powered** - Reliable audio/video processing

## Installation

```bash
go get github.com/mzgs/audio-spectrum
```

### Requirements

- Go 1.21 or later
- FFmpeg installed on your system:
  - macOS: `brew install ffmpeg`
  - Ubuntu/Debian: `sudo apt-get install ffmpeg`
  - Windows: Download from [ffmpeg.org](https://ffmpeg.org/download.html)

## Quick Start

### Basic Usage

```go
package main

import (
    "log"
    audiospectrum "github.com/mzgs/audio-spectrum"
)

func main() {
    // Generate with defaults
    err := audiospectrum.GenerateWithDefaults("input.mp3", "output.mp4")
    if err != nil {
        log.Fatal(err)
    }
}
```

### Custom Configuration

```go
package main

import (
    "log"
    audiospectrum "github.com/mzgs/audio-spectrum"
)

func main() {
    config := &audiospectrum.Config{
        InputFile:    "song.mp3",
        OutputFile:   "spectrum.mp4",
        FPS:          60,
        Duration:     30, // First 30 seconds
        BarCount:     64,
        ColorScheme:  audiospectrum.ColorSchemeFire,
        VisType:      audiospectrum.VisTypeCircular,
        BGColor:      audiospectrum.BGColorBlack,
        Width:        1920,
        Height:       1080,
        ProcessType:  audiospectrum.ProcessTypeParallel, // Use all CPU cores
    }
    
    err := audiospectrum.Generate(config)
    if err != nil {
        log.Fatal(err)
    }
}
```

## Visualization Types

- **bars** - Traditional vertical bars (default)
- **circular** - Bars radiating outward from center
- **wave** - Waveform visualization
- **radial** - Radial burst pattern
- **line** - Connected line graph spectrum
- **dots** - Particle/dots effect
- **mirror** - Mirrored bars from center
- **spiral** - Spiral pattern

## Color Schemes

- **rainbow** - Classic green to red gradient
- **fire** - Realistic fire colors (black → red → orange → yellow)
- **ocean** - Dark blue to cyan
- **purple** - Purple to pink gradient
- **neon** - Electric neon colors (cyan → magenta → green)
- **monochrome** - Grayscale gradient
- **sunset** - Deep purple → pink → orange → golden
- **forest** - Dark green to yellow-green
- **ice** - Ice blue gradient (white → light blue → cyan → deep blue)
- **lava** - Volcanic colors (black → dark red → orange → yellow → white)
- **retro** - 80s retro (purple → magenta → cyan → yellow)
- **cosmic** - Space theme (deep purple → blue → teal → pink)
- **pastel** - Soft pastel colors with low saturation
- **matrix** - Matrix green theme (dark to bright green)
- **white** - Pure white bars

## CLI Tool

The library includes a command-line tool for easy video generation:

```bash
# Build the CLI tool
cd examples/cli
go build -o audio-spectrum

# Generate video with defaults
./audio-spectrum input.mp3

# Custom options
./audio-spectrum -o output.mp4 -f 60 -b 64 -c fire -t circular input.mp3

# Try new color schemes
./audio-spectrum -c lava -bg black input.mp3
./audio-spectrum -c retro -bg black input.mp3
./audio-spectrum -c pastel -bg white input.mp3
```

## API Reference

### Functions

#### `GenerateWithDefaults(inputFile, outputFile string) error`
Generate a video with default settings.

#### `Generate(config *Config) error`
Generate a video with custom configuration.

#### `DefaultConfig() *Config`
Returns a configuration with default values.

### Configuration Options

```go
type Config struct {
    InputFile    string       // Input audio file (required)
    OutputFile   string       // Output video file (default: "spectrum_video.mp4")
    FPS          int          // Frames per second (default: 30, range: 1-120)
    Duration     float64      // Duration in seconds (default: 0 = full audio)
    BarCount     int          // Number of frequency bars (default: 32, range: 8-256)
    ColorScheme  ColorScheme  // Color scheme (default: ColorSchemeRainbow)
    VisType      VisType      // Visualization type (default: VisTypeBars)
    BGColor      BGColor      // Background color (default: BGColorGreen)
    Width        int          // Video width (default: 1280)
    Height       int          // Video height (default: 720)
    ProcessType  ProcessType  // Processing method (default: ProcessTypeFast)
}
```

### Constants

The library provides type-safe constants for all options:

```go
// Color Schemes
ColorSchemeRainbow, ColorSchemeFire, ColorSchemeOcean, ColorSchemePurple,
ColorSchemeNeon, ColorSchemeMonochrome, ColorSchemeSunset, ColorSchemeForest,
ColorSchemeIce, ColorSchemeLava, ColorSchemeRetro, ColorSchemeCosmic,
ColorSchemePastel, ColorSchemeMatrix, ColorSchemeWhite

// Visualization Types
VisTypeBars, VisTypeCircular, VisTypeWave, VisTypeRadial,
VisTypeLine, VisTypeDots, VisTypeMirror, VisTypeSpiral

// Background Colors
BGColorGreen, BGColorBlue, BGColorMagenta,
BGColorBlack, BGColorWhite, BGColorGray

// Process Types
ProcessTypeFast, ProcessTypeParallel
```

### Utility Functions

```go
GetSupportedFormats() []string         // Returns supported audio formats
GetColorSchemes() []ColorScheme        // Returns available color schemes
GetVisualizationTypes() []VisType      // Returns available visualization types
GetBackgroundColors() []BGColor        // Returns available background colors
GetProcessTypes() []ProcessType        // Returns available process types
```

## Examples

See the [examples](examples/) directory for more usage examples:
- [Simple usage](examples/simple/main.go)
- [CLI tool](examples/cli/main.go)

## Performance Tips

1. Use `ProcessType: "parallel"` for faster processing on multi-core systems
2. Lower `BarCount` for faster processing
3. Reduce resolution for quicker renders
4. Use `Duration` to limit processing time for testing

## License

MIT License - see LICENSE file for details

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## Acknowledgments

- FFmpeg for audio/video processing
- [go-dsp](https://github.com/mjibson/go-dsp) for FFT computation
- [gg](https://github.com/fogleman/gg) for 2D graphics rendering