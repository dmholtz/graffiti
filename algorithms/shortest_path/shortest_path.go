package shortest_path

import g "github.com/dmholtz/graffiti/graph"

// Encapsulates the output of a shortest path algorithm.
type ShortestPathResult[W g.Weight] struct {
	// Length stores the length of the shortest path between the source and the target node.
	// The value of Length is set to -1 iff such a path does not exist.
	Length W
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

// ShortestPathTreeNode encapsulates the directed, acyclic search graph (tree) spanned by Dijkstra's algorithm.
type ShortestPathTreeNode struct {
	// Id is the root element of this search tree and refers to a Node ID in the graph on which Dijkstra's algorithm has been applied.
	Id g.NodeId
	// Children is a list of pointers to subtrees, each referring to a node that is reached on a shortest path from the source node via this node (ID=Id).
	Children []*ShortestPathTreeNode
	// Visited is a boolean flag to mark nodes that have already been visited while traversing the search graph.
	// Caution: It is not possible to make use of this flag for more than one traversal (without resetting it).
	Visited bool
}
