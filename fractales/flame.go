package fractales

import (
	"github.com/Balise42/marzipango/fractales/orbits"
	"github.com/Balise42/marzipango/params"
	"image"
	"image/color"
	"math"
	"math/rand"
	"sync"
)

func CreateFlameComputer(params params.ImageParams) Computation {
	flameFuncs := createFlameFuncs()
	ifsMap := createFlameMap(params, flameFuncs)

	comp := func(x int, ymin int, ymax int, img *image.RGBA64, wg *sync.WaitGroup) {
		for y := ymin; y < ymax; y++ {
			val, ok := ifsMap[orbits.Coords{int64(x), int64(y)}]
			if ok {
				img.Set(x, y, val)
			} else {
				img.Set(x, y, params.Palette.Divergence)
			}
		}
		wg.Done()
	}

	return comp
}

type ifsFunc func(float64, float64) (float64, float64)

func createFlameFuncs() []ifsFunc {
	V0 := func(x float64, y float64) (float64, float64) {
		sqr := math.Sqrt(x*x + y*y)
		return 1/sqr * ((x-y)*(x+y)), 1/sqr * 2*x*y
	}

	V1 := func(x float64, y float64) (float64, float64) {
		return math.Sin(x), y
	}

	V2 := func(x float64, y float64) (float64, float64) {
		sqr := math.Sqrt(x*x + y*y)
		theta := math.Atan(x / y)
		return 1/sqr * (math.Cos(theta) + math.Sin(sqr)), 1/sqr * (math.Sin(theta) - math.Cos(sqr))
	}

	V3 := func(x float64, y float64) (float64, float64) {
		if x >= 0 && y >= 0 {
			return x, y
		} else if x < 0 && y >= 0 {
			return 2*x, y
		} else if x >= 0 && y < 0 {
			return x, y/2
		} else {
			return 2*x, y/2
		}
	}

	V4 := func(x float64, y float64) (float64, float64) {
		sqr := math.Sqrt(x*x + y*y)
		theta := math.Atan(x / y)
		return math.Sin(theta) * math.Cos(sqr), math.Cos(theta) * math.Sin(sqr)
	}

	return []ifsFunc{V0, V1, V2, V3, V4}
}

type triplet struct {
	R float64
	G float64
	B float64
	A float64
}

func createFlameMap(params params.ImageParams, funcs []ifsFunc) map[orbits.Coords]color.NRGBA {
	imgRes := make(map[orbits.Coords]color.NRGBA)
	x := float64(0)
	y := float64(0)
	rf := []float64{1.0, 1.0, 1.0, 1.0, 1.0}
	gf := []float64{0.0, 0.1, 0.2, 0.3, 0.4}
	bf := []float64{0, 0, 0, 0, 0}

	histo := make(map[orbits.Coords]int)
	cols := make(map[orbits.Coords]triplet)
	col := triplet{1.0, 0, 0, 0}

	maxValue := 0

	for i := 0; i < 500000000; i++ {
		rule := rand.Float32()
		var a, b, c, d, e, f float64
		var funcIndex int
		if rule < 0.08 {
			a = -0.98
			b = -0.12
			c = -0.6
			d = 0.01
			e = -0.028
			f = 0.07
			funcIndex = 0
		} else if rule < 0.8 {
			a = -0.5
			b = 0.43
			c = -0.06
			d = -0.44
			e = -0.09
			f = -0.88
			funcIndex = 1
		} else if rule < 0.85 {
			a = 0.18
			b = -0.12
			c = -0.18
			d = 0.04
			e = 0.18
			f = 0.40
			funcIndex = 2
		} else if rule < 0.87 {
			a = 1.62
			b = 1.03
			c = 0.59
			d = -0.66
			e = 0.25
			f = -0.72
			funcIndex = 3
		} else {
			a = 0.02
			b = 0.13
			c = -1.17
			d = -1.44
			e = -0.17
			f = -0.14
			funcIndex = 4
		}
		x1 := a*x + b*y + e
		y1 := c*x + d*y + f
		x1, y1 = funcs[funcIndex](x1, y1)

		coords := scaleFlame(x1, y1, params)
		//col := cols[coords]
		col = triplet{ (col.R + rf[funcIndex]) / 2, col.G + gf[funcIndex], col.B + bf[funcIndex], 1.0 }
		cols[coords] = col
		if i > 20 {
			histo[coords] = histo[coords] + 1
			cols[coords] = col
			if histo[coords] > maxValue {
				maxValue = histo[coords]
			}
		}
		x = x1
		y = y1
	}


	for k, col := range cols {
		alpha := math.Log(float64(histo[k]) + 1) / math.Log(float64(maxValue) + 1)
		if alpha < 0 {
			alpha = 0
		}
		tmpR := col.R//math.Pow(col.R * alpha, 0.5)
		tmpG := col.G//math.Pow(col.G * alpha, 0.5)
		tmpB := col.B//math.Pow(col.B * alpha, 0.5)
		imgRes[k] = color.NRGBA{R: uint8(tmpR * 255), G: uint8(tmpG * 255), B: uint8(tmpB * 255), A: uint8((alpha) * 255)}
	}
	return imgRes
}

func scaleFlame(x1 float64, y1 float64, imageParams params.ImageParams) orbits.Coords {
	x := (x1 + 1.0) * float64(imageParams.Width)
	y := (y1 + 1.0) * float64(imageParams.Height)
	return orbits.Coords{int64(x), int64(y)}
}
