package fractales

import (
	"math/big"
	"testing"
)

func TestMandelbrotHigh(t *testing.T) {
	MandelbrotValueHigh(&LargeComplex{big.NewFloat(0.5), big.NewFloat(0.5)}, 100)
}
