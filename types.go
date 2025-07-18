package audiospectrum

// ColorScheme represents the available color schemes for visualization
type ColorScheme string

// Available color schemes
const (
	ColorSchemeRainbow    ColorScheme = "rainbow"    // Classic green to red gradient
	ColorSchemeFire       ColorScheme = "fire"       // Dark red to bright yellow
	ColorSchemeOcean      ColorScheme = "ocean"      // Dark blue to cyan
	ColorSchemePurple     ColorScheme = "purple"     // Purple to pink gradient
	ColorSchemeNeon       ColorScheme = "neon"       // Bright full spectrum colors
	ColorSchemeMonochrome ColorScheme = "monochrome" // White gradient
	ColorSchemeSunset     ColorScheme = "sunset"     // Purple to orange
	ColorSchemeForest     ColorScheme = "forest"     // Dark green to yellow-green
	ColorSchemeWhite      ColorScheme = "white"      // Pure white bars
)

// VisType represents the available visualization types
type VisType string

// Available visualization types
const (
	VisTypeBars     VisType = "bars"     // Traditional vertical bars
	VisTypeCircular VisType = "circular" // Bars radiating outward from center
	VisTypeWave     VisType = "wave"     // Waveform visualization
	VisTypeRadial   VisType = "radial"   // Radial burst pattern
	VisTypeLine     VisType = "line"     // Connected line graph spectrum
	VisTypeDots     VisType = "dots"     // Particle/dots effect
	VisTypeMirror   VisType = "mirror"   // Mirrored bars from center
	VisTypeSpiral   VisType = "spiral"   // Spiral pattern
)

// BGColor represents the available background colors
type BGColor string

// Available background colors
const (
	BGColorGreen   BGColor = "green"   // Green chroma key
	BGColorBlue    BGColor = "blue"    // Blue chroma key
	BGColorMagenta BGColor = "magenta" // Magenta chroma key
	BGColorBlack   BGColor = "black"   // Solid black background
	BGColorWhite   BGColor = "white"   // Solid white background
	BGColorGray    BGColor = "gray"    // Solid gray background
)

// ProcessType represents the processing method
type ProcessType string

// Available processing types
const (
	ProcessTypeFast     ProcessType = "fast"     // Sequential processing
	ProcessTypeParallel ProcessType = "parallel" // Parallel processing using all CPU cores
)

// String returns the string representation of ColorScheme
func (c ColorScheme) String() string {
	return string(c)
}

// IsValid checks if the color scheme is valid
func (c ColorScheme) IsValid() bool {
	switch c {
	case ColorSchemeRainbow, ColorSchemeFire, ColorSchemeOcean, ColorSchemePurple,
		ColorSchemeNeon, ColorSchemeMonochrome, ColorSchemeSunset, ColorSchemeForest, ColorSchemeWhite:
		return true
	}
	return false
}

// String returns the string representation of VisType
func (v VisType) String() string {
	return string(v)
}

// IsValid checks if the visualization type is valid
func (v VisType) IsValid() bool {
	switch v {
	case VisTypeBars, VisTypeCircular, VisTypeWave, VisTypeRadial,
		VisTypeLine, VisTypeDots, VisTypeMirror, VisTypeSpiral:
		return true
	}
	return false
}

// String returns the string representation of BGColor
func (b BGColor) String() string {
	return string(b)
}

// IsValid checks if the background color is valid
func (b BGColor) IsValid() bool {
	switch b {
	case BGColorGreen, BGColorBlue, BGColorMagenta, BGColorBlack, BGColorWhite, BGColorGray:
		return true
	}
	return false
}

// String returns the string representation of ProcessType
func (p ProcessType) String() string {
	return string(p)
}

// IsValid checks if the process type is valid
func (p ProcessType) IsValid() bool {
	return p == ProcessTypeFast || p == ProcessTypeParallel
}