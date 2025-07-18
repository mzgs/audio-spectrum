package audiospectrum

import (
	"math"

	"github.com/fogleman/gg"
)

// drawBars draws traditional bar spectrum
func (v *Visualizer) drawBars(dc *gg.Context, magnitudes []float64) {
	for i, magnitude := range magnitudes {
		// Calculate bar height
		baseHeight := 5.0
		var barHeight float64
		var displayMagnitude float64
		
		if magnitude < 0.05 { // Increase threshold for silence
			barHeight = baseHeight + magnitude*float64(v.config.Height)*0.2 // Show very small bars for low values
			displayMagnitude = magnitude * 2
		} else {
			barHeight = baseHeight + magnitude*float64(v.config.Height)*0.7
			displayMagnitude = magnitude
		}
		
		// Get color
		color := v.getColor(displayMagnitude)
		dc.SetColor(color)
		
		// Draw bar
		x := float64(v.barPositions[i])
		barWidth := float64(v.barWidth) * 0.8
		y := float64(v.config.Height) - barHeight
		
		dc.DrawRectangle(x, y, barWidth, barHeight)
		dc.Fill()
		
		// Add glow effect for louder parts
		if magnitude > 0.5 {
			dc.SetRGBA(1, 1, 1, 0.3)
			dc.DrawRectangle(x-2, y-2, barWidth+4, barHeight+4)
			dc.Fill()
		}
	}
}

// drawCircular draws circular spectrum with bars radiating outward
func (v *Visualizer) drawCircular(dc *gg.Context, magnitudes []float64) {
	angleStep := 2 * math.Pi / float64(v.config.BarCount)
	minRadius := 80.0
	maxRadius := math.Min(float64(v.config.Width), float64(v.config.Height))/2 - 50
	
	for i, magnitude := range magnitudes {
		angle := float64(i) * angleStep
		
		// Calculate radius
		var radius float64
		var displayMagnitude float64
		
		if magnitude < 0.01 {
			radius = minRadius + 5
			displayMagnitude = 0.1
		} else {
			radius = minRadius + magnitude*(maxRadius-minRadius)
			displayMagnitude = magnitude
		}
		
		// Get color
		color := v.getColor(displayMagnitude)
		dc.SetColor(color)
		
		// Calculate line endpoints
		x1 := float64(v.centerX) + minRadius*math.Cos(angle)
		y1 := float64(v.centerY) + minRadius*math.Sin(angle)
		x2 := float64(v.centerX) + radius*math.Cos(angle)
		y2 := float64(v.centerY) + radius*math.Sin(angle)
		
		// Draw thick line
		dc.SetLineWidth(8)
		dc.DrawLine(x1, y1, x2, y2)
		dc.Stroke()
		
		// Add glow for loud parts
		if magnitude > 0.5 {
			dc.SetLineWidth(12)
			dc.SetRGBA(1, 1, 1, 0.3)
			dc.DrawLine(x1, y1, x2, y2)
			dc.Stroke()
		}
	}
}

// drawWave draws wave-form spectrum
func (v *Visualizer) drawWave(dc *gg.Context, magnitudes []float64) {
	yCenter := float64(v.config.Height) / 2
	xStep := float64(v.config.Width) / float64(len(magnitudes))
	
	for i, magnitude := range magnitudes {
		x := float64(i) * xStep
		waveHeight := 20 + magnitude*150
		
		// Get color
		color := v.getColor(magnitude)
		dc.SetColor(color)
		
		// Draw vertical line from center
		dc.SetLineWidth(3)
		dc.DrawLine(x, yCenter-waveHeight, x, yCenter+waveHeight)
		dc.Stroke()
	}
}

// drawRadial draws radial burst spectrum
func (v *Visualizer) drawRadial(dc *gg.Context, magnitudes []float64) {
	angleStep := 2 * math.Pi / float64(v.config.BarCount)
	baseRadius := 50.0
	
	for i, magnitude := range magnitudes {
		angle := float64(i) * angleStep
		
		// Create wedge shape
		var length float64
		var displayMagnitude float64
		
		if magnitude < 0.01 {
			length = 10
			displayMagnitude = 0.1
		} else {
			length = 10 + magnitude*300
			displayMagnitude = magnitude
		}
		
		// Get color
		color := v.getColor(displayMagnitude)
		dc.SetColor(color)
		
		// Calculate wedge points
		angleWidth := angleStep * 0.8
		
		dc.MoveTo(
			float64(v.centerX)+baseRadius*math.Cos(angle-angleWidth/2),
			float64(v.centerY)+baseRadius*math.Sin(angle-angleWidth/2),
		)
		dc.LineTo(
			float64(v.centerX)+baseRadius*math.Cos(angle+angleWidth/2),
			float64(v.centerY)+baseRadius*math.Sin(angle+angleWidth/2),
		)
		dc.LineTo(
			float64(v.centerX)+(baseRadius+length)*math.Cos(angle+angleWidth/2),
			float64(v.centerY)+(baseRadius+length)*math.Sin(angle+angleWidth/2),
		)
		dc.LineTo(
			float64(v.centerX)+(baseRadius+length)*math.Cos(angle-angleWidth/2),
			float64(v.centerY)+(baseRadius+length)*math.Sin(angle-angleWidth/2),
		)
		dc.ClosePath()
		dc.Fill()
		
		// Add glow for loud parts
		if magnitude > 0.5 {
			dc.SetRGBA(1, 1, 1, 0.3)
			dc.SetLineWidth(3)
			dc.Stroke()
		}
	}
}

// drawLine draws connected line spectrum
func (v *Visualizer) drawLine(dc *gg.Context, magnitudes []float64) {
	xStep := float64(v.config.Width) / float64(len(magnitudes)-1)
	
	// Start path
	dc.MoveTo(0, float64(v.config.Height)-50-magnitudes[0]*float64(v.config.Height-100))
	
	// Draw connected lines
	for i := 1; i < len(magnitudes); i++ {
		x := float64(i) * xStep
		y := float64(v.config.Height) - 50 - magnitudes[i]*float64(v.config.Height-100)
		
		// Get color for this segment
		color := v.getColor(magnitudes[i])
		dc.SetColor(color)
		dc.SetLineWidth(5)
		
		dc.LineTo(x, y)
		dc.Stroke()
		dc.MoveTo(x, y)
		
		// Add glow for loud parts
		if magnitudes[i] > 0.5 {
			dc.SetRGBA(1, 1, 1, 0.3)
			dc.SetLineWidth(8)
			dc.Stroke()
			dc.MoveTo(x, y)
		}
	}
}

// drawDots draws dots/particles spectrum
func (v *Visualizer) drawDots(dc *gg.Context, magnitudes []float64) {
	xStep := float64(v.config.Width) / float64(len(magnitudes))
	
	for i, magnitude := range magnitudes {
		x := float64(i)*xStep + xStep/2
		
		// Create multiple dots at different heights
		numDots := int(1 + magnitude*10)
		
		for j := 0; j < numDots; j++ {
			y := float64(v.config.Height) - 20 - float64(j*30) - magnitude*300
			if y < 20 {
				break
			}
			
			// Get color
			color := v.getColor(magnitude * (1 - float64(j)/10))
			dc.SetColor(color)
			
			// Draw dot
			radius := 3 + magnitude*5
			dc.DrawCircle(x, y, radius)
			dc.Fill()
			
			// Add glow
			if magnitude > 0.5 {
				dc.SetRGBA(1, 1, 1, 0.3)
				dc.DrawCircle(x, y, radius+3)
				dc.Stroke()
			}
		}
	}
}

// drawMirror draws mirror spectrum - bars from center going up and down
func (v *Visualizer) drawMirror(dc *gg.Context, magnitudes []float64) {
	yCenter := float64(v.config.Height) / 2
	
	for i, magnitude := range magnitudes {
		// Calculate bar height
		var barHeight float64
		var displayMagnitude float64
		
		if magnitude < 0.01 {
			barHeight = 5
			displayMagnitude = 0.1
		} else {
			barHeight = 5 + magnitude*float64(v.config.Height)*0.35
			displayMagnitude = magnitude
		}
		
		// Get color
		color := v.getColor(displayMagnitude)
		dc.SetColor(color)
		
		// Draw bars going up and down from center
		x := float64(v.barPositions[i])
		barWidth := float64(v.barWidth) * 0.8
		
		// Upper bar
		dc.DrawRectangle(x, yCenter-barHeight, barWidth, barHeight)
		dc.Fill()
		
		// Lower bar
		dc.DrawRectangle(x, yCenter, barWidth, barHeight)
		dc.Fill()
		
		// Add glow for loud parts
		if magnitude > 0.5 {
			dc.SetRGBA(1, 1, 1, 0.3)
			dc.DrawRectangle(x-2, yCenter-barHeight-2, barWidth+4, barHeight*2+4)
			dc.Stroke()
		}
	}
}

// drawSpiral draws spiral spectrum
func (v *Visualizer) drawSpiral(dc *gg.Context, magnitudes []float64) {
	turns := 2.0 // Number of spiral turns
	maxRadius := math.Min(float64(v.config.Width), float64(v.config.Height))/2 - 50
	
	for i := 0; i < len(magnitudes); i++ {
		magnitude := magnitudes[i]
		
		// Calculate spiral parameters
		angleStart := (float64(i) / float64(len(magnitudes))) * 2 * math.Pi * turns
		angleEnd := (float64(i+1) / float64(len(magnitudes))) * 2 * math.Pi * turns
		
		// Get color
		color := v.getColor(magnitude)
		dc.SetColor(color)
		
		// Create points along the spiral segment
		prevX, prevY := 0.0, 0.0
		for j := 0; j < 10; j++ {
			t := float64(j) / 9
			angle := angleStart + t*(angleEnd-angleStart)
			radius := 50 + (angle/(2*math.Pi*turns))*maxRadius + magnitude*50
			
			x := float64(v.centerX) + radius*math.Cos(angle)
			y := float64(v.centerY) + radius*math.Sin(angle)
			
			if j > 0 {
				thickness := 2 + magnitude*8
				dc.SetLineWidth(thickness)
				dc.DrawLine(prevX, prevY, x, y)
				dc.Stroke()
			}
			
			prevX, prevY = x, y
		}
	}
}