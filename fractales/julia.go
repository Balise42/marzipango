package fractales

import (
	"math"
	"math/big"
	"math/cmplx"

	"github.com/Balise42/marzipango/params"
)

// JuliaContinuousValueLow returns the fractional number of iterations corresponding to a complex in the Julia set in low precision
func JuliaContinuousValueLow(z complex128, maxiter int) (float64, bool) {
	c := -0.4 + 0.6i
	for i := 0; i < maxiter; i++ {
		z = z*z + c
		if absz := cmplx.Abs(z); absz > r {
			return (float64(i) + 1 - math.Log2(math.Log2(absz))), true
		}
	}
	return math.MaxInt64, false
}

// JuliaContinuousValueComputerLow returns a ValueComputation for the julia set with low precision input
func JuliaContinuousValueComputerLow(params params.ImageParams) ValueComputation {
	return func(x int, y int) (float64, bool) {
		return JuliaContinuousValueLow(scale(x, y, params), params.MaxIter)
	}
}

// JuliaContinuousValueHigh returns the fractional number of iterations corresponding to a complex in the Julia set in high precision
func JuliaContinuousValueHigh(z LargeComplex, maxiter int) (float64, bool) {
	c := LargeComplex{big.NewFloat(-0.4), big.NewFloat(0.6)}
	for i := 0; i < maxiter; i++ {
		z = z.Square().Add(&c)
		if absz := z.Abs64(); absz > r {
			return (float64(i) + 1 - math.Log2(math.Log2(absz))), true
		}
	}
	return math.MaxInt64, false
}

// JuliaContinuousValueComputerHigh returns a ValueComputation for the julia set with high precision input
func JuliaContinuousValueComputerHigh(params params.ImageParams) ValueComputation {
	return func(x int, y int) (float64, bool) {
		return JuliaContinuousValueHigh(scaleHigh(x, y, params), params.MaxIter)
	}
}

// JuliaOrbitValueLow returns the distance to the closest orbit hit by the computation of iterations corresponding to a complex in the Julia set in low precision
func JuliaOrbitValueLow(z complex128, maxiter int, orbits []Orbit) (float64, bool) {
	dist := math.MaxFloat64

	c := -0.4 + 0.6i

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

// JuliaOrbitValueComputerLow returns a ValueComputation for the julia set with orbit trapping
func JuliaOrbitValueComputerLow(params params.ImageParams, orbits []Orbit) ValueComputation {
	return func(x int, y int) (float64, bool) {
		return JuliaOrbitValueLow(scale(x, y, params), params.MaxIter, orbits)
	}
}
