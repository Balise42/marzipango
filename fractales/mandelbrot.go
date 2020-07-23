package fractales

import (
	"math"
	"math/cmplx"
)

const r = 1000

// MandelbrotValue returns a scaled value (between 0 and 1) corresponding to a complex in the Mandelbrot set
func MandelbrotValue(c complex128, maxiter int) (float64, bool) {
	z := 0 + 0i
	for i := 0; i < maxiter; i++ {
		z = z*z + c
		if absz := cmplx.Abs(z); absz > r {
			return (float64(i) + 1 - math.Log2(math.Log2(absz))) / float64(maxiter), true
		}
	}
	return math.MaxInt64 / float64(maxiter), false
}
