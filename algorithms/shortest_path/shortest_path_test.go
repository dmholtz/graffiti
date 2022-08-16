package shortest_path_test

import (
	"encoding/gob"
	"fmt"
	"os"

	g "github.com/dmholtz/graffiti/graph"
)

// Package shortest_path_test contains utility functions and testing constants

// Constants
const NUMBER_OF_RANDOM_TESTS = 2000

// Path to .gob files
const defaultGraphFile = "testdata/geo_graph_7k.gob"
const arcflag64 = "testdata/geograph_arcflag_64_7k.gob"

// Deserializer
func loadAdjacencyArrayFromGob[N any, E g.IHalfEdge](filename string) *g.AdjacencyArrayGraph[N, E] {
	var faag g.AdjacencyArrayGraph[N, E]

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
