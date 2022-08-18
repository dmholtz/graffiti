package shortest_path_test

import (
	"math"
	"math/rand"
	"testing"

	sp "github.com/dmholtz/graffiti/algorithms/shortest_path"
	g "github.com/dmholtz/graffiti/graph"
	h "github.com/dmholtz/graffiti/samples/heuristics"
)

// Differential testing: Compare the output of bidirectional A* with HaversineHeuristic with unidirectional A.
func TestBidirectionalAStar(t *testing.T) {
	aag := loadAdjacencyArrayFromGob[g.GeoPoint, g.WeightedHalfEdge[int]](defaultGraphFile) // aag is a undirected graph

	// Initialize the heurisitc
	havForwardHeuristic := h.NewHaversineHeuristic[g.WeightedHalfEdge[int], int](aag)
	havBackwardHeuristic := h.NewHaversineHeuristic[g.WeightedHalfEdge[int], int](aag)

	t.Logf("Compare %d random searches of unidirectional A* with bidirectional A*.\n", NUMBER_OF_RANDOM_TESTS)
	aStarPqPops, biAStarPqPops := 0, 0
	rand.Seed(1) // reset seed to default
	for i := 0; i < NUMBER_OF_RANDOM_TESTS; i++ {
		source := rand.Intn(aag.NodeCount())
		target := rand.Intn(aag.NodeCount())

		biAStarResult := sp.BidirectionalAStar[g.GeoPoint, g.WeightedHalfEdge[int], int](aag, aag, havForwardHeuristic, havBackwardHeuristic, source, target, false, math.MaxInt)
		aStarResult := sp.AStar[g.GeoPoint, g.WeightedHalfEdge[int], int](aag, havForwardHeuristic, source, target, false)

		if aStarResult.Length != biAStarResult.Length {
			t.Errorf("[Path(source=%d, target=%d)]: Different lengths found: A*=%d, bi-A*=%d\n", source, target, aStarResult.Length, biAStarResult.Length)
			return
		}

		// maintain PQ-pops as performance indicators
		aStarPqPops += aStarResult.PqPops
		biAStarPqPops += biAStarResult.PqPops
	}
	aStarPqPops, biAStarPqPops = aStarPqPops/NUMBER_OF_RANDOM_TESTS, biAStarPqPops/NUMBER_OF_RANDOM_TESTS
	t.Logf("Avgerage number of PQ.Pop() operations: %d (A*), %d (bi-A*)\n", aStarPqPops, biAStarPqPops)
}
