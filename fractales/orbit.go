package fractales

import (
	"image"
	"image/color"
	"image/png"
	"math"
	"os"

	"github.com/Balise42/marzipango/params"
)

type Orbit interface {
	getOrbitFastValue(z complex128) float64
	getOrbitValue(v float64) float64
}

type PointOrbit struct {
	X           float64
	Y           float64
	Translation float64
	Factor      float64
}

type LineOrbit struct {
	A           float64
	B           float64
	C           float64
	Sqrtab      float64
	Translation float64
	Factor      float64
}

type coords struct {
	X int
	Y int
}

type ImageOrbit struct {
	Distances   map[coords]float64
	Translation float64
	Factor      float64
	Width       int
	Height      int
}

func CreatePointOrbit(x float64, y float64, maxvalue float64) PointOrbit {
	orbit := PointOrbit{X: x, Y: y}
	minDist := 0.0
	maxDist := orbit.squaredDistance(-2 - 1i)

	maxDist = math.Max(maxDist, orbit.squaredDistance(-2+1i))
	maxDist = math.Max(maxDist, orbit.squaredDistance(1+1i))
	maxDist = math.Max(maxDist, orbit.squaredDistance(1-1i))

	minDist = math.Min(minDist, orbit.squaredDistance(-2+1i))
	minDist = math.Min(minDist, orbit.squaredDistance(1+1i))
	minDist = math.Min(minDist, orbit.squaredDistance(1-1i))

	maxDist = math.Sqrt(maxDist)
	minDist = math.Sqrt(minDist)

	orbit.Factor = (maxvalue - minDist) / (maxDist - minDist)
	orbit.Translation = minDist

	return orbit
}

func (p PointOrbit) getOrbitFastValue(z complex128) float64 {
	return p.squaredDistance(z)
}

func (p PointOrbit) getOrbitValue(v float64) float64 {
	return (math.Sqrt(v) - p.Translation) * p.Factor
}

func (p PointOrbit) squaredDistance(z complex128) float64 {
	return (real(z)-p.X)*(real(z)-p.X) + (imag(z)-p.Y)*(imag(z)-p.Y)
}

func CreateLineOrbit(a float64, b float64, c float64, maxvalue float64) LineOrbit {
	orbit := LineOrbit{A: a, B: b, C: c}
	orbit.Sqrtab = math.Sqrt(a*a + b*b)

	minDist := 0.0
	maxDist := orbit.getOrbitFastValue(-2 - 1i)

	maxDist = math.Max(maxDist, orbit.getOrbitFastValue(-2+1i))
	maxDist = math.Max(maxDist, orbit.getOrbitFastValue(1+1i))
	maxDist = math.Max(maxDist, orbit.getOrbitFastValue(1-1i))

	minDist = math.Min(minDist, orbit.getOrbitFastValue(-2+1i))
	minDist = math.Min(minDist, orbit.getOrbitFastValue(1+1i))
	minDist = math.Min(minDist, orbit.getOrbitFastValue(1-1i))

	maxDist = math.Sqrt(maxDist) / orbit.Sqrtab
	minDist = math.Sqrt(minDist) / orbit.Sqrtab

	orbit.Factor = (maxvalue - minDist) / (maxDist - minDist)
	orbit.Translation = minDist

	return orbit
}

func (l LineOrbit) getOrbitFastValue(z complex128) float64 {
	lineCoeff := l.A*real(z) + l.B*imag(z) + l.C
	return lineCoeff * lineCoeff
}

func (l LineOrbit) getOrbitValue(v float64) float64 {
	return (math.Sqrt(v)/l.Sqrtab - l.Translation) * l.Factor
}

func getNeighbors(x int, y int) []coords {
	return []coords{
		{x - 1, y - 1},
		{x - 1, y},
		{x - 1, y + 1},
		{x, y - 1},
		{x, y + 1},
		{x + 1, y - 1},
		{x + 1, y},
		{x + 1, y + 1}}
}

func isValid(c coords, width int, height int) bool {
	return c.X >= -100 && c.X <= width+100 && c.Y >= -100 && c.Y <= height+100
}

func isBlack(c color.Color) bool {
	r, g, b, _ := c.RGBA()
	return r == 0 && g == 0 && b == 0
}

func computeDistances(img image.Image, width int, height int) map[coords]float64 {
	distances := make(map[coords]float64)
	queue := make([]coords, 0)

	for x := 0; x < img.Bounds().Dx(); x++ {
		for y := 0; y < img.Bounds().Dy(); y++ {
			if isBlack(img.At(x, y)) {
				distances[coords{x, y}] = 0
				queue = append(queue, coords{x, y})
			}
		}
	}

	for len(queue) > 0 {
		v := queue[0]
		for _, neigh := range getNeighbors(v.X, v.Y) {
			d, ok := distances[neigh]
			alt := distances[v] + 1
			if isValid(neigh, width, height) && (!ok || d > alt) {
				distances[neigh] = alt
				if alt > 100 {
					distances[neigh] = 100
				}
				queue = append(queue, neigh)
			}
		}
		queue = queue[1:]
	}
	return distances
}

func createGrey(g float64) color.Color {
	return color.RGBA{uint8(g), uint8(g), uint8(g), 255}
}

func CreateImageOrbit(params params.ImageParams, path string, maxvalue float64) (ImageOrbit, error) {
	f, err := os.Open(path)
	if err != nil {
		return ImageOrbit{}, err
	}

	img, err := png.Decode(f)

	if err != nil {
		return ImageOrbit{}, err
	}

	distances := computeDistances(img, params.Width, params.Height)

	minDist := 0.0
	maxDist := 0.0

	for _, v := range distances {
		if maxDist < v {
			maxDist = v
		}
	}

	factor := (maxvalue - minDist) / (maxDist - minDist)
	translation := minDist

	return ImageOrbit{Distances: distances, Factor: factor, Translation: translation, Width: img.Bounds().Dx(), Height: img.Bounds().Dy()}, nil
}

func (im ImageOrbit) getOrbitFastValue(z complex128) float64 {
	x := real(z)
	y := imag(z)

	hFloat := float64(im.Height)
	wFloat := float64(im.Width)

	xOffset := 2.0
	yOffset := 1.0
	xFactor := wFloat / 3
	yFactor := hFloat / 2

	xImg := x*xFactor + xOffset/xFactor
	yImg := y*yFactor + yOffset/yFactor

	dist, ok := im.Distances[coords{int(xImg), int(yImg)}]

	if !ok {
		return math.MaxInt64
	}
	return dist
}

func (im ImageOrbit) getOrbitValue(v float64) float64 {
	if v == math.MaxInt64 {
		return math.MaxInt64
	}
	return (v - im.Translation) * im.Factor
}
