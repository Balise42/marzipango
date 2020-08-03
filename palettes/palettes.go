package palettes

import (
	"image/color"
	"math"
)

var Blue = color.RGBA64{0, 0, 0xffff, 0xffff}
var Red = color.RGBA64{0xffff, 0, 0, 0xffff}
var Green = color.RGBA64{0, 0xffff, 0, 0xffff}
var Yellow = color.RGBA64{0xffff, 0xffff, 0, 0xffff}
var Magenta = color.RGBA64{0xffff, 0, 0xffff, 0xffff}
var Cyan = color.RGBA64{0, 0xffff, 0xffff, 0xffff}
var Black = color.RGBA64{0, 0, 0, 0xffff}
var White = color.RGBA64{0xffff, 0xffff, 0xffff, 0xffff}

var ColorNames = map[string]color.Color{
	"blue":    Blue,
	"red":     Red,
	"green":   Green,
	"yellow":  Yellow,
	"magenta": Magenta,
	"cyan":    Cyan,
	"black":   Black,
	"white":   White}

// Colors used for the palette
type Colors struct {
	Divergence color.Color
	ListColors []color.Color
	MaxValue   int
}

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
