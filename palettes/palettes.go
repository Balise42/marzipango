package palettes

import (
	"image/color"
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
}

// ColorFromContinuousPalette returns the color corresponding to the value from 0 to 1 (0 is the first color of the palette, 1 is the last color of the palette)
func ColorFromContinuousPalette(value float64, converge bool, palette Colors) color.Color {
	if !converge {
		return palette.Divergence
	}

	if len(palette.ListColors) == 1 {
		return palette.ListColors[0]
	}

	colorIndex := int(float64((len(palette.ListColors) - 1)) * value)
	normalizedColor := float64((len(palette.ListColors)-1))*value - float64(colorIndex)

	c1r, c1g, c1b, _ := palette.ListColors[colorIndex].RGBA()
	c2r, c2g, c2b, _ := palette.ListColors[colorIndex+1].RGBA()

	var r, g, b uint16
	if c1r > c2r {
		r = uint16(float64(c2r) + normalizedColor*float64(c1r-c2r))
	} else {
		r = uint16(float64(c1r) + normalizedColor*float64(c2r-c1r))
	}

	if c1g > c2g {
		g = uint16(float64(c2g) + normalizedColor*float64(c1g-c2g))
	} else {
		g = uint16(float64(c1g) + normalizedColor*float64(c2g-c1g))
	}

	if c1b > c2b {
		b = uint16(float64(c2b) + normalizedColor*float64(c1b-c2b))
	} else {
		b = uint16(float64(c1b) + normalizedColor*float64(c2b-c1b))
	}

	return color.RGBA64{r, g, b, 0xffff}
}
