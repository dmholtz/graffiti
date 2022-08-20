package shortest_path_test

import (
	"encoding/gob"
	"fmt"
	"math/rand"
	"os"
	"testing"

	sp "github.com/dmholtz/graffiti/algorithms/shortest_path"
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

// DifferentialTesting compares random searches of the testedRouter with the baselineRouter and asserts same lengths.
// To obtain deterministic results, the Router's weight type is restricted to int.
func DifferentialTesting(t *testing.T, testedRouter, baselineRouter sp.Router[int], nodeCount int) {
	baselineName := "baseline"
	baselineStringer, ok := baselineRouter.(fmt.Stringer)
	if ok {
		baselineName = baselineStringer.String()
	}
	testedName := "tested"
	testedStringer, ok := testedRouter.(fmt.Stringer)
	if ok {
		testedName = testedStringer.String()
	}

	t.Logf("Compare %d random searches of %s with %s.\n", NUMBER_OF_RANDOM_TESTS, testedName, baselineName)

	rand.Seed(1) // normalize the seed
	testedPqPops, baselinePqPops := 0, 0
	for i := 0; i < NUMBER_OF_RANDOM_TESTS; i++ {
		source := rand.Intn(nodeCount)
		target := rand.Intn(nodeCount)

		testedResult := testedRouter.Route(source, target, false)
		baselineResult := baselineRouter.Route(source, target, false)

		if testedResult.Length != baselineResult.Length {
			t.Errorf("[Path(source=%d, target=%d)]: Different lengths found: [%s]=%d, [%s]=%d", source, target, testedName, testedResult.Length, baselineName, baselineResult.Length)
			return
		}

		// maintain performance indicators
		testedPqPops += testedResult.PqPops
		baselinePqPops += baselineResult.PqPops
	}
	testedPqPops = testedPqPops / NUMBER_OF_RANDOM_TESTS
	baselinePqPops = baselinePqPops / NUMBER_OF_RANDOM_TESTS
	t.Logf("Avgerage number of Pop() operations on priority queue: %d (%s), %d (%s)", testedPqPops, testedName, baselinePqPops, baselineName)
}
