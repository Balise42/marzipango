package fractales

import (
	"image"
	"math/big"
	"sync"

	"github.com/Balise42/marzipango/palettes"
	"github.com/Balise42/marzipango/params"
)

func scale(x int, y int, pos params.ImageParams) complex128 {
	re := pos.Left + float64(x)/float64(pos.Width)*(pos.Right-pos.Left)
	im := pos.Top + float64(y)/float64(pos.Height)*(pos.Bottom-pos.Top)

	return complex(re, im)
}

func scaleHigh(x int, y int, pos params.ImageParams) LargeComplex {
	ratioX := float64(x) / float64(pos.Width)

	re := big.NewFloat(0)
	re.Sub(big.NewFloat(pos.Right), big.NewFloat(pos.Left))
	re.Mul(big.NewFloat(float64(ratioX)), re)
	re.Add(re, big.NewFloat(pos.Left))

	ratioY := float64(y) / float64(pos.Height)

	im := big.NewFloat(0)
	im.Sub(big.NewFloat(pos.Bottom), big.NewFloat(pos.Top))
	im.Mul(big.NewFloat(ratioY), im)
	im.Add(im, big.NewFloat(pos.Top))

	return LargeComplex{re, im}
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
