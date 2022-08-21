package shortest_path_test

import (
	"math"
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

	testedRouter := sp.AStarRouter[g.GeoPoint, g.WeightedHalfEdge[int], int]{Graph: aag, Heuristic: altHeuristic}
	baselineRouter := sp.DijkstraRouter[g.GeoPoint, g.WeightedHalfEdge[int], int]{Graph: aag}

	DifferentialTesting(t, testedRouter, baselineRouter, aag.NodeCount())
}

// Differential testing: Compare the output of bidirectional ALT with with unidirectional ALT.
func TestBidirectionalAlt(t *testing.T) {
	aag := loadAdjacencyArrayFromGob[g.GeoPoint, g.WeightedHalfEdge[int]](defaultGraphFile) // aag is a undirected graph

	// Initialize the heurisitc
	// ALT preprocessing
	landmarks := sp.UniformLandmarks[g.GeoPoint, g.WeightedHalfEdge[int]](aag, 16)
	altForwardHeuristic := sp.NewAltHeurisitc[g.GeoPoint, g.WeightedHalfEdge[int], int](aag, aag, landmarks)
	altBackwardHeuristic := sp.NewAltHeurisitc[g.GeoPoint, g.WeightedHalfEdge[int], int](aag, aag, landmarks) // need a separate object since heuristic is stateful

	testedRouter := sp.BidirectionalAStarRouter[g.GeoPoint, g.WeightedHalfEdge[int], int]{Graph: aag, Transpose: aag, ForwardHeuristic: altForwardHeuristic, BackwardHeuristic: altBackwardHeuristic, MaxInitializerValue: math.MaxInt}
	baselineRouter := sp.AStarRouter[g.GeoPoint, g.WeightedHalfEdge[int], int]{Graph: aag, Heuristic: altForwardHeuristic}

	DifferentialTesting(t, testedRouter, baselineRouter, aag.NodeCount())
}
