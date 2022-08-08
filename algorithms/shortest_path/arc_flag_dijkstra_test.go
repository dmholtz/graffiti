package shortest_path

import (
	"encoding/gob"
	"fmt"
	"math/rand"
	"os"
	"testing"

	g "github.com/dmholtz/graffiti/graph"
)

const flaggedGraphFile = "testdata/geo_graph_arcflags_7k.gob"

func loadFlaggedGraph(filename string) g.Graph[g.PartGeoPoint, g.FlaggedHalfEdge[int, uint64]] {
	var faag g.AdjacencyArrayGraph[g.PartGeoPoint, g.FlaggedHalfEdge[int, uint64]]

	gobGraphFile, err := os.Open(filename)
	defer gobGraphFile.Close()
	if err != nil {
		panic(fmt.Sprintf("Error while reading .gob file '%s'", filename))
	}

	dec := gob.NewDecoder(gobGraphFile)
	err = dec.Decode(&faag)
	if err != nil {
		panic(fmt.Sprintf("Error while decoding .gob file '%s", filename))
	}

	return &faag
}

// Differential testing: Compare ArcFlagDijkstra's output of random searches with Dijkstra's Algorithm's output
func TestArcFlagDijkstra(t *testing.T) {
	aag := loadGraph(defaultGraphFile)         // aag is a undirected graph
	faag := loadFlaggedGraph(flaggedGraphFile) // faag is a undirected graph

	dijkstraPqPops, arcFlagDijkstraPqPops := 0, 0

	n := 2000 // number of random tests
	t.Logf("Compare %d random searches of bidirectional Dijkstra against textbook Dijkstra.\n", n)
	for i := 0; i < n; i++ {
		source := rand.Intn(faag.NodeCount())
		target := rand.Intn(faag.NodeCount())

		dijkstraResult := Dijkstra[g.GeoPoint, g.WeightedHalfEdge[int], int](aag, source, target, false)
		arcFlagDijkstraResult := ArcFlagDijkstra[g.PartGeoPoint, g.FlaggedHalfEdge[int, uint64], int](faag, source, target, false)

		if dijkstraResult.Length != arcFlagDijkstraResult.Length {
			t.Errorf("[Path(source=%d, target=%d)]: Different lengths found: Dijkstra=%d, ArcFlagDijkstra=%d", source, target, dijkstraResult.Length, arcFlagDijkstraResult.Length)
			return
		}

		// maintain performance indicators
		dijkstraPqPops += dijkstraResult.PqPops
		arcFlagDijkstraPqPops += arcFlagDijkstraResult.PqPops
	}
	t.Logf("Avgerage number of Pop() operations on priority queue: %d (Dijkstra), %d (BiDijkstra)", dijkstraPqPops/n, arcFlagDijkstraPqPops/n)
}
