package fractales

import (
	"math"
	"math/big"
	"math/cmplx"

	"github.com/Balise42/marzipango/params"
)

const r = 1000

// MandelbrotContinuousValueLow returns the fractional number of iterations corresponding to a complex in the Mandelbrot set with low precision input
func MandelbrotContinuousValueLow(c complex128, maxiter int) (float64, bool) {
	z := 0 + 0i
	for i := 0; i < maxiter; i++ {
		z = z*z + c
		if absz := cmplx.Abs(z); absz > r {
			return (float64(i) + 1 - math.Log2(math.Log2(absz))), true
		}
	}
	return math.MaxInt64, false
}

// MandelbrotContinuousValueComputerLow returns a ValueComputation for the mandelbrot set with low precision input
func MandelbrotContinuousValueComputerLow(params params.ImageParams) ValueComputation {
	return func(x int, y int) (float64, bool) {
		return MandelbrotContinuousValueLow(scale(x, y, params), params.MaxIter)
	}
}

// MandelbrotContinuousValueHigh returns the number of iterations corresponding to a complex in the Mandelbrot set with high precision input
func MandelbrotContinuousValueHigh(c *LargeComplex, maxiter int) (float64, bool) {
	z := LargeComplex{big.NewFloat(0), big.NewFloat(0)}
	for i := 0; i < maxiter; i++ {
		z = z.Square().Add(c)
		if absz := z.Abs64(); absz > r {
			return (float64(i) + 1 - math.Log2(math.Log2(absz))), true
		}
	}
	return math.MaxInt64, false
}

// MandelbrotContinuousValueComputerHigh returns a ValueComputation for the mandelbrot set with high precision input
func MandelbrotContinuousValueComputerHigh(params params.ImageParams) ValueComputation {
	return func(x int, y int) (float64, bool) {
		z := scaleHigh(x, y, params)
		return MandelbrotContinuousValueHigh(&z, params.MaxIter)
	}
}

// MandelbrotOrbitValueLow returns the distance to the closest orbit hit by the computation of iterations corresponding to a complex in the Mandelbrot set in low precision
func MandelbrotOrbitValueLow(c complex128, maxiter int, orbits []Orbit) (float64, bool) {
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

// MandelbrotOrbitValueComputerLow returns a ValueComputation for the julia set with orbit trapping
func MandelbrotOrbitValueComputerLow(params params.ImageParams, orbits []Orbit) ValueComputation {
	return func(x int, y int) (float64, bool) {
		return MandelbrotOrbitValueLow(scale(x, y, params), params.MaxIter, orbits)
	}
}

// MultibrotContinuousValueLow returns the number of iterations corresponding to a complex in the Multibrot set (with d > 2)
func MultibrotContinuousValueLow(c complex128, maxiter int, power complex128) (float64, bool) {

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

// MultibrotContinuousValueComputerLow returns a ValueComputation for the Multibrot set
func MultibrotContinuousValueComputerLow(params params.ImageParams) ValueComputation {
	return func(x int, y int) (float64, bool) {
		return MultibrotContinuousValueLow(scale(x, y, params), params.MaxIter, complex(params.Power, 0.0))
	}
}
