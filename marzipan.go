package main

import (
	"fmt"
	"image"
	"image/png"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/Balise42/marzipango.git/fractales"
)

const left = -2.0
const right = 1.0
const width = 900
const height = 600
const top = 1.0
const bottom = -1.0
const maxiter = 100

type imageParams struct {
	left   float64
	right  float64
	top    float64
	bottom float64
	width  int
	height int
}

func scale(x int, y int, pos imageParams) complex128 {
	real := pos.left + float64(x)/float64(pos.width)*(pos.right-pos.left)
	im := pos.top + float64(y)/float64(pos.height)*(pos.bottom-pos.top)
	return complex(real, im)
}

func generateImage(w io.Writer, params imageParams) error {
	img := image.NewRGBA64(image.Rect(0, 0, params.width, params.height))
	for x := 0; x < params.width; x++ {
		for y := 0; y < params.height; y++ {
			img.Set(x, y, fractales.MandelbrotColor(scale(x, y, params), maxiter))
		}
	}

	return png.Encode(w, img)
}

func parseIntParam(r *http.Request, name string, fallback int) int {
	param, err := strconv.Atoi(r.URL.Query().Get(name))
	if err != nil {
		return fallback
	}
	return param
}

func parseFloatParam(r *http.Request, name string, fallback float64) float64 {
	param, err := strconv.ParseFloat(r.URL.Query().Get(name), 64)
	if err != nil {
		return fallback
	}
	return param
}

func mandelbrot(w http.ResponseWriter, r *http.Request) {
	imgWidth := parseIntParam(r, "width", width)
	imgHeight := parseIntParam(r, "height", height)
	imgTop := parseFloatParam(r, "top", top)
	imgLeft := parseFloatParam(r, "left", left)
	imgBottom := parseFloatParam(r, "bottom", bottom)
	imgRight := parseFloatParam(r, "right", right)

	imageParams := imageParams{imgLeft, imgRight, imgTop, imgBottom, imgWidth, imgHeight}

	w.Header().Set("Content-Type", "image/png")
	err := generateImage(w, imageParams)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Println("Image served", imageParams)
}

func main() {
	http.HandleFunc("/", mandelbrot)
	err := http.ListenAndServe("localhost:8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
