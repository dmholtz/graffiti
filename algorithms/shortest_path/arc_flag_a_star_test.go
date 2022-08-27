package shortest_path_test

import (
	"testing"

	sp "github.com/dmholtz/graffiti/algorithms/shortest_path"
	g "github.com/dmholtz/graffiti/graph"
)

// Differential testing: Compare the output of the arcflag-ALT combi-router with Dijkstra's algorithm.
func TestArcFlagAlt(t *testing.T) {
	faag := loadAdjacencyArrayFromGob[g.PartGeoPoint, g.FlaggedHalfEdge[int, uint64]](arcflag64) // faag is a undirected graph

	// ALT preprocessing
	landmarks := sp.UniformLandmarks[g.PartGeoPoint, g.FlaggedHalfEdge[int, uint64]](faag, 16)
	altHeuristic := sp.NewAltHeurisitc[g.PartGeoPoint, g.FlaggedHalfEdge[int, uint64], int](faag, faag, landmarks)

	testedRouter := sp.ArcFlagAStarRouter[g.PartGeoPoint, g.FlaggedHalfEdge[int, uint64], int]{Graph: faag, Transpose: faag, Heuristic: altHeuristic}
	baselineRouter := sp.DijkstraRouter[g.PartGeoPoint, g.FlaggedHalfEdge[int, uint64], int]{Graph: faag}

	DifferentialTesting(t, testedRouter, baselineRouter, faag.NodeCount())
}
