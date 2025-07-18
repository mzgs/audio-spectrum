package main

import (
	"log"
	
	audiospectrum "github.com/mzgs/audio-spectrum"
)

func main() {
	// Example 1: Generate with all defaults
	err := audiospectrum.GenerateWithDefaults("input.mp3", "output.mp4")
	if err != nil {
		log.Fatal(err)
	}
	
	// Example 2: Custom configuration
	config := audiospectrum.DefaultConfig()
	config.InputFile = "song.mp3"
	config.OutputFile = "spectrum.mp4"
	config.ColorScheme = audiospectrum.ColorSchemeFire
	config.VisType = audiospectrum.VisTypeCircular
	config.FPS = 60
	config.BarCount = 64
	
	err = audiospectrum.Generate(config)
	if err != nil {
		log.Fatal(err)
	}
	
	// Example 3: High-quality render
	hqConfig := &audiospectrum.Config{
		InputFile:    "music.mp3",
		OutputFile:   "hq_spectrum.mp4",
		FPS:          60,
		Duration:     30, // First 30 seconds only
		BarCount:     128,
		ColorScheme:  audiospectrum.ColorSchemeOcean,
		VisType:      audiospectrum.VisTypeWave,
		BGColor:      audiospectrum.BGColorBlack,
		Width:        1920,
		Height:       1080,
		ProcessType:  audiospectrum.ProcessTypeParallel,
	}
	
	err = audiospectrum.Generate(hqConfig)
	if err != nil {
		log.Fatal(err)
	}
}