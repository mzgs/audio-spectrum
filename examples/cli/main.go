package main

import (
	"flag"
	"fmt"
	"os"
	
	audiospectrum "github.com/mzgs/audio-spectrum"
)

func main() {
	// Parse command line arguments
	var (
		outputFile   = flag.String("o", "spectrum_video.mp4", "Output video file")
		fps          = flag.Int("f", 30, "Frames per second")
		duration     = flag.Float64("d", 0, "Duration in seconds (0 for full audio)")
		bars         = flag.Int("b", 32, "Number of frequency bars")
		colorScheme  = flag.String("c", "rainbow", "Color scheme (rainbow, fire, ocean, purple, neon, monochrome, sunset, forest, ice, lava, retro, cosmic, pastel, matrix, white)")
		vizType      = flag.String("t", "bars", "Visualization type (bars, circular, wave, radial, line, dots, mirror, spiral)")
		bgColor      = flag.String("bg", "green", "Background color (green, blue, magenta, black, white, gray)")
		width        = flag.Int("w", 1280, "Video width")
		height       = flag.Int("h", 720, "Video height")
		processType  = flag.String("method", "fast", "Processing method (fast, parallel)")
	)
	
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options] input.mp3\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\nUltra-fast audio spectrum video generator\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nExamples:\n")
		fmt.Fprintf(os.Stderr, "  %s input.mp3\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -o output.mp4 -f 60 -b 64 -c fire input.mp3\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -t circular -c ocean -d 10 input.mp3\n", os.Args[0])
	}
	
	flag.Parse()
	
	// Check for input file
	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(1)
	}
	
	inputFile := flag.Arg(0)
	
	// Create configuration
	config := &audiospectrum.Config{
		InputFile:    inputFile,
		OutputFile:   *outputFile,
		FPS:          *fps,
		Duration:     *duration,
		BarCount:     *bars,
		ColorScheme:  audiospectrum.ColorScheme(*colorScheme),
		VisType:      audiospectrum.VisType(*vizType),
		BGColor:      audiospectrum.BGColor(*bgColor),
		Width:        *width,
		Height:       *height,
		ProcessType:  audiospectrum.ProcessType(*processType),
	}
	
	// Generate video
	if err := audiospectrum.Generate(config); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}