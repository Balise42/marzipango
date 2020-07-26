package fractales

import (
	"math"
	"math/cmplx"
)

// JuliaValue returns the number of iterations corresponding to a complex in the Julia set
func JuliaValue(z complex128, maxiter int) (float64, bool) {
	c := -0.4 + 0.6i
	for i := 0; i < maxiter; i++ {
		z = z*z + c
		if absz := cmplx.Abs(z); absz > r {
			return (float64(i) + 1 - math.Log2(math.Log2(absz))), true
		}
	}
	return math.MaxInt64, false
}
