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

// ValueComputation is a value computation function
type ValueComputation func(x int, y int) (float64, bool)

func CreateComputer(computeValue ValueComputation, colorPixel palettes.ColoringFunction, params params.ImageParams) Computation {
	return func(x int, ymin int, ymax int, img *image.RGBA64, wg *sync.WaitGroup) {
		for y := ymin; y < ymax; y++ {
			value, converge := computeValue(x, y)
			colorPixel(img, x, y, value, converge)
		}
		wg.Done()
	}
}
