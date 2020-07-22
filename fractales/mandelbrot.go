package fractales

import (
	"image/color"
	"math"
	"math/cmplx"
)

const r = 1000

func value(c complex128, maxiter int) (float64, bool) {
	z := 0 + 0i
	for i := 0; i < maxiter; i++ {
		z = z*z + c
		if absz := cmplx.Abs(z); absz > r {
			return (float64(i) + 1 - math.Log2(math.Log2(absz))), true
		}
	}
	return math.MaxInt64, false
}

// MandelbrotColor returns a color corresponding to a complex in the Mandelbrot set
func MandelbrotColor(z complex128, maxiter int) color.Color {
	value, converge := value(z, maxiter)
	if !converge {
		return color.Black
	}

	scaledValue := value / float64(maxiter)

	return color.RGBA64{math.MaxUint16, math.MaxUint16 - uint16(scaledValue*math.MaxUint16), math.MaxUint16 - uint16(scaledValue*math.MaxUint16), math.MaxUint16}
}
