package shortest_path

import g "github.com/dmholtz/graffiti/graph"

// Encapsulates the output of a shortest path algorithm.
type ShortestPathResult[W g.Weight] struct {
	// Distance stores the shortest path distance between the source and the target node.
	// The value of Distance is set to -1 iff such a path does not exist.
	Distance W
	// Path is a slice of NodeId values that describe the shortest path from the source to the target node.
	// The slice is empty iff such a path does not exist.
	Path []g.NodeId
	// PqPops reports the number of Pop() operations on the priority queue during the shortest path computation.
	PqPops int
	// SearchSpace reports the search space of the algorithm by enumerating all processed node IDs.
	// The slice is ordered by the time the nodes were settled.
	// Shortest path algorithms should usually not record the search space for performance reasons.
	// SearchSpace is 'nil' iff the algorithm has been instructed not to record the search space.
	SearchSpace []g.NodeId
}
