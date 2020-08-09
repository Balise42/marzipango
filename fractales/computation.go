package fractales

import (
	"image"
	"math/big"
	"sync"

	"github.com/Balise42/marzipango/palettes"
	"github.com/Balise42/marzipango/params"
)

func scale(x int, y int, pos params.ImageParams) complex128 {
	real := pos.Left + float64(x)/float64(pos.Width)*(pos.Right-pos.Left)
	im := pos.Top + float64(y)/float64(pos.Height)*(pos.Bottom-pos.Top)

	return complex(real, im)
}

func scaleHigh(x int, y int, pos params.ImageParams) LargeComplex {
	ratioX := float64(x) / float64(pos.Width)

	real := big.NewFloat(0)
	real.Sub(big.NewFloat(pos.Right), big.NewFloat(pos.Left))
	real.Mul(big.NewFloat(float64(ratioX)), real)
	real.Add(real, big.NewFloat(pos.Left))

	ratioY := float64(y) / float64(pos.Height)

	imag := big.NewFloat(0)
	imag.Sub(big.NewFloat(pos.Bottom), big.NewFloat(pos.Top))
	imag.Mul(big.NewFloat(ratioY), imag)
	imag.Add(imag, big.NewFloat(pos.Top))

	return LargeComplex{real, imag}
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

// ComputeOrbitMandelbrotHighWithContinuousPalette provides the computation for an orbit-colored Mandelbrot with the provided image parameters and orbits
func ComputeMandelbrotHighWithContinuousPalette(params params.ImageParams) Computation {
	return func(x int, ymin int, ymax int, img *image.RGBA64, wg *sync.WaitGroup) {
		for y := ymin; y < ymax; y++ {
			scaled := scaleHigh(x, y, params)
			value, converge := MandelbrotValueHigh(&scaled, params.MaxIter)
			img.Set(x, y, palettes.ColorFromContinuousPalette(value, converge, params.Palette))
		}
		wg.Done()
	}
}
