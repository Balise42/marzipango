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
	X int64
	Y int64
}

type coordsFloat struct {
	X float64
	Y float64
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

func isBlack(c color.Color) bool {
	r, g, b, _ := c.RGBA()
	return r == 0 && g == 0 && b == 0
}

// distance field computation from https://prideout.net/blog/distance_fields/
func findHullParabolas(row []float64) ([]int, []float64) {
	v := make([]int, len(row))
	z := make([]float64, len(row)+1)
	k := 0

	v[0] = 0
	z[0] = -math.MaxInt16
	z[1] = math.MaxInt16

	for i := 1; i < len(row); i++ {
		q := i
		p := v[k]
		s := intersectParabolas(p, q, row)
		for s <= z[k] {
			k = k - 1
			p = v[k]
			s = intersectParabolas(p, q, row)
		}
		k = k + 1
		v[k] = q
		z[k] = s
		z[k+1] = math.MaxInt16
	}

	return v, z
}

func intersectParabolas(p int, q int, row []float64) float64 {
	intersect := ((row[q] + float64(q*q)) - (row[p] + float64(p*p))) / (2*float64(q) - 2*float64(p))
	return intersect
}

func marchParabolas(row []float64, vertices []int, intersections []float64) {
	k := 0
	for q := range row {
		for intersections[k+1] < float64(q) {
			k = k + 1
		}
		dx := q - vertices[k]
		row[q] = float64(dx*dx) + row[vertices[k]]
	}
}

func horizontalPass(row []float64) {
	vertices, intersections := findHullParabolas(row)
	marchParabolas(row, vertices, intersections)
}

func transpose(field [][]float64) [][]float64 {
	transposed := make([][]float64, len(field[0]))
	for x := range transposed {
		transposed[x] = make([]float64, len(field))
		for y := range transposed[x] {
			transposed[x][y] = field[y][x]
		}
	}
	return transposed
}

func computeEdt(img image.Image, width int, height int, maxvalue int) map[coords]float64 {

	field := make([][]float64, width+maxvalue*2)
	for x := range field {
		field[x] = make([]float64, height+maxvalue*2)
		for y := range field[x] {
			if x > maxvalue && x < img.Bounds().Dx()+maxvalue && y > maxvalue && y < img.Bounds().Dy()+maxvalue && isBlack(img.At(x-maxvalue, y-maxvalue)) {
				field[x][y] = 0
			} else {
				field[x][y] = math.MaxInt16
			}
		}
	}

	for _, row := range field {
		horizontalPass(row)
	}

	field = transpose(field)

	for _, row := range field {
		horizontalPass(row)
	}

	field = transpose(field)

	return convertField(field, maxvalue, maxvalue)
}

func doNothing(x int, y int) {}

func convertField(field [][]float64, offsetX int, offsetY int) map[coords]float64 {
	offsetField := make(map[coords]float64)
	for x := range field {
		for y := range field[x] {
			offsetField[coords{int64(x - offsetX), int64(y - offsetY)}] = math.Sqrt(field[x][y])
		}
	}
	return offsetField
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

	distances := computeEdt(img, params.Width, params.Height, int(maxvalue))

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

	dist, ok := im.Distances[coords{int64(xImg), int64(yImg)}]

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
