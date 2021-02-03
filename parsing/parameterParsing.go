package parsing

import (
	"image/color"
	"net/http"
	"regexp"
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

	return palettes.Colors{Divergence: divergence, ListColors: listColors, MaxValue: 254}
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

func parseOrbit(rawOrbit string, defaultOrbit fractales.Orbit, imageParams params.ImageParams) fractales.Orbit {
	if strings.HasPrefix(rawOrbit, "point(") {
		paramString := strings.TrimSuffix(strings.TrimPrefix(rawOrbit, "point("), ")")
		params := strings.Split(paramString, ",")
		if len(params) != 3 {
			return defaultOrbit
		}

		x, err := strconv.ParseFloat(params[0], 64)
		if err != nil {
			return defaultOrbit
		}
		y, err := strconv.ParseFloat(params[1], 64)
		if err != nil {
			return defaultOrbit
		}
		dist, err := strconv.ParseFloat(params[2], 64)
		if err != nil {
			return defaultOrbit
		}
		return fractales.CreatePointOrbit(x, y, dist)
	} else if strings.HasPrefix(rawOrbit, "line(") {
		paramString := strings.TrimSuffix(strings.TrimPrefix(rawOrbit, "line("), ")")
		params := strings.Split(paramString, ",")
		if len(params) != 4 {
			return defaultOrbit
		}
		a, err := strconv.ParseFloat(params[0], 64)
		if err != nil {
			return defaultOrbit
		}
		b, err := strconv.ParseFloat(params[1], 64)
		if err != nil {
			return defaultOrbit
		}
		c, err := strconv.ParseFloat(params[2], 64)
		if err != nil {
			return defaultOrbit
		}
		dist, err := strconv.ParseFloat(params[3], 64)
		if err != nil {
			return defaultOrbit
		}
		return fractales.CreateLineOrbit(a, b, c, dist)
	} else if strings.HasPrefix(rawOrbit, "raster(") {
		paramString := strings.TrimSuffix(strings.TrimPrefix(rawOrbit, "raster("), ")")
		params := strings.Split(paramString, ",")

		if len(params) > 2 {
			return defaultOrbit
		}

		matched, err := regexp.MatchString(`^[a-zA-Z0-9\-]+$`, params[0])

		if err != nil || !matched {
			return defaultOrbit
		}

		path := "fractales/orbits/" + params[0] + ".png"

		dist := float64(100)

		if len(params) == 2 {
			dist, err = strconv.ParseFloat(params[1], 64)
			if err != nil {
				return defaultOrbit
			}
		}

		orbit, err := fractales.CreateImageOrbit(imageParams, path, dist)
		if err != nil {
			return defaultOrbit
		}

		return orbit
	}
	return defaultOrbit
}

func parseOrbits(r *http.Request, params params.ImageParams) ([]fractales.Orbit, bool) {
	rawOrbits, ok := r.URL.Query()["orbit"]
	defaultOrbit := fractales.CreatePointOrbit(0.5, -0.7, float64(100))

	if !ok {
		return []fractales.Orbit{defaultOrbit}, false
	}
	orbits := make([]fractales.Orbit, len(rawOrbits))

	for i, orbit := range rawOrbits {
		orbits[i] = parseOrbit(orbit, defaultOrbit, params)
	}
	return orbits, true
}

func parseFractaleType(r *http.Request, defaultType string) string {
	fractaleType := r.URL.Query().Get("type")
	if fractaleType == "mandelbrot" || fractaleType == "julia" || fractaleType == "fern" || fractaleType == "flame" || fractaleType == "sierp"{
		return fractaleType
	}
	return defaultType
}

func highPrecision(params params.ImageParams) bool {
	coords := make(map[float64]bool)
	for x := 0; x < params.Width; x++ {
		coords[params.Left+float64(x)/float64(params.Width)*(params.Right-params.Left)] = true
	}
	return len(coords) < params.Width-5
}

// ParseComputation parses the request parameters to dispatch the computation and the parameters
func ParseComputation(r *http.Request) (fractales.Computation, params.ImageParams) {
	imgWidth, imgHeight := parseImageSize(r)
	imgLeft, imgRight, imgTop, imgBottom := parseImageCoords(r)
	imgMaxIter := parseIntParam(r, "maxiter", params.Maxiter)

	listCols := color.Palette{palettes.White, palettes.Black, palettes.White}
	palette := palettes.Colors{Divergence: color.Black, ListColors: listCols, MaxValue: 500}
	imgPalette := parsePalette(r, "palette", palette)
	imgPalette.MaxValue = parseIntParam(r, "palettesize", 100)

	power := parseFloatParam(r, "power", 2)

	imageParams := params.ImageParams{Left: imgLeft, Right: imgRight, Top: imgTop, Bottom: imgBottom, Width: imgWidth, Height: imgHeight, MaxIter: imgMaxIter, Palette: imgPalette, Power: power}

	fractaleType := parseFractaleType(r, "mandelbrot")

	orbits, hasOrbits := parseOrbits(r, imageParams)

	var valueComputer fractales.ValueComputation
	var colorPixel = palettes.ContinuousColoring(imgPalette)

	if !highPrecision(imageParams) {
		if fractaleType == "julia" {
			if hasOrbits {
				valueComputer = fractales.JuliaOrbitValueComputerLow(imageParams, orbits)
			} else {
				valueComputer = fractales.JuliaContinuousValueComputerLow(imageParams)
			}
		} else if fractaleType == "mandelbrot" && power == 2 {
			if hasOrbits {
				valueComputer = fractales.MandelbrotOrbitValueComputerLow(imageParams, orbits)
			} else {
				valueComputer = fractales.MandelbrotContinuousValueComputerLow(imageParams)
			}
		} else if fractaleType == "fern" {
			valueComputer = fractales.FernValueComputeLow(imageParams)
		} else if fractaleType == "flame" {
			return fractales.CreateFlameComputer(imageParams)
		} else if fractaleType == "sierp" {
			valueComputer = fractales.SierpValueComputeLow(imageParams)
		} else if power != 2 {
			valueComputer = fractales.MultibrotContinuousValueComputerLow(imageParams)
		}
	} else {
		if fractaleType == "julia" {
			valueComputer = fractales.JuliaContinuousValueComputerHigh(imageParams)
		} else if fractaleType == "mandelbrot" && power == 2 {
			valueComputer = fractales.MandelbrotContinuousValueComputerHigh(imageParams)
		}
	}

	return fractales.CreateComputer(valueComputer, colorPixel, imageParams), imageParams
}
