package shortest_path_test

import (
	"testing"

	sp "github.com/dmholtz/graffiti/algorithms/shortest_path"
	g "github.com/dmholtz/graffiti/graph"
)

// Differential testing: Compare the output of ArcFlagDijkstra with textbook Dijkstra
func TestArcFlagDijkstra(t *testing.T) {
	faag := loadAdjacencyArrayFromGob[g.PartGeoPoint, g.FlaggedHalfEdge[int, uint64]](arcflag64) // faag is a undirected graph

	testedRouter := sp.ArcFlagRouter[g.PartGeoPoint, g.FlaggedHalfEdge[int, uint64], int]{Graph: faag}
	baselineRouter := sp.DijkstraRouter[g.PartGeoPoint, g.FlaggedHalfEdge[int, uint64], int]{Graph: faag}

	DifferentialTesting(t, testedRouter, baselineRouter, faag.NodeCount())
}
