package shortest_path_test

import (
	"math"
	"testing"

	sp "github.com/dmholtz/graffiti/algorithms/shortest_path"
	g "github.com/dmholtz/graffiti/graph"
)

// Differential testing: Compare the output of bidirectional ArcFlagDijkstra with unidirectional ArcFlagDijkstra
func TestBidirectionalArcFlagDijkstra(t *testing.T) {
	faag := loadAdjacencyArrayFromGob[g.PartGeoPoint, g.FlaggedHalfEdge[int, uint64]](arcflag64) // faag is a undirected graph

	testedRouter := sp.BidirectionalArcFlagRouter[g.PartGeoPoint, g.FlaggedHalfEdge[int, uint64], int]{Graph: faag, Transpose: faag, MaxInitializerValue: math.MaxInt}
	baselineRouter := sp.ArcFlagRouter[g.PartGeoPoint, g.FlaggedHalfEdge[int, uint64], int]{Graph: faag}

	DifferentialTesting(t, testedRouter, baselineRouter, faag.NodeCount())
}
