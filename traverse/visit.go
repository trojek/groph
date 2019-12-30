package traverse

import (
	"git.fractalqb.de/fractalqb/groph"
)

// EachOutgoing calls onDest on each node d that is a neighbour of 'from' in
// graph g. Vertex d is a neighbour of from, iff g contains the edge (d,from).
//
// For undirected graphs that are no NeighbourListers EachNeighbour
// guarantees to call WeightU with v ≥ u to detect neighbours.
func EachOutgoing(g groph.RGraph, from groph.VIdx, onDest groph.VisitVertex) {
	switch gi := g.(type) {
	case groph.OutLister:
		gi.EachOutgoing(from, onDest)
	case groph.RUndirected:
		vno := gi.VertexNo()
		n := groph.V0
		for n < from {
			if w := gi.WeightU(from, n); w != nil {
				onDest(n)
			}
			n++
		}
		for n < vno {
			if w := gi.WeightU(n, from); w != nil {
				onDest(n)
			}
			n++
		}
	default:
		vno := g.VertexNo()
		for n := groph.V0; n < vno; n++ {
			if w := g.Weight(from, n); w != nil {
				onDest(n)
			}
		}
	}
}

// EachIncoming calls onSource on each node s that is a neighbour of 'to' in
// graph g. Vertex s is a neighbour of to, iff g contains the edge (s,to).
//
// For undirected graphs that are no NeighbourListers EachNeighbour
// guarantees to call WeightU with v ≥ u to detect neighbours.
func EachIncoming(g groph.RGraph, to groph.VIdx, onSource groph.VisitVertex) {
	switch gi := g.(type) {
	case groph.InLister:
		gi.EachIncoming(to, onSource)
	case groph.RUndirected:
		vno := gi.VertexNo()
		n := groph.V0
		for n < to {
			if w := gi.WeightU(to, n); w != nil {
				onSource(n)
			}
			n++
		}
		for n < vno {
			if w := gi.WeightU(n, to); w != nil {
				onSource(n)
			}
			n++
		}
	default:
		vno := g.VertexNo()
		for n := groph.V0; n < vno; n++ {
			if w := g.Weight(n, to); w != nil {
				onSource(n)
			}
		}
	}
}
