package groph

import (
	"math"
)

const I32Del = -2147483648 // min{ int32 }

var nan32 = float32(math.NaN())

func NaN32() float32 { return nan32 }

func IsNaN32(x float32) bool { return math.IsNaN(float64(x)) }

func errState(v interface{}) error {
	if es, ok := v.(interface{ ErrState() error }); ok {
		return es.ErrState()
	}
	return nil
}

// nSum computes the sum of the n 1st integers, i.e. 1+2+3+…+n
func nSum(n VIdx) VIdx { return n * (n + 1) / 2 }

func nSumRev(n VIdx) float64 {
	r := math.Sqrt(0.25 + 2*float64(n))
	return r - 0.5
}
