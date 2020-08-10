package fractales

import (
	"math"
	"math/big"
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

// MandelbrotValue returns the number of iterations corresponding to a complex in the Mandelbrot set
func MultibrotValue(c complex128, maxiter int, power complex128) (float64, bool) {

	B := math.Pow(2, 1/(real(power)-1))

	z := 0 + 0i
	for i := 0; i < maxiter; i++ {
		z = cmplx.Pow(z, power) + c
		if absz := cmplx.Abs(z); absz > r {
			return (float64(i) + 1 - (math.Log(math.Log(absz)/math.Log(B)) / math.Log2(real(power)))), true
		}
	}
	return math.MaxInt64, false
}

// MandelbrotValue returns the number of iterations corresponding to a complex in the Mandelbrot set
func MandelbrotValueHigh(c *LargeComplex, maxiter int) (float64, bool) {
	z := LargeComplex{big.NewFloat(0), big.NewFloat(0)}
	for i := 0; i < maxiter; i++ {
		z = z.Square().Add(c)
		if absz := z.Abs64(); absz > r {
			return (float64(i) + 1 - math.Log2(math.Log2(absz))), true
		}
	}
	return math.MaxInt64, false
}

func MandelbrotOrbitValue(c complex128, maxiter int, orbits []Orbit) (float64, bool) {
	dist := math.MaxFloat64

	var z complex128
	i := 0
	for i < maxiter && cmplx.Abs(z) < 4 {
		z = z*z + c
		for _, orbit := range orbits {
			dist = math.Min(dist, orbit.getOrbitValue(orbit.getOrbitFastValue(z)))
		}
		i++
	}

	if i == maxiter {
		return math.MaxFloat64, false
	}

	return dist, true
}
