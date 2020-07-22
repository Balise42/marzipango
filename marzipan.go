package main

import (
	"fmt"
	"image"
	"image/png"
	"io"
	"log"
	"net/http"

	"github.com/Balise42/marzipango.git/fractales"
)

const left = -2.0
const right = 1.0
const width = 1800
const height = 1200
const top = 1.0
const bottom = -1.0
const maxiter = 100

func scale(x int, y int) complex128 {
	real := left + float64(x)/width*(right-left)
	im := top + float64(y)/height*(bottom-top)
	return complex(real, im)
}

func generateImage(w io.Writer) error {
	img := image.NewRGBA64(image.Rect(0, 0, width, height))
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			img.Set(x, y, fractales.MandelbrotColor(scale(x, y), maxiter))
		}
	}

	return png.Encode(w, img)
}

func mandelbrot(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "image/png")
	err := generateImage(w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Println("Image served")
}

func main() {
	http.HandleFunc("/", mandelbrot)
	err := http.ListenAndServe("localhost:8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
