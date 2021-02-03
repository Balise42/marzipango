package fractales

import (
	"github.com/Balise42/marzipango/params"
	"image"
	"image/color"
	"math"
	"math/rand"
	"sync"
)

func CreateFlameComputer(params params.ImageParams) (Computation, params.ImageParams) {
	flameFuncs := createFlameFuncs()
	ifsMap := createFlameMap(params, flameFuncs)

	comp := func(x int, ymin int, ymax int, img *image.RGBA64, wg *sync.WaitGroup) {
		for y := ymin; y < ymax; y++ {
			val, ok := ifsMap[coords{int64(x), int64(y)}]
			if ok {
				img.Set(x, y, val)
			} else {
				img.Set(x, y, params.Palette.Divergence)
			}
		}
		wg.Done()
	}

	return comp, params
}

type ifsFunc func(float64, float64) (float64, float64)

func createFlameFuncs() []ifsFunc {
	V0 := func(x float64, y float64) (float64, float64) {
		return x, y
	}

	V1 := func(x float64, y float64) (float64, float64) {
		return 3 * math.Sin(x), 3 * math.Sin(y)
	}

	V2 := func(x float64, y float64) (float64, float64) {
		r2 := x*x + y*y
		return 1 / r2 * x, 1 / r2 * y
	}

	V3 := func(x float64, y float64) (float64, float64) {
		r2 := x*x + y*y
		return x*math.Sin(r2) - y*math.Cos(r2), x*math.Cos(r2) + y*math.Sin(r^2)
	}

	return []ifsFunc{V0, V1, V2, V3}
}

type triplet struct {
	R float64
	G float64
	B float64
	A float64
}

func createFlameMap(params params.ImageParams, funcs []ifsFunc) map[coords]color.RGBA {
	res := make(map[coords]color.RGBA)
	x := float64(0)
	y := float64(0)

	resTmp := make(map[coords]triplet)
	for i := 0; i < 5000000; i++ {
		rule := rand.Float32()
		var a, b, c, d, e, f float64
		var funcIndex int
		if rule < 0.05 {
			a = 0.1
			b = 0.3
			c = 0.4
			d = 0.16
			e = 0.1
			f = 0.4
			funcIndex = 0
		} else if rule < 0.86 {
			a = 0.85
			b = 0.04
			c = -0.04
			d = 0.85
			e = -0.4
			f = 1.6
			funcIndex = 1
		} else if rule < 0.93 {
			a = -0.15
			b = 0.28
			c = 0.26
			d = 0.24
			e = 0.2
			f = 0.44
			funcIndex = 2
		} else {
			a = 0.20
			b = -0.26
			c = 0.23
			d = 0.22
			e = -0.1
			f = 1.6
			funcIndex = 3
		}
		x1 := a*x + b*y + e
		y1 := c*x + d*y + f
		x1, y1 = funcs[funcIndex](x1, y1)
		coords := scaleFlame(x1, y1, params)
		fr, fg, fb, _ := params.Palette.ListColors[funcIndex].RGBA()
		col, ok := resTmp[coords]
		if i > 20 {
			if ok {
				resTmp[coords] = triplet{(col.R + float64(fr)/256) / 2, (col.G + float64(fg)/256) / 2, (col.B + float64(fb)/256) / 2, col.A + 1}
			} else {
				resTmp[coords] = triplet{float64(fr) / 256, float64(fg) / 256, float64(fb) / 256, 1}
			}
		}
		col = resTmp[coords]
		x = x1
		y = y1
	}

	maxAlpha := 0.0
	for _, v := range resTmp {
		if v.A > maxAlpha {
			maxAlpha = v.A
		}
	}


	for k, v := range resTmp {
		scale := math.Log2(v.A) / v.A
		cr := uint8(v.R * scale * 256)
		cg := uint8(v.G * scale * 256)
		cb := uint8(v.B * scale * 256)

		alpha := uint8((math.Log2(v.A) / math.Log2(maxAlpha)) * 255)

		res[k] = color.RGBA{R: cr, G: cg, B: cb, A: alpha}
	}
	return res
}

func scaleFlame(x1 float64, y1 float64, imageParams params.ImageParams) coords {
	x := (x1 + 1.0) * float64(imageParams.Width)
	y := (y1 + 1.0) * float64(imageParams.Height)
	return coords{int64(x), int64(y)}
}
