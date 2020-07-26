package fractales

import (
	"math"
	"math/cmplx"
)

const r = 1000

// MandelbrotValue returns the number of iterations corresponding to a complex in the Mandelbrot set
func MandelbrotValue(c complex128, maxiter int) (float64, bool) {
	z := 0 + 0i
	for i := 0; i < maxiter; i++ {
		z = z*z + c
		if absz := cmplx.Abs(z); absz > r {
			return (float64(i) + 1 - math.Log2(math.Log2(absz))), true
		}
	}
	return math.MaxInt64, false
}

func MandelbrotOrbitValue(c complex128, maxiter int, orbits []Orbit) (float64, bool) {
	dist := math.MaxFloat64

	var z complex128
	for i := 0; i < maxiter; i++ {
		z = z*z + c

		for _, orbit := range orbits {
			dist = math.Min(dist, orbit.getOrbitValue(orbit.getOrbitFastValue(z)))
		}
		if cmplx.Abs(z) > 4 {
			return dist, true
		}
	}

	return math.MaxFloat64, false
}
