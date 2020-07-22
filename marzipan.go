package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"os"

	"github.com/Balise42/marzipango.git/fractales"
)

const left = -2.0
const right = 1.0
const width = 900
const height = 600
const top = 1.0
const bottom = -1.0

func scale(x int, y int) complex128 {
	real := left + float64(x)/width*(right-left)
	im := top + float64(y)/height*(bottom-top)
	return complex(real, im)
}

func main() {
	err := generateImage()
	if err != nil {
		fmt.Println("Error occured: ", err)
		os.Exit(-1)
	}
	fmt.Println("Image generated")
}

func generateImage() error {
	img := image.NewRGBA64(image.Rect(0, 0, width, height))
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			img.Set(x, y, computeColor(scale(x, y)))
		}
	}
	f, err := os.Create("mandelbrot.png")
	defer f.Close()
	if err != nil {
		return err
	}
	err = png.Encode(f, img)
	if err != nil {
		return err
	}
	return nil
}

func computeColor(z complex128) color.Color {
	return color.RGBA64{uint16((fractales.ComputeValue(z) / 100) * math.MaxUint16), 0, 0, math.MaxUint16}
}
