package audiospectrum

import (
	"fmt"
	"os"
	"time"
)

// Config holds all configuration options for the audio spectrum visualizer
type Config struct {
	InputFile    string
	OutputFile   string
	FPS          int
	Duration     float64
	BarCount     int
	ColorScheme  ColorScheme
	VisType      VisType
	BGColor      BGColor
	Width        int
	Height       int
	ProcessType  ProcessType
}

// DefaultConfig returns a Config with sensible defaults
func DefaultConfig() *Config {
	return &Config{
		OutputFile:   "spectrum_video.mp4",
		FPS:          30,
		Duration:     0, // 0 means full audio duration
		BarCount:     32,
		ColorScheme:  ColorSchemeRainbow,
		VisType:      VisTypeBars,
		BGColor:      BGColorGreen,
		Width:        1280,
		Height:       720,
		ProcessType:  ProcessTypeFast,
	}
}

// Generate creates an audio spectrum video from the given audio file
func Generate(config *Config) error {
	// Validate input
	if config.InputFile == "" {
		return fmt.Errorf("input file is required")
	}
	
	if _, err := os.Stat(config.InputFile); os.IsNotExist(err) {
		return fmt.Errorf("input file not found: %s", config.InputFile)
	}
	
	// Validate configuration
	if err := validateConfig(config); err != nil {
		return fmt.Errorf("invalid configuration: %w", err)
	}
	
	// Create visualizer config
	vizConfig := &VisualizerConfig{
		InputFile:    config.InputFile,
		OutputFile:   config.OutputFile,
		FPS:          config.FPS,
		Duration:     config.Duration,
		BarCount:     config.BarCount,
		ColorScheme:  string(config.ColorScheme),
		VizType:      string(config.VisType),
		BgColor:      string(config.BGColor),
		Width:        config.Width,
		Height:       config.Height,
		ProcessType:  string(config.ProcessType),
	}
	
	// Create and run visualizer
	visualizer := NewVisualizer(vizConfig)
	
	fmt.Printf("Processing audio file: %s\n", config.InputFile)
	startTime := time.Now()
	
	if err := visualizer.CreateVideo(); err != nil {
		return fmt.Errorf("failed to generate video: %w", err)
	}
	
	// Get file size
	fileInfo, err := os.Stat(config.OutputFile)
	if err != nil {
		return fmt.Errorf("failed to get output file info: %w", err)
	}
	
	duration := time.Since(startTime)
	fmt.Printf("\nVideo created successfully: %s\n", config.OutputFile)
	fmt.Printf("Total processing time: %.1f seconds\n", duration.Seconds())
	fmt.Printf("Output file size: %.1f MB\n", float64(fileInfo.Size())/(1024*1024))
	
	return nil
}

// GenerateWithDefaults creates a video with default settings, only requiring input/output files
func GenerateWithDefaults(inputFile, outputFile string) error {
	config := DefaultConfig()
	config.InputFile = inputFile
	config.OutputFile = outputFile
	return Generate(config)
}

func validateConfig(config *Config) error {
	// Validate FPS
	if config.FPS < 1 || config.FPS > 120 {
		return fmt.Errorf("FPS must be between 1 and 120")
	}
	
	// Validate duration
	if config.Duration < 0 {
		return fmt.Errorf("duration cannot be negative")
	}
	
	// Validate bar count
	if config.BarCount < 8 || config.BarCount > 256 {
		return fmt.Errorf("bar count must be between 8 and 256")
	}
	
	// Validate dimensions
	if config.Width < 320 || config.Height < 240 {
		return fmt.Errorf("minimum resolution is 320x240")
	}
	if config.Width > 7680 || config.Height > 4320 {
		return fmt.Errorf("maximum resolution is 7680x4320")
	}
	
	// Validate color scheme
	if !config.ColorScheme.IsValid() {
		return fmt.Errorf("invalid color scheme: %s", config.ColorScheme)
	}
	
	// Validate visualization type
	if !config.VisType.IsValid() {
		return fmt.Errorf("invalid visualization type: %s", config.VisType)
	}
	
	// Validate background color
	if !config.BGColor.IsValid() {
		return fmt.Errorf("invalid background color: %s", config.BGColor)
	}
	
	// Validate process type
	if !config.ProcessType.IsValid() {
		return fmt.Errorf("invalid process type: %s", config.ProcessType)
	}
	
	return nil
}

// GetSupportedFormats returns the supported audio formats
func GetSupportedFormats() []string {
	return []string{
		"mp3", "wav", "flac", "aac", "ogg", "m4a", "wma", "opus",
	}
}

// GetColorSchemes returns all available color schemes
func GetColorSchemes() []ColorScheme {
	return []ColorScheme{
		ColorSchemeRainbow, ColorSchemeFire, ColorSchemeOcean, ColorSchemePurple,
		ColorSchemeNeon, ColorSchemeMonochrome, ColorSchemeSunset, ColorSchemeForest,
		ColorSchemeWhite,
	}
}

// GetVisualizationTypes returns all available visualization types
func GetVisualizationTypes() []VisType {
	return []VisType{
		VisTypeBars, VisTypeCircular, VisTypeWave, VisTypeRadial,
		VisTypeLine, VisTypeDots, VisTypeMirror, VisTypeSpiral,
	}
}

// GetBackgroundColors returns all available background colors
func GetBackgroundColors() []BGColor {
	return []BGColor{
		BGColorGreen, BGColorBlue, BGColorMagenta, 
		BGColorBlack, BGColorWhite, BGColorGray,
	}
}

// GetProcessTypes returns all available process types
func GetProcessTypes() []ProcessType {
	return []ProcessType{
		ProcessTypeFast, ProcessTypeParallel,
	}
}