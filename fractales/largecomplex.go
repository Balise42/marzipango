package fractales

import (
	"math"
	"math/big"
)

type LargeComplex struct {
	real *big.Float
	imag *big.Float
}

func (z LargeComplex) Square() LargeComplex {
	realSquare := big.NewFloat(0)
	realSquare.Mul(z.real, z.real)
	imagSquare := big.NewFloat(0)
	imagSquare.Mul(z.imag, z.imag)
	doubleProd := big.NewFloat(2)
	doubleProd.Mul(doubleProd, z.real)
	doubleProd.Mul(doubleProd, z.imag)

	return LargeComplex{realSquare.Sub(realSquare, imagSquare), doubleProd}
}

func (z LargeComplex) Add(c *LargeComplex) LargeComplex {
	newReal := big.NewFloat(0)
	newImag := big.NewFloat(0)

	return LargeComplex{newReal.Add(z.real, c.real), newImag.Add(z.imag, c.imag)}
}

func (z LargeComplex) Abs64() float64 {
	realSquare := big.NewFloat(0)
	realSquare.Mul(z.real, z.real)
	imagSquare := big.NewFloat(0)
	imagSquare.Mul(z.imag, z.imag)
	abs := big.NewFloat(0)
	abs.Add(realSquare, imagSquare)

	conv, _ := abs.Float64()
	return math.Sqrt(conv)
}
