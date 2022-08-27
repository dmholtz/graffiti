package shortest_path_test

import (
	"math"
	"testing"

	sp "github.com/dmholtz/graffiti/algorithms/shortest_path"
	h "github.com/dmholtz/graffiti/examples/heuristics"
	g "github.com/dmholtz/graffiti/graph"
)

// Differential testing: Compare the output of bidirectional A* with HaversineHeuristic with unidirectional A.
func TestBidirectionalAStar(t *testing.T) {
	aag := loadAdjacencyArrayFromGob[g.GeoPoint, g.WeightedHalfEdge[int]](defaultGraphFile) // aag is a undirected graph

	// Initialize the heurisitc
	havForwardHeuristic := h.NewHaversineHeuristic[g.WeightedHalfEdge[int]](aag)
	havBackwardHeuristic := h.NewHaversineHeuristic[g.WeightedHalfEdge[int]](aag)

	testedRouter := sp.BidirectionalAStarRouter[g.GeoPoint, g.WeightedHalfEdge[int], int]{Graph: aag, Transpose: aag, ForwardHeuristic: havForwardHeuristic, BackwardHeuristic: havBackwardHeuristic, MaxInitializerValue: math.MaxInt}
	baselineRouter := sp.AStarRouter[g.GeoPoint, g.WeightedHalfEdge[int], int]{Graph: aag, Heuristic: havForwardHeuristic}

	DifferentialTesting(t, testedRouter, baselineRouter, aag.NodeCount())
}
