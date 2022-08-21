package shortest_path_test

import (
	"testing"

	sp "github.com/dmholtz/graffiti/algorithms/shortest_path"
	g "github.com/dmholtz/graffiti/graph"
)

// Differential testing: Compare the output of ArcFlagDijkstra with textbook Dijkstra
func TestTowLevelArcFlagDijkstra(t *testing.T) {
	faag := loadAdjacencyArrayFromGob[g.TwoLevelPartGeoPoint, g.TwoLevelFlaggedHalfEdge[int, uint32, uint32]](arcflag32_32) // faag is a undirected graph

	testedRouter := sp.TwoLevelArcFlagRouter[g.TwoLevelPartGeoPoint, g.TwoLevelFlaggedHalfEdge[int, uint32, uint32], int]{Graph: faag}
	baselineRouter := sp.DijkstraRouter[g.TwoLevelPartGeoPoint, g.TwoLevelFlaggedHalfEdge[int, uint32, uint32], int]{Graph: faag}

	DifferentialTesting(t, testedRouter, baselineRouter, faag.NodeCount())
}
