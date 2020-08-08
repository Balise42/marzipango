package fractales

import (
	"math/big"
	"testing"
)

func TestAbs(t *testing.T) {
	z := LargeComplex{big.NewFloat(0.3), big.NewFloat(0.4)}
	if z.Abs64() != 0.5 {
		t.Errorf("Abs(0.3 + 0.4i should be 0.5, got %f", z.Abs64())
	}
}

func TestSquare(t *testing.T) {
	z := LargeComplex{big.NewFloat(1), big.NewFloat(-2)}
	z = z.Square()
	if z.imag.Cmp(big.NewFloat(-4)) != 0 {
		t.Errorf("Square imag is dubious, wanted 0.24, got %f", z.imag)
	}
	if z.real.Cmp(big.NewFloat(-3)) != 0 {
		t.Errorf("Square real is dubious, wanted %e, got %e", big.NewFloat(-0.11), z.real)
	}
}

func TestSum(t *testing.T) {
	z := LargeComplex{big.NewFloat(0.3), big.NewFloat(0.4)}
	c := LargeComplex{big.NewFloat(1.0), big.NewFloat(2.0)}
	z = z.Add(&c)
	if z.real.Cmp(big.NewFloat(1.3)) != 0 || z.imag.Cmp(big.NewFloat(2.4)) != 0 {
		t.Errorf("Sum is dubious, wanted 1.3 + 2.4i, got %f + %fi", z.real, z.imag)
	}
}
