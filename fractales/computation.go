package fractales

import (
	"image"
	"sync"

	"github.com/Balise42/marzipango/palettes"
	"github.com/Balise42/marzipango/params"
)

func scale(x int, y int, pos params.ImageParams) complex128 {
	real := pos.Left + float64(x)/float64(pos.Width)*(pos.Right-pos.Left)
	im := pos.Top + float64(y)/float64(pos.Height)*(pos.Bottom-pos.Top)

	return complex(real, im)
}

// Computation fills in image pixels according to parameters
type Computation func(x int, ymin int, ymax int, img *image.RGBA64, wg *sync.WaitGroup)

// ComputeMandelbrotWithContinuousPalette provides the computation for a continuous-colored Mandelbrot with the provided image parameters
func ComputeMandelbrotWithContinuousPalette(params params.ImageParams) Computation {
	return func(x int, ymin int, ymax int, img *image.RGBA64, wg *sync.WaitGroup) {
		for y := ymin; y < ymax; y++ {
			value, converge := MandelbrotValue(scale(x, y, params), params.MaxIter)
			img.Set(x, y, palettes.ColorFromContinuousPalette(value, converge, params.Palette))
		}
		wg.Done()
	}
}

// ComputeJuliaWithContinuousPalette provides the computation for a continuous-colored Julia with the provided image parameters
func ComputeJuliaWithContinuousPalette(params params.ImageParams) Computation {
	return func(x int, ymin int, ymax int, img *image.RGBA64, wg *sync.WaitGroup) {
		for y := ymin; y < ymax; y++ {
			value, converge := JuliaValue(scale(x, y, params), params.MaxIter)
			img.Set(x, y, palettes.ColorFromContinuousPalette(value, converge, params.Palette))
		}
		wg.Done()
	}
}

// ComputeOrbitMandelbrotWithContinuousPalette provides the computation for an orbit-colored Mandelbrot with the provided image parameters and orbits
func ComputeOrbitMandelbrotWithContinuousPalette(params params.ImageParams, orbits []Orbit) Computation {
	return func(x int, ymin int, ymax int, img *image.RGBA64, wg *sync.WaitGroup) {
		for y := ymin; y < ymax; y++ {
			value, converge := MandelbrotOrbitValue(scale(x, y, params), params.MaxIter, orbits)
			img.Set(x, y, palettes.ColorFromContinuousPalette(value, converge, params.Palette))
		}
		wg.Done()
	}
}
