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

type LineOrbit struct {
	A           float64
	B           float64
	C           float64
	Sqrtab      float64
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

func CreateLineOrbit(a float64, b float64, c float64, maxvalue float64) LineOrbit {
	orbit := LineOrbit{A: a, B: b, C: c}
	orbit.Sqrtab = math.Sqrt(a*a + b*b)

	minDist := 0.0
	maxDist := orbit.getOrbitFastValue(-2 - 1i)

	maxDist = math.Max(maxDist, orbit.getOrbitFastValue(-2+1i))
	maxDist = math.Max(maxDist, orbit.getOrbitFastValue(1+1i))
	maxDist = math.Max(maxDist, orbit.getOrbitFastValue(1-1i))

	minDist = math.Min(minDist, orbit.getOrbitFastValue(-2+1i))
	minDist = math.Min(minDist, orbit.getOrbitFastValue(1+1i))
	minDist = math.Min(minDist, orbit.getOrbitFastValue(1-1i))

	maxDist = math.Sqrt(maxDist) / orbit.Sqrtab
	minDist = math.Sqrt(minDist) / orbit.Sqrtab

	orbit.Factor = (maxvalue - minDist) / (maxDist - minDist)
	orbit.Translation = minDist

	return orbit
}

func (l LineOrbit) getOrbitFastValue(z complex128) float64 {
	lineCoeff := l.A*real(z) + l.B*imag(z) + l.C
	return lineCoeff * lineCoeff
}

func (l LineOrbit) getOrbitValue(v float64) float64 {
	return (math.Sqrt(v)/l.Sqrtab - l.Translation) * l.Factor
}
