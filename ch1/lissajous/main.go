// lissajous.go
package main

import (
	"fmt"
	"image"
	"image/color"
	"image/gif"
	"io"
	"log"
	"math"
	"math/rand"
	"os"
	"time"
)

func main() {
	// // Create or open the out.gif file explicitly
	outFile, err := os.Create("out.gif")
	if err != nil {
		log.Fatalf("Error creating output file: %v", err)
	}
	defer outFile.Close()

	// Seed the random number generator using the current time
	rand.Seed(time.Now().UTC().UnixNano())

	// Generate the Lissajous figures and write them to out.gif
	lissajous(outFile)

}

func lissajous(out io.Writer) {
	const (
		cycles  = 5     // number of complete x oscillator revolutions
		res     = 0.001 // angular resolution
		size    = 300   // image canvas covers [-size..+size]
		nframes = 64    // number of animation frames
		delay   = 8     // delay between frames in 10ms units
	)
	freq := rand.Float64() * 3.0 // relative frequency of y oscillator
	anim := gif.GIF{LoopCount: nframes}
	phase := 0.0 // phase difference

	for i := 0; i < nframes; i++ {
		// Create a new image with a gradient palette
		rect := image.Rect(0, 0, 2*size+1, 2*size+1)
		img := image.NewPaletted(rect, generateGradientPalette(i, nframes))

		// Generate points on the Lissajous curve
		for t := 0.0; t < cycles*2*math.Pi; t += res {
			x := math.Sin(t)
			y := math.Sin(t*freq + phase)

			// Dynamically choose a color index based on the position
			colorIndex := uint8(1 + int(t*10)%len(img.Palette))
			img.SetColorIndex(size+int(x*size+0.5), size+int(y*size+0.5), colorIndex)
		}

		// Adjust the phase for the next frame
		phase += 0.1

		// Add the current frame to the animation
		anim.Delay = append(anim.Delay, delay)
		anim.Image = append(anim.Image, img)
	}

	// Write the GIF to the output file
	err := gif.EncodeAll(out, &anim)
	if err != nil {
		log.Fatalf("Error occurred while encoding the GIF: %v", err)
	}
	fmt.Println("GIF generated successfully and saved as out.gif")
}

// generateGradientPalette dynamically creates a gradient palette
func generateGradientPalette(frame, totalFrames int) []color.Color {
	colors := make([]color.Color, 256)
	for i := 0; i < 256; i++ {
		// Generate a smooth gradient
		red := uint8((i + frame*5) % 256)
		green := uint8((255 - i + frame*5) % 256)
		blue := uint8((i*2 + frame*10) % 256)
		colors[i] = color.RGBA{R: red, G: green, B: blue, A: 255}
	}
	return colors
}
