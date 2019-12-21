package tsp

import (
	"fmt"
	"testing"

	"git.fractalqb.de/fractalqb/groph"
	test "git.fractalqb.de/fractalqb/groph/internal"
)

func pathEq(p1, p2 []uint) (bool, string) {
	plen := uint(len(p1))
	if plen != uint(len(p2)) {
		return false, fmt.Sprintf("length differs: %d / %d",
			plen,
			uint(len(p2)))
	}
	s2 := uint(0)
	for s2 < plen {
		if p1[0] == p2[s2] {
			break
		}
		s2++
	}
	if s2 >= plen {
		return false, fmt.Sprintf("no start %d in p2=%v", p1[0], p2)
	}
	pidx, nidx := func(i uint) uint {
		if i == 0 {
			return plen - 1
		}
		return i - 1
	},
		func(i uint) uint {
			i++
			if i >= plen {
				i = 0
			}
			return i
		}
	s1 := uint(1)
	s2 = nidx(s2)
	if p1[s1] == p2[s2] {
		s1++
		s2 = nidx(s2)
		for s1 < plen {
			if p1[s1] != p2[s2] {
				return false, fmt.Sprintf("difference in pos %d / %d: %v %v",
					s1, s2,
					p1, p2)
			}
			s1++
			s2 = nidx(s2)
		}
	} else {
		s2 = pidx(pidx(s2))
		for s1 < plen {
			if p1[s1] != p2[s2] {
				return false, fmt.Sprintf("difference in pos %d / %d: %v %v",
					s1, s2,
					p1, p2)
			}
			s1++
			s2 = pidx(s2)
		}
	}
	return true, ""
}

func Test2OptDAgaintsGreedy(t *testing.T) {
	var points []test.Point
	var am *groph.AdjMxDf32
	for sz := uint(4); sz < 12; sz++ {
		points = test.RandomPoints(sz, points)
		am = groph.CpWeights(
			groph.NewAdjMxDf32(sz, am),
			groph.NewSliceNMeasure(points, test.Dist, false).Verify(),
		).(*groph.AdjMxDf32)
		gPath, gWeight := GreedyAdjMxDf32(am)
		tPath, tWeight := TwoOptf32(am)
		if tWeight/gWeight > 1.01 {
			t.Errorf("size %d: different path length: greedy=%f / 2-opt=%f",
				sz,
				gWeight,
				tWeight)
		}
		if ok, msg := pathEq(gPath, tPath); !ok {
			t.Logf("size %d: %s", sz, msg)
		}
	}
}

func Test2OptUAgaintsGreedy(t *testing.T) {
	var points []test.Point
	var am *groph.AdjMxUf32
	var dm *groph.AdjMxDf32
	for sz := uint(4); sz < 12; sz++ {
		points = test.RandomPoints(sz, points)
		am = groph.CpWeights(
			groph.NewAdjMxUf32(sz, am),
			groph.NewSliceNMeasure(points, test.Dist, false).Verify(),
		).(*groph.AdjMxUf32)
		dm = groph.CpWeights(
			groph.NewAdjMxDf32(sz, dm),
			groph.NewSliceNMeasure(points, test.Dist, false).Verify(),
		).(*groph.AdjMxDf32)
		gPath, gWeight := GreedyAdjMxDf32(dm)
		tPath, tWeight := TwoOptf32(am)
		if tWeight/gWeight > 1.052 {
			t.Errorf("size %d: different path length: greedy=%f / 2-opt=%f",
				sz,
				gWeight,
				tWeight)
		}
		if ok, msg := pathEq(gPath, tPath); !ok {
			t.Logf("size %d: %s", sz, msg)
		}
	}
}

const twoOptBenchSize = 120

func BenchmarkTsp2OptGenf32D(b *testing.B) {
	points := test.RandomPoints(twoOptBenchSize, nil)
	am := groph.CpWeights(
		groph.NewAdjMxDf32(twoOptBenchSize, nil),
		groph.NewSliceNMeasure(points, test.Dist, false).Verify(),
	).(*groph.AdjMxDf32)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		TwoOptf32(am)
	}
}

func BenchmarkTsp2OptGenf32U(b *testing.B) {
	points := test.RandomPoints(twoOptBenchSize, nil)
	am := groph.CpWeights(
		groph.NewAdjMxUf32(twoOptBenchSize, nil),
		groph.NewSliceNMeasure(points, test.Dist, false).Verify(),
	).(*groph.AdjMxUf32)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		TwoOptf32(am)
	}
}

// Benchmark showed that the proformance gain is not worth it
// func BenchmarkTsp2Optf32D(b *testing.B) {
// 	points := randomPoints(twoOptBenchSize, nil)
// 	am := CpWeights(
// 		NewAdjMxDf32(twoOptBenchSize, nil),
// 		NewSliceNMeasure(points, dist, false).Verify(),
// 	).(*AdjMxDf32)
// 	b.ResetTimer()
// 	for i := 0; i < b.N; i++ {
// 		am.Tsp2Opt()
// 	}
// }

// Benchmark showed that the proformance gain is not worth it
// func BenchmarkTsp2Optf32U(b *testing.B) {
// 	points := randomPoints(twoOptBenchSize, nil)
// 	am := CpWeights(
// 		NewAdjMxUf32(twoOptBenchSize, nil),
// 		NewSliceNMeasure(points, dist, false).Verify(),
// 	).(*AdjMxUf32)
// 	b.ResetTimer()
// 	for i := 0; i < b.N; i++ {
// 		am.Tsp2Opt()
// 	}
// }