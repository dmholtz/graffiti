package shortest_path_test

import (
	"math"
	"math/rand"
	"testing"

	sp "github.com/dmholtz/graffiti/algorithms/shortest_path"
	g "github.com/dmholtz/graffiti/graph"
)

// Differential testing: Compare Bidirectional Dijkstra's output of random searches with Dijkstra's Algorithm's output
func TestBidirectionalDijkstra(t *testing.T) {
	aag := loadAdjacencyArrayFromGob[g.GeoPoint, g.WeightedHalfEdge[int]](defaultGraphFile) // aag is a undirected graph

	dijkstraPqPops, biDijkstraPqPops := 0, 0

	n := 2000 // number of random tests
	t.Logf("Compare %d random searches of bidirectional Dijkstra against textbook Dijkstra.\n", n)
	for i := 0; i < n; i++ {
		source := rand.Intn(aag.NodeCount())
		target := rand.Intn(aag.NodeCount())

		dijkstraResult := sp.Dijkstra[g.GeoPoint, g.WeightedHalfEdge[int], int](aag, source, target, false)
		biDijkstraResult := sp.BidirectionalDijkstra[g.GeoPoint, g.WeightedHalfEdge[int]](aag, aag, source, target, false, math.MaxInt)

		if dijkstraResult.Length != biDijkstraResult.Length {
			t.Errorf("[Path(source=%d, target=%d)]: Different lengths found: Dijkstra=%d, BiDijkstra=%d", source, target, dijkstraResult.Length, biDijkstraResult.Length)
			return
		}

		// maintain performance indicators
		dijkstraPqPops += dijkstraResult.PqPops
		biDijkstraPqPops += biDijkstraResult.PqPops
	}
	t.Logf("Avgerage number of Pop() operations on priority queue: %d (Dijkstra), %d (BiDijkstra)", dijkstraPqPops/n, biDijkstraPqPops/n)
}
