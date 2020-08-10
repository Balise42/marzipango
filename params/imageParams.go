package params

import "github.com/Balise42/marzipango/palettes"

const Left = -2.0
const Right = 1.0
const Width = 900
const Height = 600
const Top = 1.0
const Bottom = -1.0
const Maxiter = 100

type ImageParams struct {
	Left    float64
	Right   float64
	Top     float64
	Bottom  float64
	Width   int
	Height  int
	MaxIter int
	Palette palettes.Colors
	Power   float64
}
