package shortest_path_test

import (
	"testing"

	sp "github.com/dmholtz/graffiti/algorithms/shortest_path"
	h "github.com/dmholtz/graffiti/examples/heuristics"
	g "github.com/dmholtz/graffiti/graph"
)

// Differential testing: Compare the output of A* with HaversineHeuristic with Dijkstra's algorithm.
func TestAStar(t *testing.T) {
	aag := loadAdjacencyArrayFromGob[g.GeoPoint, g.WeightedHalfEdge[int]](defaultGraphFile) // aag is a undirected graph

	// Initialize the heurisitc
	havHeuristic := h.NewHaversineHeuristic[g.WeightedHalfEdge[int]](aag)

	testedRouter := sp.AStarRouter[g.GeoPoint, g.WeightedHalfEdge[int], int]{Graph: aag, Heuristic: havHeuristic}
	baselineRouter := sp.DijkstraRouter[g.GeoPoint, g.WeightedHalfEdge[int], int]{Graph: aag}

	DifferentialTesting(t, testedRouter, baselineRouter, aag.NodeCount())
}
