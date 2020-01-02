package traverse

import (
	"sort"

	"git.fractalqb.de/fractalqb/groph"
	iutil "git.fractalqb.de/fractalqb/groph/internal/util"
)

type Traversal struct {
	g    groph.RGraph
	mem  []groph.VIdx
	tail int
	// If not nil SortBy is used to sort the neighbours v of node u. SortBy
	// returns true if the edge (u,v1) is less than (u,v2).
	SortBy  func(u, v1, v2 groph.VIdx) bool
	Visited groph.BitSet
}

func NewTraversal(g groph.RGraph) *Traversal {
	res := &Traversal{
		g:       g,
		Visited: groph.NewBitSet(g.Order()),
	}
	return res
}

func (df *Traversal) Reset(g groph.RGraph) {
	df.g = g
	df.Visited = iutil.U64Slice(df.Visited, groph.BitSetWords(g.Order()))
}

func (df *Traversal) push(v groph.VIdx) {
	df.mem = append(df.mem, v)
}

func (df *Traversal) pop() (res groph.VIdx) {
	l := len(df.mem) - 1
	res = df.mem[l]
	df.mem = df.mem[:l]
	return res
}

func (df *Traversal) take() (res groph.VIdx) {
	res = df.mem[df.tail]
	df.tail++
	return res
}

func (df *Traversal) Depth1stAt(start groph.VIdx, do groph.VisitVertex) int {
	if df.Visited.Get(start) {
		return 0
	}
	if df.mem != nil {
		df.mem = df.mem[:0]
	}
	df.push(start)
	df.Visited.Set(start)
	count := 0
	for len(df.mem) > 0 {
		start = df.pop()
		do(start)
		count++
		sortStart := len(df.mem)
		groph.EachOutgoing(df.g, start, func(n groph.VIdx) {
			if !df.Visited.Get(n) {
				df.push(n)
				df.Visited.Set(n)
			}
		})
		if df.SortBy != nil {
			sort.Slice(df.mem[sortStart:], func(v1, v2 int) bool {
				return !df.SortBy(start, v1, v2)
			})
		}
	}
	return count
}

func (df *Traversal) Depth1st(do func(n groph.VIdx, cluster int)) {
	cluster := 0
	cdo := func(n groph.VIdx) { do(n, cluster) }
	count := df.Depth1stAt(0, cdo)
	for count < df.g.Order() {
		cluster++
		start := df.Visited.FirstUnset()
		count += df.Depth1stAt(start, cdo)
	}
}

func (df *Traversal) Breadth1stAt(start groph.VIdx, do groph.VisitVertex) int {
	if df.Visited.Get(start) {
		return 0
	}
	if df.mem != nil {
		df.mem = df.mem[:0]
	}
	df.tail = 0
	df.push(start)
	df.Visited.Set(start)
	count := 0
	for df.tail < len(df.mem) {
		start = df.take()
		do(start)
		count++
		sortStart := len(df.mem)
		groph.EachOutgoing(df.g, start, func(n groph.VIdx) {
			if !df.Visited.Get(n) {
				df.push(n)
				df.Visited.Set(n)
			}
		})
		if df.SortBy != nil {
			sort.Slice(df.mem[sortStart:], func(v1, v2 int) bool {
				return df.SortBy(start, v1, v2)
			})
		}
	}
	return count
}

func (df *Traversal) Breadth1st(do func(n groph.VIdx, cluster int)) {
	cluster := 0
	cdo := func(n groph.VIdx) { do(n, cluster) }
	count := df.Breadth1stAt(0, cdo)
	for count < df.g.Order() {
		cluster++
		start := df.Visited.FirstUnset()
		count += df.Breadth1stAt(start, cdo)
	}
}
