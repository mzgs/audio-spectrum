package audiospectrum

import (
	"fmt"
	"image/color"
	"math"
	"math/cmplx"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sync"
	"time"

	"github.com/fogleman/gg"
	"github.com/mjibson/go-dsp/fft"
)

// VisualizerConfig holds all configuration for the visualizer
type VisualizerConfig struct {
	InputFile    string
	OutputFile   string
	FPS          int
	Duration     float64
	BarCount     int
	ColorScheme  string
	VizType      string
	BgColor      string
	Width        int
	Height       int
	ProcessType  string
}

// Visualizer handles the audio spectrum visualization
type Visualizer struct {
	config       *VisualizerConfig
	audioData    []float64
	sampleRate   int
	duration     float64
	totalFrames  int
	spectrumData [][]float64
	barPositions []int
	centerX      int
	centerY      int
	barWidth     int
	windowSize   int
}

// NewVisualizer creates a new visualizer instance
func NewVisualizer(config *VisualizerConfig) *Visualizer {
	v := &Visualizer{
		config:   config,
		centerX:  config.Width / 2,
		centerY:  config.Height / 2,
		barWidth: config.Width / config.BarCount,
	}
	
	// Pre-calculate bar positions
	v.barPositions = make([]int, config.BarCount)
	for i := 0; i < config.BarCount; i++ {
		v.barPositions[i] = i * v.barWidth
	}
	
	return v
}

// CreateVideo creates the spectrum visualization video
func (v *Visualizer) CreateVideo() error {
	// Load audio
	if err := v.loadAudio(); err != nil {
		return fmt.Errorf("loading audio: %w", err)
	}
	
	// Pre-compute spectrum data
	fmt.Println("Pre-computing spectrum data...")
	if err := v.precomputeSpectrum(); err != nil {
		return fmt.Errorf("computing spectrum: %w", err)
	}
	
	// Generate frames
	fmt.Printf("Generating %d frames...\n", v.totalFrames)
	if v.config.ProcessType == "parallel" {
		return v.createVideoParallel()
	}
	return v.createVideoSequential()
}

// loadAudio loads the audio file and prepares it for processing
func (v *Visualizer) loadAudio() error {
	// For now, we'll use ffmpeg to extract audio data
	// In a production version, we'd use a proper audio library
	
	// First, get audio info using ffprobe
	cmd := exec.Command("ffprobe",
		"-v", "error",
		"-show_entries", "format=duration",
		"-of", "default=noprint_wrappers=1:nokey=1",
		v.config.InputFile,
	)
	
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("getting audio duration: %w", err)
	}
	
	// Parse duration
	var fileDuration float64
	fmt.Sscanf(string(output), "%f", &fileDuration)
	
	// Set duration
	if v.config.Duration > 0 && v.config.Duration < fileDuration {
		v.duration = v.config.Duration
	} else {
		v.duration = fileDuration
	}
	
	v.totalFrames = int(v.duration * float64(v.config.FPS))
	v.sampleRate = 22050 // Standard sample rate for analysis
	
	fmt.Printf("Audio duration: %.1f seconds, %d frames\n", v.duration, v.totalFrames)
	
	// Extract raw audio data using ffmpeg
	// This is a simplified version - in production, use proper audio libraries
	return v.extractAudioData()
}

// extractAudioData extracts raw PCM data from the audio file
func (v *Visualizer) extractAudioData() error {
	// Create temp file for raw audio
	tempFile := filepath.Join(os.TempDir(), "audio_temp.raw")
	defer os.Remove(tempFile)
	
	// Convert to raw PCM using ffmpeg
	cmd := exec.Command("ffmpeg",
		"-i", v.config.InputFile,
		"-f", "f32le",
		"-acodec", "pcm_f32le",
		"-ac", "1",
		"-ar", fmt.Sprintf("%d", v.sampleRate),
		"-t", fmt.Sprintf("%.2f", v.duration),
		"-y", tempFile,
	)
	
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("converting audio: %w", err)
	}
	
	// Read the raw audio data
	data, err := os.ReadFile(tempFile)
	if err != nil {
		return fmt.Errorf("reading audio data: %w", err)
	}
	
	// Convert byte data to float32 samples
	numSamples := len(data) / 4
	v.audioData = make([]float64, numSamples)
	
	for i := 0; i < numSamples; i++ {
		// Convert 4 bytes to float32 (little endian)
		idx := i * 4
		bits := uint32(data[idx]) | uint32(data[idx+1])<<8 | uint32(data[idx+2])<<16 | uint32(data[idx+3])<<24
		v.audioData[i] = float64(math.Float32frombits(bits))
	}
	
	return nil
}

// precomputeSpectrum pre-computes all spectrum data for the video
func (v *Visualizer) precomputeSpectrum() error {
	v.windowSize = 2048
	hopLength := v.sampleRate / v.config.FPS
	
	v.spectrumData = make([][]float64, v.totalFrames)
	
	// Process each frame
	for frame := 0; frame < v.totalFrames; frame++ {
		startIdx := frame * hopLength
		endIdx := startIdx + v.windowSize
		
		if endIdx > len(v.audioData) {
			// Pad with zeros if necessary
			endIdx = len(v.audioData)
		}
		
		// Get window of audio data
		window := make([]float64, v.windowSize)
		if startIdx < len(v.audioData) {
			copy(window, v.audioData[startIdx:endIdx])
		}
		
		// Apply window function (Hamming)
		for i := range window {
			window[i] *= 0.54 - 0.46*math.Cos(2*math.Pi*float64(i)/float64(v.windowSize-1))
		}
		
		// Compute FFT
		fftData := fft.FFTReal(window)
		
		// Convert to magnitude spectrum
		magnitudes := make([]float64, len(fftData)/2)
		for i := range magnitudes {
			magnitudes[i] = cmplx.Abs(fftData[i])
		}
		
		// Create frequency bins (logarithmic scale)
		v.spectrumData[frame] = v.binFrequencies(magnitudes)
		
		// Apply smoothing with more responsive factor
		if frame > 0 {
			for i := range v.spectrumData[frame] {
				v.spectrumData[frame][i] = v.spectrumData[frame][i]*0.85 + v.spectrumData[frame-1][i]*0.15
			}
		}
	}
	
	return nil
}

// binFrequencies bins the frequency data into the desired number of bars
func (v *Visualizer) binFrequencies(magnitudes []float64) []float64 {
	bins := make([]float64, v.config.BarCount)
	
	// Create logarithmic frequency bins from 80Hz to 8000Hz
	minFreq := 80.0
	maxFreq := 8000.0
	
	freqBins := make([]float64, v.config.BarCount+1)
	for i := 0; i <= v.config.BarCount; i++ {
		freqBins[i] = minFreq * math.Pow(maxFreq/minFreq, float64(i)/float64(v.config.BarCount))
	}
	
	// Map frequency bins to FFT bins
	fftBinWidth := float64(v.sampleRate) / float64(len(magnitudes)*2)
	
	for i := 0; i < v.config.BarCount; i++ {
		startBin := int(freqBins[i] / fftBinWidth)
		endBin := int(freqBins[i+1] / fftBinWidth)
		
		if startBin >= len(magnitudes) {
			startBin = len(magnitudes) - 1
		}
		if endBin >= len(magnitudes) {
			endBin = len(magnitudes) - 1
		}
		
		// Average the magnitudes in this frequency range
		sum := 0.0
		count := 0
		for j := startBin; j <= endBin && j < len(magnitudes); j++ {
			sum += magnitudes[j]
			count++
		}
		
		if count > 0 {
			bins[i] = sum / float64(count)
		}
		
		// Normalize magnitude (FFT magnitudes can be very large)
		// First divide by window size to get proper scale
		bins[i] = bins[i] / float64(v.windowSize)
		
		// Apply logarithmic scaling for better visual response
		if bins[i] > 0 {
			// Use log scale with adjustable sensitivity
			bins[i] = math.Log10(bins[i]*1000 + 1) / 3.0
			
			// Ensure within 0-1 range
			if bins[i] < 0 {
				bins[i] = 0
			} else if bins[i] > 1 {
				bins[i] = 1
			}
		}
	}
	
	return bins
}

// createVideoSequential creates the video frame by frame
func (v *Visualizer) createVideoSequential() error {
	// Create temporary directory for frames
	tempDir := filepath.Join(os.TempDir(), fmt.Sprintf("spectrum_frames_%d", time.Now().Unix()))
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		return fmt.Errorf("creating temp dir: %w", err)
	}
	defer os.RemoveAll(tempDir)
	
	// Generate frames
	for i := 0; i < v.totalFrames; i++ {
		if i%30 == 0 {
			fmt.Printf("Processing frame %d/%d (%.1f%%)\n", i, v.totalFrames, float64(i)/float64(v.totalFrames)*100)
		}
		
		frame := v.generateFrame(i)
		filename := filepath.Join(tempDir, fmt.Sprintf("frame_%06d.png", i))
		if err := frame.SavePNG(filename); err != nil {
			return fmt.Errorf("saving frame %d: %w", i, err)
		}
	}
	
	// Create video using ffmpeg
	return v.assembleVideo(tempDir)
}

// createVideoParallel creates the video using parallel processing
func (v *Visualizer) createVideoParallel() error {
	// Create temporary directory for frames
	tempDir := filepath.Join(os.TempDir(), fmt.Sprintf("spectrum_frames_%d", time.Now().Unix()))
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		return fmt.Errorf("creating temp dir: %w", err)
	}
	defer os.RemoveAll(tempDir)
	
	// Use worker pool
	numWorkers := runtime.NumCPU()
	fmt.Printf("Using %d CPU cores for parallel processing\n", numWorkers)
	
	type job struct {
		frameIdx int
		filename string
	}
	
	jobs := make(chan job, v.totalFrames)
	errors := make(chan error, numWorkers)
	
	// Start workers
	var wg sync.WaitGroup
	for w := 0; w < numWorkers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := range jobs {
				frame := v.generateFrame(j.frameIdx)
				if err := frame.SavePNG(j.filename); err != nil {
					errors <- fmt.Errorf("saving frame %d: %w", j.frameIdx, err)
					return
				}
			}
		}()
	}
	
	// Send jobs
	for i := 0; i < v.totalFrames; i++ {
		jobs <- job{
			frameIdx: i,
			filename: filepath.Join(tempDir, fmt.Sprintf("frame_%06d.png", i)),
		}
		
		if i%30 == 0 {
			fmt.Printf("Queued frame %d/%d (%.1f%%)\n", i, v.totalFrames, float64(i)/float64(v.totalFrames)*100)
		}
	}
	close(jobs)
	
	// Wait for completion
	wg.Wait()
	close(errors)
	
	// Check for errors
	for err := range errors {
		if err != nil {
			return err
		}
	}
	
	// Create video using ffmpeg
	return v.assembleVideo(tempDir)
}

// assembleVideo uses ffmpeg to create the final video
func (v *Visualizer) assembleVideo(frameDir string) error {
	fmt.Println("Assembling video with audio...")
	
	// Create video from frames and add audio
	cmd := exec.Command("ffmpeg",
		"-framerate", fmt.Sprintf("%d", v.config.FPS),
		"-i", filepath.Join(frameDir, "frame_%06d.png"),
		"-i", v.config.InputFile,
		"-c:v", "libx264",
		"-preset", "ultrafast",
		"-pix_fmt", "yuv420p",
		"-c:a", "aac",
		"-b:a", "192k",
		"-shortest",
		"-y", v.config.OutputFile,
	)
	
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	
	return cmd.Run()
}

// generateFrame generates a single frame of the visualization
func (v *Visualizer) generateFrame(frameIdx int) *gg.Context {
	dc := gg.NewContext(v.config.Width, v.config.Height)
	
	// Set background color
	bgColor := v.getBackgroundColor()
	dc.SetColor(bgColor)
	dc.Clear()
	
	// Get spectrum data for this frame
	var magnitudes []float64
	if frameIdx < len(v.spectrumData) {
		magnitudes = v.spectrumData[frameIdx]
	} else {
		magnitudes = make([]float64, v.config.BarCount)
	}
	
	// Draw visualization based on type
	switch v.config.VizType {
	case "circular":
		v.drawCircular(dc, magnitudes)
	case "wave":
		v.drawWave(dc, magnitudes)
	case "radial":
		v.drawRadial(dc, magnitudes)
	case "line":
		v.drawLine(dc, magnitudes)
	case "dots":
		v.drawDots(dc, magnitudes)
	case "mirror":
		v.drawMirror(dc, magnitudes)
	case "spiral":
		v.drawSpiral(dc, magnitudes)
	default: // "bars"
		v.drawBars(dc, magnitudes)
	}
	
	return dc
}

// Color helper functions
func (v *Visualizer) getBackgroundColor() color.Color {
	switch v.config.BgColor {
	case "blue":
		return color.RGBA{0, 0, 255, 255}
	case "magenta":
		return color.RGBA{255, 0, 255, 255}
	case "black":
		return color.RGBA{0, 0, 0, 255}
	case "white":
		return color.RGBA{255, 255, 255, 255}
	case "gray":
		return color.RGBA{128, 128, 128, 255}
	default: // "green"
		return color.RGBA{0, 255, 0, 255}
	}
}

func (v *Visualizer) getColor(magnitude float64) color.Color {
	switch v.config.ColorScheme {
	case "fire":
		return v.getFireColor(magnitude)
	case "ocean":
		return v.getOceanColor(magnitude)
	case "purple":
		return v.getPurpleColor(magnitude)
	case "neon":
		return v.getNeonColor(magnitude)
	case "monochrome":
		return v.getMonochromeColor(magnitude)
	case "sunset":
		return v.getSunsetColor(magnitude)
	case "forest":
		return v.getForestColor(magnitude)
	case "white":
		return color.RGBA{255, 255, 255, 255}
	default: // "rainbow"
		return v.getRainbowColor(magnitude)
	}
}

func (v *Visualizer) getRainbowColor(magnitude float64) color.Color {
	// HSV to RGB conversion for rainbow effect
	h := (1.0 - magnitude) * 120.0 / 360.0 // Green to Red
	s := 1.0
	val := 0.8 + magnitude*0.2
	
	return hsvToRGB(h, s, val)
}

func (v *Visualizer) getFireColor(magnitude float64) color.Color {
	// Fire color gradient: deep red -> bright red -> orange -> yellow -> yellow-green
	if magnitude < 0.2 {
		// Deep red to bright red
		t := magnitude / 0.2
		r := uint8(180 + t*75)  // 180 to 255
		g := uint8(0)
		return color.RGBA{r, g, 0, 255}
	} else if magnitude < 0.4 {
		// Bright red to orange-red
		t := (magnitude - 0.2) / 0.2
		r := uint8(255)
		g := uint8(t * 100)  // 0 to 100
		return color.RGBA{r, g, 0, 255}
	} else if magnitude < 0.6 {
		// Orange-red to bright orange
		t := (magnitude - 0.4) / 0.2
		r := uint8(255)
		g := uint8(100 + t*80)  // 100 to 180
		return color.RGBA{r, g, 0, 255}
	} else if magnitude < 0.8 {
		// Bright orange to yellow
		t := (magnitude - 0.6) / 0.2
		r := uint8(255)
		g := uint8(180 + t*75)  // 180 to 255
		return color.RGBA{r, g, 0, 255}
	} else {
		// Yellow to yellow-green (hottest)
		t := (magnitude - 0.8) / 0.2
		r := uint8(255 - t*55)  // 255 to 200
		g := uint8(255)
		b := uint8(t * 50)  // 0 to 50 for slight green tint
		return color.RGBA{r, g, b, 255}
	}
}

func (v *Visualizer) getOceanColor(magnitude float64) color.Color {
	// Dark blue to cyan
	b := uint8(150 + magnitude*105)
	g := uint8(magnitude * 200)
	return color.RGBA{0, g, b, 255}
}

func (v *Visualizer) getPurpleColor(magnitude float64) color.Color {
	// Purple to pink
	r := uint8(180 + magnitude*75)
	b := uint8(255 - magnitude*50)
	return color.RGBA{r, 0, b, 255}
}

func (v *Visualizer) getNeonColor(magnitude float64) color.Color {
	// Bright neon colors: cyan -> blue -> magenta -> pink -> green
	if magnitude < 0.2 {
		// Cyan to electric blue
		t := magnitude / 0.2
		r := uint8(0)
		g := uint8(255 - t*155)  // 255 to 100
		b := uint8(255)
		return color.RGBA{r, g, b, 255}
	} else if magnitude < 0.4 {
		// Electric blue to magenta
		t := (magnitude - 0.2) / 0.2
		r := uint8(t * 255)
		g := uint8(100 - t*100)  // 100 to 0
		b := uint8(255)
		return color.RGBA{r, g, b, 255}
	} else if magnitude < 0.6 {
		// Magenta to hot pink
		t := (magnitude - 0.4) / 0.2
		r := uint8(255)
		g := uint8(t * 100)  // 0 to 100
		b := uint8(255 - t*55)  // 255 to 200
		return color.RGBA{r, g, b, 255}
	} else if magnitude < 0.8 {
		// Hot pink to electric green
		t := (magnitude - 0.6) / 0.2
		r := uint8(255 - t*255)  // 255 to 0
		g := uint8(100 + t*155)  // 100 to 255
		b := uint8(200 - t*200)  // 200 to 0
		return color.RGBA{r, g, b, 255}
	} else {
		// Electric green to bright cyan
		t := (magnitude - 0.8) / 0.2
		r := uint8(0)
		g := uint8(255)
		b := uint8(t * 255)  // 0 to 255
		return color.RGBA{r, g, b, 255}
	}
}

func (v *Visualizer) getMonochromeColor(magnitude float64) color.Color {
	// Grayscale
	val := uint8(50 + magnitude*205)
	return color.RGBA{val, val, val, 255}
}

func (v *Visualizer) getSunsetColor(magnitude float64) color.Color {
	if magnitude < 0.5 {
		// Purple to red
		r := uint8(magnitude * 2 * 255)
		b := uint8(255 - magnitude*2*255)
		return color.RGBA{r, 0, b, 255}
	}
	// Red to orange
	r := uint8(255)
	g := uint8((magnitude - 0.5) * 2 * 180)
	return color.RGBA{r, g, 0, 255}
}

func (v *Visualizer) getForestColor(magnitude float64) color.Color {
	// Dark green to yellow-green
	r := uint8(magnitude * 150)
	g := uint8(100 + magnitude*155)
	return color.RGBA{r, g, 0, 255}
}

// HSV to RGB conversion helper
func hsvToRGB(h, s, v float64) color.Color {
	c := v * s
	x := c * (1 - math.Abs(math.Mod(h*6, 2)-1))
	m := v - c
	
	var r, g, b float64
	switch int(h * 6) {
	case 0:
		r, g, b = c, x, 0
	case 1:
		r, g, b = x, c, 0
	case 2:
		r, g, b = 0, c, x
	case 3:
		r, g, b = 0, x, c
	case 4:
		r, g, b = x, 0, c
	default:
		r, g, b = c, 0, x
	}
	
	return color.RGBA{
		uint8((r + m) * 255),
		uint8((g + m) * 255),
		uint8((b + m) * 255),
		255,
	}
}