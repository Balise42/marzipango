package palettes

import (
	"image"
	"image/color"
	"math"
)

var Blue = color.RGBA{0, 0, 255, 255}
var Red = color.RGBA{255, 0, 0, 255}
var Green = color.RGBA{0, 255, 0, 255}
var Yellow = color.RGBA{255, 255, 0, 255}
var Magenta = color.RGBA{255, 0, 255, 255}
var Cyan = color.RGBA{0, 255, 255, 255}
var Black = color.RGBA{0, 0, 0, 255}
var White = color.RGBA{255, 255, 255, 255}
var DarkGreen = color.RGBA{0, 100, 0, 255}
var Champagne = color.RGBA{247, 231, 206, 255}
var DarkChampagne = color.RGBA{41, 25, 0, 255}
var Orange = color.RGBA{255, 127, 0, 255}
var Violet = color.RGBA{139, 0, 255, 255}
var LightPink = color.RGBA{255, 182, 193, 255}
var LightGreen = color.RGBA{172, 225, 175, 255}
var Purple = color.RGBA{148, 0, 211, 255}
var Indigo = color.RGBA{75, 0, 130, 255}
var Teal = color.RGBA{0, 128, 128, 255}
var DarkBlue = color.RGBA{0, 0, 128, 255}
var SoftPink = color.RGBA{255, 221, 244, 255}

var ColorNames = map[string]color.Color{
	"blue":          Blue,
	"red":           Red,
	"green":         Green,
	"yellow":        Yellow,
	"magenta":       Magenta,
	"cyan":          Cyan,
	"black":         Black,
	"white":         White,
	"darkgreen":     DarkGreen,
	"champagne":     Champagne,
	"darkchampagne": DarkChampagne,
	"orange":        Orange,
	"violet":        Violet,
	"lightpink":     LightPink,
	"lightgreen":    LightGreen,
	"purple":        Purple,
	"indigo":        Indigo,
	"teal":          Teal,
	"darkblue":      DarkBlue,
	"softpink":      SoftPink,
}

// Colors used for the palette
type Colors struct {
	Divergence color.Color
	ListColors []color.Color
	MaxValue   int
}

type ColoringFunction func(img *image.RGBA64, x int, y int, value float64, converge bool)

// ColorFromContinuousPalette returns the color corresponding to the value from 0 to 1 (0 is the first color of the palette, 1 is the last color of the palette)
func ColorFromContinuousPalette(rawValue float64, converge bool, palette Colors) color.Color {
	if !converge {
		return palette.Divergence
	}

	if len(palette.ListColors) == 1 {
		return palette.ListColors[0]
	}

	value := math.Mod(rawValue, float64(palette.MaxValue))

	normalized := value / float64(palette.MaxValue)
	normalized = (math.Pow(normalized-0.5, 3) + 0.125) / 0.250

	colorIndex := (int(float64((len(palette.ListColors) - 1)) * normalized)) % (len(palette.ListColors) - 1)
	normalizedColor := float64((len(palette.ListColors)-1))*normalized - float64(colorIndex)

	c1r, c1g, c1b, _ := palette.ListColors[colorIndex].RGBA()
	c2r, c2g, c2b, _ := palette.ListColors[colorIndex+1].RGBA()

	var r, g, b int16
	r = int16(float64(c2r)*normalizedColor + float64(c1r)*(1-normalizedColor))
	g = int16(float64(c2g)*normalizedColor + float64(c1g)*(1-normalizedColor))
	b = int16(float64(c2b)*normalizedColor + float64(c1b)*(1-normalizedColor))

	return color.RGBA64{uint16(r), uint16(g), uint16(b), 0xffff}
}

func ContinuousColoring(palette Colors) ColoringFunction {
	return func(img *image.RGBA64, x int, y int, value float64, converge bool) {
		img.Set(x, y, ColorFromContinuousPalette(value, converge, palette))
	}
}
