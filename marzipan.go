package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io"
	"log"
	"net/http"
	"runtime"
	"strconv"
	"strings"
	"sync"

	"github.com/Balise42/marzipango.git/fractales"
	"github.com/Balise42/marzipango.git/palettes"
)

const left = -2.0
const right = 1.0
const width = 900
const height = 600
const top = 1.0
const bottom = -1.0
const maxiter = 100

// Computation fills in image pixels according to parameters
type Computation func(x int, ymin int, ymax int, img *image.RGBA64, wg *sync.WaitGroup)

type imageParams struct {
	left    float64
	right   float64
	top     float64
	bottom  float64
	width   int
	height  int
	maxIter int
	palette palettes.Colors
}

func scale(x int, y int, pos imageParams) complex128 {
	real := pos.left + float64(x)/float64(pos.width)*(pos.right-pos.left)
	im := pos.top + float64(y)/float64(pos.height)*(pos.bottom-pos.top)

	return complex(real, im)
}

func generateImage(w io.Writer, params imageParams, comp Computation) error {
	var wg sync.WaitGroup
	img := image.NewRGBA64(image.Rect(0, 0, params.width, params.height))
	for x := 0; x < params.width; x++ {
		var numRows = params.height / runtime.NumCPU()
		for cpu := 0; cpu < runtime.NumCPU()-1; cpu++ {
			wg.Add(1)
			go comp(x, numRows*cpu, numRows*(cpu+1), img, &wg)
		}
		wg.Add(1)
		go comp(x, numRows*(runtime.NumCPU()-1), params.height, img, &wg)
	}
	wg.Wait()

	return png.Encode(w, img)
}

func computeMandelbrotWithContinuousPalette(params imageParams) Computation {
	return func(x int, ymin int, ymax int, img *image.RGBA64, wg *sync.WaitGroup) {
		for y := ymin; y < ymax; y++ {
			value, converge := fractales.MandelbrotValue(scale(x, y, params), params.maxIter)
			img.Set(x, y, palettes.ColorFromContinuousPalette(value, converge, params.palette))
		}
		wg.Done()
	}
}

func computeJuliaWithContinuousPalette(params imageParams) Computation {
	return func(x int, ymin int, ymax int, img *image.RGBA64, wg *sync.WaitGroup) {
		for y := ymin; y < ymax; y++ {
			value, converge := fractales.JuliaValue(scale(x, y, params), params.maxIter)
			img.Set(x, y, palettes.ColorFromContinuousPalette(value, converge, params.palette))
		}
		wg.Done()
	}
}

func computeColumn(x int, ymin int, ymax int, params imageParams, img *image.RGBA64, wg *sync.WaitGroup) {
	for y := ymin; y < ymax; y++ {
		value, converge := fractales.MandelbrotValue(scale(x, y, params), params.maxIter)
		img.Set(x, y, palettes.ColorFromContinuousPalette(value, converge, params.palette))
	}
	wg.Done()
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

func parseColor(name string) color.Color {
	color, ok := palettes.ColorNames[name]
	if ok {
		return color
	}
	return palettes.Black
}

func parsePalette(r *http.Request, name string, fallback palettes.Colors) palettes.Colors {
	param := r.URL.Query().Get(name)
	if len(param) < 1 {
		return fallback
	}

	paramList := strings.Split(param, ",")
	divergence := parseColor(paramList[0])
	listColors := make([]color.Color, len(paramList)-1)
	for i := 1; i < len(paramList); i++ {
		listColors[i-1] = parseColor(paramList[i])
	}

	return palettes.Colors{Divergence: divergence, ListColors: listColors, MaxValue: 100}
}

func parseImageSize(r *http.Request) (int, int) {
	if r.URL.Query().Get("size") != "" {
		param, err := strconv.Atoi(r.URL.Query().Get("size"))
		if err == nil {
			return param, param
		}
	}
	return parseIntParam(r, "width", width), parseIntParam(r, "height", height)
}

func parseImageCoords(r *http.Request) (float64, float64, float64, float64) {
	if r.URL.Query().Get("x") != "" && r.URL.Query().Get("y") != "" && r.URL.Query().Get("window") != "" {
		x := parseFloatParam(r, "x", 0)
		y := parseFloatParam(r, "y", 0)
		space := parseFloatParam(r, "window", 1)
		return x - space, x + space, y - space, y + space
	}
	return parseFloatParam(r, "left", left), parseFloatParam(r, "right", right), parseFloatParam(r, "top", top), parseFloatParam(r, "bottom", bottom)
}

func parseComputation(r *http.Request) (Computation, imageParams) {
	imgWidth, imgHeight := parseImageSize(r)
	imgLeft, imgRight, imgTop, imgBottom := parseImageCoords(r)
	imgMaxIter := parseIntParam(r, "maxiter", maxiter)

	listCols := color.Palette{palettes.White, palettes.Black, palettes.White}
	palette := palettes.Colors{Divergence: color.Black, ListColors: listCols, MaxValue: 100}
	imgPalette := parsePalette(r, "palette", palette)

	imageParams := imageParams{imgLeft, imgRight, imgTop, imgBottom, imgWidth, imgHeight, imgMaxIter, imgPalette}

	if r.URL.Query().Get("type") == "julia" {
		return computeJuliaWithContinuousPalette(imageParams), imageParams
	} else {
		return computeMandelbrotWithContinuousPalette(imageParams), imageParams
	}
}

func fractale(w http.ResponseWriter, r *http.Request) {
	comp, imageParams := parseComputation(r)

	w.Header().Set("Content-Type", "image/png")
	err := generateImage(w, imageParams, comp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Println("Image served", imageParams)
}

func main() {
	http.HandleFunc("/", fractale)
	err := http.ListenAndServe("localhost:8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
