package shortest_path_test

import (
	"math"
	"math/rand"
	"testing"

	sp "github.com/dmholtz/graffiti/algorithms/shortest_path"
	g "github.com/dmholtz/graffiti/graph"
)

// Differential testing: Compare the output of ALT (A*, Landmarks and Triangular Inequalities) with Dijkstra's algorithm.
func TestAlt(t *testing.T) {
	aag := loadAdjacencyArrayFromGob[g.GeoPoint, g.WeightedHalfEdge[int]](defaultGraphFile) // aag is a undirected graph

	// ALT preprocessing
	landmarks := sp.UniformLandmarks[g.GeoPoint, g.WeightedHalfEdge[int]](aag, 16)
	altHeuristic := sp.NewAltHeurisitc[g.GeoPoint, g.WeightedHalfEdge[int], int](aag, aag, landmarks)

	t.Logf("Compare %d random searches of ALT with Dijkstra's algorithm.\n", NUMBER_OF_RANDOM_TESTS)
	altPqPops, dijkstraPqPops := 0, 0
	rand.Seed(1) // reset seed to default
	for i := 0; i < NUMBER_OF_RANDOM_TESTS; i++ {
		source := rand.Intn(aag.NodeCount())
		target := rand.Intn(aag.NodeCount())

		altResult := sp.AStar[g.GeoPoint, g.WeightedHalfEdge[int], int](aag, altHeuristic, source, target, false)
		dijkstraResult := sp.Dijkstra[g.GeoPoint, g.WeightedHalfEdge[int], int](aag, source, target, false)

		if altResult.Length != dijkstraResult.Length {
			t.Errorf("[Path(source=%d, target=%d)]: Different lengths found: ALT=%d, Dijkstra=%d\n", source, target, altResult.Length, dijkstraResult.Length)
			return
		}

		// maintain PQ-pops as performance indicators
		altPqPops += altResult.PqPops
		dijkstraPqPops += dijkstraResult.PqPops
	}
	altPqPops, dijkstraPqPops = altPqPops/NUMBER_OF_RANDOM_TESTS, dijkstraPqPops/NUMBER_OF_RANDOM_TESTS
	t.Logf("Avgerage number of PQ.Pop() operations: %d (ALT), %d (Dijkstra)\n", altPqPops, dijkstraPqPops)
}

// Differential testing: Compare the output of bidirectional ALT with with unidirectional ALT.
func TestBidirectionalAlt(t *testing.T) {
	aag := loadAdjacencyArrayFromGob[g.GeoPoint, g.WeightedHalfEdge[int]](defaultGraphFile) // aag is a undirected graph

	// Initialize the heurisitc
	// ALT preprocessing
	landmarks := sp.UniformLandmarks[g.GeoPoint, g.WeightedHalfEdge[int]](aag, 16)
	altForwardHeuristic := sp.NewAltHeurisitc[g.GeoPoint, g.WeightedHalfEdge[int], int](aag, aag, landmarks)
	altBackwardHeuristic := sp.NewAltHeurisitc[g.GeoPoint, g.WeightedHalfEdge[int], int](aag, aag, landmarks) // need a separate object since heuristic is stateful

	t.Logf("Compare %d random searches of unidirectional ALT with bidirectional ALT.\n", NUMBER_OF_RANDOM_TESTS)
	altPqPops, biAltPqPops := 0, 0
	rand.Seed(1) // reset seed to default
	for i := 0; i < NUMBER_OF_RANDOM_TESTS; i++ {
		source := rand.Intn(aag.NodeCount())
		target := rand.Intn(aag.NodeCount())

		biAltResult := sp.BidirectionalAStar[g.GeoPoint, g.WeightedHalfEdge[int], int](aag, aag, altForwardHeuristic, altBackwardHeuristic, source, target, false, math.MaxInt)
		altResult := sp.AStar[g.GeoPoint, g.WeightedHalfEdge[int], int](aag, altForwardHeuristic, source, target, false)

		if altResult.Length != biAltResult.Length {
			t.Errorf("[Path(source=%d, target=%d)]: Different lengths found: ALT=%d, bi-ALT=%d\n", source, target, altResult.Length, biAltResult.Length)
			return
		}

		// maintain PQ-pops as performance indicators
		altPqPops += altResult.PqPops
		biAltPqPops += biAltResult.PqPops
	}
	altPqPops, biAltPqPops = altPqPops/NUMBER_OF_RANDOM_TESTS, biAltPqPops/NUMBER_OF_RANDOM_TESTS
	t.Logf("Avgerage number of PQ.Pop() operations: %d (ALT), %d (bi-ALT)\n", altPqPops, biAltPqPops)
}
