package fractales

import (
	"github.com/Balise42/marzipango/params"
	"math/rand"
)

func SierpValueComputeLow(params params.ImageParams) ValueComputation {
	sierpFuncs := createSierpFuncs()
	ifsMap := createSierpMap(params, sierpFuncs)

	return func(x int, y int) (float64, bool) {
		val, ok := ifsMap[coords{int64(x),int64(y)}]
		if !ok {
			return 0, false
		} else {
			return val, true
		}
	}
}

func createSierpMap(params params.ImageParams, funcs []ifsFunc) map[coords]float64 {
	res := make(map[coords]float64)
	x := float64(0)
	y := float64(0)

	for i := 0; i < 50000000; i++ {
		rule := rand.Intn(3)
		x1, y1 := funcs[rule](x, y)
		coords := scaleFlame(x1, y1, params)
		res[coords] = res[coords] + 1
		x, y = x1, y1
	}
	return res
}

func createSierpFuncs() []ifsFunc {
	F0 := func(x float64, y float64) (float64, float64) {
		return x/2, y/2
	}

	F1 := func(x float64, y float64) (float64, float64) {
		return (x+1)/2, y/2
	}

	F2 := func(x float64, y float64) (float64, float64) {
		return x/2, (y+1) / 2
	}

	return []ifsFunc{F0, F1, F2}
}

