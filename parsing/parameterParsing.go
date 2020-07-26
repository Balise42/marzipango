package parsing

import (
	"image/color"
	"net/http"
	"strconv"
	"strings"

	"github.com/Balise42/marzipango/fractales"
	"github.com/Balise42/marzipango/palettes"
	"github.com/Balise42/marzipango/params"
)

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
	return parseIntParam(r, "width", params.Width), parseIntParam(r, "height", params.Height)
}

func parseImageCoords(r *http.Request) (float64, float64, float64, float64) {
	if r.URL.Query().Get("x") != "" && r.URL.Query().Get("y") != "" && r.URL.Query().Get("window") != "" {
		x := parseFloatParam(r, "x", 0)
		y := parseFloatParam(r, "y", 0)
		space := parseFloatParam(r, "window", 1)
		return x - space, x + space, y - space, y + space
	}
	return parseFloatParam(r, "left", params.Left), parseFloatParam(r, "right", params.Right), parseFloatParam(r, "top", params.Top), parseFloatParam(r, "bottom", params.Bottom)
}

// ParseComputation parses the request parameters to dispatch the computation and the parameters
func ParseComputation(r *http.Request) (fractales.Computation, params.ImageParams) {
	imgWidth, imgHeight := parseImageSize(r)
	imgLeft, imgRight, imgTop, imgBottom := parseImageCoords(r)
	imgMaxIter := parseIntParam(r, "maxiter", params.Maxiter)

	listCols := color.Palette{palettes.White, palettes.Black, palettes.White}
	palette := palettes.Colors{Divergence: color.Black, ListColors: listCols, MaxValue: 100}
	imgPalette := parsePalette(r, "palette", palette)

	imageParams := params.ImageParams{imgLeft, imgRight, imgTop, imgBottom, imgWidth, imgHeight, imgMaxIter, imgPalette}

	if r.URL.Query().Get("type") == "julia" {
		return fractales.ComputeJuliaWithContinuousPalette(imageParams), imageParams
	} else {
		return fractales.ComputeMandelbrotWithContinuousPalette(imageParams), imageParams
	}
}