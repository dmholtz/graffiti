package shortest_path

import (
	"encoding/gob"
	"fmt"
	"os"
	"testing"

	g "github.com/dmholtz/graffiti/graph"
)

const defaultGraphFile = "testdata/geo_graph_7k.gob"

func loadGraph(filename string) g.Graph[g.GeoPoint, g.WeightedHalfEdge[int]] {
	var aag g.AdjacencyArrayGraph[g.GeoPoint, g.WeightedHalfEdge[int]]

	gobGraphFile, err := os.Open(filename)
	defer gobGraphFile.Close()
	if err != nil {
		panic(fmt.Sprintf("Error while reading .gob file '%s'", filename))
	}

	dec := gob.NewDecoder(gobGraphFile)
	err = dec.Decode(&aag)
	if err != nil {
		panic(fmt.Sprintf("Error while decoding .gob file '%s", filename))
	}

	return &aag
}

func TestDijkstra(t *testing.T) {
	aag := loadGraph(defaultGraphFile)
	t.Log(aag.NodeCount())
}
