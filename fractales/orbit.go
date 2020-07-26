package fractales

import "math"

type Orbit interface {
	getOrbitFastValue(z complex128) float64
	getOrbitValue(v float64) float64
}
type PointOrbit struct {
	X           float64
	Y           float64
	Translation float64
	Factor      float64
}

func CreatePointOrbit(x float64, y float64, maxvalue float64) PointOrbit {
	orbit := PointOrbit{X: x, Y: y}
	minDist := 0.0
	maxDist := orbit.squaredDistance(-2 - 1i)

	maxDist = math.Max(maxDist, orbit.squaredDistance(-2+1i))
	maxDist = math.Max(maxDist, orbit.squaredDistance(1+1i))
	maxDist = math.Max(maxDist, orbit.squaredDistance(1-1i))

	minDist = math.Min(minDist, orbit.squaredDistance(-2+1i))
	minDist = math.Min(minDist, orbit.squaredDistance(1+1i))
	minDist = math.Min(minDist, orbit.squaredDistance(1-1i))

	maxDist = math.Sqrt(maxDist)
	minDist = math.Sqrt(minDist)

	orbit.Factor = (maxvalue - minDist) / (maxDist - minDist)
	orbit.Translation = minDist

	return orbit
}

func (p PointOrbit) getOrbitFastValue(z complex128) float64 {
	return p.squaredDistance(z)
}

func (p PointOrbit) getOrbitValue(v float64) float64 {
	return (math.Sqrt(v) - p.Translation) * p.Factor
}

func (p PointOrbit) squaredDistance(z complex128) float64 {
	return (real(z)-p.X)*(real(z)-p.X) + (imag(z)-p.Y)*(imag(z)-p.Y)
}
