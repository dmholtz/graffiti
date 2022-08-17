package shortest_path_test

import (
	"math/rand"
	"testing"

	sp "github.com/dmholtz/graffiti/algorithms/shortest_path"
	g "github.com/dmholtz/graffiti/graph"
	h "github.com/dmholtz/graffiti/samples/heuristics"
)

// Differential testing: Compare the output of A* with HaversineHeuristic with Dijkstra's algorithm.
func TestAStar(t *testing.T) {
	aag := loadAdjacencyArrayFromGob[g.GeoPoint, g.WeightedHalfEdge[int]](defaultGraphFile) // aag is a undirected graph

	// Initialize the heurisitc
	havHeurisitc := h.NewHaversineHeuristic[g.WeightedHalfEdge[int], int](aag)

	t.Logf("Compare %d random searches of A* with Dijkstra's algorithm.\n", NUMBER_OF_RANDOM_TESTS)
	aStarPqPops, dijkstraPqPops := 0, 0
	for i := 0; i < NUMBER_OF_RANDOM_TESTS; i++ {
		source := rand.Intn(aag.NodeCount())
		target := rand.Intn(aag.NodeCount())

		aStarResult := sp.AStar[g.GeoPoint, g.WeightedHalfEdge[int], int](aag, havHeurisitc, source, target, false)
		dijkstraResult := sp.Dijkstra[g.GeoPoint, g.WeightedHalfEdge[int], int](aag, source, target, false)

		if aStarResult.Length != dijkstraResult.Length {
			t.Errorf("[Path(source=%d, target=%d)]: Different lengths found: A*=%d, Dijkstra=%d\n", source, target, aStarResult.Length, dijkstraResult.Length)
			return
		}

		// maintain PQ-pops as performance indicators
		aStarPqPops += aStarResult.PqPops
		dijkstraPqPops += dijkstraResult.PqPops
	}
	aStarPqPops, dijkstraPqPops = aStarPqPops/NUMBER_OF_RANDOM_TESTS, dijkstraPqPops/NUMBER_OF_RANDOM_TESTS
	t.Logf("Avgerage number of PQ.Pop() operations: %d (A*), %d (Dijkstra)\n", aStarPqPops, dijkstraPqPops)
}
