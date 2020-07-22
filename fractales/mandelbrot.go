package fractales

import (
	"math"
	"math/cmplx"
)

const maxiter = 100
const R = 1000

func ComputeValue(c complex128) float64 {
	z := 0 + 0i
	for i := 0; i < maxiter; i++ {
		z = z*z + c
		absz := cmplx.Abs(z)
		if absz > R {
			return (float64(i) + 1 - math.Log2(math.Log2(absz)))
		}
	}
	return math.MaxInt64
}
