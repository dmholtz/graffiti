package shortest_path

import g "github.com/dmholtz/graffiti/graph"

// Router is the interface that wraps the basic Route method of a shortest path algorithm.
type Router[W g.Weight] interface {
	// Route computes the shortest path from the source node to the target node of the underlying graph.
	//
	// The search space of the algorithm's execution is reported iff recordSearchSpace is true.
	// Note that recording the search space will decrease the performance of Route significantly.
	Route(source, target g.NodeId, recordSearchSpace bool) ShortestPathResult[W]
}
