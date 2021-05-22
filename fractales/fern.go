package fractales

import (
	"github.com/Balise42/marzipango/fractales/orbits"
	"github.com/Balise42/marzipango/params"
	"math/rand"
)



func FernValueComputeLow(params params.ImageParams) ValueComputation {
	ifsMap := createFernMap(params)
	
	return func(x int, y int) (float64, bool) {
		val, ok := ifsMap[orbits.Coords{int64(x), int64(y)}]
		if ok {
			return float64(val), true
		}
		return float64(0), false
	}
}

func createFernMap(params params.ImageParams) map[orbits.Coords]int {
	res := make(map[orbits.Coords]int)
	x := float64(0)
	y := float64(0)
	res[scaleFern(x, y, params)] = 1

	for i := 0; i < 100000000; i++ {
		rule := rand.Float32()
		var a, b, c, d, e, f float64
		if rule < 0.05 {
			a = 0
			b = 0
			c = 0
			d = 0.16
			e = 0
			f = 0
		} else if rule < 0.86 {
			a = 0.85
			b = 0.04
			c = -0.04
			d = 0.85
			e = 0
			f = 1.6
		} else if rule < 0.93 {
			a = -0.15
			b = 0.28
			c = 0.26
			d = 0.24
			e = 0
			f = 0.44
		} else {
			a = 0.20
			b = -0.26
			c = 0.23
			d = 0.22
			e = 0
			f = 1.6
		}
		x1 := a * x + b * y + e
		y1 := c * x + d * y + f
		coords := scaleFern(x1, y1, params)
		_, ok := res[coords]
		if !ok {
			res[coords] = 1
		} else {
			res[coords] = res[coords] + 1
		}
		x = x1
		y = y1
	}
	return res
}

func scaleFern(x1 float64, y1 float64, imageParams params.ImageParams) orbits.Coords {
	x := (x1 + 2.1820) * float64(imageParams.Width) / (2.1820 + 2.6558)
	y := (9.9983 - y1) * float64(imageParams.Height) / 9.9983
	return orbits.Coords{int64(x), int64(y)}
}
