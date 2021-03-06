package util

import (
	"errors"
	"fmt"

	"git.fractalqb.de/fractalqb/groph"
	iutil "git.fractalqb.de/fractalqb/groph/internal/util"
)

// ClipError is returned when the order of src and dst in Cp*Weights does not
// match. If err is < 0 there were -err vertices ignored from src. If err > 0
// then err vertices in dst were not covered.
type ClipError int

func (err ClipError) Error() string {
	if err < 0 {
		return fmt.Sprintf("%d vertices ignored from source", -err)
	}
	return fmt.Sprintf("%d vertices not covered in destination", err)
}

// Clipped returns the clipping if err is a ClipError. Otherwise it returns 0.
func Clipped(err error) int {
	if ce, ok := err.(ClipError); ok {
		return int(ce)
	}
	return 0
}

// SrcDstError is returned when an error occurred in the src or dst graph during
// Cp*Weights.
type SrcDstError struct {
	Src error
	Dst error
}

func (err SrcDstError) Error() string {
	if err.Src != nil {
		if err.Dst != nil {
			return fmt.Sprintf("source: %s; dest: %s", err.Src, err.Dst)
		}
		return fmt.Sprintf("source: %s", err.Src)
	} else if err.Dst != nil {
		return fmt.Sprintf("dest: %s", err.Dst)
	}
	return "unspecific error"
}

// MustCp panics when err is not nil. Otherwiese it returns g.
func MustCp(g groph.WGraph, err error) groph.WGraph {
	if err != nil {
		panic(err)
	}
	return g
}

// CpWeights copies the edge weights from one graph to another.
// Vertices are identified by their index, i.e. the user has to take care of
// the vertex order. If the number of vertices in the graph differs the smaller
// graph determines how many edge weights are copied.
func CpWeights(dst groph.WGraph, src groph.RGraph) (dstout groph.WGraph, err error) {
	sz := dst.Order()
	if src.Order() < sz {
		sz = src.Order()
	}
	if udst, ok := dst.(groph.WUndirected); ok {
		if usrc, ok := src.(groph.RUndirected); ok {
			for i := 0; i < sz; i++ {
				udst.SetWeightU(i, i, usrc.WeightU(i, i))
				for j := 0; j < i; j++ {
					udst.SetWeightU(i, j, usrc.WeightU(i, j))
				}
			}
		} else {
			return dst, errors.New("cannot copy from directed to undirected graph")
		}
	} else if usrc, ok := src.(groph.RUndirected); ok {
		for i := 0; i < sz; i++ {
			dst.SetWeight(i, i, usrc.WeightU(i, i))
			for j := 0; j < i; j++ {
				w := usrc.WeightU(i, j)
				dst.SetWeight(i, j, w)
				dst.SetWeight(j, i, w)
			}
		}
	} else {
		for i := 0; i < sz; i++ {
			for j := 0; j < sz; j++ {
				dst.SetWeight(i, j, src.Weight(i, j))
			}
		}
	}
	sderr := SrcDstError{Src: iutil.ErrState(src), Dst: iutil.ErrState(dst)}
	if sderr.Src != nil || sderr.Dst != nil {
		return dst, sderr
	}
	vnd := dst.Order() - src.Order()
	if vnd == 0 {
		return dst, nil
	}
	return dst, ClipError(vnd)
}

// CpXWeights “transfers” the edge weights from src Graph to dst Graph
// with the same vertex restrictions as CpWeights. CpXWeights applies
// the transformation function xf() to each edge weight.
//
// Panic of xf will be recovered and returned as error.
func CpXWeights(
	dst groph.WGraph,
	src groph.RGraph,
	xf func(in interface{}) interface{},
) (dstout groph.WGraph, err error) {
	defer func() {
		if p := recover(); p != nil {
			switch e := p.(type) {
			case error:
				err = e
			default:
				err = fmt.Errorf("panic: %v", p)
			}
		}
	}()
	sz := dst.Order()
	if src.Order() < sz {
		sz = src.Order()
	}
	var w interface{}
	if udst, ok := dst.(groph.WUndirected); ok {
		if usrc, ok := src.(groph.RUndirected); ok {
			for i := 0; i < sz; i++ {
				w = usrc.WeightU(i, i)
				udst.SetWeightU(i, i, xf(w))
				for j := 0; j < i; j++ {
					w = usrc.WeightU(i, j)
					udst.SetWeightU(i, j, xf(w))
				}
			}
		} else {
			return dst, errors.New("cannot copy from directed to undirected graph")
		}
	} else if usrc, ok := src.(groph.RUndirected); ok {
		for i := 0; i < sz; i++ {
			w = usrc.WeightU(i, i)
			dst.SetWeight(i, i, xf(w))
			for j := 0; j < i; j++ {
				w := xf(usrc.WeightU(i, j))
				dst.SetWeight(i, j, w)
				dst.SetWeight(j, i, w)
			}
		}
	} else {
		for i := 0; i < sz; i++ {
			for j := 0; j < sz; j++ {
				w = src.Weight(i, j)
				dst.SetWeight(i, j, xf(w))
			}
		}
	}
	vnd := dst.Order() - src.Order()
	if vnd == 0 {
		return dst, nil
	}
	return dst, ClipError(vnd)
}
