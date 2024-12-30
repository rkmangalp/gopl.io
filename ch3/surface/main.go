package main

import (
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"
)

const (
	defaultWidth, defaultHeight = 600, 320 // Default canvas size in pixels
	cells                       = 100      // Number of grid cells
	xyrange                     = 30.0     // Axis ranges (-xyrange..+xyrange)
	angle                       = math.Pi / 6
)

var sin30, cos30 = math.Sin(angle), math.Cos(angle) // sin(30°), cos(30°)

func main() {
	http.HandleFunc("/surface", surfaceHandler)
	log.Println("Starting server on :8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func surfaceHandler(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters for width, height, and color
	width := parseQueryParam(r, "width", defaultWidth)
	height := parseQueryParam(r, "height", defaultHeight)
	color := r.URL.Query().Get("color")
	if color == "" {
		color = "blue" // Default color
	}

	xyscale := float64(width) / 2 / xyrange // Pixels per x or y unit
	zscale := float64(height) * 0.4         // Pixels per z unit

	// Set Content-Type header
	w.Header().Set("Content-Type", "image/svg+xml")

	// Write SVG to the response
	fmt.Fprintf(w, "<svg xmlns='http://www.w3.org/2000/svg' "+
		"style='stroke: grey; fill: %s; stroke-width: 0.7' "+
		"width='%d' height='%d'>", color, width, height)
	for i := 0; i < cells; i++ {
		for j := 0; j < cells; j++ {
			ax, ay, validA := corner(i+1, j, width, height, xyscale, zscale)
			bx, by, validB := corner(i, j, width, height, xyscale, zscale)
			cx, cy, validC := corner(i, j+1, width, height, xyscale, zscale)
			dx, dy, validD := corner(i+1, j+1, width, height, xyscale, zscale)

			// Skip polygon if any corner is invalid
			if !(validA && validB && validC && validD) {
				continue
			}

			fmt.Fprintf(w, "<polygon points='%g,%g %g,%g %g,%g %g,%g'/>\n",
				ax, ay, bx, by, cx, cy, dx, dy)
		}
	}
	fmt.Fprintln(w, "</svg>")
}

func parseQueryParam(r *http.Request, param string, defaultValue int) int {
	// Parse a query parameter as an integer, return the default value if not provided or invalid
	valueStr := r.URL.Query().Get(param)
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return defaultValue
	}
	return value
}

func corner(i, j, width, height int, xyscale, zscale float64) (float64, float64, bool) {
	// Find point (x, y) at corner of cell (i, j)
	x := xyrange * (float64(i)/cells - 0.5)
	y := xyrange * (float64(j)/cells - 0.5)

	// Compute surface height z
	z := f(x, y)
	if math.IsNaN(z) || math.IsInf(z, 0) {
		return 0, 0, false // Return invalid flag if z is not finite
	}

	// Project (x, y, z) isometrically onto 2D SVG canvas (sx, sy)
	sx := float64(width)/2 + (x-y)*cos30*xyscale
	sy := float64(height)/2 + (x+y)*sin30*xyscale - z*zscale
	return sx, sy, true
}

func f(x, y float64) float64 {
	r := math.Hypot(x, y) // distance from (0,0)
	return math.Sin(r) / r

	// Egg Box
	// return math.Sin(x) * math.Sin(y)

	// Moguls
	// return math.Cos(x) + math.Cos(y)

	// Saddle
	// return x*x - y*y
}
