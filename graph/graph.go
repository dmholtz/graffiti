package graph

// A node of a graph is identified by an nonnegative 32 bit integer. Depending on the application, negative node ids might represent errors.
type NodeId = int

// Simplest capability description of an outgoing edge without any annotations such as weight.
type IHalfEdge interface {
	// To() returns the head of the edge, i.e. the node id to which this edge points.
	To() NodeId
}

// The weight of an edge can be of any number type that supports addition and subtraction.
type Weight interface {
	int | float64
}

// Capabilities description of weighted half edges
type IWeightedHalfEdge[W Weight] interface {
	// IWeightedHalfEdge inherits all capabilities of IHalfEdge.
	IHalfEdge
	// Weight() returns the weight of type W associated with this edge.
	Weight() W
}

// Generic interface of a graph
type Graph[N any, E IHalfEdge] interface {
	// NodeCount() returns the number of nodes in the graph.
	NodeCount() int
	// EdgeCount() returns the number of edges in the graph.
	EdgeCount() int
	// GetNode(id) returns the node n of the graph, which is uniquely identified ID 'id'.
	// The method panics iff the graph does not contain a node with ID 'id'.
	GetNode(id NodeId) N
	// GetHalfEdgesFrom(id) returns a slice of edges, which leave the node with ID=id.
	// If the node referred by 'id' does not have any leaving edges, an empty slice is returned.
	// The method panics iff the graph does not contain a node with ID 'id'.
	GetHalfEdgesFrom(id NodeId) []E
}
