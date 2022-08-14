package shortest_path_test

import (
	"math/rand"
	"testing"

	sp "github.com/dmholtz/graffiti/algorithms/shortest_path"
	g "github.com/dmholtz/graffiti/graph"
)

// Differential testing: Compare the output of ArcFlagDijkstra with textbook Dijkstra
func TestArcFlagDijkstra(t *testing.T) {
	faag := loadAdjacencyArrayFromGob[g.PartGeoPoint, g.FlaggedHalfEdge[int, uint64]](arcflag64) // faag is a undirected graph

	t.Logf("Compare %d random searches of ArcFlagDijkstra with textbook Dijkstra.\n", NUMBER_OF_RANDOM_TESTS)
	biPqPops, uniPqPops := 0, 0
	for i := 0; i < NUMBER_OF_RANDOM_TESTS; i++ {
		source := rand.Intn(faag.NodeCount())
		target := rand.Intn(faag.NodeCount())

		arcFlagResult := sp.ArcFlagDijkstra[g.PartGeoPoint, g.FlaggedHalfEdge[int, uint64], int](faag, source, target, false)
		dijkstraResult := sp.Dijkstra[g.PartGeoPoint, g.FlaggedHalfEdge[int, uint64], int](faag, source, target, false)

		if arcFlagResult.Length != dijkstraResult.Length {
			t.Errorf("[Path(source=%d, target=%d)]: Different lengths found: ArcFlagDijkstra=%d, Dijkstra=%d\n", source, target, arcFlagResult.Length, dijkstraResult.Length)
			return
		}

		// maintain PQ-pops as performance indicators
		biPqPops += arcFlagResult.PqPops
		uniPqPops += dijkstraResult.PqPops
	}
	biPqPops, uniPqPops = biPqPops/NUMBER_OF_RANDOM_TESTS, uniPqPops/NUMBER_OF_RANDOM_TESTS
	t.Logf("Avgerage number of PQ.Pop() operations: %d (ArcFlagDijkstra), %d (Dijkstra)\n", biPqPops, uniPqPops)
}
