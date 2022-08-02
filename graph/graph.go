package graph

// A node of a graph is identified by an nonnegative 32 bit integer. Depending on the application, negative node ids might represent errors.
type NodeId = int

// Simplest capability description of an outgoing edge without any annotations such as weight.
type IHalfEdge interface {
	To() NodeId // returns the nodeId to which the head of this half edge points to
}

// The weight of an edge can be of any number type that supports addition and subtraction.
type Weight interface {
	int | float64
}

// A weighted half edge combines the capabilities of an unweighted half edge and a query to retrieve the weight.
type IWeightedHalfEdge[W Weight] interface {
	IHalfEdge  // capabilities of an unweighted half edge
	Weight() W // returns the weight of the edge
}

// Generic interface of a graph
type Graph[N any, E IHalfEdge] interface {
	NodeCount() int
	EdgeCount() int
	GetNode(id NodeId) N
	GetHalfEdgesFrom(id NodeId) []E
}
