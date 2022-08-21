package shortest_path_test

import (
	"math"
	"testing"

	sp "github.com/dmholtz/graffiti/algorithms/shortest_path"
	g "github.com/dmholtz/graffiti/graph"
)

// Differential testing: Compare Bidirectional Dijkstra's output of random searches with Dijkstra's Algorithm's output
func TestBidirectionalDijkstra(t *testing.T) {
	aag := loadAdjacencyArrayFromGob[g.GeoPoint, g.WeightedHalfEdge[int]](defaultGraphFile) // aag is a undirected graph

	testedRouter := sp.BiDijkstraRouter[g.GeoPoint, g.WeightedHalfEdge[int], int]{Graph: aag, Transpose: aag, MaxInitializerValue: math.MaxInt}
	baselineRouter := sp.DijkstraRouter[g.GeoPoint, g.WeightedHalfEdge[int], int]{Graph: aag}

	DifferentialTesting(t, testedRouter, baselineRouter, aag.NodeCount())
}
