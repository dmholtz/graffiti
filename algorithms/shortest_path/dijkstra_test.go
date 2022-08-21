package shortest_path_test

import (
	"math"
	"math/rand"
	"testing"

	sp "github.com/dmholtz/graffiti/algorithms/shortest_path"
	g "github.com/dmholtz/graffiti/graph"
)

// Differential testing: Compare the output of one-to-all Dijkstra to one-to-one Dijkstra
func TestOneToAllDijkstra(t *testing.T) {
	aag := loadAdjacencyArrayFromGob[g.GeoPoint, g.WeightedHalfEdge[int]](defaultGraphFile)

	t.Logf("Compare %d random searches of one-to-all Dijkstra with one-to-one Dijkstra.\n", NUMBER_OF_RANDOM_TESTS)
	dijkstraRouter := sp.DijkstraRouter[g.GeoPoint, g.WeightedHalfEdge[int], int]{Graph: aag}
	for i := 0; i < int(math.Sqrt(NUMBER_OF_RANDOM_TESTS)); i++ {
		source := rand.Intn(aag.NodeCount())

		one2AllResult := sp.DijkstraOneToAll[g.GeoPoint, g.WeightedHalfEdge[int], int](aag, source)

		for j := 0; j < int(math.Sqrt(NUMBER_OF_RANDOM_TESTS)); j++ {
			target := rand.Intn(aag.NodeCount())
			one2OneResult := dijkstraRouter.Route(source, target, false)

			if one2AllResult.Lengths[target] != one2OneResult.Length {
				t.Errorf("[Path(source=%d, target=%d)]: Different lengths found: one-to-all Dijkstra=%d, one-to-one Dijkstra=%d\n", source, target, one2AllResult.Lengths[target], one2OneResult.Length)
				return
			}
		}
	}
}
