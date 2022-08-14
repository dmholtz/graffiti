package shortest_path_test

import (
	"math"
	"math/rand"
	"testing"

	sp "github.com/dmholtz/graffiti/algorithms/shortest_path"
	g "github.com/dmholtz/graffiti/graph"
)

// Differential testing: Compare the output of bidirectional ArcFlagDijkstra with unidirectional ArcFlagDijkstra
func TestBidirectionalArcFlagDijkstra(t *testing.T) {
	faag := loadAdjacencyArrayFromGob[g.PartGeoPoint, g.FlaggedHalfEdge[int, uint64]](arcflag64) // faag is a undirected graph

	t.Logf("Compare %d random searches of bidirectional ArcFlagDijkstra with unidirectional ArcFlagDijkstra.\n", NUMBER_OF_RANDOM_TESTS)
	biPqPops, uniPqPops := 0, 0
	for i := 0; i < NUMBER_OF_RANDOM_TESTS; i++ {
		source := rand.Intn(faag.NodeCount())
		target := rand.Intn(faag.NodeCount())

		biResult := sp.BidirectionalArcFlagDijkstra[g.PartGeoPoint, g.FlaggedHalfEdge[int, uint64]](faag, faag, source, target, false, math.MaxInt)
		uniResult := sp.ArcFlagDijkstra[g.PartGeoPoint, g.FlaggedHalfEdge[int, uint64], int](faag, source, target, false)

		if biResult.Length != uniResult.Length {
			t.Errorf("[Path(source=%d, target=%d)]: Different lengths found: Bi-ArcFlagDijkstra=%d, ArcFlagDijkstra=%d\n", source, target, biResult.Length, uniResult.Length)
			return
		}

		// maintain PQ-pops as performance indicators
		biPqPops += biResult.PqPops
		uniPqPops += uniResult.PqPops
	}
	biPqPops, uniPqPops = biPqPops/NUMBER_OF_RANDOM_TESTS, uniPqPops/NUMBER_OF_RANDOM_TESTS
	t.Logf("Avgerage number of PQ.Pop() operations: %d (Bi-ArcFlagDijkstra), %d (ArcFlagDijkstra)\n", biPqPops, uniPqPops)
}
