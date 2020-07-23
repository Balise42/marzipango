package palettes

import (
	"image/color"
)

// Colors used for the palette
type Colors struct {
	Divergence color.Color
	ListColors color.Palette
}

// ColorFromContinuousPalette returns the color corresponding to the value from 0 to 1 (0 is the first color of the palette, 1 is the last color of the palette)
func ColorFromContinuousPalette(value float64, converge bool, palette Colors) color.Color {
	if !converge {
		return palette.Divergence
	}

	colorIndex := int(float64((len(palette.ListColors) - 1)) * value)
	normalizedColor := float64((len(palette.ListColors)-1))*value - float64(colorIndex)

	c1r, c1g, c1b, c1a := palette.ListColors[colorIndex].RGBA()
	c2r, c2g, c2b, c2a := palette.ListColors[colorIndex+1].RGBA()

	return color.RGBA64{
		uint16(float64(c1r) + normalizedColor*float64(c2r-c1r)),
		uint16(float64(c1g) + normalizedColor*float64(c2g-c1g)),
		uint16(float64(c1b) + normalizedColor*float64(c2b-c1b)),
		uint16(float64(c1a) + normalizedColor*float64(c2a-c1a)),
	}
}
