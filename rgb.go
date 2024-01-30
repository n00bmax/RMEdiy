package main

import (
	"image"
	"image/color"
	"math"
	"time"

	"github.com/kbinani/screenshot"
	"github.com/lucasb-eyer/go-colorful"
	"gitlab.com/gomidi/midi/v2"
	"k8s.io/klog/v2"
)

// ColorInfo represents a color and its associated name and index
type ColorInfo struct {
	Name  string
	Index int
	R     uint8
	G     uint8
	B     uint8
}

func captureScreenArea(x, y, width, height int) (*image.RGBA, error) {
	// screenshot.CaptureDisplay()
	var all image.Rectangle = image.Rect(0, 0, 0, 0)

	bounds := screenshot.GetDisplayBounds(0)
	all = bounds.Union(all)
	img, err := screenshot.CaptureRect(bounds)
	if err != nil {
		return nil, err
	}

	return img, nil
}

func calculateAverageColor(img *image.RGBA, darkThreshold float64) colorful.Color {
	var totalR, totalG, totalB float64
	var pixelCount int

	for i := 0; i < len(img.Pix); i += 4 {
		px := colorful.Color{
			R: float64(img.Pix[i]) / 255,
			G: float64(img.Pix[i+1]) / 255,
			B: float64(img.Pix[i+2]) / 255,
		}

		// Use the Clamped method to ensure the RGB values are in the range [0, 1]
		clampedPx := px.Clamped()

		// Skip dark colors
		if (clampedPx.R+clampedPx.G+clampedPx.B)/3 > darkThreshold {
			totalR += clampedPx.R
			totalG += clampedPx.G
			totalB += clampedPx.B
			pixelCount++
		}
	}

	if pixelCount > 0 {
		return colorful.Color{
			R: totalR / float64(pixelCount),
			G: totalG / float64(pixelCount),
			B: totalB / float64(pixelCount),
		}
	}

	return colorful.Color{}
}

var predefinedColors = map[string]ColorInfo{
	"Green": {Name: "Green", Index: 0, R: 0, G: 255, B: 0},
	"Cyan":  {Name: "Cyan", Index: 1, R: 0, G: 255, B: 255},
	"Amber": {Name: "Amber", Index: 2, R: 200, G: 191, B: 0},
	// "Monochrome":     {Name: "Monochrome", Index: 3, R: 255, G: 255, B: 255},
	"Red":    {Name: "Red", Index: 4, R: 255, G: 0, B: 0},
	"Orange": {Name: "Orange", Index: 5, R: 255, G: 165, B: 0},
}

var RGBStop = make(chan struct{})

var RGBticker *time.Ticker // time.NewTicker(time.Duration(config.Sync.RGB.RefreshRate) * time.Millisecond)

func SyncRGB() {
	mess := []byte{}
	go func() {

		for range RGBticker.C {

			x, y, width, height := 200, 200, 200, 200 // Adjust these coordinates and dimensions as needed
			img, err := captureScreenArea(x, y, width, height)
			if err != nil {
				klog.Fatal(err)
				continue
			}

			mess = append(rmeSysExCommandBase, generateDeviceCommand(33, rgbToClosestColorIndex(calculateDominantColor(img)))...)
			send, err := midi.SendTo(out)
			mess = append(mess, SysExClose...)

			if err != nil {
				klog.Infof("ERROR: %s\n", err)
				return
			}

			klog.Infof("%#v", mess)

			send(mess)

		}

	}()
}

func meterColorChangeRequest(color int) {
	mess := append(rmeSysExCommandBase, generateDeviceCommand(33, color)...)
	send, err := midi.SendTo(out)
	mess = append(mess, SysExClose...)
	if err != nil {
		klog.Infof("ERROR: %s\n", err)
		return
	}
	klog.V(2).Infof("%#v", mess)
	send(mess)
}
func calculateDominantColor(img *image.RGBA) (int, int, int) {
	colorProminence := calculateColorProminence(img)
	dominantColor := findDominantColor(colorProminence)
	return int(dominantColor.R), int(dominantColor.G), int(dominantColor.B)
}

func calculateColorProminence(img *image.RGBA) map[color.RGBA]float64 {
	colorProminence := make(map[color.RGBA]float64)
	bounds := img.Bounds()

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			c := color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)}
			brightness := calculateBrightness(c)
			colorProminence[c] += brightness
		}
	}

	return colorProminence
}

func calculateBrightness(c color.RGBA) float64 {
	return 0.299*float64(c.R) + 0.587*float64(c.G) + 0.114*float64(c.B)
}

func findDominantColor(colorProminence map[color.RGBA]float64) color.RGBA {
	var dominantColor color.RGBA
	maxProminence := 0.0

	for c, prominence := range colorProminence {
		if prominence > maxProminence {
			maxProminence = prominence
			dominantColor = c
		}
	}

	return dominantColor
}

func rgbToClosestColorIndex(r, g, b int) int {
	minDistance := math.MaxFloat64
	closestIndex := -1
	// monochrome selection
	klog.Info(r, g, b)
	if (r <= 20 && g <= 20 && b <= 20) || (r >= 200 && g >= 200 && b >= 200) {
		return 3
	}

	targetColor := color.RGBA{uint8(r), uint8(g), uint8(b), 255}

	for _, colorInfo := range predefinedColors {
		predefinedColor := color.RGBA{colorInfo.R, colorInfo.G, colorInfo.B, 255}

		distance := colorDistance(targetColor, predefinedColor)
		if distance < minDistance {
			minDistance = distance
			closestIndex = colorInfo.Index
		}
	}

	return closestIndex
}

func colorDistance(c1, c2 color.RGBA) float64 {

	// Calculate Euclidean distance between colors
	rDiff := int(c2.R) - int(c1.R)
	gDiff := int(c2.G) - int(c1.G)
	bDiff := int(c2.B) - int(c1.B)

	distance := math.Sqrt(float64(rDiff*rDiff + gDiff*gDiff + bDiff*bDiff))

	return distance
}
