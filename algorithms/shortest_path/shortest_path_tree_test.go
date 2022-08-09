package shortest_path

import (
	"testing"

	g "github.com/dmholtz/graffiti/graph"
)

func TestMultiplePathsInShortestPathTree(t *testing.T) {
	t.Parallel()

	// prepare the test case: diamond graph
	alg := &g.AdjacencyListGraph[struct{}, g.WeightedHalfEdge[int]]{}
	// insert four (empty) nodes
	alg.AppendNode(struct{}{})
	alg.AppendNode(struct{}{})
	alg.AppendNode(struct{}{})
	alg.AppendNode(struct{}{})
	// insert diamond edges with equal path length from node 0 to node 3
	alg.InsertHalfEdge(0, g.WeightedHalfEdge[int]{To_: 1, Weight_: 7})
	alg.InsertHalfEdge(0, g.WeightedHalfEdge[int]{To_: 2, Weight_: 8})
	alg.InsertHalfEdge(1, g.WeightedHalfEdge[int]{To_: 3, Weight_: 8})
	alg.InsertHalfEdge(2, g.WeightedHalfEdge[int]{To_: 3, Weight_: 7})

	tree := ShortestPathTree[struct{}, g.WeightedHalfEdge[int], int](alg, 0)

	for _, child1 := range tree.Children {
		//t.Logf("Level 1: %d", child1.Id)
		if len(child1.Children) < 1 {
			t.Error("No path to level 2")
		}
		for _, child2 := range child1.Children {
			//t.Logf("Level 2: %d", child2.Id)
			if child2.Id != 3 {
				t.Errorf("No path to node 3.")
			}
		}
	}
}
